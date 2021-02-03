package trades_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/stretchr/testify/assert"
)

func TestTradeFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.Trade
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{401597393},
			expected: trades.Trade{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{401597393, 1574694475039, 0.005, 7244.9},
			expected: trades.Trade{
				Pair:   "tBTCUSD",
				ID:     401597393,
				MTS:    1574694475039,
				Amount: 0.005,
				Price:  7244.9,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.TFromRaw("tBTCUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestTradeExecutionFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.TradeExecuted
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{401597395},
			expected: trades.TradeExecuted{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{401597395, 1574694478808, 0.005, 7245.3},
			expected: trades.TradeExecuted{
				Pair:   "tBTCUSD",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.TEFromRaw("tBTCUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
