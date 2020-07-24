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
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, urlPath, data)
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
	item := raw[0].([]interface{})
	// [1] == success, [] || [0] == false
	if len(item) > 0 && item[0].(int) == 1 {
		return true, nil
	}
	return false, nil
}
