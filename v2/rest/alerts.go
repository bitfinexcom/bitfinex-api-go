package rest

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// AlertService manages the Alert endpoint.
type AlertService struct {
	requestFactory
	Synchronous
}

// All - Get all active alerts
// see https://docs.bitfinex.com/reference/rest-auth-alert-list
func (s *AlertService) All(alertType string) ([]interface{}, error) {
	if alertType == "" {
		alertType = "price"
	}
	data := map[string]interface{}{"type": alertType}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, "alerts", data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Set - Set a new price alert
// see https://docs.bitfinex.com/reference/rest-auth-alert-set
func (s *AlertService) Set(alertType, symbol string, price float64) ([]interface{}, error) {
	if alertType == "" {
		alertType = "price"
	}
	data := map[string]interface{}{
		"type":   alertType,
		"symbol": symbol,
		"price":  price,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithBytes(common.PermissionWrite, "alert/set", b)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Delete - Delete a price alert
// see https://docs.bitfinex.com/reference/rest-auth-alert-delete
func (s *AlertService) Delete(symbol string, price float64) ([]interface{}, error) {
	endpoint := fmt.Sprintf("alert/price:%s:%v/del", symbol, price)
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionWrite, endpoint)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
