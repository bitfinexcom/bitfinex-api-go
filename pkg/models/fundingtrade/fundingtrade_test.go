package fundingtrade_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingtrade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingTradeFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{636040}

		got, err := fundingtrade.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			636040,
			"fUST",
			1574077528000,
			41237922,
			-100,
			0.0024,
			2,
			nil,
		}

		got, err := fundingtrade.FromRaw(payload)
		require.Nil(t, err)

		expected := &fundingtrade.FundingTrade{
			ID:         636040,
			Symbol:     "fUST",
			MTSCreated: 1574077528000,
			OfferID:    41237922,
			Amount:     -100,
			Rate:       0.0024,
			Period:     2,
			Maker:      0,
		}
		assert.Equal(t, expected, got)
	})
}

func TestFundingTradeSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := fundingtrade.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				636040,
				"fUST",
				1574077528000,
				41237922,
				-100,
				0.0024,
				2,
				nil,
			},
			[]interface{}{636041},
		}
		got, err := fundingtrade.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				636040,
				"fUST",
				1574077528000,
				41237922,
				-100,
				0.0024,
				2,
				nil,
			},
			[]interface{}{
				636041,
				"fUST",
				1574077528000,
				41237922,
				-100,
				0.0025,
				2,
				nil,
			},
		}

		got, err := fundingtrade.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &fundingtrade.Snapshot{
			Snapshot: []*fundingtrade.FundingTrade{
				{
					ID:         636040,
					Symbol:     "fUST",
					MTSCreated: 1574077528000,
					OfferID:    41237922,
					Amount:     -100,
					Rate:       0.0024,
					Period:     2,
					Maker:      0,
				},
				{
					ID:         636041,
					Symbol:     "fUST",
					MTSCreated: 1574077528000,
					OfferID:    41237922,
					Amount:     -100,
					Rate:       0.0025,
					Period:     2,
					Maker:      0,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}
