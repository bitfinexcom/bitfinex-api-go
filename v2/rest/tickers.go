package rest

import (
	"strings"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// TickerService manages the Tickers endpoint.
type TickerService struct {	
	Synchronous
}

// Return Tickers for the public account.
func (s *TickerService) Get(symbols string) (bitfinex.TickerSnapshot, error) {

	path := []string {"tickers?symbols=", symbols}
	raw, err := s.Request(NewRequestWithMethod(strings.Join(path,""), "GET"))
	
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewTickerSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
