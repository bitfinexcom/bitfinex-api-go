package websocket

import (
	"bytes"
	"context"
	"fmt"
	"github.com/op/go-logging"
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

var productionBaseURL = "wss://api-pub.bitfinex.com/ws/2"

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

// DMSCancelOnDisconnect cancels session orders on disconnect.
const DMSCancelOnDisconnect int = 4

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
	return newWs(w.parameters.URL, w.parameters.LogTransport, w.parameters.Logger)
}

// Client provides a unified interface for users to interact with the Bitfinex V2 Websocket API.
// nolint:megacheck,structcheck
type Client struct {
	asyncFactory       AsynchronousFactory // for re-creating transport during reconnects

	timeout            int64 // read timeout
	apiKey             string
	apiSecret          string
	cancelOnDisconnect bool
	Authentication     AuthState
	asynchronous       Asynchronous
	nonce              utils.NonceGenerator
	isConnected        bool
	terminal           bool
	resetSubscriptions []*subscription
	init               bool
	log                *logging.Logger

	// connection & operational behavior
	parameters         *Parameters

	// subscription manager
	subscriptions      *subscriptions
	factories          map[string]messageFactory
	orderbooks         map[string]*Orderbook

	// close signal sent to user on shutdown
	shutdown           chan bool
	resetWebsocket     chan bool

	// downstream listener channel to deliver API objects
	listener           chan interface{}

	// race management
	lock               sync.Mutex
	waitGroup          sync.WaitGroup
}

// Credentials assigns authentication credentials to a connection request.
func (c *Client) Credentials(key string, secret string) *Client {
	c.apiKey = key
	c.apiSecret = secret
	return c
}

// CancelOnDisconnect ensures all orders will be canceled if this API session is disconnected.
func (c *Client) CancelOnDisconnect(cxl bool) *Client {
	c.cancelOnDisconnect = cxl
	return c
}

func (c *Client) sign(msg string) (string, error) {
	sig := hmac.New(sha512.New384, []byte(c.apiSecret))
	_, err := sig.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig.Sum(nil)), nil
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
		subscriptions:  newSubscriptions(params.HeartbeatTimeout, params.Logger),
		orderbooks:     make(map[string]*Orderbook),
		nonce:          nonce,
		isConnected:    false,
		parameters:     params,
		listener:       make(chan interface{}),
		terminal:       false,
		resetWebsocket: make(chan bool),
		shutdown:       make(chan bool),
		asynchronous:   async.Create(),
		log:            params.Logger,
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
	c.registerFactory(ChanBook, newBookFactory(c.subscriptions, c.orderbooks, c.parameters.ManageOrderbook))
	c.registerFactory(ChanCandles, newCandlesFactory(c.subscriptions))
}

// IsConnected returns true if the underlying asynchronous transport is connected to an endpoint.
func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) listenDisconnect() {
	for {
		select {
		case <- c.resetWebsocket:
			// transport websocket is shutting down
			c.lock.Lock()
			c.isConnected = false
			c.lock.Unlock()
			err := c.reconnect(fmt.Errorf("reconnecting"))
			if err != nil {
				c.killListener(err)
				c.log.Warningf("socket disconnect: %s", err.Error())
				// exit routine if failed to reconnect
				return
			}
		case <- c.shutdown:
			// websocket client killed completely
			return
		case e := <- c.subscriptions.ListenDisconnect(): // subscription heartbeat timeout
			if e != nil {
				c.log.Warningf("heartbeat disconnect: %s", e.Error())
			}
			c.lock.Lock()
			c.isConnected = false
			c.lock.Unlock()
			if e != nil {
				c.closeAsyncAndWait(c.parameters.ShutdownTimeout)
				err := c.reconnect(e)
				if err != nil {
					c.log.Warningf("socket disconnect: %s", err.Error())
					// exit routine if failed to reconnect
					return
				}
			}
		}
	}
}

func (c *Client) dumpParams() {
	c.log.Debug("----Bitfinex Client Parameters----")
	c.log.Debugf("AutoReconnect=%t", c.parameters.AutoReconnect)
	c.log.Debugf("ReconnectInterval=%s", c.parameters.ReconnectInterval)
	c.log.Debugf("ReconnectAttempts=%d", c.parameters.ReconnectAttempts)
	c.log.Debugf("ShutdownTimeout=%s", c.parameters.ShutdownTimeout)
	c.log.Debugf("ResubscribeOnReconnect=%t", c.parameters.ResubscribeOnReconnect)
	c.log.Debugf("HeartbeatTimeout=%s", c.parameters.HeartbeatTimeout)
	c.log.Debugf("URL=%s", c.parameters.URL)
	c.log.Debugf("ManageOrderbook=%t", c.parameters.ManageOrderbook)
}

// Connect to the Bitfinex API, this should only be called once.
func (c *Client) Connect() error {
	c.dumpParams()
	c.terminal = false
	// wait for reset websocket signals
	go c.listenDisconnect()
	c.reset()
	return c.connect()
}

