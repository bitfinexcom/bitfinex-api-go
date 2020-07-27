package stats_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatsFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{1573554000000}

		got, err := stats.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{1573554000000, 25957.94278561}

		got, err := stats.FromRaw(payload)
		require.Nil(t, err)

		expected := &stats.Stat{
			Period: 1573554000000,
			Volume: 25957.94278561,
		}
		assert.Equal(t, expected, got)
	})
}

func TestNewStatsSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments #1", func(t *testing.T) {
		payload := []interface{}{}

		got, err := stats.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("invalid arguments #2", func(t *testing.T) {
		payload := []interface{}{1573554000000}

		got, err := stats.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{1573554000000, 25957.94278561},
			[]interface{}{1573553940000, 25964.29056204},
		}

		got, err := stats.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := []*stats.Stat{
			{
				Period: 1573554000000,
				Volume: 25957.94278561,
			},
			{
				Period: 1573553940000,
				Volume: 25964.29056204,
			},
		}
		assert.Equal(t, expected, got)
	})
}
