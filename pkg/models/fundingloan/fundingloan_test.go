package fundingloan_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/stretchr/testify/assert"
)

func TestFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingloan.Loan
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{"exchange"},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest funding loans item": {
			pld: []interface{}{
				2995368, "fUST", 0, 1574077517000, 1574077517000, 100, nil, "ACTIVE", "FIXED",
				nil, nil, 0.0024, 2, 1574077517000, 1574077517000, 0, nil, nil, 0, nil, 0,
			},
			expected: &fundingloan.Loan{
				ID:            2995368,
				Symbol:        "fUST",
				Side:          0,
				MTSCreated:    1574077517000,
				MTSUpdated:    1574077517000,
				Amount:        100,
				Status:        "ACTIVE",
				RateType:      "FIXED",
				Rate:          0.0024,
				Period:        2,
				MTSOpened:     1574077517000,
				MTSLastPayout: 1574077517000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest funding loans history item": {
			pld: []interface{}{
				13683223, "fBTC", -1, 1575446268000, 1575446644000, 0.02, nil, "CLOSED (used)", "FIXED",
				nil, nil, 0, 7, 1575446268000, 1575446643000, nil, 0, nil, 0, nil, 0,
			},
			expected: &fundingloan.Loan{
				ID:            13683223,
				Symbol:        "fBTC",
				Side:          -1,
				MTSCreated:    1575446268000,
				MTSUpdated:    1575446644000,
				Amount:        0.02,
				Status:        "CLOSED (used)",
				RateType:      "FIXED",
				Rate:          0,
				Period:        7,
				MTSOpened:     1575446268000,
				MTSLastPayout: 1575446643000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fls item": {
			pld: []interface{}{
				2995442, "fUSD", -1, 1575291961000, 1575295850000, 820, 0, "ACTIVE", nil,
				nil, nil, 0.002, 7, 1575282446000, 1575295850000, 0, 0, nil, 0, nil, 0,
			},
			expected: &fundingloan.Loan{
				ID:            2995442,
				Symbol:        "fUSD",
				Side:          -1,
				MTSCreated:    1575291961000,
				MTSUpdated:    1575295850000,
				Amount:        820,
				Status:        "ACTIVE",
				Rate:          0.002,
				Period:        7,
				MTSOpened:     1575282446000,
				MTSLastPayout: 1575295850000,
				Notify:        false,
				Hidden:        false,
				Insure:        false,
				Renew:         false,
				RateReal:      0,
				NoClose:       false,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := fundingloan.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *fundingloan.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest funding loans": {
			pld: []interface{}{
				[]interface{}{
					2995368, "fUST", 0, 1574077517000, 1574077517000, 100, nil, "ACTIVE", "FIXED",
					nil, nil, 0.0024, 2, 1574077517000, 1574077517000, 0, nil, nil, 0, nil, 0,
				},
			},
			expected: &fundingloan.Snapshot{
				Snapshot: []*fundingloan.Loan{
					{
						ID:            2995368,
						Symbol:        "fUST",
						Side:          0,
						MTSCreated:    1574077517000,
						MTSUpdated:    1574077517000,
						Amount:        100,
						Status:        "ACTIVE",
						RateType:      "FIXED",
						Rate:          0.0024,
						Period:        2,
						MTSOpened:     1574077517000,
						MTSLastPayout: 1574077517000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest funding loan history": {
			pld: []interface{}{
				[]interface{}{
					13683223, "fBTC", -1, 1575446268000, 1575446644000, 0.02, nil, "CLOSED (used)", "FIXED",
					nil, nil, 0, 7, 1575446268000, 1575446643000, nil, 0, nil, 0, nil, 0,
				},
			},
			expected: &fundingloan.Snapshot{
				Snapshot: []*fundingloan.Loan{
					{
						ID:            13683223,
						Symbol:        "fBTC",
						Side:          -1,
						MTSCreated:    1575446268000,
						MTSUpdated:    1575446644000,
						Amount:        0.02,
						Status:        "CLOSED (used)",
						RateType:      "FIXED",
						Rate:          0,
						Period:        7,
						MTSOpened:     1575446268000,
						MTSLastPayout: 1575446643000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws fls": {
			pld: []interface{}{
				[]interface{}{
					2995442, "fUSD", -1, 1575291961000, 1575295850000, 820, 0, "ACTIVE", nil,
					nil, nil, 0.002, 7, 1575282446000, 1575295850000, 0, 0, nil, 0, nil, 0,
				},
			},
			expected: &fundingloan.Snapshot{
				Snapshot: []*fundingloan.Loan{
					{
						ID:            2995442,
						Symbol:        "fUSD",
						Side:          -1,
						MTSCreated:    1575291961000,
						MTSUpdated:    1575295850000,
						Amount:        820,
						Status:        "ACTIVE",
						Rate:          0.002,
						Period:        7,
						MTSOpened:     1575282446000,
						MTSLastPayout: 1575295850000,
						Notify:        false,
						Hidden:        false,
						Insure:        false,
						Renew:         false,
						RateReal:      0,
						NoClose:       false,
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
			got, err := fundingloan.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestNewFromRaw(t *testing.T) {
	pld := []interface{}{
		2995368, "fUST", 0, 1574077517000, 1574077517000, 100, nil, "ACTIVE", nil,
		nil, nil, 0.0024, 2, 1574077517000, 1574077517000, 0, nil, nil, 1, nil, 0,
	}

	expected := "fundingloan.New"
	o, err := fundingloan.NewFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestUpdateFromRaw(t *testing.T) {
	pld := []interface{}{
		2995368, "fUST", 0, 1574077517000, 1574077517000, 100, nil, "ACTIVE", nil,
		nil, nil, 0.0024, 2, 1574077517000, 1574077517000, 0, nil, nil, 1, nil, 0,
	}

	expected := "fundingloan.Update"
	o, err := fundingloan.UpdateFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}

func TestCancelFromRaw(t *testing.T) {
	pld := []interface{}{
		2995368, "fUST", 0, 1574077517000, 1574077517000, 100, nil, "ACTIVE", nil,
		nil, nil, 0.0024, 2, 1574077517000, 1574077517000, 0, nil, nil, 1, nil, 0,
	}

	expected := "fundingloan.Cancel"
	o, err := fundingloan.CancelFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(o).String()
	assert.Equal(t, expected, got)
}
