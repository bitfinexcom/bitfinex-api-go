package mux

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/msg"
)

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages and reconnect client with all its subscriptions
// in case of a failure
type Mux struct {
	cid            int
	publicInbound  chan msg.Msg
	privateInbound chan msg.Msg
	publicClients  map[int]*client.Client
	privateClient  *client.Client
	mtx            *sync.RWMutex
	Err            error
	transform      bool
	apikey         string
	apisec         string
	subInfo        map[int64]event.Info
	authenticated  bool
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		publicInbound:  make(chan msg.Msg),
		privateInbound: make(chan msg.Msg),
		publicClients:  make(map[int]*client.Client),
		mtx:            &sync.RWMutex{},
		subInfo:        map[int64]event.Info{},
	}
}

// TransformRaw enables data transformation and mapping to appropriate
// models before sending it to consumer
func (m *Mux) TransformRaw() *Mux {
	m.transform = true
	return m
}

// WithApiKEY accepts and persists api key
func (m *Mux) WithApiKEY(key string) *Mux {
	m.apikey = key
	return m
}

// WithApiSEC accepts and persists api sec
func (m *Mux) WithApiSEC(sec string) *Mux {
	m.apisec = sec
	return m
}

// Subscribe - given the details in form of event.Subscribe, subscribes client
func (m *Mux) Subscribe(sub event.Subscribe) *Mux {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.Err != nil {
		return m
	}

	if alreadySubscribed := m.publicClients[m.cid].Subs.Added(sub); alreadySubscribed {
		return m
	}

	m.publicClients[m.cid].Subscribe(sub)

	if limitReached := m.publicClients[m.cid].Subs.LimitReached(); limitReached {
		log.Printf("30 subs limit is reached on cid: %d, spawning new conn\n", m.cid)
		m.addClient()
	}
	return m
}

// Start creates initial clients for accepting connections
func (m *Mux) Start() *Mux {
	return m.addClient()
}

// Listen accepts a callback func that will get called each time mux receives a
// message from any of its clients/subscriptions. It should be called last, after
// all setup calls are made as it's blocking
func (m *Mux) Listen(cb func(interface{}, error)) error {
	if m.Err != nil {
		return m.Err
	}

	for {
		select {
		case ms, ok := <-m.publicInbound:
			if !ok {
				return errors.New("channel has closed unexpectedly")
			}

			if ms.Err != nil {
				cb(nil, fmt.Errorf("conn:%d has failed | err:%s | reconnecting", ms.CID, ms.Err))
				m.reconnect(ms.CID)
				continue
			}

			// return raw payload data if transform is off
			if !m.transform {
				cb(ms.Data, nil)
				continue
			}

			// handle event type message
			if ms.IsEvent() {
				cb(m.trackSub(ms.ProcessEvent()))
				continue
			}

			// handle data type message
			if ms.IsRaw() {
				cb(ms.ProcessRaw(m.subInfo))
				continue
			}

			cb(nil, fmt.Errorf("unrecognized msg signature: %s", ms.Data))
		}
	}
}

func (m *Mux) hasAPIKeys() bool {
	return len(m.apikey) != 0 && len(m.apisec) != 0
}

func (m *Mux) addClient() *Mux {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.Err != nil {
		return m
	}

	if m.hasAPIKeys() && m.privateClient == nil {
		m.addPrivateClient()
	}

	return m.addPublicClient()
}

func (m *Mux) trackSub(i event.Info, err error) (event.Info, error) {
	// keep track of chanID to subInfo mapping
	if i.Event == "subscribed" || i.Event == "auth" {
		m.subInfo[i.ChanID] = i
	}
	m.authenticated = i.Event == "auth"
	return i, err
}

func (m *Mux) reconnect(cid int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	// pull old client subscriptions
	subs := m.publicClients[cid].Subs.GetAll()
	// add fresh client
	m.addClient()
	// resubscribe old events
	for _, sub := range subs {
		log.Printf("resubscribing: %+v\n", sub)
		m.Subscribe(sub)
	}
	// remove old, closed channel from the lost
	delete(m.publicClients, cid)
}

func (m *Mux) addPublicClient() *Mux {
	// adding new client so making sure we increment cid
	m.cid++
	// create new public client and pass error to mux if any
	c := client.New(m.cid).Public()
	if c.Err != nil {
		m.Err = c.Err
		return m
	}
	// add new client to list for later reference
	m.publicClients[m.cid] = c
	// start listening for incoming client messages
	go c.Read(m.publicInbound)
	return m
}

func (m *Mux) addPrivateClient() *Mux {
	// create new private client and pass error to mux if any
	c := client.
		New(m.cid).
		Private(m.apikey, m.apisec)

	if c.Err != nil {
		m.Err = c.Err
		return m
	}

	m.privateClient = c
	go c.Read(m.privateInbound)
	return m
}
