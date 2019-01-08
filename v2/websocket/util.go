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

//Misc Funding Callback Convenience Methods

func (c *Client) OnFundingInfo(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingInfo{}, callback)
}

func (c *Client) OnFundingTrade(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingTrade{}, callback)
}

func (c *Client) OnFundingTradeSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingTradeSnapshot{}, callback)
}

func (c *Client) OnFundingTradeExecution(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingTradeExecution{}, callback)
}

func (c *Client) OnFundingTradeUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.FundingTradeUpdate{}, callback)
}

//Balance Callback Convenience Methods

func (c *Client) OnBalanceInfo(callback ClientCallback) {
	c.RegisterCallback(bitfinex.BalanceInfo{}, callback)
}

func (c *Client) OnBalanceUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.BalanceUpdate{}, callback)
}

//Margin Callback Convenience Methods

func (c *Client) OnMarginInfoBase(callback ClientCallback) {
	c.RegisterCallback(bitfinex.MarginInfoBase{}, callback)
}

func (c *Client) OnMarginInfoUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.MarginInfoUpdate{}, callback)
}

func (c *Client) OnNotification(callback ClientCallback) {
	c.RegisterCallback(bitfinex.Notification{}, callback)
}

//Wallet Callback Convenience Methods

func (c *Client) OnWalletSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.WalletSnapshot{}, callback)
}

func (c *Client) OnWalletUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.WalletUpdate{}, callback)
}

//Public Book Callback Convenience Methods

func (c *Client) OnBookUpdateSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.BookUpdateSnapshot{}, callback)
}

func (c *Client) OnBookUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.BookUpdate{}, callback)
}

//Public Candle Callback Convenience Methods

func (c *Client) OnCandleSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.CandleSnapshot{}, callback)
}

func (c *Client) OnCandle(callback ClientCallback) {
	c.RegisterCallback(bitfinex.Candle{}, callback)
}

//Public Trade Callback Convenience Methods

func (c *Client) OnTradeSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TradeSnapshot{}, callback)
}

func (c *Client) OnTrade(callback ClientCallback) {
	c.RegisterCallback(bitfinex.Trade{}, callback)
}

//Public Ticker Callback Convenience Methods

func (c *Client) OnTicker(callback ClientCallback) {
	c.RegisterCallback(bitfinex.Ticker{}, callback)
}

func (c *Client) OnTickerSnapshot(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TickerSnapshot{}, callback)
}

func (c *Client) OnTickerUpdate(callback ClientCallback) {
	c.RegisterCallback(bitfinex.TickerUpdate{}, callback)
}
