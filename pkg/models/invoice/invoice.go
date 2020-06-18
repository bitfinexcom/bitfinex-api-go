package invoice

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// Invoice data structure
type Invoice struct {
	InvoiceHash string
	Invoice     string
	Amount      string
}

var invoiceFields = map[string]int{
	"InvoiceHash": 0,
	"Invoice":     1,
	"Amount":      4,
}

// NewFromRaw takes in slice of interfaces and converts them to
// pointer to Invoice
func NewFromRaw(raw []interface{}) (*Invoice, error) {
	if len(raw) < 5 {
		return nil, fmt.Errorf("data slice too short for Invoice: %#v", raw)
	}

	invc := &Invoice{}

	invc.InvoiceHash = convert.SValOrEmpty(raw[invoiceFields["InvoiceHash"]])
	invc.Invoice = convert.SValOrEmpty(raw[invoiceFields["Invoice"]])
	invc.Amount = convert.SValOrEmpty(raw[invoiceFields["Amount"]])

	return invc, nil
}
