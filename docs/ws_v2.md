# websocket
--
    import "github.com/bitfinexcom/bitfinex-api-go/v2/websocket"


## Usage

```go
const (
	ChanBook    = "book"
	ChanTrades  = "trades"
	ChanTicker  = "ticker"
	ChanCandles = "candles"
	ChanStatus  = "status"
)
```
Available channels

```go
const (
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventPing        = "ping"
)
```
Events

```go
const (
	ErrorCodeUnknownEvent         int = 10000
	ErrorCodeUnknownPair          int = 10001
	ErrorCodeUnknownBookPrecision int = 10011
	ErrorCodeUnknownBookLength    int = 10012
	ErrorCodeSubscriptionFailed   int = 10300
	ErrorCodeAlreadySubscribed    int = 10301
	ErrorCodeUnknownChannel       int = 10302
	ErrorCodeUnsubscribeFailed    int = 10400
	ErrorCodeNotSubscribed        int = 10401
)
```
error codes pulled from v2 docs & API usage

```go
const DMSCancelOnDisconnect int = 4
```
DMSCancelOnDisconnect cancels session orders on disconnect.

```go
const KEEP_ALIVE_TIMEOUT = 10
```
seconds to wait in between re-sending the keep alive ping

```go
const MaxChannels = 25
```

```go
const WS_READ_CAPACITY = 10
```
size of channel that the websocket reader routine pushes websocket updates into

```go
const WS_WRITE_CAPACITY = 5000
```
size of channel that the websocket writer routine pulls from

```go
var (
	ErrWSNotConnected     = fmt.Errorf("websocket connection not established")
	ErrWSAlreadyConnected = fmt.Errorf("websocket connection already established")
)
```
ws-specific errors

#### func  ConvertBytesToJsonNumberArray

```go
func ConvertBytesToJsonNumberArray(raw_bytes []byte) ([]interface{}, error)
```

#### type Asynchronous

```go
type Asynchronous interface {
	Connect() error
	Send(ctx context.Context, msg interface{}) error
	Listen() <-chan []byte
	Close()
	Done() <-chan error
}
```

Asynchronous interface decouples the underlying transport from API logic.

#### type AsynchronousFactory

```go
type AsynchronousFactory interface {
	Create() Asynchronous
}
```

AsynchronousFactory provides an interface to re-create asynchronous transports
during reconnect events.

#### func  NewWebsocketAsynchronousFactory

```go
func NewWebsocketAsynchronousFactory(parameters *Parameters) AsynchronousFactory
```
NewWebsocketAsynchronousFactory creates a new websocket factory with a given
URL.

#### type AuthEvent

```go
type AuthEvent struct {
	Event   string       `json:"event"`
	Status  string       `json:"status"`
	ChanID  int64        `json:"chanId,omitempty"`
	UserID  int64        `json:"userId,omitempty"`
	SubID   string       `json:"subId"`
	AuthID  string       `json:"auth_id,omitempty"`
	Message string       `json:"msg,omitempty"`
	Caps    Capabilities `json:"caps"`
}
```


#### type AuthState

```go
type AuthState authState // prevent user construction of authStates

```

AuthState provides a typed authentication state.

```go
const (
	NoAuthentication         AuthState = 0
	PendingAuthentication    AuthState = 1
	SuccessfulAuthentication AuthState = 2
	RejectedAuthentication   AuthState = 3
)
```
Authentication states

#### type BookFactory

```go
type BookFactory struct {
}
```


#### func (*BookFactory) Build

