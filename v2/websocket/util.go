package websocket

import "github.com/bitfinexcom/bitfinex-api-go/v2"

func (c *Client) OnOrderSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.OrderSnapshot{}, callback)
}

func (c *Client) OnOrderNew(callback ClientCallback) {
	c.RegisterCallback(bitfinex.OrderNew{}, callback)
}

func (c *Client) OnOrderUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.OrderUpdate{}, callback)
}

func (c *Client) OnOrderCancel(callback ClientCallback) {
	c.RegisterCallback(bitfinex.OrderCancel{}, callback)
}

func (c *Client) OnPositionSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.PositionSnapshot{}, callback)
}

func (c *Client) OnPositionNew(callback ClientCallback) {
	c.RegisterCallback(bitfinex.PositionNew{}, callback)
}

func (c *Client) OnPositionUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.PositionUpdate{}, callback)
}

func (c *Client) OnPositionCancel(callback ClientCallback) {
	c.RegisterCallback(bitfinex.PositionCancel{}, callback)
}

func (c *Client) OnTradeExecution(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TradeExecution{}, callback)
}

func (c *Client) OnTradeExecutionUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TradeExecutionUpdate{}, callback)
}
