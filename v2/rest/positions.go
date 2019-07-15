package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// PositionService manages the Position endpoint.
type PositionService struct {
	requestFactory
	Synchronous
}

// All returns all positions for the authenticated account.
func (s *PositionService) All() (*bitfinex.PositionSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, "positions")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)

	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewPositionSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
