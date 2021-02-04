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
