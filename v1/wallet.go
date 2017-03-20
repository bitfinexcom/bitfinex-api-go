package bitfinex

import "strconv"

const (
	WALLET_TRADING  = "trading"
	WALLET_EXCHANGE = "exchange"
	WALLET_DEPOSIT  = "deposit"
)

type WalletService struct {
	client *Client
}

type TransferStatus struct {
	Status  string
	Message string
}

// Transfer funds between wallets
func (c *WalletService) Transfer(amount float64, currency, from, to string) ([]TransferStatus, error) {

	payload := map[string]interface{}{
		"amount":     strconv.FormatFloat(amount, 'f', -1, 32),
		"currency":   currency,
		"walletfrom": from,
		"walletto":   to,
	}

	req, err := c.client.newAuthenticatedRequest("GET", "transfer", payload)

	if err != nil {
		return nil, err
	}

	status := make([]TransferStatus, 0)

	_, err = c.client.do(req, &status)

	return status, err
}

type WithdrawStatus struct {
	Status       string
	Message      string
	WithdrawalID int `json:"withdrawal_id"`
}

// Withdraw a cryptocurrency to a digital wallet
func (c *WalletService) WithdrawCrypto(amount float64, currency, wallet, destinationAddress string) ([]WithdrawStatus, error) {

	payload := map[string]interface{}{
		"amount":         strconv.FormatFloat(amount, 'f', -1, 32),
		"walletselected": wallet,
		"withdraw_type":  currency,
		"address":        destinationAddress,
	}

	req, err := c.client.newAuthenticatedRequest("GET", "withdraw", payload)

	if err != nil {
		return nil, err
	}

	status := make([]WithdrawStatus, 0)

	_, err = c.client.do(req, &status)

	return status, err

}

type BankAccount struct {
	AccountName   string // Account name
	AccountNumber string // Account number or IBAN
	BankName      string // Bank Name
	BankAddress   string // Bank Address
	BankCity      string // Bank City
	BankCountry   string // Bank Country
	SwiftCode     string // SWIFT Code
}

func (c *WalletService) WithdrawWire(amount float64, expressWire bool, wallet string, beneficiaryBank, intermediaryBank BankAccount, message string) ([]WithdrawStatus, error) {

	var express int
	if expressWire {
		express = 1
	} else {
		express = 0
	}

	payload := map[string]interface{}{
		"amount":                    strconv.FormatFloat(amount, 'f', -1, 32),
		"walletselected":            wallet,
		"withdraw_type":             "wire",
		"expressWire":               express,
		"account_name":              beneficiaryBank.AccountName,
		"account_number":            beneficiaryBank.AccountNumber,
		"bank_name":                 beneficiaryBank.BankName,
		"bank_address":              beneficiaryBank.BankAddress,
		"bank_city":                 beneficiaryBank.BankCity,
		"bank_country":              beneficiaryBank.BankCountry,
		"swift":                     beneficiaryBank.SwiftCode,
		"detail_payment":            message,
		"intermediary_bank_account": intermediaryBank.AccountNumber,
		"intermediary_bank_address": intermediaryBank.BankAddress,
		"intermediary_bank_city":    intermediaryBank.BankCity,
		"intermediary_bank_country": intermediaryBank.BankCountry,
		"intermediary_bank_swift":   intermediaryBank.SwiftCode,
	}

	req, err := c.client.newAuthenticatedRequest("GET", "withdraw", payload)

	if err != nil {
		return nil, err
	}

	status := make([]WithdrawStatus, 0)

	_, err = c.client.do(req, &status)

	return status, err

}
