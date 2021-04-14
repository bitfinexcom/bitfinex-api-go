package wallet_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *wallet.Wallet
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{"exchange"},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"valid rest pld with meta": {
			pld: []interface{}{
				"exchange", "UST", 19788.6529257, 0, 19788.6529257,
				"Exchange 2.0 UST for USD @ 11.696",
				map[string]interface{}{
					"reason":        "TRADE",
					"order_id":      1189740779,
					"order_id_oppo": 1189785673,
					"trade_price":   "11.696",
					"trade_amount":  "-2.0",
					"order_cid":     1598516362757,
					"order_gid":     1598516362629,
				},
			},
			expected: &wallet.Wallet{
				Type:              "exchange",
				Currency:          "UST",
				Balance:           19788.6529257,
				UnsettledInterest: 0,
				BalanceAvailable:  19788.6529257,
				LastChange:        "Exchange 2.0 UST for USD @ 11.696",
				TradeDetails: map[string]interface{}{
					"order_cid":     1598516362757,
					"order_gid":     1598516362629,
					"order_id":      1189740779,
					"order_id_oppo": 1189785673,
					"reason":        "TRADE",
					"trade_amount":  "-2.0",
					"trade_price":   "11.696",
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"valid rest pld no meta": {
			pld: []interface{}{
				"exchange", "UST", 19788.6529257, 0, 19788.6529257,
				"Exchange 2.0 UST for USD @ 11.696", nil,
			},
			expected: &wallet.Wallet{
				Type:              "exchange",
				Currency:          "UST",
				Balance:           19788.6529257,
				UnsettledInterest: 0,
				BalanceAvailable:  19788.6529257,
				LastChange:        "Exchange 2.0 UST for USD @ 11.696",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"valid ws wallet with meta": {
			pld: []interface{}{
				"exchange", "BTC", 1.61169184, 0, nil,
				"Exchange 0.01 BTC for USD @ 7804.6",
				map[string]interface{}{
					"reason":        "TRADE",
					"order_id":      34988418651,
					"order_id_oppo": 34990541044,
					"trade_price":   "7804.6",
					"trade_amount":  "0.01",
				},
			},
			expected: &wallet.Wallet{
				Type:              "exchange",
				Currency:          "BTC",
				Balance:           1.61169184,
				UnsettledInterest: 0,
				BalanceAvailable:  0,
				LastChange:        "Exchange 0.01 BTC for USD @ 7804.6",
				TradeDetails: map[string]interface{}{
					"order_id":      34988418651,
					"order_id_oppo": 34990541044,
					"reason":        "TRADE",
					"trade_amount":  "0.01",
					"trade_price":   "7804.6",
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"valid ws wallet no meta": {
			pld: []interface{}{
				"exchange", "BTC", 1.61169184, 0, nil,
				"Exchange 0.01 BTC for USD @ 7804.6", nil,
			},
			expected: &wallet.Wallet{
				Type:              "exchange",
				Currency:          "BTC",
				Balance:           1.61169184,
				UnsettledInterest: 0,
				BalanceAvailable:  0,
				LastChange:        "Exchange 0.01 BTC for USD @ 7804.6",
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := wallet.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *wallet.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"valid rest pld": {
			pld: []interface{}{
				[]interface{}{
					"exchange", "UST", 19788.6529257, 0, 19788.6529257,
					"Exchange 2.0 UST for USD @ 11.696",
					map[string]interface{}{
						"reason":        "TRADE",
						"order_id":      1189740779,
						"order_id_oppo": 1189785673,
						"trade_price":   "11.696",
						"trade_amount":  "-2.0",
						"order_cid":     1598516362757,
						"order_gid":     1598516362629,
					},
				},
				[]interface{}{
					"exchange", "UST", 19788.6529257, 0, 19788.6529257,
					"Exchange 2.0 UST for USD @ 11.696", nil,
				},
			},
			expected: &wallet.Snapshot{
				Snapshot: []*wallet.Wallet{
					{
						Type:              "exchange",
						Currency:          "UST",
						Balance:           19788.6529257,
						UnsettledInterest: 0,
						BalanceAvailable:  19788.6529257,
						LastChange:        "Exchange 2.0 UST for USD @ 11.696",
						TradeDetails: map[string]interface{}{
							"order_cid":     1598516362757,
							"order_gid":     1598516362629,
							"order_id":      1189740779,
							"order_id_oppo": 1189785673,
							"reason":        "TRADE",
							"trade_amount":  "-2.0",
							"trade_price":   "11.696",
						},
					},
					{
						Type:              "exchange",
						Currency:          "UST",
						Balance:           19788.6529257,
						UnsettledInterest: 0,
						BalanceAvailable:  19788.6529257,
						LastChange:        "Exchange 2.0 UST for USD @ 11.696",
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"valid ws pld": {
			pld: []interface{}{
				[]interface{}{"exchange", "SAN", 19.76, 0, nil, nil, nil},
			},
			expected: &wallet.Snapshot{
				Snapshot: []*wallet.Wallet{
					{
						Type:              "exchange",
						Currency:          "SAN",
						Balance:           19.76,
						UnsettledInterest: 0,
						BalanceAvailable:  0,
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
			got, err := wallet.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
