package rest

import (
	"strings"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// CandleService manages the Candles endpoint.
type CandleService struct {	
	Synchronous
}

// Return Candles for the public account.
func (s *CandleService) Get(symbol string, resolution bitfinex.CandleResolution, param string) (bitfinex.CandleSnapshot, error) {

	path := []string {"candles/trade:", string(resolution) ,":t",symbol,"/",param}
	raw, err := s.Request(NewRequestWithMethod(strings.Join(path,""), "GET"))
	
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewCandleSnapshotFromRaw(symbol,resolution,raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
