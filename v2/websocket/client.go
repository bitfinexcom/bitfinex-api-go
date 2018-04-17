package websocket

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
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
	EventPing        = "ping"
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

// AsynchronousFactory provides an interface to re-create asynchronous transports during reconnect events.
type AsynchronousFactory interface {
	Create() Asynchronous
}

// WebsocketAsynchronousFactory creates a websocket-based asynchronous transport.
type WebsocketAsynchronousFactory struct {
	parameters *Parameters
}

// NewWebsocketAsynchronousFactory creates a new websocket factory with a given URL.
func NewWebsocketAsynchronousFactory(parameters *Parameters) AsynchronousFactory {
	return &WebsocketAsynchronousFactory{
		parameters: parameters,
	}
}

// Create returns a new websocket transport.
func (w *WebsocketAsynchronousFactory) Create() Asynchronous {
	return newWs(w.parameters.URL, w.parameters.LogTransport)
}

// Client provides a unified interface for users to interact with the Bitfinex V2 Websocket API.
type Client struct {
	asyncFactory AsynchronousFactory // for re-creating transport during reconnects

	timeout            int64 // read timeout
	apiKey             string
	apiSecret          string
	Authentication     AuthState
	asynchronous       Asynchronous
	nonce              utils.NonceGenerator
	isConnected        bool
	terminal           bool
	resetSubscriptions []*subscription
	init               bool

	// connection & operational behavior
	parameters *Parameters

	// subscription manager
	subscriptions *subscriptions
	factories     map[string]messageFactory

	// close signal sent to user on shutdown
	shutdown chan bool

	// downstream listener channel to deliver API objects
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

// New creates a default client.
func New() *Client {
	return NewWithParams(NewDefaultParameters())
}

// NewWithAsyncFactory creates a new default client with a given asynchronous transport factory interface.
func NewWithAsyncFactory(async AsynchronousFactory) *Client {
	return NewWithParamsAsyncFactory(NewDefaultParameters(), async)
}

// NewWithParams creates a new default client with a given set of parameters.
func NewWithParams(params *Parameters) *Client {
	return NewWithParamsAsyncFactory(params, NewWebsocketAsynchronousFactory(params))
}

// NewWithAsyncFactoryNonce creates a new default client with a given asynchronous transport factory and nonce generator.
func NewWithAsyncFactoryNonce(async AsynchronousFactory, nonce utils.NonceGenerator) *Client {
	return NewWithParamsAsyncFactoryNonce(NewDefaultParameters(), async, nonce)
}

// NewWithParamsNonce creates a new default client with a given set of parameters and nonce generator.
func NewWithParamsNonce(params *Parameters, nonce utils.NonceGenerator) *Client {
	return NewWithParamsAsyncFactoryNonce(params, NewWebsocketAsynchronousFactory(params), nonce)
}

// NewWithParamsAsyncFactory creates a new default client with a given set of parameters and asynchronous transport factory interface.
func NewWithParamsAsyncFactory(params *Parameters, async AsynchronousFactory) *Client {
	return NewWithParamsAsyncFactoryNonce(params, async, utils.NewEpochNonceGenerator())
}

// NewWithParamsAsyncFactoryNonce creates a new client with a given set of parameters, asynchronous transport factory, and nonce generator interfaces.
func NewWithParamsAsyncFactoryNonce(params *Parameters, async AsynchronousFactory, nonce utils.NonceGenerator) *Client {
	c := &Client{
		asyncFactory:   async,
		Authentication: NoAuthentication,
		factories:      make(map[string]messageFactory),
		subscriptions:  newSubscriptions(params.HeartbeatTimeout),
		nonce:          nonce,
		isConnected:    false,
		parameters:     params,
		listener:       make(chan interface{}),
		terminal:       false,
	}
	c.registerPublicFactories()
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
	c.registerFactory(ChanTicker, newTickerFactory(c.subscriptions))
	c.registerFactory(ChanTrades, newTradeFactory(c.subscriptions))
	c.registerFactory(ChanBook, newBookFactory(c.subscriptions))
	c.registerFactory(ChanCandles, newCandlesFactory(c.subscriptions))
}

// IsConnected returns true if the underlying asynchronous transport is connected to an endpoint.
func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) listenDisconnect() {
	select {
	case e := <-c.asynchronous.Done(): // transport shutdown
		if e != nil {
			log.Printf("socket disconnect: %s", e.Error())
		}
		c.isConnected = false
		c.reconnect(e)
	case e := <-c.subscriptions.ListenDisconnect(): // subscription heartbeat timeout
		if e != nil {
			log.Printf("heartbeat disconnect: %s", e.Error())
		}
		c.isConnected = false
		if e != nil {
			c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
			c.reconnect(e)
		}
	case <-c.shutdown: // normal shutdown
		c.isConnected = false
	}
}

