package rest

import (
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// OrderService manages data flow for the Order API endpoint
type DerivativesService struct {
	requestFactory
	Synchronous
}

// Update the amount of collateral for a Derivative position
// see https://docs.bitfinex.com/reference#rest-auth-deriv-pos-collateral-set for more info
func (s *WalletService) SetCollateral(symbol string, amount float64) (bool, error) {
	urlPath := path.Join("deriv", "collateral", "set")
	data := map[string]interface{}{
		"symbol":     symbol,
		"collateral": amount,
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionWrite, urlPath, data)
	if err != nil {
		return false, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return false, err
	}
	// [[1]] == success, [] || [[0]] == false
	if len(raw) <= 0 {
		return false, nil
	}
	item, ok := raw[0].([]interface{})
	if !ok || len(item) == 0 {
		return false, nil
	}
	// [1] == success, [0] == false
	if val, ok := item[0].(float64); ok && int(val) == 1 {
		return true, nil
	}
	return false, nil
}

// CollateralLimits - get collateral limits for a derivative position
// see https://docs.bitfinex.com/reference#rest-auth-deriv-pos-collateral-limits
func (s *DerivativesService) CollateralLimits(symbol string) ([]interface{}, error) {
	data := map[string]interface{}{
		"symbol": symbol,
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, path.Join("calc", "deriv", "collateral", "limits"), data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
