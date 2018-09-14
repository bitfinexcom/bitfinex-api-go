package websocket

import "github.com/bitfinexcom/bitfinex-api-go/v2"

//Order Callback Convenience Methods

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

//Position Callback Convenience Methods

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

//TradeExecution Callback Convenience Methods

func (c *Client) OnTradeExecution(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TradeExecution{}, callback)
}

func (c *Client) OnTradeExecutionUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TradeExecutionUpdate{}, callback)
}

//FundingOffer Callback Convenience Methods

func (c *Client) OnFundingOfferSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingOfferSnapshot{}, callback)
}

func (c *Client) OnFundingOfferNew(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingOfferNew{}, callback)
}

func (c *Client) OnFundingOfferUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingOfferUpdate{}, callback)
}

func (c *Client) OnFundingOfferCancel(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingOfferCancel{}, callback)
}

//FundingCredit Callback Convenience Methods

func (c *Client) OnFundingCreditSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingCreditSnapshot{}, callback)
}

func (c *Client) OnFundingCreditNew(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingCreditNew{}, callback)
}

func (c *Client) OnFundingCreditUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingCreditUpdate{}, callback)
}

func (c *Client) OnFundingCreditCancel(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingCreditCancel{}, callback)
}

//FundingLoan Callback Convenience Methods

func (c *Client) OnFundingLoanSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingLoanSnapshot{}, callback)
}

func (c *Client) OnFundingLoanNew(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingLoanNew{}, callback)
}

func (c *Client) OnFundingLoanUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingLoanUpdate{}, callback)
}

func (c *Client) OnFundingLoanCancel(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingLoanCancel{}, callback)
}
