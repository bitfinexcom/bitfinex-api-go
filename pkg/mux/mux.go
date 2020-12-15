package mux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"unicode"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
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
	chanIdToName  map[int64]string
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		publicInbound: make(chan client.Msg),
		clients:       make(map[int]*client.Client),
		mtx:           &sync.RWMutex{},
		chanIdToName:  map[int64]string{},
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

			// return raw payload data if transform is off
			if !m.transform {
				cb(msg.Data, nil)
				continue
			}

			t := bytes.TrimLeftFunc(msg.Data, unicode.IsSpace)
			if bytes.HasPrefix(t, []byte("[")) {
				_, err := m.processRaw(msg.Data)
				if err != nil {
					cb(nil, fmt.Errorf("parsing msg: %s, err: %s", msg.Data, err))
				}
				continue
			}

			if bytes.HasPrefix(t, []byte("{")) {
				e, err := m.processEvent(msg.Data)
				if err != nil {
					cb(nil, fmt.Errorf("parsing msg: %s, err: %s", msg.Data, err))
				}
				// keep track of chanID:chanName mapping to know what data
				// type to transform raw payload to
				if e.Event == "subscribed" {
					m.chanIdToName[e.ChanID] = e.Channel
				}
				cb(e, nil)
				continue
			}

			cb(nil, fmt.Errorf("unrecognized msg signature: %s", msg.Data))
		}
	}
}

func (m *Mux) processEvent(in []byte) (e event.Info, err error) {
	err = json.Unmarshal(in, &e)
	return
}

func (m *Mux) processRaw(in []byte) (raw []interface{}, err error) {
	if err = json.Unmarshal(in, &raw); err != nil {
		return
	}

	// payload data is always last element of the slice
	pld := raw[len(raw)-1]
	// chanID is alwaus 1st element of the slice
	chID := convert.I64ValOrZero(raw[0])
	// allocate channel name by id to know how to transform raw data
	channel, ok := m.chanIdToName[chID]
	if !ok {
		err = fmt.Errorf("unrecognized chanId:%d", chID)
		return
	}

	switch data := pld.(type) {
	case string:
		log.Printf("%s chan string pld: %s\n", channel, data)
	case []interface{}:
		if _, ok := data[0].([]interface{}); ok {
			log.Printf("%s chan snapshot pld: %+v\n", channel, data)
		} else {
			log.Printf("%s chan update: %+v\n", channel, data)
		}
	}

	return
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
