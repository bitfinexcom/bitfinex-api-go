package trades_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticatedTradeExecutionUpdateFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.AuthTradeExecutionUpdate
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{402088407},
			expected: trades.AuthTradeExecutionUpdate{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{
				402088407, "tETHUST", 1574963975602, 34938060782,
				-0.2, 153.57, "MARKET", 0, -1, -0.061668, "USD",
			},
			expected: trades.AuthTradeExecutionUpdate{
				ID:          402088407,
				Pair:        "tETHUST",
				MTS:         1574963975602,
				OrderID:     34938060782,
				ExecAmount:  -0.2,
				ExecPrice:   153.57,
				OrderType:   "MARKET",
				OrderPrice:  0,
				Maker:       -1,
				Fee:         -0.061668,
				FeeCurrency: "USD",
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.ATEUFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
