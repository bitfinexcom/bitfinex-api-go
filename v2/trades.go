package bitfinex

import (
	"path"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	client *Client
}

// All returns all orders for the authenticated account.
func (s *TradeService) All(symbol string) (TradeSnapshot, error) {
	req, err := s.client.newAuthenticatedRequest("POST", path.Join("trades", symbol, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil})
	if err != nil {
		return nil, err
	}

	var raw []interface{}
	_, err = s.client.do(req, &raw)
	if err != nil {
		return nil, err
	}

	os, err := tradeSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
