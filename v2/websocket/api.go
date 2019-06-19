package websocket

import (
	"context"
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

type FlagRequest struct {
	Event string `json:"event"`
	Flags int `json:"flags"`
}

// API for end-users to interact with Bitfinex.

// Send publishes a generic message to the Bitfinex API.
func (c *Client) Send(ctx context.Context, msg interface{}) error {
	socket, err := c.getSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, msg)
}

func (c *Client) EnableFlag(ctx context.Context, flag int) (string, error) {
	req := &FlagRequest{
		Event: "conf",
		Flags: flag,
	}
	// TODO enable flag on reconnect?
	// send to all sockets
	for _, socket := range c.sockets {
		err := socket.Asynchronous.Send(ctx, req)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// returns the count of websocket connections that are currently active
func (c *Client) ConnectionCount() int {
	return len(c.sockets)
}

func (c *Client) AvailableCapacity() int {
	return c.getTotalAvailableSocketCapacity()
}

// starts a new websocket connection. This function is only exposed in case you want to
// implicitly add new connections otherwise connection management is already handled for you.
func (c *Client) StartNewConnection() error {
	return c.connectSocket(SocketId(len(c.sockets)))
}

func (c *Client) subscribeBySocket(ctx context.Context, socket *Socket, req *SubscriptionRequest) (string, error) {
	c.subscriptions.add(socket.Id, req)
	err := socket.Asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
}

// Subscribe sends a subscription request to the Bitfinex API and tracks the subscription status by ID.
func (c *Client) Subscribe(ctx context.Context, req *SubscriptionRequest) (string, error) {
	if c.getTotalAvailableSocketCapacity() <= 1 {
		err := c.StartNewConnection()
		if err != nil {
			return "", err
		}
	}
	// get socket with the highest available capacity
	socket, err := c.getMostAvailableSocket()
	if err != nil {
		return "", err
	}
	return c.subscribeBySocket(ctx, socket, req)
}

// SubscribeTicker sends a subscription request for the ticker.
func (c *Client) SubscribeTicker(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTicker,
		Symbol:  symbol,
	}
	return c.Subscribe(ctx, req)
}

// SubscribeTrades sends a subscription request for the trade feed.
func (c *Client) SubscribeTrades(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTrades,
		Symbol:  symbol,
	}
	return c.Subscribe(ctx, req)
}

// SubscribeBook sends a subscription request for market data for a given symbol, at a given frequency, with a given precision, returning no more than priceLevels price entries.
// Default values are Precision0, Frequency0, and priceLevels=25.
func (c *Client) SubscribeBook(ctx context.Context, symbol string, precision bitfinex.BookPrecision, frequency bitfinex.BookFrequency, priceLevel int) (string, error) {
	if priceLevel < 0 {
		return "", fmt.Errorf("negative price levels not supported: %d", priceLevel)
	}
	req := &SubscriptionRequest{
		SubID:     c.nonce.GetNonce(),
		Event:     EventSubscribe,
		Channel:   ChanBook,
		Symbol:    symbol,
		Precision: string(precision),
		Len:       fmt.Sprintf("%d", priceLevel), // needed for R0?
	}
	if !bitfinex.IsRawBook(string(precision)) {
		req.Frequency = string(frequency)
	}
	return c.Subscribe(ctx, req)
}

// SubscribeCandles sends a subscription request for OHLC candles.
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, resolution bitfinex.CandleResolution) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanCandles,
		Key:     fmt.Sprintf("trade:%s:%s", resolution, symbol),
	}
	return c.Subscribe(ctx, req)
}

func (c *Client) GetOrderbook(symbol string) (*Orderbook, error) {
	if val, ok := c.orderbooks[symbol]; ok {
		// take dereferenced copy of orderbook
		return val, nil
	}
	return nil, fmt.Errorf("Orderbook %s does not exist", symbol)
}

// SubmitOrder sends an order request.
func (c *Client) SubmitOrder(ctx context.Context, order *bitfinex.OrderNewRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, order)
}

func (c *Client) SubmitUpdateOrder(ctx context.Context, orderUpdate *bitfinex.OrderUpdateRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, orderUpdate)
}

// SubmitCancel sends a cancel request.
func (c *Client) SubmitCancel(ctx context.Context, cancel *bitfinex.OrderCancelRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, cancel)
}

// LookupSubscription looks up a subscription request by ID
func (c *Client) LookupSubscription(subID string) (*SubscriptionRequest, error) {
	s, err := c.subscriptions.lookupBySubscriptionID(subID)
	if err != nil {
		return nil, err
	}
	return s.Request, nil
}
