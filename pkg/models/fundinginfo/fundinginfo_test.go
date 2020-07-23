package fundinginfo_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundinginfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingInfoFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{"sym"}

		got, err := fundinginfo.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"sym",
			"fUST",
			[]interface{}{
				0.0024,
				0.0024,
				1.9522164351851852,
				1.4818560606060607,
			},
		}

		got, err := fundinginfo.FromRaw(payload)
		require.Nil(t, err)

		expected := &fundinginfo.FundingInfo{
			Symbol:       "fUST",
			YieldLoan:    0.0024,
			YieldLend:    0.0024,
			DurationLoan: 1.9522164351851852,
			DurationLend: 1.4818560606060607,
		}
		assert.Equal(t, expected, got)
	})
}
