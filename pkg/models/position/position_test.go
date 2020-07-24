package position_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPositionFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{"tBTCUSD"},
		}

		p, err := position.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, p)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"tBTCUSD",
			"ACTIVE",
			0.0195,
			8565.0267019,
			0,
			0,
			-0.33455568705000516,
			-0.0003117550117425625,
			7045.876419249083,
			3.0673001895895604,
			nil,
			142355652,
		}

		p, err := position.FromRaw(payload)
		require.Nil(t, err)

		expected := &position.Position{
			Id:                   142355652,
			Symbol:               "tBTCUSD",
			Status:               "ACTIVE",
			Amount:               0.0195,
			BasePrice:            8565.0267019,
			MarginFunding:        0,
			MarginFundingType:    0,
			ProfitLoss:           -0.33455568705000516,
			ProfitLossPercentage: -0.0003117550117425625,
			LiquidationPrice:     7045.876419249083,
			Leverage:             3.0673001895895604,
		}

		assert.Equal(t, expected, p)
	})
}

func TestNewPositionSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{"tBTCUSD"},
			[]interface{}{
				"tBTCUSD",
				"ACTIVE",
				0.0195,
				8565.0267019,
				0,
				0,
				-0.33455568705000516,
				-0.0003117550117425625,
				7045.876419249083,
				3.0673001895895604,
				nil,
				142355652,
			},
		}

		p, err := position.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, p)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				"tBTCUSD",
				"ACTIVE",
				0.0195,
				8565.0267019,
				0,
				0,
				-0.33455568705000516,
				-0.0003117550117425625,
				7045.876419249083,
				3.0673001895895604,
				nil,
				142355652,
			},
			[]interface{}{
				"tBTCUSD",
				"ACTIVE",
				0.0195,
				8575,
				0,
				0,
				-0.33455568705000516,
				-0.0003117550117425625,
				7045.876419249083,
				3.0673001895895604,
				nil,
				142355653,
			},
		}

		p, err := position.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &position.Snapshot{
			Snapshot: []*position.Position{
				{
					Id:                   142355652,
					Symbol:               "tBTCUSD",
					Status:               "ACTIVE",
					Amount:               0.0195,
					BasePrice:            8565.0267019,
					MarginFunding:        0,
					MarginFundingType:    0,
					ProfitLoss:           -0.33455568705000516,
					ProfitLossPercentage: -0.0003117550117425625,
					LiquidationPrice:     7045.876419249083,
					Leverage:             3.0673001895895604,
				},
				{
					Id:                   142355653,
					Symbol:               "tBTCUSD",
					Status:               "ACTIVE",
					Amount:               0.0195,
					BasePrice:            8575,
					MarginFunding:        0,
					MarginFundingType:    0,
					ProfitLoss:           -0.33455568705000516,
					ProfitLossPercentage: -0.0003117550117425625,
					LiquidationPrice:     7045.876419249083,
					Leverage:             3.0673001895895604,
				},
			},
		}

		assert.Equal(t, expected, p)
	})
}
