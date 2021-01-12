package order_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrderFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{33950998275}

		got, err := order.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
		}

		got, err := order.FromRaw(payload)
		require.Nil(t, err)

		expected := &order.Order{
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
		}
		assert.Equal(t, expected, got)
	})
}

func TestOrdersSnapshotFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{}
		got, err := order.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("partially valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
			[]interface{}{33950998276},
		}
		got, err := order.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
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
		}

		got, err := order.SnapshotFromRaw(payload)
		require.Nil(t, err)

		expected := &order.Snapshot{
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
				{
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
			},
		}

		assert.Equal(t, expected, got)
	})
}

func TestNewFromRaw(t *testing.T) {
	pld := []interface{}{
		33950998276, nil, 1573476747887, "tETHUSD", 1573476748000, 1573476748000, -0.5,
		-0.5, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 220, 0, 0, 0, nil, nil,
		nil, 0, 1, nil, nil, nil, "BFX", nil, nil, nil,
	}
	expected := "order.New"
	o, err := order.NewFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestUpdateFromRaw(t *testing.T) {
	pld := []interface{}{
		33950998276, nil, 1573476747887, "tETHUSD", 1573476748000, 1573476748000, -0.5,
		-0.5, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 220, 0, 0, 0, nil, nil,
		nil, 0, 1, nil, nil, nil, "BFX", nil, nil, nil,
	}
	expected := "order.Update"
	o, err := order.UpdateFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestCancelFromRaw(t *testing.T) {
	pld := []interface{}{
		33950998276, nil, 1573476747887, "tETHUSD", 1573476748000, 1573476748000, -0.5,
		-0.5, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 220, 0, 0, 0, nil, nil,
		nil, 0, 1, nil, nil, nil, "BFX", nil, nil, nil,
	}
	expected := "order.Cancel"
	o, err := order.CancelFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}
