package mux

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
)

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages and reconnect client with all its subscriptions
// in case of a failure
type Mux struct {
	CID     int
	Inbound chan client.Msg
	Clients map[int]*client.Client
	mtx     *sync.RWMutex
	Err     error
	APIKey  string
	APISec  string
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		Inbound: make(chan client.Msg),
		Clients: make(map[int]*client.Client),
		mtx:     &sync.RWMutex{},
	}
}

// Subscribe - given the details in form of hash table, subscribes client
func (m *Mux) Subscribe(sub map[string]string) *Mux {
	if m.Err != nil {
		return m
	}

	if alreadySubscribed := m.Clients[m.CID].Subs.Added(sub); alreadySubscribed {
		return m
	}

	m.Clients[m.CID].Subscribe(sub)

	if limitReached := m.Clients[m.CID].Subs.LimitReached(); limitReached {
		log.Printf("30 subs limit is reached on cID: %d, spawning new conn\n", m.CID)
		m.AddClient()
	}
	return m
}

// AddClient adds public or authenticated client depending
// on mux api keys presence
func (m *Mux) AddClient() *Mux {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.APIKey) == 0 && len(m.APISec) == 0 {
		return m.addPublicClient()
	}

	return m.addPrivateClient()
}

// Listen accepts a callback func that will get called each time mux receives a
// message from any of its clients/subscriptions. It should be called last, after
// all setup calls are made as it's blocking
func (m *Mux) Listen(cb func([]byte, error)) {
	if m.Err != nil {
		cb(nil, m.Err)
		return
	}

	log.Println("starting to listen...")

	for {
		select {
		case msg, ok := <-m.Inbound:
			log.Printf("m:%s, e:%v, chan:%t\n", msg.Msg, msg.Err, ok)
			if !ok {
				cb(nil, errors.New("channel has closed unexpectedly, restart"))
				return
			}

			if msg.Err != nil {
				cb(nil, fmt.Errorf("conn:%d has failed | err:%s | reconnecting", msg.CID, msg.Err))
				m.reconnect(msg.CID)
				continue
			}

			cb(msg.Msg, nil)
		}
	}
}

func (m *Mux) reconnect(cID int) {
	// pull old client subscriptions
	subs := m.Clients[cID].Subs.GetAll()
	// add fresh client
	m.AddClient()
	// resubscribe old events
	for subID, sub := range subs {
		log.Printf("resubscribing: %s\n", subID)
		m.Subscribe(sub)
	}
	// remove old, closed channel from the lost
	delete(m.Clients, cID)
}

func (m *Mux) addPublicClient() *Mux {
	if m.Err != nil {
		return m
	}

	// adding new client so making sure we increment cid
	m.CID++

	c := client.New(m.CID).Public()
	if c.Err != nil {
		m.Err = c.Err
		return m
	}

	m.Clients[m.CID] = c
	// start listening for incoming messages
	go c.Read(m.Inbound)
	return m
}

func (m *Mux) addPrivateClient() *Mux {
	if m.Err != nil {
		return m
	}

	// TODO: implement auth channel handler
	return m
}
