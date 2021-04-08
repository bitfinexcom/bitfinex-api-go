package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/msg"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	id        int
	conn      net.Conn
	nonceGen  *utils.EpochNonceGenerator
	subsLimit int
	subs      map[event.Subscribe]bool
}

// New returns pointer to Client instance
func New() *Client {
	return &Client{
		subs:     make(map[event.Subscribe]bool),
		nonceGen: utils.NewEpochNonceGenerator(),
	}
}

// WithID assigns clinet ID
func (c *Client) WithID(ID int) *Client {
	c.id = ID
	return c
}

// WithSubsLimit sets limit of subscriptions on the instance
func (c *Client) WithSubsLimit(limit int) *Client {
	c.subsLimit = limit
	return c
}

// Public creates and returns client to interact with public channels
func (c *Client) Public(url string) (*Client, error) {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return c, nil
}

// Private creates and returns client to interact with private channels
func (c *Client) Private(key, sec, url string, dms int) (*Client, error) {
	nonce := c.nonceGen.GetNonce()
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	payload := "AUTH" + nonce
	sig := hmac.New(sha512.New384, []byte(sec))
	if _, err := sig.Write([]byte(payload)); err != nil {
		return nil, err
	}

	pldSign := hex.EncodeToString(sig.Sum(nil))
	sub := event.Subscribe{
		Event:       "auth",
		APIKEY:      key,
		AuthSig:     pldSign,
		AuthPayload: payload,
		AuthNonce:   nonce,
		DMS:         dms,
	}

	if err := c.Subscribe(sub); err != nil {
		return nil, err
	}
	return c, nil
}

// Subscribe takes subscription payload as per docs and subscribes client to it.
// We keep track of subscriptions so that when client failes, we can resubscribe.
func (c *Client) Subscribe(sub event.Subscribe) error {
	sub.Event = "subscribe"
	if err := c.Send(sub); err != nil {
		return err
	}

	c.AddSub(sub)
	return nil
}

// Unsubscribe takes channel id and unsubscribes client from it.
func (c *Client) Unsubscribe(chanID int64) error {
	pld := struct {
		Event  string `json:"event"`
		ChanID int64  `json:"chanId"`
	}{
		Event:  "unsubscribe",
		ChanID: chanID,
	}

	if err := c.Send(pld); err != nil {
		return err
	}

	return nil
}

// Send takes payload in form of interface and sends it to api
func (c *Client) Send(pld interface{}) error {
	b, err := json.Marshal(pld)
	if err != nil {
		return err
	}

	return wsutil.WriteClientBinary(c.conn, b)
}

// Close closes the socket connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Read starts consuming data stream
func (c *Client) Read(ch chan<- msg.Msg) {
	defer c.conn.Close()

	for {
		ms, opCode, err := wsutil.ReadServerData(c.conn)
		m := msg.Msg{Data: ms, CID: c.id}

		if err != nil {
			m.Err = err
			ch <- m
			return
		}

		if opCode == ws.OpClose {
			m.Err = errors.New("client has closed unexpectedly")
			ch <- m
			return
		}

		ch <- m
	}
}

// SubsLimitReached returns true if number of subs > subsLimit
func (c *Client) SubsLimitReached() bool {
	if c.subsLimit == 0 {
		return false
	}
	return len(c.subs) == c.subsLimit
}

// SubAdded checks if given subscription is already added. Used to
// avoid duplicate subscriptions per client
func (c *Client) SubAdded(sub event.Subscribe) (isAdded bool) {
	_, isAdded = c.subs[sub]
	return
}

// AddSub adds new subscription to the list
func (c *Client) AddSub(sub event.Subscribe) {
	c.subs[sub] = true
}

// RemoveSub removes new subscription to the list
func (c *Client) RemoveSub(sub event.Subscribe) {
	delete(c.subs, sub)
}

// GetAllSubs returns all subscriptions
func (c *Client) GetAllSubs() (res []event.Subscribe) {
	for sub := range c.subs {
		res = append(res, sub)
	}
	return
}
