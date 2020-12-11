package mux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"unicode"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
)

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages and reconnect client with all its subscriptions
// in case of a failure
type Mux struct {
	cid           int
	publicInbound chan client.Msg
	authInbound   chan client.Msg
	clients       map[int]*client.Client
	mtx           *sync.RWMutex
	Err           error
	transform     bool
	apiKey        string
	apiSec        string
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		publicInbound: make(chan client.Msg),
		clients:       make(map[int]*client.Client),
		mtx:           &sync.RWMutex{},
	}
}

func (m *Mux) TransformRaw() *Mux {
	m.transform = true
	return m
}

// Subscribe - given the details in form of hash table, subscribes client
func (m *Mux) Subscribe(sub event.Subscribe) *Mux {
	if m.Err != nil {
		return m
	}

	if alreadySubscribed := m.clients[m.cid].Subs.Added(sub); alreadySubscribed {
		return m
	}

	m.clients[m.cid].Subscribe(sub)

	if limitReached := m.clients[m.cid].Subs.LimitReached(); limitReached {
		log.Printf("30 subs limit is reached on cid: %d, spawning new conn\n", m.cid)
		m.AddClient()
	}
	return m
}

// AddClient adds public or authenticated client depending
// on mux api keys presence
func (m *Mux) AddClient() *Mux {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.apiKey) == 0 && len(m.apiSec) == 0 {
		return m.addPublicClient()
	}

	return m.addPrivateClient()
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
		case msg, ok := <-m.authInbound:
			if !ok {
				return errors.New("authenticated channel has closed unexpectedly")
			}
			cb(msg.Data, nil)
		case msg, ok := <-m.publicInbound:
			if !ok {
				return errors.New("channel has closed unexpectedly")
			}

			if msg.Err != nil {
				cb(nil, fmt.Errorf("conn:%d has failed | err:%s | reconnecting", msg.CID, msg.Err))
				m.reconnect(msg.CID)
				continue
			}

			if !m.transform {
				cb(msg.Data, nil)
				continue
			}

			t := bytes.TrimLeftFunc(msg.Data, unicode.IsSpace)
			if bytes.HasPrefix(t, []byte("[")) {
				cb(msg.Data, nil)
				continue
			}

			if bytes.HasPrefix(t, []byte("{")) {
				e := event.Info{}
				if err := json.Unmarshal(msg.Data, &e); err != nil {
					cb(nil, fmt.Errorf("failed parsing msg: %s", err))
					continue
				}
				cb(e, nil)
				continue
			}

			cb(nil, fmt.Errorf("unrecognized msg signature: %s", msg.Data))
		}
	}
}

func (m *Mux) reconnect(cid int) {
	// pull old client subscriptions
	subs := m.clients[cid].Subs.GetAll()
	// add fresh client
	m.AddClient()
	// resubscribe old events
	for _, sub := range subs {
		log.Printf("resubscribing: %+v\n", sub)
		m.Subscribe(sub)
	}
	// remove old, closed channel from the lost
	delete(m.clients, cid)
}

func (m *Mux) addPublicClient() *Mux {
	if m.Err != nil {
		return m
	}
	// adding new client so making sure we increment cid
	m.cid++

	c := client.New(m.cid).Public()
	if c.Err != nil {
		m.Err = c.Err
		return m
	}

	m.clients[m.cid] = c
	// start listening for incoming messages
	go c.Read(m.publicInbound)
	return m
}

func (m *Mux) addPrivateClient() *Mux {
	if m.Err != nil {
		return m
	}

	// TODO: implement auth channel handler
	return m
}
