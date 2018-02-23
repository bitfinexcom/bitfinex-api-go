package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	Synchronous
}

// All returns all orders for the authenticated account.
func (s *TradeService) All(symbol string) (*bitfinex.TradeSnapshot, error) {

	raw, err := s.Request(NewRequestWithData(path.Join("trades", symbol, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil}))

	if err != nil {
		return nil, err
	}

	dat := make([][]float64, 0)
	for _, r := range raw {
		if f, ok := r.([]float64); ok {
			dat = append(dat, f)
		}
	}

	os, err := bitfinex.NewTradeSnapshotFromRaw(symbol, dat)
	if err != nil {
		return nil, err
	}
	return os, nil
}
