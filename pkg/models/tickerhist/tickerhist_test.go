package tickerhist_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tickerhist"
	"github.com/stretchr/testify/assert"
)

func TestTickerHistFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected tickerhist.TickerHist
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{402088407},
			expected: tickerhist.TickerHist{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{
				"tBTCUSD", 54281, nil, 54282, nil, nil, nil, nil, nil, nil, nil, nil, 1619769715000,
			},
			expected: tickerhist.TickerHist{
				Symbol: "tBTCUSD",
				Bid:    54281,
				Ask:    54282,
				MTS:    1619769715000,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := tickerhist.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestTickerHistSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      [][]interface{}
		expected tickerhist.Snapshot
	}{
		"invalid payload": {
			pld:      [][]interface{}{},
			expected: tickerhist.Snapshot{},
		},
		"valid payload": {
			pld: [][]interface{}{
				{"tBTCUSD", 54281, nil, 54282, nil, nil, nil, nil, nil, nil, nil, nil, 1619769715000},
				{"tLTCUSD", 264.66, nil, 264.9, nil, nil, nil, nil, nil, nil, nil, nil, 1619770205000},
			},
			expected: tickerhist.Snapshot{
				Snapshot: []tickerhist.TickerHist{
					{
						Symbol: "tBTCUSD",
						Bid:    54281,
						Ask:    54282,
						MTS:    1619769715000,
					},
					{
						Symbol: "tLTCUSD",
						Bid:    264.66,
						Ask:    264.9,
						MTS:    1619770205000,
					},
				},
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got := tickerhist.SnapshotFromRaw(v.pld)
			assert.Equal(t, v.expected, got)
		})
	}
}
