package bitfinex

import "errors"

type DepositService struct {
	client *Client
}

type DepositResponse struct {
	Result   string
	Method   string
	Currency string
	Address  string
}

func (d *DepositResponse) Success() (bool, error) {
	if d.Result == "success" {
		return true, nil
	} else {
		err := errors.New(d.Address)
		return false, err
	}
}

func (s *DepositService) New(method, walletName string, renew int) (DepositResponse, error) {

	payload := map[string]interface{}{
		"method":      method,
		"wallet_name": walletName,
		"renew":       renew,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "deposit/new", payload)

	if err != nil {
		return DepositResponse{}, err
	}

	var v DepositResponse
	_, err = s.client.do(req, &v)

	if err != nil {
		return DepositResponse{}, err
	}
	return v, nil
}
