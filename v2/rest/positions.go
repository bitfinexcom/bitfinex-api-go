package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/notification"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
)

// PositionService manages the Position endpoint.
type PositionService struct {
	requestFactory
	Synchronous
}

// All - retrieves all of the active positions
// see https://docs.bitfinex.com/reference#rest-auth-positions for more info
func (s *PositionService) All() (*position.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "positions")
	if err != nil {
		return nil, err
	}

	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	pss, err := position.SnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pss, nil
}

// Claim - submits a request to claim an active position with the given id
// see https://docs.bitfinex.com/reference#claim-position for more info
func (s *PositionService) Claim(cp *position.ClaimRequest) (*notification.Notification, error) {
	bytes, err := cp.ToJSON()
	if err != nil {
		return nil, err
	}

	req, err := s.requestFactory.NewAuthenticatedRequestWithBytes(common.PermissionWrite, "position/claim", bytes)
	if err != nil {
		return nil, err
	}

	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	return notification.FromRaw(raw)
}
