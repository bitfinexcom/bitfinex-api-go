package status_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/stretchr/testify/assert"
)

func TestFromWSRaw(t *testing.T) {
	testCases := map[string]struct {
		key     string
		data    []interface{}
		err     func(*testing.T, error)
		success func(*testing.T, interface{})
	}{
		"empty slice": {
			key:  "",
			data: []interface{}{},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"invalid key": {
			key:  "foo",
			data: []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"deriv snapshot invalid pld": {
			key:  "deriv:tBTCF0:USTF0",
			data: []interface{}{[]interface{}{1}},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"deriv raw invalid pld": {
			key:  "deriv:tBTCF0:USTF0",
			data: []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"liq snapshot invalid pld": {
			key:  "liq:global",
			data: []interface{}{[]interface{}{1}},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"liq raw invalid pld": {
			key:  "liq:global",
			data: []interface{}{1},
			err: func(t *testing.T, got error) {
				assert.Error(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Nil(t, got)
			},
		},
		"deriv snapshot valid pld": {
			key: "deriv:tBTCF0:USTF0",
			data: []interface{}{
				[]interface{}{
					1596124822000, nil, 0.896, 0.771995, nil, 1396531.67460709,
					nil, 1596153600000, 0.0001056, 6, nil, -0.01381344, nil, nil,
					0.7664, nil, nil, 246502.09001551, nil, nil, nil, nil, 0.3,
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
					},
				})
			},
		},
		"deriv raw valid pld": {
			key: "deriv:tBTCF0:USTF0",
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
		"liquidation snapshot valid pld": {
			key: "liq:global",
			data: []interface{}{
				[]interface{}{
					"pos", 145400868, 1609144352338, nil, "tETHF0:USTF0",
					-1.67288094, 730.96, nil, 1, 1, nil, 736.13,
				},
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.LiquidationsSnapshot{
					Snapshot: []*status.Liquidation{
						{
							Symbol:        "tETHF0:USTF0",
							PositionID:    145400868,
							MTS:           1609144352338,
							Amount:        -1.67288094,
							BasePrice:     730.96,
							IsMatch:       1,
							IsMarketSold:  1,
							PriceAcquired: 736.13,
						},
					},
				})
			},
		},
		"liquidation raw valid pld": {
			key: "liq:global",
			data: []interface{}{
				"pos", 145400868, 1609144352338, nil, "tETHF0:USTF0",
				-1.67288094, 730.96, nil, 1, 1, nil, 736.13,
			},
			err: func(t *testing.T, got error) {
				assert.Nil(t, got)
			},
			success: func(t *testing.T, got interface{}) {
				assert.Equal(t, got, &status.Liquidation{
					Symbol:        "tETHF0:USTF0",
					PositionID:    145400868,
					MTS:           1609144352338,
					Amount:        -1.67288094,
					BasePrice:     730.96,
					IsMatch:       1,
					IsMarketSold:  1,
					PriceAcquired: 736.13,
				})
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			got, err := status.FromWSRaw(v.key, v.data)
			v.err(t, err)
			v.success(t, got)
		})
	}
}
