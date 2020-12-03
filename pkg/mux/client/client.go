package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Msg struct {
	Msg []byte
	Err error
	CID int
}

type Client struct {
	Conn net.Conn
	Err  error
	ID   int
}

// New returns pointer to Client instance
func New(ID int) *Client {
	return &Client{ID: ID}
}

func (c *Client) Public() *Client {
	if c.Err != nil {
		return c
	}

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "wss://api-pub.bitfinex.com/ws/2")
	if err != nil {
		c.Err = err
		return c
	}

	c.Conn = conn
	return c
}

// Subscribe takes subscription payload as per DOCS and subscribes connection to it
func (c *Client) Subscribe(pld map[string]string) *Client {
	if c.Err != nil {
		return c
	}

	b, err := json.Marshal(pld)
	if err != nil {
		c.Err = fmt.Errorf("creating msg payload: %s, msg: %+v", err, pld)
		return c
	}

	if err = wsutil.WriteClientMessage(c.Conn, ws.OpText, b); err != nil {
		c.Err = fmt.Errorf("sending msg: %s, pld: %s", err, b)
		return c
	}

	return c
}

// Close closes socket connection
func (c *Client) Close() error {
	return c.Conn.Close()
}

func (c *Client) Read(ch chan<- Msg) {
	for {
		msg, _, err := wsutil.ReadServerData(c.Conn)
		if err != nil {
			ch <- Msg{nil, err, c.ID}
		}

		ch <- Msg{msg, nil, c.ID}
	}
}
