package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// PositionService manages the Position endpoint.
type PositionService struct {
	requestFactory
	Synchronous
}

// Retrieves all of the active positions
// see https://docs.bitfinex.com/reference#rest-auth-positions for more info
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

// Submits a request to claim an active position with the given id
// see https://docs.bitfinex.com/reference#claim-position for more info
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
	return bitfinex.NewNotificationFromRaw(raw)
}
