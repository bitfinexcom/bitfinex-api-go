package bitfinex

import (
	"net/url"
)

type SymbolsDetailsService struct {
	client *Client
}

type SymbolDetail struct {
	Pair             string `json:"pair"`
	PricePrecision   int    `json:"price_precision"`
	InitialMargin    string `json:"initial_margin"`
	MinimumMargin    string `json:"minimum_margin"`
	MaximumOrderSize string `json:"maximum_order_size"`
	MinimumOrderSize string `json:"minimum_order_size"`
	Expiration       string `json:"expiration"`
}

func (s *SymbolsDetailsService) GetSymbolsDetails() ([]SymbolDetail, error) {
	params := url.Values{}
	req, err := s.client.newRequest("GET", "symbols_details", params)
	if err != nil {
		return nil, err
	}

	var v []SymbolDetail

	_, err = s.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
