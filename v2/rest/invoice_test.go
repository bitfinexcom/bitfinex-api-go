package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/invoice"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateInvoice(t *testing.T) {
	t.Run("unsupported currency", func(t *testing.T) {
		c := rest.NewClient()

		args := rest.DepositInvoiceRequest{
			Currency: "ETH",
			Wallet:   "exchange",
			Amount:   "0.0001",
		}

		invc, err := c.Invoice.GenerateInvoice(args)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("amount too small", func(t *testing.T) {
		c := rest.NewClient()

		args := rest.DepositInvoiceRequest{
			Currency: "LNX",
			Wallet:   "exchange",
			Amount:   "0.0000001",
		}

		invc, err := c.Invoice.GenerateInvoice(args)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("amount too large", func(t *testing.T) {
		c := rest.NewClient()

		args := rest.DepositInvoiceRequest{
			Currency: "LNX",
			Wallet:   "exchange",
			Amount:   "0.03",
		}

		invc, err := c.Invoice.GenerateInvoice(args)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{"invoicehash"}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)

		pld := rest.DepositInvoiceRequest{
			Currency: "LNX",
			Wallet:   "exchange",
			Amount:   "0.0001",
		}

		invc, err := c.Invoice.GenerateInvoice(pld)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("valid response data", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/deposit/invoice", r.RequestURI)
			assert.Equal(t, "POST", r.Method)
			respMock := []interface{}{
				"invoicehash",
				"invoice",
				nil,
				nil,
				"0.002",
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)

		pld := rest.DepositInvoiceRequest{
			Currency: "LNX",
			Wallet:   "exchange",
			Amount:   "0.002",
		}

		invc, err := c.Invoice.GenerateInvoice(pld)
		require.Nil(t, err)

		expected := &invoice.Invoice{
			InvoiceHash: "invoicehash",
			Invoice:     "invoice",
			Amount:      "0.002",
		}
		assert.Equal(t, expected, invc)
	})
}
