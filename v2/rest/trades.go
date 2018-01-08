package rest

import (
	"path"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	Synchronous
}

// All returns all orders for the authenticated account.
func (s *TradeService) All(symbol string) (bitfinex.TradeSnapshot, error) {

	raw, err := s.Request(NewRequestWithData(path.Join("trades", symbol, "hist"), map[string]interface{}{"start":nil, "end": nil, "limit": nil}))

	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewTradeSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
