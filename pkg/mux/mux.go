package mux

import (
	"errors"
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
)

const clientSubsLimit = 30

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages and reconnect client with all its subscriptions
// in case of a failure
type Mux struct {
	CID     int
	Inbound chan client.Msg
	Clients map[int]*client.Client
	Subs    map[int]map[string]map[string]string
	Err     error
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		CID:     0,
		Inbound: make(chan client.Msg),
		Clients: make(map[int]*client.Client),
		Subs:    make(map[int]map[string]map[string]string),
	}
}

// Subscribe - given the details in form of hash table, subscribes client
func (m *Mux) Subscribe(sub map[string]string) *Mux {
	if m.Err != nil {
		return m
	}

	subID := m.getSubID(sub)

	// check if already subscribed
	if _, ok := m.Subs[m.CID][subID]; ok {
		return m
	}

	// keep track of subscriptions
	if _, ok := m.Subs[m.CID]; !ok {
		m.Subs[m.CID] = make(map[string]map[string]string)
	}

	if _, ok := m.Subs[m.CID][subID]; !ok {
		m.Subs[m.CID][subID] = sub
	}

	// check if new subscription will not exceed the subscriptions limit per client
	// if it does, create new client and call Subscribe recursively with same payload
	if len(m.Subs[m.CID]) == clientSubsLimit {
		log.Printf("%d subs limit is reached on cID: %d, spawning new conn\n", clientSubsLimit, m.CID)
		m.AddPublicChan()
		return m.Subscribe(sub)
	}

	// subscribe and keep track of subscription
	m.Clients[m.CID].Subscribe(sub)
	m.Subs[m.CID][subID] = sub
	return m
}

// AddPublicChan adds public cannel to mux
func (m *Mux) AddPublicChan() *Mux {
	if m.Err != nil {
		return m
	}

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

// private methods
func (m *Mux) getSubID(sub map[string]string) (key string) {
	for _, v := range sub {
		key = key + "#" + v
	}
	return
}

func (m *Mux) reconnect(cID int) {
	// get client subscriptions
	oldSubs, ok := m.Subs[cID]
	if ok {
		delete(m.Subs, cID)
	}

	if _, ok := m.Clients[cID]; ok {
		if err := m.Clients[cID].Conn.Close(); err != nil {
			log.Printf("failed closing client: %s\n", err)
		}

		delete(m.Clients, cID)
	}

	for subID, sub := range oldSubs {
		log.Printf("resubscribing: %s\n", subID)
		m.Subscribe(sub)
	}
}
