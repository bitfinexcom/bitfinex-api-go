package websocket

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/bitfinexcom/bitfinex-api-go/utils"

	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

var productionBaseURL = "wss://api.bitfinex.com/ws/2"

// ws-specific errors
var (
	ErrWSNotConnected     = fmt.Errorf("websocket connection not established")
	ErrWSAlreadyConnected = fmt.Errorf("websocket connection already established")
)

// Available channels
const (
	ChanBook    = "book"
	ChanTrades  = "trades"
	ChanTicker  = "ticker"
	ChanCandles = "candles"
)

// Events
const (
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
)

// Authentication states
const (
	NoAuthentication         AuthState = 0
	PendingAuthentication    AuthState = 1
	SuccessfulAuthentication AuthState = 2
	RejectedAuthentication   AuthState = 3
)

// private type--cannot instantiate.
type authState byte

// AuthState provides a typed authentication state.
type AuthState authState // prevent user construction of authStates

// Asynchronous interface decouples the underlying transport from API logic.
type Asynchronous interface {
	Connect() error
	Send(ctx context.Context, msg interface{}) error
	Listen() <-chan []byte
	Close()
	Done() <-chan error
}

// Client provides a unified interface for users to interact with the Bitfinex V2 Websocket API.
type Client struct {
	timeout        int64 // read timeout
	apiKey         string
	apiSecret      string
	Authentication AuthState
	asynchronous   Asynchronous
	nonce          utils.NonceGenerator

	// subscription manager
	subscriptions *subscriptions
	factories     map[string]messageFactory

	// close signal sent to user on shutdown
	shutdown chan bool

	listener chan interface{}
}

// Credentials assigns authentication credentials to a connection request.
func (c *Client) Credentials(key string, secret string) *Client {
	c.apiKey = key
	c.apiSecret = secret
	return c
}

func (c *Client) sign(msg string) string {
	sig := hmac.New(sha512.New384, []byte(c.apiSecret))
	sig.Write([]byte(msg))
	return hex.EncodeToString(sig.Sum(nil))
}

func (c *Client) registerFactory(channel string, factory messageFactory) {
	c.factories[channel] = factory
}

// NewClientWithURL creates a new default client with a given API endpoint.
func NewClientWithURL(url string) *Client {
	return NewClientWithAsync(newWs(url))
}

// NewClientWithAsync creates a new default client with a given asynchronous transport interface.
func NewClientWithAsync(async Asynchronous) *Client {
	return NewClientWithAsyncNonce(async, utils.NewEpochNonceGenerator())
}

// NewClientWithAsync creates a new client with a given asynchronous transport and nonce generator interfaces.
func NewClientWithAsyncNonce(async Asynchronous, nonce utils.NonceGenerator) *Client {
	c := &Client{
		asynchronous:   async,
		shutdown:       make(chan bool),
		Authentication: NoAuthentication,
		factories:      make(map[string]messageFactory),
		listener:       make(chan interface{}),
		subscriptions:  newSubscriptions(),
		nonce:          nonce,
	}
	c.registerPublicFactories()
	// wait for shutdown signals from child & caller
	go c.listenDisconnect()
	return c
}

func extractSymbolResolutionFromKey(subscription string) (symbol string, resolution bitfinex.CandleResolution, err error) {
	var res, sym string
	str := strings.Split(subscription, ":")
	if len(str) < 3 {
		return "", resolution, fmt.Errorf("could not convert symbol resolution for %s: len %d", subscription, len(str))
	}
	res = str[1]
	sym = str[2]
	resolution, err = bitfinex.CandleResolutionFromString(res)
	if err != nil {
		return "", resolution, err
	}
	return sym, resolution, nil
}

