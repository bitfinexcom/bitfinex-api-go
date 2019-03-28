package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// LedgerService manages the Ledger endpoint.
type LedgerService struct {
	requestFactory
	Synchronous
}

// returns the ledgers from the selected currency for authenticated account.
func (s *LedgerService) All(symbol string) (*bitfinex.LedgerSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(path.Join("ledgers", symbol, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil})
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)

	if err != nil {
		return nil, err
	}

	dat := make([][]float64, 0)
	for _, r := range raw {
		if f, ok := r.([]float64); ok {
			dat = append(dat, f)
		}
	}

	os, err := bitfinex.NewLedgerSnapshotFromRaw(symbol, dat)
	if err != nil {
		return nil, err
	}
	return os, nil
}
