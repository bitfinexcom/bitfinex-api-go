package fundingoffer_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingOfferFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{652606505}

		got, err := fundingoffer.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
		}

		got, err := fundingoffer.FromRaw(payload)
		require.Nil(t, err)

		expected := &fundingoffer.Offer{
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
		}
		assert.Equal(t, expected, got)
	})
}

func TestFundingOfferSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := fundingoffer.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
			[]interface{}{652606506},
		}
		got, err := fundingoffer.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
			[]interface{}{
				652606506,
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
				1,
				nil,
				0,
				nil,
			},
		}

		got, err := fundingoffer.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &fundingoffer.Snapshot{
			Snapshot: []*fundingoffer.Offer{
				{
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
				{
					ID:         652606506,
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
					Hidden:     true,
					Insure:     false,
					Renew:      false,
					RateReal:   0,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestFundingOfferCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		focr := fundingoffer.CancelRequest{ID: 123}

		got, err := focr.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"foc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}

func TestFundingOfferSubmitRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		fosr := fundingoffer.SubmitRequest{
			Type:   "LIMIT",
			Symbol: "fETH",
			Amount: 0.29797676,
			Rate:   0.0002,
			Period: 2,
		}
		got, err := fosr.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"fon\", null, {\"type\":\"LIMIT\",\"symbol\":\"fETH\",\"amount\":\"0.29797676\",\"rate\":\"0.0002\",\"period\":2}]"
		assert.Equal(t, expected, string(got))
	})
}
