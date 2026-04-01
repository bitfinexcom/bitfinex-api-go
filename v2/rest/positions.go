package rest

import (
	"path"

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

// History - retrieves past in-active positions
// see https://docs.bitfinex.com/reference#rest-auth-positions-hist
func (s *PositionService) History() (*position.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("positions", "hist"))
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

// Audit - retrieves positions audit
// see https://docs.bitfinex.com/reference#rest-auth-positions-audit
func (s *PositionService) Audit() (*position.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("positions", "audit"))
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

// Snapshot - retrieves positions snapshot
// see https://docs.bitfinex.com/reference#rest-auth-positions-snap
func (s *PositionService) Snapshot() (*position.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("positions", "snap"))
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

// Increase - increase the size of a position
// see https://docs.bitfinex.com/reference#rest-auth-position-increase
func (s *PositionService) Increase(symbol string, amount float64) (*notification.Notification, error) {
	data := map[string]interface{}{
		"symbol": symbol,
		"amount": amount,
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionWrite, path.Join("position", "increase"), data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return notification.FromRaw(raw)
}

// IncreaseInfo - get info about increasing a position
// see https://docs.bitfinex.com/reference#rest-auth-position-increase-info
func (s *PositionService) IncreaseInfo(symbol string, amount float64) ([]interface{}, error) {
	data := map[string]interface{}{
		"symbol": symbol,
		"amount": amount,
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, path.Join("position", "increase", "info"), data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