// reset assumes transport has already died or been closed
func (c *Client) reset() {
	subs := c.subscriptions.Reset()
	if subs != nil {
		c.resetSubscriptions = subs
	}
	c.init = true
	ws := c.asyncFactory.Create()
	c.asynchronous = ws

	// listen to data from async
	go c.listenUpstream(ws)
}

func (c *Client) connect() error {
	err := c.asynchronous.Connect()
	if err != nil {
		return err
	}
	c.isConnected = true
	// enable flag
	if c.parameters.ManageOrderbook {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err_flag := c.EnableFlag(ctx, bitfinex.Checksum)
		if err_flag != nil {
			return err_flag
		}
	}
	return nil
}

func (c *Client) reconnect(err error) error {
	if c.terminal {
		// dont attempt to reconnect if terminal
		return err
	}
	if !c.parameters.AutoReconnect {
		err := fmt.Errorf("AutoReconnect setting is disabled, do not reconnect: %s", err.Error())
		return err
	}
	reconnectTry := 0
	for ; reconnectTry < c.parameters.ReconnectAttempts; reconnectTry++ {
		c.log.Debugf("waiting %s until reconnect...", c.parameters.ReconnectInterval)
		time.Sleep(c.parameters.ReconnectInterval)
		c.log.Infof("reconnect attempt %d/%d", reconnectTry+1, c.parameters.ReconnectAttempts)
		c.reset()
		err = c.connect()
		if err == nil {
			c.log.Debugf("reconnect OK")
			reconnectTry = 0
			return nil
		}
		c.log.Warningf("reconnect failed: %s", err.Error())
	}
	if err != nil {
		c.log.Errorf("could not reconnect: %s", err.Error())
	}
	return err
}


// start this goroutine before connecting, but this should die during a connection failure
func (c *Client) listenUpstream(ws Asynchronous) {
	for {
		select {
		case <- ws.Done(): // transport shutdown
			c.resetWebsocket <- true
			return
		case msg := <- ws.Listen():
			if msg != nil {
				// Errors here should be non critical so we just log them.
				// log.Printf("[DEBUG]: %s\n", msg)
				err := c.handleMessage(msg)
				if err != nil {
					c.log.Warning(err)
				}
			}
		}
	}
}

// terminal, unrecoverable state. called after async is closed.
func (c *Client) killListener(e error) {
	if c.listener != nil {
		if e != nil {
			c.listener <- e
		}
		close(c.listener)
	}
}

func (c *Client) closeAsyncAndWait(t time.Duration) {
	if !c.init {
		return
	}
	timeout := make(chan bool)
	c.waitGroup.Add(1)
	go func() {
		select {
		case <-c.asynchronous.Done():
			c.waitGroup.Done()
		case <-timeout:
			c.waitGroup.Done()
		}
	}()
	c.asynchronous.Close()
	go func() {
		time.Sleep(t)
		close(timeout)
	}()
	c.waitGroup.Wait()
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
// Close will close the Done() channel.
func (c *Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
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
	case <-c.asynchronous.Done():
		close(c.shutdown) // kill reset socket listener
		return // successful cleanup
	case <-timeout:
		c.log.Debug("shutdown timed out")
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
			c.log.Debugf("resubscribing to %s with nonce %s", sub.Request.String(), sub.Request.SubID)
			_, err := c.Subscribe(context.Background(), sub.Request)
			if err != nil {
				c.log.Errorf("could not resubscribe: %s", err.Error())
			}
		}
		c.resetSubscriptions = nil
	}
}

// called when an info event is received
func (c *Client) handleOpen() error {
	if c.hasCredentials() {
		err_auth := c.authenticate(context.Background())
		if err_auth != nil {
			return err_auth
		}
	} else {
		c.checkResubscription()
	}
	return nil
}

// called when an auth event is received
func (c *Client) handleAuthAck(auth *AuthEvent) {
	if c.Authentication == SuccessfulAuthentication {
		err := c.subscriptions.activate(auth.SubID, auth.ChanID)
		if err != nil {
			c.log.Errorf("could not activate auth subscription: %s", err.Error())
		}
		c.checkResubscription()
	} else {
		c.log.Error("authentication failed")
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
	sig, err := c.sign(payload)
	if err != nil {
		return err
	}
	s := &SubscriptionRequest{
		Event:       "auth",
		APIKey:      c.apiKey,
		AuthSig:     sig,
		AuthPayload: payload,
		AuthNonce:   nonce,
		Filter:      filter,
		SubID:       nonce,
	}
	if c.cancelOnDisconnect {
		s.DMS = DMSCancelOnDisconnect
	}
	c.subscriptions.add(s)

	if err := c.asynchronous.Send(ctx, s); err != nil {
		return err
	}
	c.Authentication = PendingAuthentication

	return nil
}
