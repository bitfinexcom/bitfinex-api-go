package trades_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/stretchr/testify/assert"
)

func TestAuthFundingTradeFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.AuthFundingTrade
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{636854},
			expected: trades.AuthFundingTrade{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{636854, "fUSD", 1575282446000, 41238905, -1000, 0.002, 7, nil},
			expected: trades.AuthFundingTrade{
				ID:         636854,
				Symbol:     "fUSD",
				MTSCreated: 1575282446000,
				OfferID:    41238905,
				Amount:     -1000,
				Rate:       0.002,
				Period:     7,
				Maker:      0,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.AFTFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestAuthFundingTradeUpdateFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.AuthFundingTradeUpdate
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{636854},
			expected: trades.AuthFundingTradeUpdate{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{636854, "fUSD", 1575282446000, 41238905, -1000, 0.002, 7, nil},
			expected: trades.AuthFundingTradeUpdate{
				ID:         636854,
				Symbol:     "fUSD",
				MTSCreated: 1575282446000,
				OfferID:    41238905,
				Amount:     -1000,
				Rate:       0.002,
				Period:     7,
				Maker:      0,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.AFTUFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestAuthFundingTradeExecutionFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected trades.AuthFundingTradeExecuted
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      []interface{}{636854},
			expected: trades.AuthFundingTradeExecuted{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: []interface{}{636854, "fUSD", 1575282446000, 41238905, -1000, 0.002, 7, nil},
			expected: trades.AuthFundingTradeExecuted{
				ID:         636854,
				Symbol:     "fUSD",
				MTSCreated: 1575282446000,
				OfferID:    41238905,
				Amount:     -1000,
				Rate:       0.002,
				Period:     7,
				Maker:      0,
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.AFTEFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestAuthFundingTradeSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      [][]interface{}
		expected trades.AuthFundingTradeSnapshot
		err      func(*testing.T, error)
	}{
		"invalid payload": {
			pld:      [][]interface{}{},
			expected: trades.AuthFundingTradeSnapshot{},
			err: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		"valid payload": {
			pld: [][]interface{}{{
				636854, "fUSD", 1575282446000, 41238905, -1000, 0.002, 7, nil,
			}},
			expected: trades.AuthFundingTradeSnapshot{
				Snapshot: []trades.AuthFundingTrade{
					{
						ID:         636854,
						Symbol:     "fUSD",
						MTSCreated: 1575282446000,
						OfferID:    41238905,
						Amount:     -1000,
						Rate:       0.002,
						Period:     7,
						Maker:      0,
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := trades.AFTSnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
