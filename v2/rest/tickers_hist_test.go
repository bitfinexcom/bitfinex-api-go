package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tickerhist"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicTickerHistoryResource(t *testing.T) {
	t.Run("missing arguments", func(t *testing.T) {
		c := rest.NewClient()
		th, err := c.TickersHistory.Get(rest.GetTickerHistPayload{})
		require.NotNil(t, err)
		require.Nil(t, th)
	})

	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{"abc123"}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)
		th, err := c.TickersHistory.Get(rest.GetTickerHistPayload{})
		require.NotNil(t, err)
		require.Nil(t, th)
	})

	t.Run("valid req payload and resp data", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/tickers/hist?end=123456&limit=5&start=123456&symbols=tBTCUSD", r.RequestURI)
			respMock := []interface{}{
				[]interface{}{
					"tBTCUSD", 54281, nil, 54282, nil, nil, nil, nil, nil, nil, nil, nil, 1619769715000,
				},
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)
		th, err := c.TickersHistory.Get(rest.GetTickerHistPayload{
			Symbols: []string{"tBTCUSD"},
			Limit:   5,
			Start:   123456,
			End:     123456,
		})
		require.Nil(t, err)

		expected := []tickerhist.TickerHist{
			{
				Symbol: "tBTCUSD",
				Bid:    54281,
				Ask:    54282,
				MTS:    1619769715000,
			},
		}
		assert.Equal(t, expected, th)
	})
}
