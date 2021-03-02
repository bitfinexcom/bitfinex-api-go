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
	cid           int
	dms           int
	publicChan    chan msg.Msg
	publicClients map[int]*client.Client
	privateChan   chan msg.Msg
	privateClient *client.Client
	mtx           *sync.RWMutex
	Err           error
	transform     bool
	apikey        string
	apisec        string
	subInfo       map[int64]event.Info
	authenticated bool
	publicURL     string
	authURL       string
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		publicChan:    make(chan msg.Msg),
		privateChan:   make(chan msg.Msg),
		publicClients: make(map[int]*client.Client),
		mtx:           &sync.RWMutex{},
		subInfo:       map[int64]event.Info{},
		publicURL:     "wss://api-pub.bitfinex.com/ws/2",
		authURL:       "wss://api.staging.bitfinex.com/ws/2",
	}
}

// TransformRaw enables data transformation and mapping to appropriate
// models before sending it to consumer
func (m *Mux) TransformRaw() *Mux {
	m.transform = true
	return m
}

// WithAPIKEY accepts and persists api key
func (m *Mux) WithAPIKEY(key string) *Mux {
	m.apikey = key
	return m
}

// WithDeadManSwitch - when socket is closed, cancel all account orders
func (m *Mux) WithDeadManSwitch() *Mux {
	m.dms = 4
	return m
}

// WithAPISEC accepts and persists api sec
func (m *Mux) WithAPISEC(sec string) *Mux {
	m.apisec = sec
	return m
}

// WithAuthURL accepts and persists auth url
func (m *Mux) WithAuthURL(url string) *Mux {
	m.authURL = url
	return m
}

// Subscribe - given the details in form of event.Subscribe,
// subscribes client to public channels
func (m *Mux) Subscribe(sub event.Subscribe) *Mux {
	if m.Err != nil {
		return m
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	if subscribed := m.publicClients[m.cid].SubAdded(sub); subscribed {
		return m
	}

	if m.Err = m.publicClients[m.cid].Subscribe(sub); m.Err != nil {
		return m
	}

	if limitReached := m.publicClients[m.cid].SubsLimitReached(); limitReached {
		log.Printf("subs limit is reached on cid: %d, spawning new conn\n", m.cid)
		m.addPublicClient()
	}
	return m
}

// Start creates initial clients for accepting connections
func (m *Mux) Start() *Mux {
	if m.hasAPIKeys() && m.privateClient == nil {
		m.addPrivateClient()
	}

	return m.addPublicClient()
}

// Listen accepts a callback func that will get called each time mux
// receives a message from any of its clients/subscriptions. It
// should be called last, after all setup calls are made
func (m *Mux) Listen(cb func(interface{}, error)) error {
	if m.Err != nil {
		return m.Err
	}

	for {
		select {
		case ms, ok := <-m.publicChan:
			if !ok {
				return errors.New("channel has closed unexpectedly")
			}
			if ms.Err != nil {
				cb(nil, fmt.Errorf("conn:%d has failed | err:%s | reconnecting", ms.CID, ms.Err))
				m.resetPublicClient(ms.CID)
				continue
			}
			// return raw payload data if transform is off
			if !m.transform {
				cb(ms.Data, nil)
				continue
			}
			// handle event type message
			if ms.IsEvent() {
				cb(m.recordEvent(ms.ProcessEvent()))
				continue
			}
			// handle data type message
			if ms.IsRaw() {
				raw, pld, chID, _, err := ms.PreprocessRaw()
				if err != nil {
					cb(nil, err)
					continue
				}

				inf, ok := m.subInfo[chID]
				if !ok {
					cb(nil, fmt.Errorf("unrecognized chanId:%d", chID))
					continue
				}
				cb(ms.ProcessPublic(raw, pld, chID, inf))
				continue
			}
			cb(nil, fmt.Errorf("unrecognized msg signature: %s", ms.Data))
		case ms, ok := <-m.privateChan:
			if !ok {
				return errors.New("channel has closed unexpectedly")
			}
			if ms.Err != nil {
				cb(nil, fmt.Errorf("err: %s | reconnecting", ms.Err))
				m.resetPrivateClient()
				continue
			}
			// return raw payload data if transform is off
			if !m.transform {
				cb(ms.Data, nil)
				continue
			}
			// handle event type message
			if ms.IsEvent() {
				cb(m.recordEvent(ms.ProcessEvent()))
				continue
			}
			// handle data type message
			if ms.IsRaw() {
				raw, pld, chID, msgType, err := ms.PreprocessRaw()
				if err != nil {
					cb(nil, err)
					continue
				}
				cb(ms.ProcessPrivate(raw, pld, chID, msgType))
				continue
			}
			cb(nil, fmt.Errorf("unrecognized msg signature: %s", ms.Data))
		}
	}
}

// Send meant for authenticated input, takes payload in form of interface
// and calls client with it
func (m *Mux) Send(pld interface{}) error {
	if !m.authenticated || m.privateClient == nil {
		return errors.New("not authorized")
	}
	return m.privateClient.Send(pld)
}

func (m *Mux) hasAPIKeys() bool {
	return len(m.apikey) != 0 && len(m.apisec) != 0
}

func (m *Mux) recordEvent(i event.Info, err error) (event.Info, error) {
	switch i.Event {
	case "subscribed":
		m.subInfo[i.ChanID] = i
	case "auth":
		if i.Status == "OK" {
			m.subInfo[i.ChanID] = i
			m.authenticated = true
		}
	}
	// add more cases if/when needed
	return i, err
}

func (m *Mux) resetPublicClient(cid int) {
	// pull old client subscriptions
	subs := m.publicClients[cid].GetAllSubs()
	// add fresh client
	m.addPublicClient()
	// resubscribe old events
	for _, sub := range subs {
		log.Printf("resubscribing: %+v\n", sub)
		m.Subscribe(sub)
	}
	// remove old, closed channel from the list
	delete(m.publicClients, cid)
}

func (m *Mux) resetPrivateClient() {
	m.authenticated = false
	m.privateClient = nil
	m.addPrivateClient()
}

func (m *Mux) addPublicClient() *Mux {
	// adding new client so making sure we increment cid
	m.cid++
	// create new public client and pass error to mux if any
	c, err := client.
		New().
		WithID(m.cid).
		WithSubsLimit(25).
		Public(m.publicURL)
	if err != nil {
		m.Err = err
		return m
	}
	// add new client to list for later reference
	m.publicClients[m.cid] = c
	// start listening for incoming client messages
	go c.Read(m.publicChan)
	return m
}

func (m *Mux) addPrivateClient() *Mux {
	// create new private client and pass error to mux if any
	c, err := client.New().Private(m.apikey, m.apisec, m.authURL, m.dms)
	if err != nil {
		m.Err = err
		return m
	}

	m.privateClient = c
	go c.Read(m.privateChan)
	return m
}
