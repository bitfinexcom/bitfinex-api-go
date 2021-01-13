package fundingcredit_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingCreditFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{2995368}

		got, err := fundingcredit.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			26222883,
			"fUST",
			1,
			1574013661000,
			1574079687000,
			350,
			nil,
			"ACTIVE",
			nil,
			nil,
			nil,
			0.0024,
			2,
			1574013661000,
			1574078487000,
			1,
			nil,
			nil,
			0,
			nil,
			1,
			"tBTCUST",
		}

		got, err := fundingcredit.FromRaw(payload)
		require.Nil(t, err)

		expected := &fundingcredit.Credit{
			ID:            26222883,
			Symbol:        "fUST",
			Side:          "",
			MTSCreated:    1574013661000,
			MTSUpdated:    1574079687000,
			Amount:        350,
			Status:        "ACTIVE",
			Rate:          0.0024,
			Period:        2,
			MTSOpened:     1574013661000,
			MTSLastPayout: 1574078487000,
			Notify:        true,
			Hidden:        false,
			Insure:        false,
			Renew:         false,
			RateReal:      0,
			NoClose:       true,
			PositionPair:  "tBTCUST",
		}
		assert.Equal(t, expected, got)
	})
}

func TestFundingCreditSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := fundingcredit.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				26222883,
				"fUST",
				1,
				1574013661000,
				1574079687000,
				350,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574013661000,
				1574078487000,
				1,
				nil,
				nil,
				0,
				nil,
				1,
				"tBTCUST",
			},
			[]interface{}{26222883},
		}
		got, err := fundingcredit.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				26222883,
				"fUST",
				1,
				1574013661000,
				1574079687000,
				350,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574013661000,
				1574078487000,
				1,
				nil,
				nil,
				0,
				nil,
				1,
				"tBTCUST",
			},
			[]interface{}{
				26222884,
				"fUST",
				1,
				1574013661000,
				1574079687000,
				350,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574013661000,
				1574078487000,
				1,
				nil,
				nil,
				1,
				nil,
				1,
				"tBTCUST",
			},
		}

		got, err := fundingcredit.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &fundingcredit.Snapshot{
			Snapshot: []*fundingcredit.Credit{
				{
					ID:            26222883,
					Symbol:        "fUST",
					Side:          "",
					MTSCreated:    1574013661000,
					MTSUpdated:    1574079687000,
					Amount:        350,
					Status:        "ACTIVE",
					Rate:          0.0024,
					Period:        2,
					MTSOpened:     1574013661000,
					MTSLastPayout: 1574078487000,
					Notify:        true,
					Hidden:        false,
					Insure:        false,
					Renew:         false,
					RateReal:      0,
					NoClose:       true,
					PositionPair:  "tBTCUST",
				},
				{
					ID:            26222884,
					Symbol:        "fUST",
					Side:          "",
					MTSCreated:    1574013661000,
					MTSUpdated:    1574079687000,
					Amount:        350,
					Status:        "ACTIVE",
					Rate:          0.0024,
					Period:        2,
					MTSOpened:     1574013661000,
					MTSLastPayout: 1574078487000,
					Notify:        true,
					Hidden:        false,
					Insure:        false,
					Renew:         true,
					RateReal:      0,
					NoClose:       true,
					PositionPair:  "tBTCUST",
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestFundingCreditCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		flcr := fundingcredit.CancelRequest{ID: 123}
		got, err := flcr.MarshalJSON()

		require.Nil(t, err)

		expected := "[0, \"fcc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}

func TestNewFromRaw(t *testing.T) {
	pld := []interface{}{
		26222883, "fUST", 1, 1574013661000, 1574079687000, 350, nil, "ACTIVE", nil, nil,
		nil, 0.0024, 2, 1574013661000, 1574078487000, 1, nil, nil, 0, nil, 1, "tBTCUST",
	}

	expected := "fundingcredit.New"
	o, err := fundingcredit.NewFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestUpdateFromRaw(t *testing.T) {
	pld := []interface{}{
		26222883, "fUST", 1, 1574013661000, 1574079687000, 350, nil, "ACTIVE", nil, nil,
		nil, 0.0024, 2, 1574013661000, 1574078487000, 1, nil, nil, 0, nil, 1, "tBTCUST",
	}

	expected := "fundingcredit.Update"
	o, err := fundingcredit.UpdateFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}