func (c *Client) dumpParams() {
	log.Print("----Bitfinex Client Parameters----")
	log.Printf("AutoReconnect=%t", c.parameters.AutoReconnect)
	log.Printf("ReconnectInterval=%s", c.parameters.ReconnectInterval)
	log.Printf("ReconnectAttempts=%d", c.parameters.ReconnectAttempts)
	log.Printf("ShutdownTimeout=%s", c.parameters.ShutdownTimeout)
	log.Printf("ResubscribeOnReconnect=%t", c.parameters.ResubscribeOnReconnect)
	log.Printf("HeartbeatTimeout=%s", c.parameters.HeartbeatTimeout)
	log.Printf("URL=%s", c.parameters.URL)
}

// Connect to the Bitfinex API, this should only be called once.
func (c *Client) Connect() error {
	c.dumpParams()
	c.reset()
	return c.connect()
}

// reset assumes transport has already died or been closed
func (c *Client) reset() {
	subs := c.subscriptions.Reset()
	if subs != nil {
		c.resetSubscriptions = subs
	}
	c.shutdown = make(chan bool)
	c.init = true
	c.asynchronous = c.asyncFactory.Create()
	// wait for shutdown signals from child & caller
	go c.listenDisconnect()
	// listen to data from async
	go c.listenUpstream()
}

func (c *Client) connect() error {
	err := c.asynchronous.Connect()
	if err == nil {
		c.isConnected = true
	}
	return err
}

func (c *Client) reconnect(err error) error {
	if c.terminal {
		c.exit(err)
		return err
	}
	if !c.parameters.AutoReconnect {
		err := fmt.Errorf("AutoReconnect setting is disabled, do not reconnect: %s", err.Error())
		c.exit(err)
		return err
	}
	for ; c.parameters.reconnectTry < c.parameters.ReconnectAttempts; c.parameters.reconnectTry++ {
		log.Printf("waiting %s until reconnect...", c.parameters.ReconnectInterval)
		time.Sleep(c.parameters.ReconnectInterval)
		log.Printf("reconnect attempt %d/%d", c.parameters.reconnectTry+1, c.parameters.ReconnectAttempts)
		c.reset()
		err = c.connect()
		if err == nil {
			log.Print("reconnect OK")
			c.parameters.reconnectTry = 0
			return nil
		}
		log.Printf("reconnect failed: %s", err.Error())
	}
	if err != nil {
		log.Printf("could not reconnect: %s", err.Error())
	}
	return c.exit(err)
}

func (c *Client) exit(err error) error {
	c.terminal = true
	c.close(err)
	return err
}

// start this goroutine before connecting, but this should die during a connection failure
func (c *Client) listenUpstream() {
	for {
		select {
		case <-c.shutdown:
			return // only exit point
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

// terminal, unrecoverable state. called after async is closed.
func (c *Client) close(e error) {
	if c.listener != nil {
		if e != nil {
			c.listener <- e
		}
		close(c.listener)
	}
	// shutdowns goroutines
	close(c.shutdown)
}

func (c *Client) closeAsyncAndWait(t time.Duration) {
	if !c.init {
		return
	}
	timeout := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		select {
		case <-c.asynchronous.Done(): // will this work?
			wg.Done()
		case <-timeout:
			log.Print("blocking async shutdown timed out")
			wg.Done()
		}
	}()
	c.asynchronous.Close()
	go func() {
		time.Sleep(t)
		close(timeout)
	}()
	wg.Wait()
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
// Close will close the Done() channel.
func (c *Client) Close() {
	c.terminal = true
	c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
	c.subscriptions.Close()

	// clean shutdown waits on shutdown channel, which is triggered by cascading resource
	// cleanups after a closed asynchronous transport
	timeout := make(chan bool)
	go func() {
		time.Sleep(c.parameters.ShutdownTimeout)
		close(timeout)
	}()
	select {
	case <-c.shutdown:
		return // successful cleanup
	case <-timeout:
		log.Print("shutdown timed out")
		return
	}
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

func (c *Client) sendUnsubscribeMessage(ctx context.Context, chanID int64) error {
	return c.asynchronous.Send(ctx, unsubscribeMsg{Event: "unsubscribe", ChanID: chanID})
}

func (c *Client) checkResubscription() {
	if c.parameters.ResubscribeOnReconnect && c.resetSubscriptions != nil {
		for _, sub := range c.resetSubscriptions {
			if sub.Request.Event == "auth" {
				continue
			}
			sub.Request.SubID = c.nonce.GetNonce() // new nonce
			log.Printf("resubscribing to %s with nonce %s", sub.Request.String(), sub.Request.SubID)
			_, err := c.Subscribe(context.Background(), sub.Request)
			if err != nil {
				log.Printf("could not resubscribe: %s", err.Error())
			}
		}
	}
}

// called when an info event is received
func (c *Client) handleOpen() {
	if c.hasCredentials() {
		c.authenticate(context.Background())
	} else {
		c.checkResubscription()
	}
}

// called when an auth event is received
func (c *Client) handleAuthAck(auth *AuthEvent) {
	if c.Authentication == SuccessfulAuthentication {
		err := c.subscriptions.activate(auth.SubID, auth.ChanID)
		if err != nil {
			log.Printf("could not activate auth subscription: %s", err.Error())
		}
		c.checkResubscription()
	} else {
		log.Print("authentication failed")
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
	// sub is removed from manager on ack from API
	return c.sendUnsubscribeMessage(ctx, sub.ChanID)
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
