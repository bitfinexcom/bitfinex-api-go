package status_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/stretchr/testify/assert"
)

func TestDerivFromRaw(t *testing.T) {
	testCases := map[string]struct {
		symbol  string
		data    []interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty slice": {
			symbol: "tBTCF0:USTF0",
			data:   []interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid pld": {
			symbol: "tBTCF0:USTF0",
			data:   []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"valid pld": {
			symbol: "tBTCF0:USTF0",
			data: []interface{}{
				1596124822000, nil, 0.896, 0.771995, nil, 1396531.67460709,
				nil, 1596153600000, 0.0001056, 6, nil, -0.01381344, nil, nil,
				0.7664, nil, nil, 246502.09001551, nil, nil, nil, nil, 0.3,
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.Derivative{
					Symbol:               "tBTCF0:USTF0",
					MTS:                  1596124822000,
					Price:                0.896,
					SpotPrice:            0.771995,
					InsuranceFundBalance: 1396531.67460709,
					FundingEventMTS:      1596153600000,
					FundingAccrued:       0.0001056,
					FundingStep:          6,
					CurrentFunding:       -0.01381344,
					MarkPrice:            0.7664,
					OpenInterest:         246502.09001551,
					ClampMIN:             0,
					ClampMAX:             0.3,
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.DerivFromRaw(v.symbol, v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}

func TestDerivFromRestRaw(t *testing.T) {
	testCases := map[string]struct {
		data    []interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty slice": {
			data: []interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid pld": {
			data: []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"valid pld": {
			data: []interface{}{
				"tBTCF0:USTF0", 1596124822000, nil, 0.896, 0.771995, nil,
				1396531.67460709, nil, 1596153600000, 0.0001056, 6, nil,
				-0.01381344, nil, nil, 0.7664, nil, nil, 246502.09001551,
				nil, nil, nil, nil, 0.3,
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.Derivative{
					Symbol:               "tBTCF0:USTF0",
					MTS:                  1596124822000,
					Price:                0.896,
					SpotPrice:            0.771995,
					InsuranceFundBalance: 1396531.67460709,
					FundingEventMTS:      1596153600000,
					FundingAccrued:       0.0001056,
					FundingStep:          6,
					CurrentFunding:       -0.01381344,
					MarkPrice:            0.7664,
					OpenInterest:         246502.09001551,
					ClampMIN:             0,
					ClampMAX:             0.3,
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.DerivFromRestRaw(v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}

func TestDerivSnapshotFromRaw(t *testing.T) {
	testCases := map[string]struct {
		symbol  string
		data    [][]interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty slice": {
			symbol: "tBTCF0:USTF0",
			data:   [][]interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid pld": {
			symbol: "tBTCF0:USTF0",
			data:   [][]interface{}{{1}},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"valid pld": {
			symbol: "tBTCF0:USTF0",
			data: [][]interface{}{
				{
					1596124822000, nil, 0.896, 0.771995, nil, 1396531.67460709,
					nil, 1596153600000, 0.0001056, 6, nil, -0.01381344, nil, nil,
					0.7664, nil, nil, 246502.09001551, nil, nil, nil, nil, 0.3,
				},
				{
					1596124822001, nil, 0.896, 0.771995, nil, 1396531.67460709,
					nil, 1596153600000, 0.0001056, 6, nil, -0.01381344, nil, nil,
					0.7664, nil, nil, 246502.09001551, nil, nil, nil, nil, 0.3, nil, 123,
				},
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.DerivativesSnapshot{
					Snapshot: []*status.Derivative{
						{
							Symbol:               "tBTCF0:USTF0",
							MTS:                  1596124822000,
							Price:                0.896,
							SpotPrice:            0.771995,
							InsuranceFundBalance: 1396531.67460709,
							FundingEventMTS:      1596153600000,
							FundingAccrued:       0.0001056,
							FundingStep:          6,
							CurrentFunding:       -0.01381344,
							MarkPrice:            0.7664,
							OpenInterest:         246502.09001551,
							ClampMIN:             0,
							ClampMAX:             0.3,
						},
						{
							Symbol:               "tBTCF0:USTF0",
							MTS:                  1596124822001,
							Price:                0.896,
							SpotPrice:            0.771995,
							InsuranceFundBalance: 1396531.67460709,
							FundingEventMTS:      1596153600000,
							FundingAccrued:       0.0001056,
							FundingStep:          6,
							CurrentFunding:       -0.01381344,
							MarkPrice:            0.7664,
							OpenInterest:         246502.09001551,
							ClampMIN:             0,
							ClampMAX:             0.3,
						},
					},
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.DerivSnapshotFromRaw(v.symbol, v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}
