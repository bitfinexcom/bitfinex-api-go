package bitfinex

import (
	"net/url"
)

type SymbolsService struct {
	client *Client
}

func (s *SymbolsService) GetSymbols() ([]string, error) {
	params := url.Values{}
	req, err := s.client.newRequest("GET", "symbols", params)
	if err != nil {
		return nil, err
	}

	var v []string

	_, err = s.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
