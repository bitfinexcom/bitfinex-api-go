package rest

import (
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// MovementsService manages the Movements endpoint.
type MovementsService struct {
	requestFactory
	Synchronous
}

// History - Get deposit/withdrawal history for given currency
// see https://docs.bitfinex.com/reference/rest-auth-movements
func (s *MovementsService) History(currency string) ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("movements", currency, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// AllHistory - Get deposit/withdrawal history for all currencies
// see https://docs.bitfinex.com/reference/rest-auth-movements
func (s *MovementsService) AllHistory() ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("movements", "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
