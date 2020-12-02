package mux

import (
	"errors"
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
)

const clientSubsLimit = 30

type clientID int

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages
type Mux struct {
	Clients map[clientID]*client.Client
	Inbound chan client.Msg
	Subs    map[clientID]int
	Err     error
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		Clients: make(map[clientID]*client.Client),
		Inbound: make(chan client.Msg),
		Subs:    make(map[clientID]int),
	}
}

// Subscribe checks to see how many subscriptions last added client has
// and if subscribing new channel would cross the clientSubsLimit, it will
// create new client connection and call itself again recursively with same payload
func (m *Mux) Subscribe(sub map[string]string) *Mux {
	if m.Err != nil {
		return m
	}

	idx := clientID(len(m.Clients) - 1)

	if m.Subs[idx] >= clientSubsLimit {
		log.Printf("%d subs limit is reached on conn: %d, spawning new conn\n", clientSubsLimit, idx)
		m.AddPublicChan()
		m.Subscribe(sub)
		return m
	}

	m.Clients[idx].Subscribe(sub)
	m.Subs[idx]++
	return m
}

// AddPublicChan adds public cannel to mux
func (m *Mux) AddPublicChan() *Mux {
	if m.Err != nil {
		return m
	}

	c := client.New().Public()
	if c.Err != nil {
		m.Err = c.Err
		return m
	}

	m.Clients[clientID(len(m.Clients))] = c
	return m
}

// Listen accepts a callback func that will get called each time mux receives a
// message from any of its clients/subscriptions. It should be called last, after
// all setup calls are made
func (m *Mux) Listen(cb func([]byte, error)) {
	if m.Err != nil {
		cb(nil, m.Err)
		return
	}

	go m.listen()

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
func (m *Mux) listen() {
	if m.Err != nil {
		return
	}

	for _, c := range m.Clients {
		go c.Read(m.Inbound)
	}
}

func (m *Mux) reconnect(cid int) {

}
