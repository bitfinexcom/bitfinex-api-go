package fundingloan_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFundingLoanFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{2995368}

		got, err := fundingloan.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			2995368,
			"fUST",
			0,
			1574077517000,
			1574077517000,
			100,
			nil,
			"ACTIVE",
			nil,
			nil,
			nil,
			0.0024,
			2,
			1574077517000,
			1574077517000,
			0,
			nil,
			nil,
			1,
			nil,
			0,
		}

		got, err := fundingloan.FromRaw(payload)
		require.Nil(t, err)

		expected := &fundingloan.Loan{
			ID:            2995368,
			Symbol:        "fUST",
			Side:          "",
			MTSCreated:    1574077517000,
			MTSUpdated:    1574077517000,
			Amount:        100,
			Status:        "ACTIVE",
			Rate:          0.0024,
			Period:        2,
			MTSOpened:     1574077517000,
			MTSLastPayout: 1574077517000,
			Notify:        false,
			Hidden:        false,
			Insure:        false,
			Renew:         true,
			RateReal:      0,
			NoClose:       false,
		}
		assert.Equal(t, expected, got)
	})
}

func TestFundingLoanSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := fundingloan.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				2995368,
				"fUST",
				0,
				1574077517000,
				1574077517000,
				100,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574077517000,
				1574077517000,
				0,
				nil,
				nil,
				1,
				nil,
				0,
			},
			[]interface{}{2995369},
		}
		got, err := fundingloan.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				2995368,
				"fUST",
				0,
				1574077517000,
				1574077517000,
				100,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574077517000,
				1574077517000,
				0,
				nil,
				nil,
				1,
				nil,
				0,
			},
			[]interface{}{
				2995369,
				"fUST",
				0,
				1574077517000,
				1574077517000,
				100,
				nil,
				"ACTIVE",
				nil,
				nil,
				nil,
				0.0024,
				2,
				1574077517000,
				1574077517000,
				1,
				1,
				nil,
				1,
				nil,
				1,
			},
		}

		got, err := fundingloan.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &fundingloan.Snapshot{
			Snapshot: []*fundingloan.Loan{
				{
					ID:            2995368,
					Symbol:        "fUST",
					Side:          "",
					MTSCreated:    1574077517000,
					MTSUpdated:    1574077517000,
					Amount:        100,
					Status:        "ACTIVE",
					Rate:          0.0024,
					Period:        2,
					MTSOpened:     1574077517000,
					MTSLastPayout: 1574077517000,
					Notify:        false,
					Hidden:        false,
					Insure:        false,
					Renew:         true,
					RateReal:      0,
					NoClose:       false,
				},
				{
					ID:            2995369,
					Symbol:        "fUST",
					Side:          "",
					MTSCreated:    1574077517000,
					MTSUpdated:    1574077517000,
					Amount:        100,
					Status:        "ACTIVE",
					Rate:          0.0024,
					Period:        2,
					MTSOpened:     1574077517000,
					MTSLastPayout: 1574077517000,
					Notify:        true,
					Hidden:        true,
					Insure:        false,
					Renew:         true,
					RateReal:      0,
					NoClose:       true,
				},
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestFundingLoanCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		flcr := fundingloan.CancelRequest{ID: 123}
		got, err := flcr.MarshalJSON()

		require.Nil(t, err)

		expected := "[0, \"flc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}
