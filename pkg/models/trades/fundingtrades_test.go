package trades_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/stretchr/testify/assert"
)

func TestFundingTradeFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.FundingTrade
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{401597393},
			expected: trades.FundingTrade{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{133323072, 1574694245478, -258.7458086, 0.0002587, 2},
			expected: trades.FundingTrade{
				Symbol: "fUSD",
				ID:     133323072,
				MTS:    1574694245478,
				Amount: -258.7458086,
				Rate:   0.0002587,
				Period: 2,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.FTFromRaw("fUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestFundingTradeExecutionFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.FundingTradeExecuted
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{401597393},
			expected: trades.FundingTradeExecuted{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2},
			expected: trades.FundingTradeExecuted{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.FTEFromRaw("fUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestFundingTradeExecutionUpdateFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.FundingTradeExecutionUpdate
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{401597393},
			expected: trades.FundingTradeExecutionUpdate{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{133323543, 1574694605000, -59.84, 0.00023647, 2},
			expected: trades.FundingTradeExecutionUpdate{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.FTEUFromRaw("fUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestFundingTradeSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      [][]interface{}
		expected trades.FundingTradeSnapshot
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      [][]interface{}{},
			expected: trades.FundingTradeSnapshot{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: [][]interface{}{
				{133323543, 1574694605000, -59.84, 0.00023647, 2},
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
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.FTSnapshotFromRaw("fUSD", v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
