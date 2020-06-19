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
		args := rest.AveragePriceArgs{
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
