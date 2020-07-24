package rest

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/invoice"
)

// InvoiceService manages Invoice endpoint
type InvoiceService struct {
	requestFactory
	Synchronous
}

// DepositInvoiceRequest - data structure for constructing deposit invoice request payload
type DepositInvoiceRequest struct {
	Currency string `json:"currency,omitempty"`
	Wallet   string `json:"wallet,omitempty"`
	Amount   string `json:"amount,omitempty"`
}

var validCurrencies = map[string]struct {
	name string
	min  float64
	max  float64
}{
	"LNX": {
		name: "Bitcoin Lightning Network",
		min:  0.000001,
		max:  0.02,
	},
}

func validCurrency(currency string) error {
	if _, ok := validCurrencies[currency]; !ok {
		var sb strings.Builder
		sb.WriteString(currency)
		sb.WriteString(" is not supported currency. Supported currencies: [")

		for sc := range validCurrencies {
			sb.WriteString(fmt.Sprintf(" %s(%s) ", sc, validCurrencies[sc].name))
		}

		sb.WriteString("]")

		return fmt.Errorf(sb.String())
	}

	return nil
}

func validAmount(currency, amount string) error {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}

	if f < validCurrencies[currency].min {
		return fmt.Errorf(
			"Minimum allowed amount for %s is %f. Got: %f",
			currency,
			validCurrencies[currency].min,
			f,
		)
	}

	if f > validCurrencies[currency].max {
		return fmt.Errorf(
			"Maximum allowed amount for %s is %f. Got: %f",
			currency,
			validCurrencies[currency].max,
			f,
		)
	}

	return nil
}

// GenerateInvoice generates a Lightning Network deposit invoice
// Accepts DepositInvoiceRequest type as argument
// https://docs.bitfinex.com/reference#rest-auth-deposit-invoice
func (is *InvoiceService) GenerateInvoice(payload DepositInvoiceRequest) (*invoice.Invoice, error) {
	if err := validCurrency(payload.Currency); err != nil {
		return nil, err
	}

	if err := validAmount(payload.Currency, payload.Amount); err != nil {
		return nil, err
	}

	pldBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := is.NewAuthenticatedRequestWithBytes(
		common.PermissionWrite,
		path.Join("deposit", "invoice"),
		pldBytes,
	)
	if err != nil {
		return nil, err
	}

	raw, err := is.Request(req)
	if err != nil {
		return nil, err
	}

	invc, err := invoice.NewFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return invc, nil
}
