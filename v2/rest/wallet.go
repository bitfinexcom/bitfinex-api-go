package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"strconv"
)

// WalletService manages data flow for the Wallet API endpoint
type WalletService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
func (s *WalletService) Wallet() (*bitfinex.WalletSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, "wallets")
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewWalletSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}

func (ws *WalletService) Transfer(from, to, currency, currencyTo string, amount float64) (*bitfinex.Notification, error) {
	// `/v2/auth/w/transfer` (params: `from`, `to`, `currency`, `currency_to`, `amount`)
	body := map[string]interface{}{
		"from": from,
		"to": to,
		"currency": currency,
		"currency_to": currencyTo,
		"amount": strconv.FormatFloat(amount, 'f', -1, 64),
	}
	req, err := ws.requestFactory.NewAuthenticatedRequestWithData(bitfinex.PermissionWrite, "transfer", body)
	if err != nil {
		return nil, err
	}
	raw, err := ws.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewNotificationFromRaw(raw)
}

func (ws *WalletService) depositAddress(wallet string, method string, renew int) (*bitfinex.Notification, error) {
	// `/v2/auth/w/deposit/address` (params: `wallet`, `method`, `op_renew`(=1 to regenerate the wallet address))
	body := map[string]interface{}{
		"wallet": wallet,
		"method": method,
		"op_renew": renew,
	}
	req, err := ws.requestFactory.NewAuthenticatedRequestWithData(bitfinex.PermissionWrite, "deposit/address", body)
	if err != nil {
		return nil, err
	}
	raw, err := ws.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewNotificationFromRaw(raw)
}

func (ws *WalletService) DepositAddress(wallet, method string) (*bitfinex.Notification, error) {
	// `/v2/auth/w/deposit/address` (params: `wallet`, `method`, `op_renew`(=1 to regenerate the wallet address))
	return ws.depositAddress(wallet, method, 0)
}

func (ws *WalletService) CreateDepositAddress(wallet, method string) (*bitfinex.Notification, error) {
	// `/v2/auth/w/deposit/address` (params: `wallet`, `method`, `op_renew`(=1 to regenerate the wallet address))
	return ws.depositAddress(wallet, method, 1)
}

func (ws *WalletService) Withdraw(wallet, method string, amount float64, address string) (*bitfinex.Notification, error) {
	// `/v2/auth/w/withdraw` (params: `wallet`, `method`, `amount`, `address`
	body := map[string]interface{}{
		"wallet": wallet,
		"method": method,
		"amount": strconv.FormatFloat(amount, 'f', -1, 64),
		"address": address,
	}
	req, err := ws.requestFactory.NewAuthenticatedRequestWithData(bitfinex.PermissionWrite, "withdraw", body)
	if err != nil {
		return nil, err
	}
	raw, err := ws.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewNotificationFromRaw(raw)
}
