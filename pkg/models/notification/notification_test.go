package notification_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/notification"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNotificationFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{1575282446000}

		got, err := notification.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("missing notification info", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"on-req",
			1234,
			nil,
			nil,
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:        1575282446000,
			Type:       "on-req",
			MessageID:  1234,
			NotifyInfo: nil,
			Code:       4567,
			Status:     "SUCCESS",
			Text:       "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("notification info present but empty", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"on-req",
			1234,
			nil,
			[]interface{}{},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:        1575282446000,
			Type:       "on-req",
			MessageID:  1234,
			NotifyInfo: nil,
			Code:       4567,
			Status:     "SUCCESS",
			Text:       "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'on-req' raw snapshot", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"on-req",
			1234,
			nil,
			[]interface{}{
				[]interface{}{
					33950998275,
					nil,
					1573476747887,
					"tETHUSD",
					1573476748000,
					1573476748000,
					-0.5,
					-0.5,
					"LIMIT",
					nil,
					nil,
					nil,
					0,
					"ACTIVE",
					nil,
					nil,
					220,
					0,
					0,
					0,
					nil,
					nil,
					nil,
					0,
					1,
					nil,
					nil,
					nil,
					"BFX",
					nil,
					nil,
					nil,
				},
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "on-req",
			MessageID: 1234,
			NotifyInfo: &order.Snapshot{
				Snapshot: []*order.Order{
					{
						ID:            33950998275,
						GID:           0,
						CID:           1573476747887,
						Symbol:        "tETHUSD",
						MTSCreated:    1573476748000,
						MTSUpdated:    1573476748000,
						Amount:        -0.5,
						AmountOrig:    -0.5,
						Type:          "LIMIT",
						TypePrev:      "",
						MTSTif:        0,
						Flags:         0,
						Status:        "ACTIVE",
						Price:         220,
						PriceAvg:      0,
						PriceTrailing: 0,
						PriceAuxLimit: 0,
						Notify:        false,
						Hidden:        true,
						PlacedID:      0,
						Meta:          map[string]interface{}{},
					},
				},
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'on-req' single raw", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"on-req",
			1234,
			nil,
			[]interface{}{
				33950998275,
				nil,
				1573476747887,
				"tETHUSD",
				1573476748000,
				1573476748000,
				-0.5,
				-0.5,
				"LIMIT",
				nil,
				nil,
				nil,
				0,
				"ACTIVE",
				nil,
				nil,
				220,
				0,
				0,
				0,
				nil,
				nil,
				nil,
				0,
				1,
				nil,
				nil,
				nil,
				"BFX",
				nil,
				nil,
				nil,
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "on-req",
			MessageID: 1234,
			NotifyInfo: &order.New{
				ID:            33950998275,
				GID:           0,
				CID:           1573476747887,
				Symbol:        "tETHUSD",
				MTSCreated:    1573476748000,
				MTSUpdated:    1573476748000,
				Amount:        -0.5,
				AmountOrig:    -0.5,
				Type:          "LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "ACTIVE",
				Price:         220,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        true,
				PlacedID:      0,
				Meta:          map[string]interface{}{},
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'ou-req' single raw", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"ou-req",
			1234,
			nil,
			[]interface{}{
				33950998276,
				nil,
				1573476747887,
				"tETHUSD",
				1573476748000,
				1573476748000,
				-0.5,
				-0.5,
				"LIMIT",
				nil,
				nil,
				nil,
				0,
				"ACTIVE",
				nil,
				nil,
				220,
				0,
				0,
				0,
				nil,
				nil,
				nil,
				0,
				1,
				nil,
				nil,
				nil,
				"BFX",
				nil,
				nil,
				nil,
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "ou-req",
			MessageID: 1234,
			NotifyInfo: &order.Update{
				ID:            33950998276,
				GID:           0,
				CID:           1573476747887,
				Symbol:        "tETHUSD",
				MTSCreated:    1573476748000,
				MTSUpdated:    1573476748000,
				Amount:        -0.5,
				AmountOrig:    -0.5,
				Type:          "LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "ACTIVE",
				Price:         220,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        true,
				PlacedID:      0,
				Meta:          map[string]interface{}{},
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'oc-req' single raw", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"oc-req",
			1234,
			nil,
			[]interface{}{
				33950998277,
				nil,
				1573476747887,
				"tETHUSD",
				1573476748000,
				1573476748000,
				-0.5,
				-0.5,
				"LIMIT",
				nil,
				nil,
				nil,
				0,
				"ACTIVE",
				nil,
				nil,
				220,
				0,
				0,
				0,
				nil,
				nil,
				nil,
				0,
				1,
				nil,
				nil,
				nil,
				"BFX",
				nil,
				nil,
				nil,
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "oc-req",
			MessageID: 1234,
			NotifyInfo: &order.Cancel{
				ID:            33950998277,
				GID:           0,
				CID:           1573476747887,
				Symbol:        "tETHUSD",
				MTSCreated:    1573476748000,
				MTSUpdated:    1573476748000,
				Amount:        -0.5,
				AmountOrig:    -0.5,
				Type:          "LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "ACTIVE",
				Price:         220,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        true,
				PlacedID:      0,
				Meta:          map[string]interface{}{},
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'fon-req' single raw", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"fon-req",
			1234,
			nil,
			[]interface{}{
				652606505,
				"fETH",
				1574000611000,
				1574000611000,
				0.29797676,
				0.29797676,
				"LIMIT",
				nil,
				nil,
				0,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0002,
				2,
				1,
				nil,
				nil,
				0,
				nil,
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "fon-req",
			MessageID: 1234,
			NotifyInfo: &fundingoffer.New{
				ID:         652606505,
				Symbol:     "fETH",
				MTSCreated: 1574000611000,
				MTSUpdated: 1574000611000,
				Amount:     0.29797676,
				AmountOrig: 0.29797676,
				Type:       "LIMIT",
				Flags:      0,
				Status:     "ACTIVE",
				Rate:       0.0002,
				Period:     2,
				Notify:     true,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'foc-req' single raw", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"foc-req",
			1234,
			nil,
			[]interface{}{
				652606505,
				"fETH",
				1574000611000,
				1574000611000,
				0.29797676,
				0.29797676,
				"LIMIT",
				nil,
				nil,
				0,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0002,
				2,
				1,
				nil,
				nil,
				0,
				nil,
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "foc-req",
			MessageID: 1234,
			NotifyInfo: &fundingoffer.Cancel{
				ID:         652606505,
				Symbol:     "fETH",
				MTSCreated: 1574000611000,
				MTSUpdated: 1574000611000,
				Amount:     0.29797676,
				AmountOrig: 0.29797676,
				Type:       "LIMIT",
				Flags:      0,
				Status:     "ACTIVE",
				Rate:       0.0002,
				Period:     2,
				Notify:     true,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'uca'", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"uca",
			1234,
			nil,
			[]interface{}{
				"fETH",
				1574000611000,
				"LIMIT",
			},
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "uca",
			MessageID: 1234,
			NotifyInfo: []interface{}{
				"fETH",
				1574000611000,
				"LIMIT",
			},
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})

	t.Run("'pc'", func(t *testing.T) {
		payload := []interface{}{
			1575282446000,
			"pc",
			1234,
			nil,
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
			4567,
			"SUCCESS",
			"foo bar",
		}

		got, err := notification.FromRaw(payload)
		require.Nil(t, err)

		expected := &notification.Notification{
			MTS:       1575282446000,
			Type:      "pc",
			MessageID: 1234,
			NotifyInfo: &position.Cancel{
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
			Code:   4567,
			Status: "SUCCESS",
			Text:   "foo bar",
		}

		assert.Equal(t, expected, got)
	})
}
