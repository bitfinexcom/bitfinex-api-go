package websocket

import (
	"context"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
)

type FlagRequest struct {
	Event string `json:"event"`
	Flags int    `json:"flags"`
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

// Gen the count of currently active websocket connections
func (c *Client) ConnectionCount() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return len(c.sockets)
}

// Get the available capacity of the current
// websocket connections
func (c *Client) AvailableCapacity() int {
	return c.getTotalAvailableSocketCapacity()
}

// Start a new websocket connection. This function is only exposed in case you want to
// implicitly add new connections otherwise connection management is already handled for you.
func (c *Client) StartNewConnection() error {
	return c.connectSocket(SocketId(c.ConnectionCount()))
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

// Submit a request to subscribe to the given SubscriptionRequuest
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

// Submit a request to receive ticker updates
func (c *Client) SubscribeTicker(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTicker,
		Symbol:  symbol,
	}
	return c.Subscribe(ctx, req)
}

// Submit a request to receive trade updates
func (c *Client) SubscribeTrades(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTrades,
		Symbol:  symbol,
	}
	return c.Subscribe(ctx, req)
}

// Submit a  subscription request for market data for the given symbol, at the given frequency, with the given precision, returning no more than priceLevels price entries.
// Default values are Precision0, Frequency0, and priceLevels=25.
func (c *Client) SubscribeBook(ctx context.Context, symbol string, precision common.BookPrecision, frequency common.BookFrequency, priceLevel int) (string, error) {
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
	if !book.IsRawBook(string(precision)) {
		req.Frequency = string(frequency)
	}
	return c.Subscribe(ctx, req)
}

// Submit a subscription request to receive candle updates
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, resolution common.CandleResolution) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanCandles,
		Key:     fmt.Sprintf("trade:%s:%s", resolution, symbol),
	}
	return c.Subscribe(ctx, req)
}

// Submit a subscription request for status updates
func (c *Client) SubscribeStatus(ctx context.Context, symbol string, sType common.StatusType) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanStatus,
		Key:     fmt.Sprintf("%s:%s", string(sType), symbol),
	}
	return c.Subscribe(ctx, req)
}

// Retrieve the Orderbook for the given symbol which is managed locally.
// This requires ManageOrderbook=True and an active chanel subscribed to the given
// symbols orderbook
func (c *Client) GetOrderbook(symbol string) (*Orderbook, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	if val, ok := c.orderbooks[symbol]; ok {
		// take dereferenced copy of orderbook
		return val, nil
	}
	return nil, fmt.Errorf("Orderbook %s does not exist", symbol)
}

// Submit a request to create a new order
func (c *Client) SubmitOrder(ctx context.Context, onr *order.NewRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, onr)
}

// Submit and update request to change an existing orders values
func (c *Client) SubmitUpdateOrder(ctx context.Context, our *order.UpdateRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, our)
}

// Submit a cancel request for an existing order
func (c *Client) SubmitCancel(ctx context.Context, ocr *order.CancelRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, ocr)
}

// Get a subscription request using a subscription ID
func (c *Client) LookupSubscription(subID string) (*SubscriptionRequest, error) {
	s, err := c.subscriptions.lookupBySubscriptionID(subID)
	if err != nil {
		return nil, err
	}
	return s.Request, nil
}

// Submit a new funding offer request
func (c *Client) SubmitFundingOffer(ctx context.Context, fundingOffer *fundingoffer.SubmitRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, fundingOffer)
}

// Submit a request to cancel and existing funding offer
func (c *Client) SubmitFundingCancel(ctx context.Context, fundingOffer *fundingoffer.CancelRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, fundingOffer)
}

// CloseFundingLoan - cancels funding loan by ID. Emits an error if not authenticated.
func (c *Client) CloseFundingLoan(ctx context.Context, flcr *fundingloan.CancelRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, flcr)
}

// CloseFundingCredit - cancels funding credit by ID. Emits an error if not authenticated.
func (c *Client) CloseFundingCredit(ctx context.Context, fundingOffer *fundingcredit.CancelRequest) error {
	socket, err := c.GetAuthenticatedSocket()
	if err != nil {
		return err
	}
	return socket.Asynchronous.Send(ctx, fundingOffer)
}
