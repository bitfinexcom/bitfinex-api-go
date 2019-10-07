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

// Retrieves all of the wallets for the account
// see https://docs.bitfinex.com/reference#rest-auth-wallets for more info
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

// Submits a request to transfer funds from one Bitfinex wallet to another
// see https://docs.bitfinex.com/reference#transfer-between-wallets for more info
func (ws *WalletService) Transfer(from, to, currency, currencyTo string, amount float64) (*bitfinex.Notification, error) {
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

// Retrieves the deposit address for the given Bitfinex wallet
// see https://docs.bitfinex.com/reference#deposit-address for more info
func (ws *WalletService) DepositAddress(wallet, method string) (*bitfinex.Notification, error) {
	return ws.depositAddress(wallet, method, 0)
}

// Submits a request to create a new deposit address for the give Bitfinex wallet. Old addresses are still valid.
// See https://docs.bitfinex.com/reference#deposit-address for more info
func (ws *WalletService) CreateDepositAddress(wallet, method string) (*bitfinex.Notification, error) {
	return ws.depositAddress(wallet, method, 1)
}

// Submits a request to withdraw funds from the given Bitfinex wallet to the given address
// See https://docs.bitfinex.com/reference#withdraw for more info
func (ws *WalletService) Withdraw(wallet, method string, amount float64, address string) (*bitfinex.Notification, error) {
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