```go
func (f *BookFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (*BookFactory) BuildSnapshot

```go
func (f *BookFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (BookFactory) Close

```go
func (s BookFactory) Close()
```
Close is terminal. Do not call heartbeat after close.

#### func (BookFactory) ListenDisconnect

```go
func (s BookFactory) ListenDisconnect() <-chan HeartbeatDisconnect
```
ListenDisconnect returns an error channel which receives a message when a
heartbeat has expired a channel.

#### func (BookFactory) ResetAll

```go
func (s BookFactory) ResetAll()
```
Removes all tracked subscriptions

#### func (BookFactory) ResetSocketSubscriptions

```go
func (s BookFactory) ResetSocketSubscriptions(socketId SocketId) []*subscription
```
Reset clears all subscriptions assigned to the given socket ID, and returns a
slice of the existing subscriptions prior to reset

#### type CandlesFactory

```go
type CandlesFactory struct {
}
```


#### func (*CandlesFactory) Build

```go
func (f *CandlesFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (*CandlesFactory) BuildSnapshot

```go
func (f *CandlesFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (CandlesFactory) Close

```go
func (s CandlesFactory) Close()
```
Close is terminal. Do not call heartbeat after close.

#### func (CandlesFactory) ListenDisconnect

```go
func (s CandlesFactory) ListenDisconnect() <-chan HeartbeatDisconnect
```
ListenDisconnect returns an error channel which receives a message when a
heartbeat has expired a channel.

#### func (CandlesFactory) ResetAll

```go
func (s CandlesFactory) ResetAll()
```
Removes all tracked subscriptions

#### func (CandlesFactory) ResetSocketSubscriptions

```go
func (s CandlesFactory) ResetSocketSubscriptions(socketId SocketId) []*subscription
```
Reset clears all subscriptions assigned to the given socket ID, and returns a
slice of the existing subscriptions prior to reset

#### type Capabilities

```go
type Capabilities struct {
	Orders    Capability `json:"orders"`
	Account   Capability `json:"account"`
	Funding   Capability `json:"funding"`
	History   Capability `json:"history"`
	Wallets   Capability `json:"wallets"`
	Withdraw  Capability `json:"withdraw"`
	Positions Capability `json:"positions"`
}
```


#### type Capability

```go
type Capability struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}
```


#### type Client

```go
type Client struct {
	Authentication AuthState
}
```

Client provides a unified interface for users to interact with the Bitfinex V2
Websocket API. nolint:megacheck,structcheck

#### func  New

```go
func New() *Client
```
New creates a default client.

#### func  NewWithAsyncFactory

```go
func NewWithAsyncFactory(async AsynchronousFactory) *Client
```
NewWithAsyncFactory creates a new default client with a given asynchronous
transport factory interface.

#### func  NewWithAsyncFactoryNonce

```go
func NewWithAsyncFactoryNonce(async AsynchronousFactory, nonce utils.NonceGenerator) *Client
```
NewWithAsyncFactoryNonce creates a new default client with a given asynchronous
transport factory and nonce generator.

#### func  NewWithParams

```go
func NewWithParams(params *Parameters) *Client
```
NewWithParams creates a new default client with a given set of parameters.

#### func  NewWithParamsAsyncFactory

```go
func NewWithParamsAsyncFactory(params *Parameters, async AsynchronousFactory) *Client
```
NewWithParamsAsyncFactory creates a new default client with a given set of
parameters and asynchronous transport factory interface.

#### func  NewWithParamsAsyncFactoryNonce

```go
func NewWithParamsAsyncFactoryNonce(params *Parameters, async AsynchronousFactory, nonce utils.NonceGenerator) *Client
```
NewWithParamsAsyncFactoryNonce creates a new client with a given set of
parameters, asynchronous transport factory, and nonce generator interfaces.

#### func  NewWithParamsNonce

```go
func NewWithParamsNonce(params *Parameters, nonce utils.NonceGenerator) *Client
```
NewWithParamsNonce creates a new default client with a given set of parameters
and nonce generator.

#### func (*Client) AvailableCapacity

```go
func (c *Client) AvailableCapacity() int
```
Get the available capacity of the current websocket connections

#### func (*Client) CancelOnDisconnect

```go
func (c *Client) CancelOnDisconnect(cxl bool) *Client
```
CancelOnDisconnect ensures all orders will be canceled if this API session is
disconnected.

#### func (*Client) Close

```go
func (c *Client) Close()
```
Close the websocket client which will cause for all active sockets to be exited
and the Done() function to be called

#### func (*Client) Connect

```go
func (c *Client) Connect() error
```
Connect to the Bitfinex API, this should only be called once.

#### func (*Client) ConnectionCount

```go
func (c *Client) ConnectionCount() int
```
Gen the count of currently active websocket connections

#### func (*Client) Credentials

```go
func (c *Client) Credentials(key string, secret string) *Client
```
Credentials assigns authentication credentials to a connection request.

#### func (*Client) EnableFlag

```go
func (c *Client) EnableFlag(ctx context.Context, flag int) (string, error)
```
Submit a request to enable the given flag

#### func (*Client) GetAuthenticatedSocket

```go
func (c *Client) GetAuthenticatedSocket() (*Socket, error)
```
Get the authenticated socket. Due to rate limitations there can only be one
authenticated socket active at a time

#### func (*Client) GetOrderbook

```go
func (c *Client) GetOrderbook(symbol string) (*Orderbook, error)
```
Retrieve the Orderbook for the given symbol which is managed locally. This
requires ManageOrderbook=True and an active chanel subscribed to the given
symbols orderbook

#### func (*Client) IsConnected

```go
func (c *Client) IsConnected() bool
```
Returns true if the underlying asynchronous transport is connected to an
endpoint.

#### func (*Client) Listen

```go
func (c *Client) Listen() <-chan interface{}
```
Listen for all incoming api websocket messages When a websocket connection is
terminated, the publisher channel will close.

#### func (*Client) LookupSubscription

```go
func (c *Client) LookupSubscription(subID string) (*SubscriptionRequest, error)
```
Get a subscription request using a subscription ID

#### func (*Client) Send

```go
func (c *Client) Send(ctx context.Context, msg interface{}) error
```
Send publishes a generic message to the Bitfinex API.

#### func (*Client) StartNewConnection

```go
func (c *Client) StartNewConnection() error
```
Start a new websocket connection. This function is only exposed in case you want
to implicitly add new connections otherwise connection management is already
handled for you.

#### func (*Client) SubmitCancel

```go
func (c *Client) SubmitCancel(ctx context.Context, cancel *bitfinex.OrderCancelRequest) error
```
Submit a cancel request for an existing order

#### func (*Client) SubmitFundingCancel

```go
func (c *Client) SubmitFundingCancel(ctx context.Context, fundingOffer *bitfinex.FundingOfferCancelRequest) error
```
Submit a request to cancel and existing funding offer

#### func (*Client) SubmitFundingOffer

```go
func (c *Client) SubmitFundingOffer(ctx context.Context, fundingOffer *bitfinex.FundingOfferRequest) error
```
Submit a new funding offer request

#### func (*Client) SubmitOrder

```go
func (c *Client) SubmitOrder(ctx context.Context, order *bitfinex.OrderNewRequest) error
```
Submit a request to create a new order

#### func (*Client) SubmitUpdateOrder

```go
func (c *Client) SubmitUpdateOrder(ctx context.Context, orderUpdate *bitfinex.OrderUpdateRequest) error
```
Submit and update request to change an existing orders values

#### func (*Client) Subscribe

```go
func (c *Client) Subscribe(ctx context.Context, req *SubscriptionRequest) (string, error)
```
Submit a request to subscribe to the given SubscriptionRequuest

#### func (*Client) SubscribeBook

```go
func (c *Client) SubscribeBook(ctx context.Context, symbol string, precision bitfinex.BookPrecision, frequency bitfinex.BookFrequency, priceLevel int) (string, error)
```
Submit a subscription request for market data for the given symbol, at the given
frequency, with the given precision, returning no more than priceLevels price
entries. Default values are Precision0, Frequency0, and priceLevels=25.

#### func (*Client) SubscribeCandles

```go
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, resolution bitfinex.CandleResolution) (string, error)
```
Submit a subscription request to receive candle updates

#### func (*Client) SubscribeStatus

```go
func (c *Client) SubscribeStatus(ctx context.Context, symbol string, sType bitfinex.StatusType) (string, error)
```
Submit a subscription request for status updates

#### func (*Client) SubscribeTicker

```go
func (c *Client) SubscribeTicker(ctx context.Context, symbol string) (string, error)
```
Submit a request to receive ticker updates

#### func (*Client) SubscribeTrades

```go
func (c *Client) SubscribeTrades(ctx context.Context, symbol string) (string, error)
```
Submit a request to receive trade updates

#### func (*Client) Unsubscribe

```go
func (c *Client) Unsubscribe(ctx context.Context, id string) error
```
Unsubscribe from the existing subscription with the given id

#### type ConfEvent

```go
type ConfEvent struct {
	Flags int `json:"flags"`
}
```


#### type ErrorEvent

```go
type ErrorEvent struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`

	// also contain members related to subscription reject
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}
```


#### type FlagRequest

```go
type FlagRequest struct {
	Event string `json:"event"`
	Flags int    `json:"flags"`
}
```


#### type HeartbeatDisconnect

```go
type HeartbeatDisconnect struct {
	Subscription *subscription
	Error        error
}
```


#### type InfoEvent

```go
type InfoEvent struct {
	Version  float64      `json:"version"`
	ServerId string       `json:"serverId"`
	Platform PlatformInfo `json:"platform"`
	Code     int          `json:"code"`
	Msg      string       `json:"msg"`
}
```


#### type Orderbook

```go
type Orderbook struct {
}
```


#### func (*Orderbook) Asks

```go
func (ob *Orderbook) Asks() []bitfinex.BookUpdate
```

#### func (*Orderbook) Bids

```go
func (ob *Orderbook) Bids() []bitfinex.BookUpdate
```

#### func (*Orderbook) Checksum

```go
func (ob *Orderbook) Checksum() uint32
```

#### func (*Orderbook) SetWithSnapshot

```go
func (ob *Orderbook) SetWithSnapshot(bs *bitfinex.BookUpdateSnapshot)
```

#### func (*Orderbook) Symbol

```go
func (ob *Orderbook) Symbol() string
```

#### func (*Orderbook) UpdateWith

```go
func (ob *Orderbook) UpdateWith(bu *bitfinex.BookUpdate)
```

#### type Parameters

```go
type Parameters struct {
	AutoReconnect     bool
	ReconnectInterval time.Duration
	ReconnectAttempts int

	ShutdownTimeout       time.Duration
	CapacityPerConnection int
	Logger                *logging.Logger

	ResubscribeOnReconnect bool

	HeartbeatTimeout time.Duration
	LogTransport     bool

	URL             string
	ManageOrderbook bool
}
```

Parameters defines adapter behavior.

#### func  NewDefaultParameters

```go
func NewDefaultParameters() *Parameters
```

#### type PlatformInfo

```go
type PlatformInfo struct {
	Status int `json:"status"`
}
```


#### type RawEvent

```go
type RawEvent struct {
	Data interface{}
}
```


#### type Socket

```go
type Socket struct {
	Id SocketId
	Asynchronous
	IsConnected        bool
	ResetSubscriptions []*subscription
	IsAuthenticated    bool
}
```


#### type SocketId

```go
type SocketId int
```


#### type StatsFactory

```go
type StatsFactory struct {
}
```


#### func (*StatsFactory) Build

```go
func (f *StatsFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (*StatsFactory) BuildSnapshot

```go
func (f *StatsFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (StatsFactory) Close

```go
func (s StatsFactory) Close()
```
Close is terminal. Do not call heartbeat after close.

#### func (StatsFactory) ListenDisconnect

```go
func (s StatsFactory) ListenDisconnect() <-chan HeartbeatDisconnect
```
ListenDisconnect returns an error channel which receives a message when a
heartbeat has expired a channel.

#### func (StatsFactory) ResetAll

```go
func (s StatsFactory) ResetAll()
```
Removes all tracked subscriptions

#### func (StatsFactory) ResetSocketSubscriptions

```go
func (s StatsFactory) ResetSocketSubscriptions(socketId SocketId) []*subscription
```
Reset clears all subscriptions assigned to the given socket ID, and returns a
slice of the existing subscriptions prior to reset

#### type SubscribeEvent

```go
type SubscribeEvent struct {
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}
```


#### type SubscriptionRequest

```go
type SubscriptionRequest struct {
	SubID string `json:"subId"`
	Event string `json:"event"`

	// authenticated
	APIKey      string   `json:"apiKey,omitempty"`
	AuthSig     string   `json:"authSig,omitempty"`
	AuthPayload string   `json:"authPayload,omitempty"`
	AuthNonce   string   `json:"authNonce,omitempty"`
	Filter      []string `json:"filter,omitempty"`
	DMS         int      `json:"dms,omitempty"` // dead man switch

	// unauthenticated
	Channel   string `json:"channel,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair,omitempty"`
}
```


#### func (*SubscriptionRequest) String

```go
func (s *SubscriptionRequest) String() string
```

#### type SubscriptionSet

```go
type SubscriptionSet []*subscription
```

SubscriptionSet is a typed version of an array of subscription pointers,
intended to meet the sortable interface. We need to sort Reset()'s return values
for tests with more than 1 subscription (range map order is undefined)

#### func (SubscriptionSet) Len

```go
func (s SubscriptionSet) Len() int
```

#### func (SubscriptionSet) Less

```go
func (s SubscriptionSet) Less(i, j int) bool
```

#### func (SubscriptionSet) RemoveByChannelId

```go
func (s SubscriptionSet) RemoveByChannelId(chanId int64) SubscriptionSet
```

#### func (SubscriptionSet) RemoveBySubscriptionId

```go
func (s SubscriptionSet) RemoveBySubscriptionId(subID string) SubscriptionSet
```

#### func (SubscriptionSet) Swap

```go
func (s SubscriptionSet) Swap(i, j int)
```

#### type TickerFactory

```go
type TickerFactory struct {
}
```


#### func (*TickerFactory) Build

```go
func (f *TickerFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (*TickerFactory) BuildSnapshot

```go
func (f *TickerFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (TickerFactory) Close

```go
func (s TickerFactory) Close()
```
Close is terminal. Do not call heartbeat after close.

#### func (TickerFactory) ListenDisconnect

```go
func (s TickerFactory) ListenDisconnect() <-chan HeartbeatDisconnect
```
ListenDisconnect returns an error channel which receives a message when a
heartbeat has expired a channel.

#### func (TickerFactory) ResetAll

```go
func (s TickerFactory) ResetAll()
```
Removes all tracked subscriptions

#### func (TickerFactory) ResetSocketSubscriptions

```go
func (s TickerFactory) ResetSocketSubscriptions(socketId SocketId) []*subscription
```
Reset clears all subscriptions assigned to the given socket ID, and returns a
slice of the existing subscriptions prior to reset

#### type TradeFactory

```go
type TradeFactory struct {
}
```


#### func (*TradeFactory) Build

```go
func (f *TradeFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (*TradeFactory) BuildSnapshot

```go
func (f *TradeFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
```

#### func (TradeFactory) Close

```go
func (s TradeFactory) Close()
```
Close is terminal. Do not call heartbeat after close.

#### func (TradeFactory) ListenDisconnect

```go
func (s TradeFactory) ListenDisconnect() <-chan HeartbeatDisconnect
```
ListenDisconnect returns an error channel which receives a message when a
heartbeat has expired a channel.

#### func (TradeFactory) ResetAll

```go
func (s TradeFactory) ResetAll()
```
Removes all tracked subscriptions

#### func (TradeFactory) ResetSocketSubscriptions

```go
func (s TradeFactory) ResetSocketSubscriptions(socketId SocketId) []*subscription
```
Reset clears all subscriptions assigned to the given socket ID, and returns a
slice of the existing subscriptions prior to reset

#### type UnsubscribeEvent

```go
type UnsubscribeEvent struct {
	Status string `json:"status"`
	ChanID int64  `json:"chanId"`
}
```


#### type UnsubscribeRequest

```go
type UnsubscribeRequest struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}
```


#### type WebsocketAsynchronousFactory

```go
type WebsocketAsynchronousFactory struct {
}
```

WebsocketAsynchronousFactory creates a websocket-based asynchronous transport.

#### func (*WebsocketAsynchronousFactory) Create

```go
func (w *WebsocketAsynchronousFactory) Create() Asynchronous
```
Create returns a new websocket transport.
