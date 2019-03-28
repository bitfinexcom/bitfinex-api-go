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
func (s *LedgerService) Ledger(currency string, start int32, end int32, limit int32) (string, error) {
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(path.Join("ledgers", currency, "hist"), map[string]interface{}{"start": start, "end": end, "limit": limit})
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)

	if err != nil {
		return nil, err
	}
	
	return raw, nil
}
