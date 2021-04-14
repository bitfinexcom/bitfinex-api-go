package fundingcredit_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/stretchr/testify/assert"
)

func TestFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingcredit.Credit
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{"exchange"},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest funding credits item": {
			pld: []interface{}{
				26222883, "fUST", 1, 1574013661000, 1574079687000, 350, nil, "ACTIVE", "FIXED", nil,
				nil, 0.0024, 2, 1574013661000, 1574078487000, 0, nil, nil, 0, nil, 0, "tBTCUST",
			},
			expected: &fundingcredit.Credit{
				ID:            26222883,
				Symbol:        "fUST",
				Side:          1,
				MTSCreated:    1574013661000,
				MTSUpdated:    1574079687000,
				Amount:        350,
				Status:        "ACTIVE",
				RateType:      "FIXED",
				Rate:          0.0024,
				Period:        2,
				MTSOpened:     1574013661000,
				MTSLastPayout: 1574078487000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
				PositionPair:  "tBTCUST",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest funding credits history item": {
			pld: []interface{}{
				171988300, "fUSD", 1, 1574230085000, 1574402835000, 50.70511182, nil,
				"CLOSED (expired)", "FIXED", nil, nil, 0.00024799, 2, 1574230085000,
				1574403364000, nil, 0, nil, 0, nil, 0, "tEOSUSD",
			},
			expected: &fundingcredit.Credit{
				ID:            171988300,
				Symbol:        "fUSD",
				Side:          1,
				MTSCreated:    1574230085000,
				MTSUpdated:    1574402835000,
				Amount:        50.70511182,
				Status:        "CLOSED (expired)",
				RateType:      "FIXED",
				Rate:          0.00024799,
				Period:        2,
				MTSOpened:     1574230085000,
				MTSLastPayout: 1574403364000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
				PositionPair:  "tEOSUSD",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fcs item": {
			pld: []interface{}{
				26223578, "fUST", 1, 1575052261000, 1575296187000, 350, 0, "ACTIVE", nil, nil,
				nil, 0, 30, 1575052261000, 1575293487000, 0, 0, nil, 0, nil, 0, "tBTCUST",
			},
			expected: &fundingcredit.Credit{
				ID:            26223578,
				Symbol:        "fUST",
				Side:          1,
				MTSCreated:    1575052261000,
				MTSUpdated:    1575296187000,
				Amount:        350,
				Status:        "ACTIVE",
				Rate:          0,
				Period:        30,
				MTSOpened:     1575052261000,
				MTSLastPayout: 1575293487000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
				PositionPair:  "tBTCUST",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := fundingcredit.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingcredit.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest funding credits": {
			pld: []interface{}{
				[]interface{}{
					26222883, "fUST", 1, 1574013661000, 1574079687000, 350, nil, "ACTIVE", "FIXED", nil,
					nil, 0.0024, 2, 1574013661000, 1574078487000, 0, nil, nil, 0, nil, 0, "tBTCUST",
				},
			},
			expected: &fundingcredit.Snapshot{
				Snapshot: []*fundingcredit.Credit{
					{
						ID:            26222883,
						Symbol:        "fUST",
						Side:          1,
						MTSCreated:    1574013661000,
						MTSUpdated:    1574079687000,
						Amount:        350,
						Status:        "ACTIVE",
						RateType:      "FIXED",
						Rate:          0.0024,
						Period:        2,
						MTSOpened:     1574013661000,
						MTSLastPayout: 1574078487000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
						PositionPair:  "tBTCUST",
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest funding credits history": {
			pld: []interface{}{
				[]interface{}{
					171988300, "fUSD", 1, 1574230085000, 1574402835000, 50.70511182, nil,
					"CLOSED (expired)", "FIXED", nil, nil, 0.00024799, 2, 1574230085000,
					1574403364000, nil, 0, nil, 0, nil, 0, "tEOSUSD",
				},
			},
			expected: &fundingcredit.Snapshot{
				Snapshot: []*fundingcredit.Credit{
					{
						ID:            171988300,
						Symbol:        "fUSD",
						Side:          1,
						MTSCreated:    1574230085000,
						MTSUpdated:    1574402835000,
						Amount:        50.70511182,
						Status:        "CLOSED (expired)",
						RateType:      "FIXED",
						Rate:          0.00024799,
						Period:        2,
						MTSOpened:     1574230085000,
						MTSLastPayout: 1574403364000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
						PositionPair:  "tEOSUSD",
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fcs": {
			pld: []interface{}{
				[]interface{}{
					26223578, "fUST", 1, 1575052261000, 1575296187000, 350, 0, "ACTIVE", nil, nil,
					nil, 0, 30, 1575052261000, 1575293487000, 0, 0, nil, 0, nil, 0, "tBTCUST",
				},
			},
			expected: &fundingcredit.Snapshot{
				Snapshot: []*fundingcredit.Credit{
					{
						ID:            26223578,
						Symbol:        "fUST",
						Side:          1,
						MTSCreated:    1575052261000,
						MTSUpdated:    1575296187000,
						Amount:        350,
						Status:        "ACTIVE",
						Rate:          0,
						Period:        30,
						MTSOpened:     1575052261000,
						MTSLastPayout: 1575293487000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
						PositionPair:  "tBTCUST",
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
			got, err := fundingcredit.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
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

func TestCancelFromRaw(t *testing.T) {
	pld := []interface{}{
		26222883, "fUST", 1, 1574013661000, 1574079687000, 350, nil, "ACTIVE", nil, nil,
		nil, 0.0024, 2, 1574013661000, 1574078487000, 1, nil, nil, 0, nil, 1, "tBTCUST",
	}

	expected := "fundingcredit.Cancel"
	o, err := fundingcredit.CancelFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}
