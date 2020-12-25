package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/msg"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/subs"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/utils"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	Conn     net.Conn
	Subs     *subs.Subs
	Err      error
	ID       int
	isPublic bool
	nonceGen *utils.EpochNonceGenerator
}

// New returns pointer to Client instance
func New(ID int) *Client {
	return &Client{
		ID:       ID,
		Subs:     subs.New(),
		nonceGen: utils.NewEpochNonceGenerator(),
	}
}

// Public creates and returns public client to interact with public channels
func (c *Client) Public() *Client {
	if c.Err != nil {
		return c
	}

	c.isPublic = true
	c.Subs.SubsLimit = 25
	c.Conn, _, _, c.Err = ws.DefaultDialer.Dial(context.Background(), "wss://api-pub.bitfinex.com/ws/2")
	return c
}

// Private creates and returns private client to interact with private channels
func (c *Client) Private(key, sec string) *Client {
	if c.Err != nil {
		return c
	}

	// nonce := fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1000)
	nonce := c.nonceGen.GetNonce()
	c.Subs.SubsLimit = 20
	c.Conn, _, _, c.Err = ws.DefaultDialer.Dial(context.Background(), "wss://api.staging.bitfinex.com/ws/2")
	if c.Err != nil {
		return c
	}

	payload := "AUTH" + nonce
	sig := hmac.New(sha512.New384, []byte(sec))
	sig.Write([]byte(payload))
	pldSign := hex.EncodeToString(sig.Sum(nil))
	sub := event.Subscribe{
		Event:       "auth",
		APIKEY:      key,
		AuthSig:     pldSign,
		AuthPayload: payload,
		AuthNonce:   nonce,
	}

	c.Subscribe(sub)
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

func (c *Client) Read(ch chan<- msg.Msg) {
	for {
		ms, _, err := wsutil.ReadServerData(c.Conn)
		if err != nil {
			c.Conn.Close()
			ch <- msg.Msg{
				Data:     nil,
				Err:      err,
				CID:      c.ID,
				IsPublic: c.isPublic,
			}
			return
		}

		ch <- msg.Msg{
			Data:     ms,
			Err:      nil,
			CID:      c.ID,
			IsPublic: c.isPublic,
		}
	}
}
