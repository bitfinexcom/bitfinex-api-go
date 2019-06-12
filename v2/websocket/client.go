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

type SocketId int
type Socket struct {
	Id                 SocketId
	Asynchronous
	IsConnected        bool
	ResetSubscriptions []*subscription
	IsAuthenticated    bool
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
	sockets            map[SocketId]*Socket
	nonce              utils.NonceGenerator
	terminal           bool
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
	resetWebsocket     chan SocketId

	// downstream listener channel to deliver API objects
	listener           chan interface{}

	// race management
	mtx                *sync.RWMutex
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
		parameters:     params,
		listener:       make(chan interface{}),
		terminal:       false,
		resetWebsocket: make(chan SocketId),
		shutdown:       nil,
		sockets:        make(map[SocketId]*Socket),
		mtx:            &sync.RWMutex{},
		log:            params.Logger,
	}
	c.registerPublicFactories()
	return c
}

// Connect to the Bitfinex API, this should only be called once.
func (c *Client) Connect() error {
	c.dumpParams()
	c.terminal = false
	go c.listenDisconnect()
	return c.connectSocket(SocketId(len(c.sockets)))
}


// IsConnected returns true if the underlying asynchronous transport is connected to an endpoint.
func (c *Client) IsConnected() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	for _, socket := range c.sockets {
		if socket.IsConnected {
			return true
		}
	}
	return false
}

// Listen provides an atomic interface for receiving API messages.
// When a websocket connection is terminated, the publisher channel will close.
func (c *Client) Listen() <-chan interface{} {
	return c.listener
}

// Close provides an interface for a user initiated shutdown.
// Close will close the Done() channel.
func (c *Client) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.terminal = true
	var tmpWaitGroup sync.WaitGroup
	socketCount := len(c.sockets)
	if socketCount > 0 {
		tmpWaitGroup.Add(socketCount)
		for _, socket := range c.sockets {
			go func(s *Socket) {
				c.closeAsyncAndWait(s, c.parameters.ShutdownTimeout)
				tmpWaitGroup.Done()
			}(socket)
		}
		tmpWaitGroup.Wait()
	}
	c.subscriptions.Close()
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

