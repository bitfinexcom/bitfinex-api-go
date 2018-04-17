package websocket

import (
	"context"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// API for end-users to interact with Bitfinex.

// Send publishes a generic message to the Bitfinex API.
func (c *Client) Send(ctx context.Context, msg interface{}) error {
	return c.asynchronous.Send(ctx, msg)
}

// Subscribe sends a subscription request to the Bitfinex API and tracks the subscription status by ID.
func (c *Client) Subscribe(ctx context.Context, req *SubscriptionRequest) (string, error) {
	c.subscriptions.add(req)
	err := c.asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
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

// SubmitOrder sends an order request.
func (c *Client) SubmitOrder(ctx context.Context, order *bitfinex.OrderNewRequest) error {
	return c.asynchronous.Send(ctx, order)
}

// SubmitCancel sends a cancel request.
func (c *Client) SubmitCancel(ctx context.Context, cancel *bitfinex.OrderCancelRequest) error {
	return c.asynchronous.Send(ctx, cancel)
}
