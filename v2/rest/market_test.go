package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketAveragePrice(t *testing.T) {
	t.Run("calls with valid query params", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/calc/trade/avg?amount=100&period=2&rate_limit=1000.5&symbol=fUSD", r.RequestURI)
			assert.Equal(t, "POST", r.Method)
			respMock := []interface{}{
				123.123,
				100,
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)
		args := rest.AveragePriceRequest{
			Symbol:    "fUSD",
			Amount:    "100",
			RateLimit: "1000.5",
			Period:    2,
		}

		avgPrice, err := c.Market.AveragePrice(args)
		require.Nil(t, err)

		expected := []float64{
			123.123,
			100,
		}
		assert.Equal(t, expected, avgPrice)
	})
}

func TestForeignExchangeRate(t *testing.T) {
	t.Run("missing arguments", func(t *testing.T) {
		c := rest.NewClient()
		rsp, err := c.Market.ForeignExchangeRate(rest.ForeignExchangeRateRequest{
			FirstCurrency: "BTC",
		})

		require.NotNil(t, err)
		require.Nil(t, rsp)
	})

	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/calc/fx", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := rest.ForeignExchangeRateRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := rest.ForeignExchangeRateRequest{
				FirstCurrency:  "BTC",
				SecondCurrency: "USD",
			}
			assert.Equal(t, expectedReqPld, gotReqPld)

			respMock := []interface{}{9151.5}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)
		rsp, err := c.Market.ForeignExchangeRate(rest.ForeignExchangeRateRequest{
			FirstCurrency:  "BTC",
			SecondCurrency: "USD",
		})

		require.Nil(t, err)
		assert.Equal(t, float64(9151.5), rsp[0])
	})
}