func (c *Client) listenDisconnect() {
	for {
		select {
		case socketId := <- c.resetWebsocket:
			if socket, ok := c.sockets[socketId]; ok {
				c.mtx.Lock()
				socket.IsConnected = false
				c.mtx.Unlock()
				err := c.reconnect(socket, fmt.Errorf("reconnecting"))
				if err != nil {
					c.killListener(err)
					c.log.Warningf("socket disconnect: %s", err.Error())
					// exit routine if failed to reconnect
					return
				}
			}
		case <- c.shutdown:
			return
		case hbErr := <- c.subscriptions.ListenDisconnect(): // subscription heartbeat timeout
			c.log.Warningf("heartbeat disconnect: %s", hbErr.Error.Error())
			if socket, ok := c.sockets[hbErr.Subscription.SocketId]; ok {
				c.mtx.Lock()
				socket.IsConnected = false
				c.mtx.Unlock()
				// reconnect to the socket
				c.closeAsyncAndWait(socket, c.parameters.ShutdownTimeout)
				err := c.reconnect(socket, hbErr.Error)
				if err != nil {
					c.log.Warningf("socket disconnect: %s", err.Error())
					return
				}
			}
		}
	}
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

func (c *Client) reconnect(socket *Socket, err error) error {
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
		err := c.reconnectSocket(socket)
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

func (c *Client) dumpParams() {
	c.log.Debug("----Bitfinex Client Parameters----")
	c.log.Debugf("AutoReconnect=%t", c.parameters.AutoReconnect)
	c.log.Debugf("CapacityPerConnection=%t", c.parameters.CapacityPerConnection)
	c.log.Debugf("ReconnectInterval=%s", c.parameters.ReconnectInterval)
	c.log.Debugf("ReconnectAttempts=%d", c.parameters.ReconnectAttempts)
	c.log.Debugf("ShutdownTimeout=%s", c.parameters.ShutdownTimeout)
	c.log.Debugf("ResubscribeOnReconnect=%t", c.parameters.ResubscribeOnReconnect)
	c.log.Debugf("HeartbeatTimeout=%s", c.parameters.HeartbeatTimeout)
	c.log.Debugf("URL=%s", c.parameters.URL)
	c.log.Debugf("ManageOrderbook=%t", c.parameters.ManageOrderbook)
}

func (c *Client) connectSocket(socketId SocketId) error {
	async := c.asyncFactory.Create()
	// connect socket
	err := async.Connect()
	if err != nil {
		// unable to establish connection
		return err
	}
	socket := &Socket{
		Id: socketId,
		Asynchronous: async,
		IsConnected: true,
		ResetSubscriptions: make([]*subscription, 0),
		IsAuthenticated: false,
	}
	c.mtx.Lock()
	// add socket to managed map
	c.sockets[socket.Id] = socket
	c.mtx.Unlock()
	// enable orderbook flag if set in params
	// TODO - find better way to enable flags on websocket creation
	// TODO!!
	if c.parameters.ManageOrderbook {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err_flag := c.EnableFlag(ctx, bitfinex.Checksum)
		if err_flag != nil {
			return err_flag
		}
	}
	go c.listenUpstream(socket)
	return nil
}

func (c *Client) reconnectSocket(socket *Socket) error {
	// make sure the socket is closed
	socket.Close()
	// remove subscriptions from manager but keep a copy so we can resubscribe once
	// a new connection is established
	oldSubscriptions := c.subscriptions.ResetSocketSubscriptions(socket.Id)
	// establish a new connection
	err := c.connectSocket(socket.Id)
	if err != nil {
		return err
	}
	// check re-subscription on reconnect is true in the client params
	if !c.parameters.ResubscribeOnReconnect  {
		return nil
	}
	// set the resubscriptions of the socket. This will be used later on when the auth/info
	// event is passed i
	socket.ResetSubscriptions = oldSubscriptions
	return nil
}

// start this goroutine before connecting, but this should die during a connection failure
func (c *Client) listenUpstream(socket *Socket) {
	for {
		select {
		case <- socket.Asynchronous.Done(): // transport shutdown
			c.resetWebsocket <- socket.Id
			return
		case msg := <- socket.Asynchronous.Listen():
			if msg != nil {
				err := c.handleMessage(socket.Id, msg)
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

func (c *Client) closeAsyncAndWait(socket *Socket, t time.Duration) {
	if !c.init {
		return
	}
	timeout := make(chan bool)
	c.waitGroup.Add(1)
	go func() {
		select {
		case <-socket.Asynchronous.Done():
			c.waitGroup.Done()
		case <-timeout:
			c.waitGroup.Done()
		}
	}()
	socket.Asynchronous.Close()
	go func() {
		time.Sleep(t)
		close(timeout)
	}()
	c.waitGroup.Wait()
}

func (c *Client) handleMessage(socketId SocketId, msg []byte) error {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	err := error(nil)
	// either a channel data array or an event object, raw json encoding
	if bytes.HasPrefix(t, []byte("[")) {
		err = c.handleChannel(socketId, msg)
	} else if bytes.HasPrefix(t, []byte("{")) {
		err = c.handleEvent(socketId, msg)
	} else {
		return fmt.Errorf("unexpected message in socket (id=%d): %s", socketId, msg)
	}
	return err
}

func (c *Client) sendUnsubscribeMessage(ctx context.Context, chanID int64) error {
	// get the socket that the channel is assigned to
	sub, err := c.subscriptions.lookupByChannelID(chanID)
	if err != nil {
		return err
	}
	socket, err := c.socketById(sub.SocketId)
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, unsubscribeMsg{Event: "unsubscribe", ChanID: chanID})
}

func (c *Client) checkResubscription(socketId SocketId) {
	socket, err := c.socketById(socketId)
	if err != nil {
		panic(err)
	}
	if c.parameters.ResubscribeOnReconnect && socket.ResetSubscriptions != nil {
		for _, sub := range socket.ResetSubscriptions {
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
		socket.ResetSubscriptions = nil
	}
}

// called when an info event is received
func (c *Client) handleOpen(socketId SocketId) error {
	autSocket, _ := c.getAuthenticatedSocket()
	// if we have auth credentials and there is currently no authenticated
	// sockets (we are only allowed one)
	if c.hasCredentials() && autSocket == nil {
		err_auth := c.authenticate(context.Background(), socketId)
		if err_auth != nil {
			return err_auth
		}
	} else {
		c.checkResubscription(socketId)
	}
	return nil
}

// called when an auth event is received
func (c *Client) handleAuthAck(socketId SocketId, auth *AuthEvent) {
	if c.Authentication == SuccessfulAuthentication {
		// set socket to authenticated
		socket, err := c.socketById(socketId)
		if err != nil {
			panic(err)
		}
		socket.IsAuthenticated = true
		err = c.subscriptions.activate(auth.SubID, auth.ChanID)
		if err != nil {
			c.log.Errorf("could not activate auth subscription: %s", err.Error())
		}
		c.checkResubscription(socketId)
	} else {
		c.log.Error("authentication failed")
	}
}

func (c *Client) hasCredentials() bool {
	return c.apiKey != "" && c.apiSecret != ""
}

// Authenticate creates the payload for the authentication request and sends it
// to the API. The filters will be applied to the authenticated channel, i.e.
// only subscribe to the filtered messages.
func (c *Client) authenticate(ctx context.Context, socketId SocketId, filter ...string) error {
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
	c.subscriptions.add(socketId, s)
	socket, err := c.socketById(socketId)
	if err != nil {
		return err
	}
	if err = socket.Asynchronous.Send(ctx, s); err != nil {
		return err
	}
	c.Authentication = PendingAuthentication
	return nil
}

// get a random socket
func (c *Client) getSocket() (*Socket, error) {
	if len(c.sockets) <= 0 {
		return nil, fmt.Errorf("no socket found")
	}
	return c.sockets[0], nil
}

func (c *Client) getLowestCapacitySocket() (*Socket, error) {
	var lowestSocket *Socket
	lowestCapacity := -1
	for _, socket := range c.sockets {
		capac := c.getSocketCapacity(socket.Id)
		if lowestSocket == nil {
			lowestSocket = socket
			lowestCapacity = capac
			continue
		}
		if capac < lowestCapacity {
			lowestSocket = socket
			lowestCapacity = capac
		}
	}
	if lowestSocket == nil {
		return nil, fmt.Errorf("no socket found")
	}
	return lowestSocket, nil
}

// lookup the socket with the given Id, throw error if not found
func (c *Client) socketById(socketId SocketId) (*Socket, error) {
	if socket, ok := c.sockets[socketId]; ok {
		return socket, nil
	}
	return nil, fmt.Errorf("could not find socket with ID %d", socketId)
}

// calculates how many free channels are available across all of the sockets
func (c *Client) getAvailableSocketCapacity() int {
	freeCapacity := 0
	for _, socket := range c.sockets {
		freeCapacity += c.getSocketCapacity(socket.Id)
	}
	return freeCapacity
}

// calculates how many free channels are available on the given socket
func (c *Client) getSocketCapacity(socketId SocketId) int {
	subs, err := c.subscriptions.lookupBySocketId(socketId)
	if err == nil {
		return c.parameters.CapacityPerConnection - subs.Len()
	}
	return 0
}

// get the authenticated socket
func (c *Client) getAuthenticatedSocket() (*Socket, error) {
	for _, socket := range c.sockets {
		if socket.IsAuthenticated {
			return socket, nil
		}
	}
	return nil, fmt.Errorf("no authenticated socket found")
}
