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
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/client"
)

// Mux will manage all connections and subscriptions. Will check if subscriptions
// limit is reached and spawn new connection when that happens. It will also listen
// to all incomming client messages and reconnect client with all its subscriptions
// in case of a failure
type Mux struct {
	cid           int
	publicInbound chan client.Msg
	clients       map[int]*client.Client
	mtx           *sync.RWMutex
	Err           error
	transform     bool
	apiKey        string
	apiSec        string
	chanInfo      map[int64]event.Info
}

// New returns pointer to instance of mux
func New() *Mux {
	return &Mux{
		publicInbound: make(chan client.Msg),
		clients:       make(map[int]*client.Client),
		mtx:           &sync.RWMutex{},
		chanInfo:      map[int64]event.Info{},
	}
}

// TransformRaw enables data transformation and mapping to appropriate
// models before sending it to consumer
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

// AddClient adds public or authenticated client depending on mux api keys presence
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
				cb(m.processRaw(msg.Data))
				continue
			}

			if bytes.HasPrefix(t, []byte("{")) {
				cb(m.processEvent(msg.Data))
				continue
			}

			cb(nil, fmt.Errorf("unrecognized msg signature: %s", msg.Data))
		}
	}
}

func (m *Mux) processEvent(in []byte) (i event.Info, err error) {
	if err = json.Unmarshal(in, &i); err != nil {
		return i, fmt.Errorf("parsing msg: %s, err: %s", in, err)
	}
	// keep track of chanID to chanInfo mapping
	if i.Event == "subscribed" {
		m.chanInfo[i.ChanID] = i
	}
	return
}

func (m *Mux) processRaw(in []byte) (interface{}, error) {
	var raw []interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		return nil, fmt.Errorf("parsing msg: %s, err: %s", in, err)
	}
	// payload data is always last element of the slice
	pld := raw[len(raw)-1]
	// chanID is always 1st element of the slice
	chID := convert.I64ValOrZero(raw[0])
	// allocate channel name by id to know how to transform raw data
	inf, ok := m.chanInfo[chID]
	if !ok {
		return nil, fmt.Errorf("unrecognized chanId:%d", chID)
	}

	switch data := pld.(type) {
	case string:
		log.Printf("%d string pld: %s\n", inf.ChanID, data)
	case []interface{}:
		switch inf.Channel {
		case "trades":
			return trade.FromWSRaw(inf.Pair, data)
		case "ticker":
			return ticker.FromWSRaw(inf.Symbol, data)
		case "book":
			return book.FromWSRaw(inf.Symbol, inf.Precision, data)
		case "candles":
			return candle.FromWSRaw(inf.Key, data)
		}
	}

	return raw, nil
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
	// create new public client and pass error to mux if any
	c := client.New(m.cid).Public()
	if c.Err != nil {
		m.Err = c.Err
		return m
	}
	// add new client to list for later reference
	m.clients[m.cid] = c
	// start listening for incoming client messages
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
