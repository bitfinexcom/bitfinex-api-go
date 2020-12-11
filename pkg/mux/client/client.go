package client

import (
	"context"
	"encoding/json"
	"net"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/subs"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Msg struct {
	Data []byte
	Err  error
	CID  int
}

type Client struct {
	Conn net.Conn
	Subs *subs.Subs
	Err  error
	ID   int
}

// New returns pointer to Client instance
func New(ID int) *Client {
	return &Client{
		ID:   ID,
		Subs: subs.New(),
	}
}

// Public creates and returns public client to interact with public channels
func (c *Client) Public() *Client {
	if c.Err != nil {
		return c
	}

	c.Conn, _, _, c.Err = ws.DefaultDialer.Dial(context.Background(), "wss://api-pub.bitfinex.com/ws/2")
	return c
}

// Subscribe takes subscription payload as per docs and subscribes connection to it
func (c *Client) Subscribe(sub event.Subscribe) *Client {
	if c.Err != nil {
		return c
	}

	c.Subs.Add(sub)
	b, _ := json.Marshal(sub)
	if c.Err = wsutil.WriteClientBinary(c.Conn, b); c.Err != nil {
		c.Subs.Remove(sub)
		return c
	}

	return c
}

func (c *Client) Read(ch chan<- Msg) {
	for {
		msg, _, err := wsutil.ReadServerData(c.Conn)
		if err != nil {
			c.Conn.Close()
			ch <- Msg{nil, err, c.ID}
			return
		}

		ch <- Msg{msg, nil, c.ID}
	}
}
