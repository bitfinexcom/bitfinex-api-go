package trade_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTradeFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{399251013}

		got, err := trade.FromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid trading arguments", func(t *testing.T) {
		payload := []interface{}{
			388063448,
			1567526214876,
			1.918524,
			10682,
		}

		got, err := trade.FromRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &trade.Trade{
			Pair:   "tBTCUSD",
			ID:     388063448,
			MTS:    1567526214876,
			Amount: 1.918524,
			Price:  10682,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("valid funding arguments", func(t *testing.T) {
		payload := []interface{}{
			124486873,
			1567526287066,
			-210.69675707,
			0.00034369,
			2,
		}

		got, err := trade.FromRaw("fUSD", payload)
		require.Nil(t, err)

		expected := &trade.Trade{
			Pair:   "fUSD",
			ID:     124486873,
			MTS:    1567526287066,
			Amount: -210.69675707,
			Rate:   0.00034369,
			Period: 2,
		}

		assert.Equal(t, expected, got)
	})
}

func TestTradeSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := [][]interface{}{}
		got, err := trade.SnapshotFromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := [][]interface{}{
			{
				124486873,
				1567526287066,
				-210.69675707,
				0.00034369,
				2,
			},
			{124486874},
		}
		got, err := trade.SnapshotFromRaw("fUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := [][]interface{}{
			{
				124486873,
				1567526287066,
				-210.69675707,
				0.00034369,
				2,
			},
			{
				124486874,
				1567526287066,
				-210.69675707,
				0.00034369,
				3,
			},
		}

		got, err := trade.SnapshotFromRaw("fUSD", payload)
		require.Nil(t, err)

		expected := &trade.Snapshot{
			Snapshot: []*trade.Trade{
				{
					Pair:   "fUSD",
					ID:     124486873,
					MTS:    1567526287066,
					Amount: -210.69675707,
					Rate:   0.00034369,
					Period: 2,
				},
				{
					Pair:   "fUSD",
					ID:     124486874,
					MTS:    1567526287066,
					Amount: -210.69675707,
					Rate:   0.00034369,
					Period: 3,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestFromWSRaw(t *testing.T) {
	t.Run("missing arguments", func(t *testing.T) {
		payload := []interface{}{}

		got, err := trade.FromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{399251013}

		got, err := trade.FromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid update arguments", func(t *testing.T) {
		payload := []interface{}{
			388063448,
			1567526214876,
			1.918524,
			10682,
		}

		got, err := trade.FromRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &trade.Trade{
			Pair:   "tBTCUSD",
			ID:     388063448,
			MTS:    1567526214876,
			Amount: 1.918524,
			Price:  10682,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("valid snapshot", func(t *testing.T) {
		payload := [][]interface{}{
			{
				124486873,
				1567526287066,
				-210.69675707,
				0.00034369,
				2,
			},
			{
				124486874,
				1567526287066,
				-210.69675707,
				0.00034369,
				3,
			},
		}

		got, err := trade.SnapshotFromRaw("fUSD", payload)
		require.Nil(t, err)

		expected := &trade.Snapshot{
			Snapshot: []*trade.Trade{
				{
					Pair:   "fUSD",
					ID:     124486873,
					MTS:    1567526287066,
					Amount: -210.69675707,
					Rate:   0.00034369,
					Period: 2,
				},
				{
					Pair:   "fUSD",
					ID:     124486874,
					MTS:    1567526287066,
					Amount: -210.69675707,
					Rate:   0.00034369,
					Period: 3,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}
