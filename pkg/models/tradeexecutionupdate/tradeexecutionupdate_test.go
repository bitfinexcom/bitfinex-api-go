package tradeexecutionupdate_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecutionupdate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTradeExecutionUpdateFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{402088407}

		got, err := tradeexecutionupdate.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid trading arguments", func(t *testing.T) {
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
			-0.061668,
			"USD",
		}

		got, err := tradeexecutionupdate.FromRaw(payload)
		require.Nil(t, err)

		expected := &tradeexecutionupdate.TradeExecutionUpdate{
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
		}

		assert.Equal(t, expected, got)
	})
}

func TestTradeExecutionUpdateSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := tradeexecutionupdate.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				402088407,
				"tETHUST",
				1574963975602,
				34938060782,
				-0.2,
				153.57,
				"MARKET",
				0,
				-1,
				-0.061668,
				"USD",
			},
			[]interface{}{402088408},
		}
		got, err := tradeexecutionupdate.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				402088407,
				"tETHUST",
				1574963975602,
				34938060782,
				-0.2,
				153.57,
				"MARKET",
				0,
				-1,
				-0.061668,
				"USD",
			},
			[]interface{}{
				402088408,
				"tETHUST",
				1574963975602,
				34938060782,
				-0.2,
				153.57,
				"MARKET",
				0,
				-1,
				-0.061668,
				"USD",
			},
		}

		got, err := tradeexecutionupdate.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &tradeexecutionupdate.Snapshot{
			Snapshot: []*tradeexecutionupdate.TradeExecutionUpdate{
				{
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
				{
					ID:          402088408,
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
			},
		}

		assert.Equal(t, expected, got)
	})
}
