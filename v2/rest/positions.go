package rest

import (
	"fmt"
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

func (s *PositionService) Claim(cp *bitfinex.ClaimPositionRequest) (*bitfinex.Notification, error) {
	bytes, err := cp.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithBytes(bitfinex.PermissionWrite, "position/claim", bytes)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(raw)
	return bitfinex.NewNotificationFromRaw(raw)
}
