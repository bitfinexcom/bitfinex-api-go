package margin_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/margin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingTradeFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{636040}

		got, err := margin.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid type", func(t *testing.T) {
		payload := []interface{}{"foo", []interface{}{}}

		got, err := margin.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid 'base' type payload", func(t *testing.T) {
		payload := []interface{}{"base", []interface{}{}}

		got, err := margin.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid 'base' type payload", func(t *testing.T) {
		payload := []interface{}{
			"base",
			[]interface{}{
				-13.014640000000007,
				0,
				49331.70267297,
				49318.68803297,
				27,
			},
		}

		got, err := margin.FromRaw(payload)
		require.Nil(t, err)

		expected := &margin.InfoBase{
			UserProfitLoss: -13.014640000000007,
			UserSwaps:      0,
			MarginBalance:  49331.70267297,
			MarginNet:      49318.68803297,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("invalid 'sym' type payload #1", func(t *testing.T) {
		payload := []interface{}{"sym", []interface{}{}}

		got, err := margin.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid 'sym' type payload #2", func(t *testing.T) {
		payload := []interface{}{"sym", "tETHUSD", []interface{}{}}

		got, err := margin.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid 'sym' type payload", func(t *testing.T) {
		payload := []interface{}{
			"sym",
			"tETHUSD",
			[]interface{}{
				149361.09689202666,
				149639.26293509,
				830.0182168075556,
				895.0658432466332,
				nil,
				nil,
				nil,
				nil,
			},
		}

		got, err := margin.FromRaw(payload)
		require.Nil(t, err)

		expected := &margin.InfoUpdate{
			Symbol:          "tETHUSD",
			TradableBalance: 149361.09689202666,
		}

		assert.Equal(t, expected, got)
	})
}
