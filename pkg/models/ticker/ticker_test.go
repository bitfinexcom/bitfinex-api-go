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

func TestSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := [][]interface{}{}
		got, err := ticker.SnapshotFromRaw("tBTCUSD", payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
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
				9205.8,
				21.038892079999997,
				9205.9,
				31.41755015,
				38.97852525,
				0.0043,
				9205.87852525,
				1815.68824558,
				9299.1234,
				9111.8,
			},
		}

		got, err := ticker.SnapshotFromRaw("tBTCUSD", payload)
		require.Nil(t, err)

		expected := &ticker.Snapshot{
			Snapshot: []*ticker.Ticker{
				{
					Symbol:          "tBTCUSD",
					Frr:             0,
					Bid:             9206.8,
					BidPeriod:       0,
					BidSize:         21.038892079999997,
					Ask:             9205.9,
					AskPeriod:       0,
					AskSize:         31.41755015,
					DailyChange:     38.97852525,
					DailyChangePerc: 0.0043,
					LastPrice:       9205.87852525,
					Volume:          1815.68824558,
					High:            9299,
					Low:             9111.8,
				},
				{
					Symbol:          "tBTCUSD",
					Frr:             0,
					Bid:             9205.8,
					BidPeriod:       0,
					BidSize:         21.038892079999997,
					Ask:             9205.9,
					AskPeriod:       0,
					AskSize:         31.41755015,
					DailyChange:     38.97852525,
					DailyChangePerc: 0.0043,
					LastPrice:       9205.87852525,
					Volume:          1815.68824558,
					High:            9299.1234,
					Low:             9111.8,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}
