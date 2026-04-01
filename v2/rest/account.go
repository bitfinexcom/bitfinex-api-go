package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// AccountService manages the Account endpoint.
type AccountService struct {
	requestFactory
	Synchronous
}

// UserInfo - Get account info
// see https://docs.bitfinex.com/reference/rest-auth-info-user
func (s *AccountService) UserInfo() ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "info/user")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Summary - Get account summary (fees, volume)
// see https://docs.bitfinex.com/reference/rest-auth-summary
func (s *AccountService) Summary() ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "summary")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// MarginInfo - Get margin info for a given key (e.g. "base", symbol)
// see https://docs.bitfinex.com/reference/rest-auth-info-margin
func (s *AccountService) MarginInfo(key string) ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "info/margin/"+key)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// AvailableBalance - Calculate available balance for order/offer
// see https://docs.bitfinex.com/reference/rest-auth-calc-avail-balance
func (s *AccountService) AvailableBalance(symbol string, dir int, rate float64, orderType string) ([]interface{}, error) {
	data := map[string]interface{}{
		"symbol": symbol,
		"dir":    dir,
		"rate":   rate,
		"type":   orderType,
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, "calc/order/avail", data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// LoginsHistory - Get login history
// see https://docs.bitfinex.com/reference/rest-auth-logins-hist
func (s *AccountService) LoginsHistory() ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "logins/hist")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// AuditHistory - Get account changelog/audit
// see https://docs.bitfinex.com/reference/rest-auth-audit-hist
func (s *AccountService) AuditHistory() ([]interface{}, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, "audit/hist")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
