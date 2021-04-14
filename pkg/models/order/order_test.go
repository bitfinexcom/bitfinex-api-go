package order_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/stretchr/testify/assert"
)

func TestFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *order.Order
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{"exchange"},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest retire orders": {
			pld: []interface{}{
				33950998275, nil, 1573476747887, "tETHUSD", 1573476748000, 1573476748000,
				-0.5, -0.5, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 220, 0, 0, 0, nil, nil,
				nil, 0, 0, nil, nil, nil, "BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         33950998275,
				CID:        1573476747887,
				Symbol:     "tETHUSD",
				MTSCreated: 1573476748000,
				MTSUpdated: 1573476748000,
				Amount:     -0.5,
				AmountOrig: -0.5,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Price:      220,
				Notify:     false,
				Hidden:     false,
				Routing:    "BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest submit order": {
			pld: []interface{}{
				30630788061, nil, 1567590617439, "tBTCUSD", 1567590617439, 1567590617439,
				0.001, 0.001, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 15, 0, 0, 0, nil, nil,
				nil, 0, nil, nil, nil, nil, "API>BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         30630788061,
				CID:        1567590617439,
				Symbol:     "tBTCUSD",
				MTSCreated: 1567590617439,
				MTSUpdated: 1567590617439,
				Amount:     0.001,
				AmountOrig: 0.001,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Price:      15,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest update order": {
			pld: []interface{}{
				30854813589, nil, 1568109670135, "tBTCUSD", 1568109673000, 1568109866000,
				0.002, 0.002, "LIMIT", "LIMIT", nil, nil, 0, "ACTIVE", nil, nil, 20, 0, 0,
				0, nil, nil, nil, 0, 0, nil, nil, nil, "API>BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         30854813589,
				CID:        1568109670135,
				Symbol:     "tBTCUSD",
				MTSCreated: 1568109673000,
				MTSUpdated: 1568109866000,
				Amount:     0.002,
				AmountOrig: 0.002,
				Type:       "LIMIT",
				TypePrev:   "LIMIT",
				Status:     "ACTIVE",
				Price:      20,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest cancel order": {
			pld: []interface{}{
				30937950333, nil, 1568298279766, "tBTCUSD", 1568298281000,
				1568298281000, 0.001, 0.001, "LIMIT", nil, nil, nil, 0,
				"ACTIVE", nil, nil, 15, 0, 0, 0, nil, nil, nil, 0, 0, nil,
				nil, nil, "API>BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         30937950333,
				CID:        1568298279766,
				Symbol:     "tBTCUSD",
				MTSCreated: 1568298281000,
				MTSUpdated: 1568298281000,
				Amount:     0.001,
				AmountOrig: 0.001,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Price:      15,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest order multi op": {
			pld: []interface{}{
				31492970019, nil, 1569347312970, "tBTCUSD", 1569347312971, 1569347312971, 0.001,
				0.001, "LIMIT", nil, nil, nil, 4096, "ACTIVE", nil, nil, 15, 0, 0, 0, nil, nil,
				nil, 0, nil, nil, nil, nil, "API>BFX", nil, nil, map[string]interface{}{"$F7": 1},
			},
			expected: &order.Order{
				ID:         31492970019,
				CID:        1569347312970,
				Symbol:     "tBTCUSD",
				MTSCreated: 1569347312971,
				MTSUpdated: 1569347312971,
				Amount:     0.001,
				AmountOrig: 0.001,
				Type:       "LIMIT",
				Flags:      4096,
				Status:     "ACTIVE",
				Price:      15,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
				Meta:       map[string]interface{}{"$F7": 1},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest cancel order multi": {
			pld: []interface{}{
				31123704044, nil, 1568711144715, "tBTCUSD", 1568711147000, 1568711147000,
				0.001, 0.001, "LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 15, 0, 0, 0, nil, nil,
				nil, 0, 0, nil, nil, nil, "API>BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         31123704044,
				CID:        1568711144715,
				Symbol:     "tBTCUSD",
				MTSCreated: 1568711147000,
				MTSUpdated: 1568711147000,
				Amount:     0.001,
				AmountOrig: 0.001,
				Type:       "LIMIT",
				Status:     "ACTIVE",
				Price:      15,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"rest orders history item": {
			pld: []interface{}{
				33961681942, "1227", 1337, "tBTCUSD", 1573482478000, 1573485373000,
				0.001, 0.001, "EXCHANGE LIMIT", nil, nil, nil, "0", "CANCELED", nil, nil,
				15, 0, 0, 0, nil, nil, nil, 0, 0, nil, nil, nil, "API>BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         33961681942,
				CID:        1337,
				Symbol:     "tBTCUSD",
				MTSCreated: 1573482478000,
				MTSUpdated: 1573485373000,
				Amount:     0.001,
				AmountOrig: 0.001,
				Type:       "EXCHANGE LIMIT",
				Status:     "CANCELED",
				Price:      15,
				Notify:     false,
				Hidden:     false,
				Routing:    "API>BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws os item": {
			pld: []interface{}{
				34930659963, nil, 1574955083558, "tETHUSD", 1574955083558, 1574955083573,
				0.201104, 0.201104, "EXCHANGE LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil, 120,
				0, 0, 0, nil, nil, nil, 0, 0, nil, nil, nil, "BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         34930659963,
				CID:        1574955083558,
				Symbol:     "tETHUSD",
				MTSCreated: 1574955083558,
				MTSUpdated: 1574955083573,
				Amount:     0.201104,
				AmountOrig: 0.201104,
				Type:       "EXCHANGE LIMIT",
				Status:     "ACTIVE",
				Price:      120,
				Notify:     false,
				Hidden:     false,
				Routing:    "BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws on ou oc": {
			pld: []interface{}{
				34930659963, nil, 1574955083558, "tETHUSD", 1574955083558, 1574955354487,
				0.201104, 0.201104, "EXCHANGE LIMIT", nil, nil, nil, 0, "CANCELED", nil, nil,
				120, 0, 0, 0, nil, nil, nil, 0, 0, nil, nil, nil, "BFX", nil, nil, nil,
			},
			expected: &order.Order{
				ID:         34930659963,
				CID:        1574955083558,
				Symbol:     "tETHUSD",
				MTSCreated: 1574955083558,
				MTSUpdated: 1574955354487,
				Amount:     0.201104,
				AmountOrig: 0.201104,
				Type:       "EXCHANGE LIMIT",
				Status:     "CANCELED",
				Price:      120,
				Notify:     false,
				Hidden:     false,
				Routing:    "BFX",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := order.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *order.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest orders hist": {
			pld: []interface{}{
				[]interface{}{
					33961681942, "1227", 1337, "tBTCUSD", 1573482478000, 1573485373000,
					0.001, 0.001, "EXCHANGE LIMIT", nil, nil, nil, "0", "CANCELED", nil, nil,
					15, 0, 0, 0, nil, nil, nil, 0, 0, nil, nil, nil, "API>BFX", nil, nil, nil,
				},
			},
			expected: &order.Snapshot{
				Snapshot: []*order.Order{
					{
						ID:         33961681942,
						CID:        1337,
						Symbol:     "tBTCUSD",
						MTSCreated: 1573482478000,
						MTSUpdated: 1573485373000,
						Amount:     0.001,
						AmountOrig: 0.001,
						Type:       "EXCHANGE LIMIT",
						Status:     "CANCELED",
						Price:      15,
						Notify:     false,
						Hidden:     false,
						Routing:    "API>BFX",
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws orders os": {
			pld: []interface{}{
				[]interface{}{
					34930659963, nil, 1574955083558, "tETHUSD", 1574955083558, 1574955083573,
					0.201104, 0.201104, "EXCHANGE LIMIT", nil, nil, nil, 0, "ACTIVE", nil, nil,
					120, 0, 0, 0, nil, nil, nil, 0, 0, nil, nil, nil, "BFX", nil, nil, nil,
				},
			},
			expected: &order.Snapshot{
				Snapshot: []*order.Order{
					{
						ID:         34930659963,
						CID:        1574955083558,
						Symbol:     "tETHUSD",
						MTSCreated: 1574955083558,
						MTSUpdated: 1574955083573,
						Amount:     0.201104,
						AmountOrig: 0.201104,
						Type:       "EXCHANGE LIMIT",
						Status:     "ACTIVE",
						Price:      120,
						Notify:     false,
						Hidden:     false,
						Routing:    "BFX",
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
			got, err := order.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
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
