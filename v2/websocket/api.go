package websocket

import (
	"context"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// API for end-users to interact with Bitfinex.

// SubscribeTicker sends a subscription request for the ticker.
func (c *Client) SubscribeTicker(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTicker,
		Symbol:  symbol,
	}
	c.subscriptions.add(req)
	err := c.asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
}

// SubscribeTrades sends a subscription request for the trade feed.
func (c *Client) SubscribeTrades(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanTrades,
		Symbol:  symbol,
	}
	c.subscriptions.add(req)
	err := c.asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
}

// SubscribeBook sends a subscription request for market data.
func (c *Client) SubscribeBook(ctx context.Context, symbol string) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanBook,
		Symbol:  symbol,
	}
	c.subscriptions.add(req)
	err := c.asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
}

// SubscribeCandles sends a subscription request for OHLC candles.
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, resolution bitfinex.CandleResolution) (string, error) {
	req := &SubscriptionRequest{
		SubID:   c.nonce.GetNonce(),
		Event:   EventSubscribe,
		Channel: ChanCandles,
		Key:     fmt.Sprintf("trade:%s:%s", resolution, symbol),
	}
	c.subscriptions.add(req)
	err := c.asynchronous.Send(ctx, req)
	if err != nil {
		// propagate send error
		return "", err
	}
	return req.SubID, nil
}

// SubmitOrder sends an order request.
func (c *Client) SubmitOrder(ctx context.Context, order *bitfinex.OrderNewRequest) error {
	return c.asynchronous.Send(ctx, order)
}

// SubmitCancel sends a cancel request.
func (c *Client) SubmitCancel(ctx context.Context, cancel *bitfinex.OrderCancelRequest) error {
	return c.asynchronous.Send(ctx, cancel)
}
