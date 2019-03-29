package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// LedgerService manages the Ledgers endpoint.
type LedgerService struct {
	requestFactory
	Synchronous
}

// Ledgers returns ledgers for the selected currency and authenticated account.
func (s *LedgerService) Ledgers(currency string) (*bitfinex.LedgerSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(path.Join("ledgers", currency, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil})
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

	os, err := bitfinex.NewLedgerSnapshotFromRaw(currency, dat)
	if err != nil {
		return nil, err
	}
	return os, nil
}
