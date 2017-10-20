package bitfinex

import ()

// PositionService manages the Position endpoint.
type PositionService struct {
	client *Client
}

// All returns all positions for the authenticated account.
func (s *PositionService) All() (PositionSnapshot, error) {
	req, err := s.client.newAuthenticatedRequest("POST", "positions", nil)
	if err != nil {
		return nil, err
	}

	var raw []interface{}
	_, err = s.client.do(req, &raw)
	if err != nil {
		return nil, err
	}

	os, err := positionSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