func (c *Client) registerPublicFactories() {
	c.registerFactory(ChanTicker, func(chanID int64, raw []interface{}) (msg interface{}, err error) {
		sub, err := c.subscriptions.lookupByChannelID(chanID)
		if err == nil {
			tick, err := bitfinex.NewTickerFromRaw(sub.Request.Symbol, raw)
			return &tick, err
		}
		return nil, err
	})
	c.registerFactory(ChanTrades, func(chanID int64, raw []interface{}) (msg interface{}, err error) {
		return bitfinex.NewTradeSnapshotFromRaw(raw)
	})
	c.registerFactory(ChanBook, func(chanID int64, raw []interface{}) (msg interface{}, err error) {
		sub, err := c.subscriptions.lookupByChannelID(chanID)
		if err == nil {
			update, err := bitfinex.NewBookUpdateFromRaw(sub.Request.Symbol, raw)
			return &update, err
		}
		return nil, err
	})
	c.registerFactory(ChanCandles, func(chanID int64, raw []interface{}) (msg interface{}, err error) {
		sub, err := c.subscriptions.lookupByChannelID(chanID)
		if err != nil {
			return nil, err
		}
		sym, res, err := extractSymbolResolutionFromKey(sub.Request.Key)
		if err != nil {
			return nil, err
		}
		book, err := bitfinex.NewCandleFromRaw(sym, res, raw)
		return &book, err
	})
}

// NewClient creates a new default client.
func NewClient() *Client {
	return NewClientWithURL(productionBaseURL)
}

// Connect to the Bitfinex API.
func (c *Client) Connect() error {
	err := c.asynchronous.Connect()
	if err == nil {
		go c.listenUpstream()
	}
	return err
}

func (c *Client) listenDisconnect() {
	// block until finished
	select {
	case err := <-c.asynchronous.Done(): // child shutdown
		c.close(err)
		return
	case <-c.shutdown: // normal shutdown
		return
	}
}

func (c *Client) listenUpstream() {
	for {
		select {
		case <-c.shutdown:
			return
		case msg := <-c.asynchronous.Listen():
			if msg != nil {
				// Errors here should be non critical so we just log them.
				err := c.handleMessage(msg)
				if err != nil {
					log.Printf("[WARN]: %s\n", err)
				}
			}
		}
	}
}

// cleanly dispose of resources & signal we are finished
func (c *Client) close(e error) {
	// internal goroutine shutdown
	close(c.shutdown)

	if c.listener != nil {
		if e != nil {
			c.listener <- e
		}
		close(c.listener)
	}
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the listen channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
// Close will close the Done() channel.
func (c *Client) Close() {
	// close transport
	c.asynchronous.Close() // will trigger a close()
}

func (c *Client) handleMessage(msg []byte) error {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	err := error(nil)
	// either a channel data array or an event object, raw json encoding
	if bytes.HasPrefix(t, []byte("[")) {
		err = c.handleChannel(msg)
	} else if bytes.HasPrefix(t, []byte("{")) {
		err = c.handleEvent(msg)
	} else {
		return fmt.Errorf("unexpected message: %s", msg)
	}
	return err
}

func (c *Client) sendUnsubscribeMessage(ctx context.Context, id int64) error {
	return c.asynchronous.Send(ctx, unsubscribeMsg{Event: "unsubscribe", ChanID: id})
}

func (c *Client) unsubscribeByChanID(ctx context.Context, id int64) error {
	return c.sendUnsubscribeMessage(ctx, id)
}

func (c *Client) handleOpen() {
	if c.hasCredentials() {
		c.authenticate(context.Background())
	}
}

func (c *Client) hasCredentials() bool {
	return c.apiKey != "" && c.apiSecret != ""
}

// Unsubscribe looks up an existing subscription by ID and sends an unsubscribe request.
func (c *Client) Unsubscribe(ctx context.Context, id string) error {
	sub, err := c.subscriptions.lookupBySubscriptionID(id)
	if err != nil {
		return err
	}
	return c.unsubscribeByChanID(ctx, sub.ChanID)
}

// Authenticate creates the payload for the authentication request and sends it
// to the API. The filters will be applied to the authenticated channel, i.e.
// only subscribe to the filtered messages.
func (c *Client) authenticate(ctx context.Context, filter ...string) error {
	nonce := c.nonce.GetNonce()

	payload := "AUTH" + nonce
	s := &SubscriptionRequest{
		Event:       "auth",
		APIKey:      c.apiKey,
		AuthSig:     c.sign(payload),
		AuthPayload: payload,
		AuthNonce:   nonce,
		Filter:      filter,
		SubID:       nonce,
	}
	c.subscriptions.add(s)

	if err := c.asynchronous.Send(ctx, s); err != nil {
		return err
	}
	c.Authentication = PendingAuthentication

	return nil
}

// SetReadTimeout sets the read timeout for the underlying websocket connections.
func (c *Client) SetReadTimeout(t time.Duration) {
	atomic.StoreInt64(&c.timeout, t.Nanoseconds())
}
