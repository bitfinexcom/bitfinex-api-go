package tradeexecution_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecution"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTradeExecutionFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{402088407}

		got, err := tradeexecution.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid trading arguments #1", func(t *testing.T) {
		payload := []interface{}{
			402088407,
			"tETHUST",
			1574963975602,
			34938060782,
			-0.2,
			153.57,
			"MARKET",
			0,
			-1,
			nil,
			nil,
			0,
		}

		got, err := tradeexecution.FromRaw(payload)
		require.Nil(t, err)

		expected := &tradeexecution.TradeExecution{
			ID:         402088407,
			Pair:       "tETHUST",
			MTS:        1574963975602,
			OrderID:    34938060782,
			ExecAmount: -0.2,
			ExecPrice:  153.57,
			OrderType:  "MARKET",
			OrderPrice: 0,
			Maker:      -1,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("valid trading arguments #2", func(t *testing.T) {
		payload := []interface{}{
			402088407,
			"tETHUST",
			1574963975602,
			34938060782,
			-0.2,
			153.57,
		}

		got, err := tradeexecution.FromRaw(payload)
		require.Nil(t, err)

		expected := &tradeexecution.TradeExecution{
			ID:         402088407,
			Pair:       "tETHUST",
			MTS:        1574963975602,
			OrderID:    34938060782,
			ExecAmount: -0.2,
			ExecPrice:  153.57,
		}

		assert.Equal(t, expected, got)
	})
}
