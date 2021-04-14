package trades_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTradesFromWSRaw(t *testing.T) {
	cases := map[string]struct {
		raw      []interface{}
		pld      []interface{}
		expected interface{}
		pair     string
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			raw:      []interface{}{},
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid funding pair snapshot": {
			raw: []interface{}{},
			pld: []interface{}{
				[]interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2},
			},
			expected: trades.FundingTradeSnapshot{
				Snapshot: []trades.FundingTrade{
					{
						Symbol: "fUSD",
						ID:     133323543,
						MTS:    1574694605000,
						Amount: -59.84,
						Rate:   0.00023647,
						Period: 2,
					},
				},
			},
			pair: "fUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"invalid funding pair snapshot": {
			raw: []interface{}{},
			pld: []interface{}{
				[]interface{}{},
			},
			expected: trades.FundingTradeSnapshot{},
			pair:     "fUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid trading pair snapshot": {
			raw: []interface{}{},
			pld: []interface{}{
				[]interface{}{401597395, 1574694478808, 0.005, 7245.3},
			},
			expected: trades.TradeSnapshot{
				Snapshot: []trades.Trade{
					{
						Pair:   "tBTCUSD",
						ID:     401597395,
						MTS:    1574694478808,
						Amount: 0.005,
						Price:  7245.3,
					},
				},
			},
			pair: "tBTCUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"invalid trading pair snapshot": {
			raw: []interface{}{},
			pld: []interface{}{
				[]interface{}{},
			},
			expected: trades.TradeSnapshot{},
			pair:     "tBTCUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid trading execution update": {
			raw: []interface{}{17470, "tu", []interface{}{401597395, 1574694478808, 0.005, 7245.3}},
			pld: []interface{}{401597395, 1574694478808, 0.005, 7245.3},
			expected: trades.TradeExecutionUpdate{
				Pair:   "tBTCUSD",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
			pair: "tBTCUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"invalid trading execution update": {
			raw:      []interface{}{17470, "tu", []interface{}{}},
			pld:      []interface{}{},
			expected: nil,
			pair:     "tBTCUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid trading execution": {
			raw: []interface{}{17470, "te", []interface{}{401597395, 1574694478808, 0.005, 7245.3}},
			pld: []interface{}{401597395, 1574694478808, 0.005, 7245.3},
			expected: trades.TradeExecuted{
				Pair:   "tBTCUSD",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
			pair: "tBTCUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"invalid trading execution": {
			raw:      []interface{}{17470, "te", []interface{}{}},
			pld:      []interface{}{},
			expected: nil,
			pair:     "tBTCUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid funding trading execution": {
			raw: []interface{}{337371, "fte", []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2}},
			pld: []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2},
			expected: trades.FundingTradeExecuted{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
			pair: "fUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"ivalid funding trading execution": {
			raw:      []interface{}{337371, "fte", []interface{}{}},
			pld:      []interface{}{},
			expected: nil,
			pair:     "fUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid funding trading execution update": {
			raw: []interface{}{337371, "ftu", []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2}},
			pld: []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2},
			expected: trades.FundingTradeExecutionUpdate{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
			pair: "fUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"ivalid funding trading execution update": {
			raw:      []interface{}{337371, "ftu", []interface{}{}},
			pld:      []interface{}{},
			expected: nil,
			pair:     "fUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid trade": {
			raw: []interface{}{17470, []interface{}{401597395, 1574694478808, 0.005, 7245.3}},
			pld: []interface{}{401597395, 1574694478808, 0.005, 7245.3},
			expected: trades.Trade{
				Pair:   "tBTCUSD",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
			pair: "tBTCUSD",
			err: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"ivalid trade": {
			raw:      []interface{}{17470, []interface{}{}},
			pld:      []interface{}{},
			expected: nil,
			pair:     "tBTCUSD",
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.FromWSRaw(v.pair, v.raw, v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
