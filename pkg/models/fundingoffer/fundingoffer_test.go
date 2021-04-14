package fundingoffer_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/stretchr/testify/assert"
)

func TestFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingoffer.Offer
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{"exchange"},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest active funding offer": {
			pld: []interface{}{
				652606505, "fETH", 1574000611000, 1574000611000, 0.29797676, 0.29797676,
				"LIMIT", nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0002, 2, 0, nil, nil, 0, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         652606505,
				Symbol:     "fETH",
				MTSCreated: 1574000611000,
				MTSUpdated: 1574000611000,
				Amount:     0.29797676,
				AmountOrig: 0.29797676,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Rate:       0.0002,
				Period:     2,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest submit funding offer": {
			pld: []interface{}{
				604366339, "fUSD", 1568713496502, 1568713496502, 50, 50, "LIMIT", nil,
				nil, nil, "ACTIVE", nil, nil, nil, 0.00002, 2, false, nil, nil, false, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         604366339,
				Symbol:     "fUSD",
				MTSCreated: 1568713496502,
				MTSUpdated: 1568713496502,
				Amount:     50,
				AmountOrig: 50,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Rate:       2e-05,
				Period:     2,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest cancel funding offer": {
			pld: []interface{}{
				604393839, "fUSD", 1568716545000, 1568716545000, 50, 50, "LIMIT", nil,
				nil, nil, "ACTIVE", nil, nil, nil, 0.06, 2, false, nil, nil, false, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         604393839,
				Symbol:     "fUSD",
				MTSCreated: 1568716545000,
				MTSUpdated: 1568716545000,
				Amount:     50,
				AmountOrig: 50,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Rate:       0.06,
				Period:     2,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest funding offer hist item": {
			pld: []interface{}{
				653170899, "fUSD", 1574072620000, 1574072620000, 0,
				-57.9, nil, nil, nil, nil, "EXECUTED at 0.0368% (57.9)",
				nil, nil, nil, 0.000369, 2, 0, 0, nil, nil, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         653170899,
				Symbol:     "fUSD",
				MTSCreated: 1574072620000,
				MTSUpdated: 1574072620000,
				Amount:     0,
				AmountOrig: -57.9,
				Status:     "EXECUTED at 0.0368% (57.9)",
				Rate:       0.000369,
				Period:     2,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fon fou foc": {
			pld: []interface{}{
				41238747, "fUST", 1575026670000, 1575026670000, 5000, 5000, "LIMIT", nil,
				nil, 0, "ACTIVE", nil, nil, nil, 0.006000000000000001, 30, 0, 0, nil, 0, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         41238747,
				Symbol:     "fUST",
				MTSCreated: 1575026670000,
				MTSUpdated: 1575026670000,
				Amount:     5000,
				AmountOrig: 5000,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Rate:       0.006000000000000001,
				Period:     30,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fos item": {
			pld: []interface{}{
				41237920, "fETH", 1573912039000, 1573912039000, 0.5, 0.5, "LIMIT",
				nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0024, 2, 0, 0, nil, 0, nil,
			},
			expected: &fundingoffer.Offer{
				ID:         41237920,
				Symbol:     "fETH",
				MTSCreated: 1573912039000,
				MTSUpdated: 1573912039000,
				Amount:     0.5,
				AmountOrig: 0.5,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Rate:       0.0024,
				Period:     2,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := fundingoffer.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingoffer.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest funding offer hist": {
			pld: []interface{}{
				[]interface{}{
					653170899, "fUSD", 1574072620000, 1574072620000, 0, -57.9, nil, nil, nil, nil,
					"EXECUTED at 0.0368% (57.9)", nil, nil, nil, 0.000369, 2, 0, 0, nil, nil, nil,
				},
			},
			expected: &fundingoffer.Snapshot{
				Snapshot: []*fundingoffer.Offer{
					{
						ID:         653170899,
						Symbol:     "fUSD",
						MTSCreated: 1574072620000,
						MTSUpdated: 1574072620000,
						Amount:     0,
						AmountOrig: -57.9,
						Status:     "EXECUTED at 0.0368% (57.9)",
						Rate:       0.000369,
						Period:     2,
						Notify:     false,
						Hidden:     false,
						Insure:     false,
						Renew:      false,
						RateReal:   0,
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fos": {
			pld: []interface{}{
				[]interface{}{
					41237920, "fETH", 1573912039000, 1573912039000, 0.5, 0.5, "LIMIT",
					nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0024, 2, 0, 0, nil, 0, nil,
				},
			},
			expected: &fundingoffer.Snapshot{
				Snapshot: []*fundingoffer.Offer{
					{
						ID:         41237920,
						Symbol:     "fETH",
						MTSCreated: 1573912039000,
						MTSUpdated: 1573912039000,
						Amount:     0.5,
						AmountOrig: 0.5,
						Type:       "LIMIT",
						Status:     "ACTIVE",
						Rate:       0.0024,
						Period:     2,
						Notify:     false,
						Hidden:     false,
						Insure:     false,
						Renew:      false,
						RateReal:   0,
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := fundingoffer.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestCancelFromRaw(t *testing.T) {
	pld := []interface{}{
		652606505, "fETH", 1574000611000, 1574000611000, 0.29797676, 0.29797676, "LIMIT",
		nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0002, 2, 1, nil, nil, 0, nil,
	}

	expected := "fundingoffer.Cancel"
	o, err := fundingoffer.CancelFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestNewFromRaw(t *testing.T) {
	pld := []interface{}{
		652606505, "fETH", 1574000611000, 1574000611000, 0.29797676, 0.29797676, "LIMIT",
		nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0002, 2, 1, nil, nil, 0, nil,
	}

	expected := "fundingoffer.New"
	o, err := fundingoffer.NewFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestUpdateFromRaw(t *testing.T) {
	pld := []interface{}{
		652606505, "fETH", 1574000611000, 1574000611000, 0.29797676, 0.29797676, "LIMIT",
		nil, nil, 0, "ACTIVE", nil, nil, nil, 0.0002, 2, 1, nil, nil, 0, nil,
	}

	expected := "fundingoffer.Update"
	o, err := fundingoffer.UpdateFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}
