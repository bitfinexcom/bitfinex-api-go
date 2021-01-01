package ticker_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTickerFromRestRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{"tBTCUSD"}

		got, err := ticker.FromRestRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"tBTCUSD",
			10654,
			53.62425959,
			10655,
			76.68743116,
			745.1,
			0.0752,
			10655,
			14420.34271146,
			10766,
			9889.1449809,
		}

		got, err := ticker.FromRestRaw(payload)
		require.Nil(t, err)

		expected := &ticker.Ticker{
			Symbol:          "tBTCUSD",
			Frr:             0,
			Bid:             10654,
			BidPeriod:       0,
			BidSize:         53.62425959,
			Ask:             10655,
			AskPeriod:       0,
			AskSize:         76.68743116,
			DailyChange:     745.1,
			DailyChangePerc: 0.0752,
			LastPrice:       10655,
			Volume:          14420.34271146,
			High:            10766,
			Low:             9889.1449809,
		}

		assert.Equal(t, expected, got)
	})
}

func TestFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{7616.5, 31.89055171}

		got, err := ticker.FromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("trading pair update / snapshot ticker", func(t *testing.T) {
		payload := []interface{}{
			7616.5,
			31.89055171,
			7617.5,
			43.358118629999986,
			-550.8,
			-0.0674,
			7617.1,
			8314.71200815,
			8257.8,
			7500,
		}

		got, err := ticker.FromRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Ticker{
			Symbol:          "tBTCUSD",
			Bid:             7616.5,
			BidSize:         31.89055171,
			Ask:             7617.5,
			AskSize:         43.358118629999986,
			DailyChange:     -550.8,
			DailyChangePerc: -0.0674,
			LastPrice:       7617.1,
			Volume:          8314.71200815,
			High:            8257.8,
			Low:             7500,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("funding pair spanshot ticker", func(t *testing.T) {
		payload := []interface{}{
			0.0003447013698630137,
			0.000316,
			30,
			1682003.0922634401,
			0.00031783,
			4,
			23336.958053110004,
			0.00000707,
			0.0209,
			0.0003446,
			156123478.78447533,
			0.000347,
			0.00009,
			nil,
			nil,
			146247919.52883354,
		}

		got, err := ticker.FromRaw("fUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Ticker{
			Symbol:             "fUSD",
			Frr:                0.0003447013698630137,
			Bid:                0.000316,
			BidPeriod:          30,
			BidSize:            1682003.0922634401,
			Ask:                0.00031783,
			AskPeriod:          4,
			AskSize:            23336.958053110004,
			DailyChange:        0.00000707,
			DailyChangePerc:    0.0209,
			LastPrice:          0.0003446,
			Volume:             156123478.78447533,
			High:               0.000347,
			Low:                0.00009,
			FrrAmountAvailable: 146247919.52883354,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("funding pair update ticker", func(t *testing.T) {
		payload := []interface{}{
			0.0003447013698630137,
			0.000316,
			30,
			1682003.0922634401,
			0.00031783,
			4,
			23336.958053110004,
			0.00000707,
			0.0209,
			0.0003446,
			156123478.78447533,
			0.000347,
			0.00009,
		}

		got, err := ticker.FromRaw("fUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Ticker{
			Symbol:          "fUSD",
			Frr:             0.0003447013698630137,
			Bid:             0.000316,
			BidPeriod:       30,
			BidSize:         1682003.0922634401,
			Ask:             0.00031783,
			AskPeriod:       4,
			AskSize:         23336.958053110004,
			DailyChange:     0.00000707,
			DailyChangePerc: 0.0209,
			LastPrice:       0.0003446,
			Volume:          156123478.78447533,
			High:            0.000347,
			Low:             0.00009,
		}

		assert.Equal(t, expected, got)
	})
}

func TestSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := [][]interface{}{}
		got, err := ticker.SnapshotFromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := [][]interface{}{
			{
				9206.8,
				21.038892079999997,
				9205.9,
				31.41755015,
				38.97852525,
				0.0043,
				9205.87852525,
				1815.68824558,
				9299,
				9111.8,
			},
			{
				9206.8,
			},
		}
		got, err := ticker.SnapshotFromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("trading pair snapshot", func(t *testing.T) {
		payload := [][]interface{}{
			{
				9206.8,
				21.038892079999997,
				9205.9,
				31.41755015,
				38.97852525,
				0.0043,
				9205.87852525,
				1815.68824558,
				9299,
				9111.8,
			},
		}

		got, err := ticker.SnapshotFromRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Snapshot{
			Snapshot: []*ticker.Ticker{
				{
					Symbol:          "tBTCUSD",
					Bid:             9206.8,
					BidSize:         21.038892079999997,
					Ask:             9205.9,
					AskSize:         31.41755015,
					DailyChange:     38.97852525,
					DailyChangePerc: 0.0043,
					LastPrice:       9205.87852525,
					Volume:          1815.68824558,
					High:            9299,
					Low:             9111.8,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestFromWSRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{7616.5, 31.89055171}
		got, err := ticker.FromWSRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("missing arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := ticker.FromWSRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("trading pair ticker", func(t *testing.T) {
		payload := []interface{}{
			7616.5,
			31.89055171,
			7617.5,
			43.358118629999986,
			-550.8,
			-0.0674,
			7617.1,
			8314.71200815,
			8257.8,
			7500,
		}

		got, err := ticker.FromWSRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Ticker{
			Symbol:          "tBTCUSD",
			Bid:             7616.5,
			BidSize:         31.89055171,
			Ask:             7617.5,
			AskSize:         43.358118629999986,
			DailyChange:     -550.8,
			DailyChangePerc: -0.0674,
			LastPrice:       7617.1,
			Volume:          8314.71200815,
			High:            8257.8,
			Low:             7500,
		}

		assert.Equal(t, expected, got)
	})

	t.Run("valid trading pair ticker snapshot", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				9206.8,
				21.038892079999997,
				9205.9,
				31.41755015,
				38.97852525,
				0.0043,
				9205.87852525,
				1815.68824558,
				9299,
				9111.8,
			},
		}

		got, err := ticker.FromWSRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Snapshot{
			Snapshot: []*ticker.Ticker{
				{
					Symbol:          "tBTCUSD",
					Bid:             9206.8,
					BidSize:         21.038892079999997,
					Ask:             9205.9,
					AskSize:         31.41755015,
					DailyChange:     38.97852525,
					DailyChangePerc: 0.0043,
					LastPrice:       9205.87852525,
					Volume:          1815.68824558,
					High:            9299,
					Low:             9111.8,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}
