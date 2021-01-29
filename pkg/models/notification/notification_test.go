package notification_test

import (
	"encoding/json"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/notification"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/stretchr/testify/assert"
)

func TestNotificationMapping(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected interface{}
		err      func(*testing.T, error)
	}{
		"on-req": {
			pld: []byte(`[
				0,
				"n",
				[
					1611922089,"on-req",null,null,
					[
						1201469553,0,788,"tBTCUSD",1611922089073,1611922089073,0.001,0.001,"EXCHANGE LIMIT",
						null,null,null,0,"ACTIVE",null,null,33,0,0,0,null,null,null,0,0,null,null,null,
						"API>BFX",null,null,null
					],
					null,"SUCCESS","Submitting exchange limit buy order for 0.001 BTC."
				]
			]`),
			expected: &notification.Notification{
				MTS:       1611922089,
				Type:      "on-req",
				MessageID: 0,
				NotifyInfo: order.New{
					ID:            1201469553,
					GID:           0,
					CID:           788,
					Symbol:        "tBTCUSD",
					MTSCreated:    1611922089073,
					MTSUpdated:    1611922089073,
					Amount:        0.001,
					AmountOrig:    0.001,
					Type:          "EXCHANGE LIMIT",
					TypePrev:      "",
					MTSTif:        0,
					Flags:         0,
					Status:        "ACTIVE",
					Price:         33,
					PriceAvg:      0,
					PriceTrailing: 0,
					PriceAuxLimit: 0,
					Notify:        false,
					Hidden:        false,
					PlacedID:      0,
					Routing:       "API>BFX",
					Meta:          nil,
				},
				Code:   0,
				Status: "SUCCESS",
				Text:   "Submitting exchange limit buy order for 0.001 BTC.",
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"on-req snapshot": {
			pld: []byte(`[
				0,
				"n",
				[
					1611922089,"on-req",null,null,
					[[
						1201469553,0,788,"tBTCUSD",1611922089073,1611922089073,0.001,0.001,"EXCHANGE LIMIT",
						null,null,null,0,"ACTIVE",null,null,33,0,0,0,null,null,null,0,0,null,null,null,
						"API>BFX",null,null,null
					]],
					null,"SUCCESS","Submitting exchange limit buy order for 0.001 BTC."
				]
			]`),
			expected: &notification.Notification{
				MTS:       1611922089,
				Type:      "on-req",
				MessageID: 0,
				NotifyInfo: &order.Snapshot{
					Snapshot: []*order.Order{
						{
							ID:            1201469553,
							GID:           0,
							CID:           788,
							Symbol:        "tBTCUSD",
							MTSCreated:    1611922089073,
							MTSUpdated:    1611922089073,
							Amount:        0.001,
							AmountOrig:    0.001,
							Type:          "EXCHANGE LIMIT",
							TypePrev:      "",
							MTSTif:        0,
							Flags:         0,
							Status:        "ACTIVE",
							Price:         33,
							PriceAvg:      0,
							PriceTrailing: 0,
							PriceAuxLimit: 0,
							Notify:        false,
							Hidden:        false,
							PlacedID:      0,
							Routing:       "API>BFX",
							Meta:          nil,
						},
					},
				},
				Code:   0,
				Status: "SUCCESS",
				Text:   "Submitting exchange limit buy order for 0.001 BTC.",
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"ou-req": {
			pld: []byte(`[
				0,
				"n",
				[
					1575289447641,"ou-req",null,null,
					[
						1185815100,null,1575289350475,"tETHUSD",1575289351944,1575289351948,-3,
						-3,"LIMIT",null,null,null,0,"ACTIVE",null,null,240,0,0,0,null,null,null,
						0,0,null,null,null,"API>BFX",null,null,null
					],
					null,"SUCCESS","Submitting update to limit sell order for 3 ETH."
				]
			]`),
			expected: &notification.Notification{
				MTS:       1575289447641,
				Type:      "ou-req",
				MessageID: 0,
				NotifyInfo: order.Update{
					ID:            1185815100,
					GID:           0,
					CID:           1575289350475,
					Symbol:        "tETHUSD",
					MTSCreated:    1575289351944,
					MTSUpdated:    1575289351948,
					Amount:        -3,
					AmountOrig:    -3,
					Type:          "LIMIT",
					TypePrev:      "",
					MTSTif:        0,
					Flags:         0,
					Status:        "ACTIVE",
					Price:         240,
					PriceAvg:      0,
					PriceTrailing: 0,
					PriceAuxLimit: 0,
					Notify:        false,
					Hidden:        false,
					PlacedID:      0,
					Routing:       "API>BFX",
					Meta:          nil,
				},
				Code:   0,
				Status: "SUCCESS",
				Text:   "Submitting update to limit sell order for 3 ETH.",
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"fon-req": {
			pld: []byte(`[
				0,
				"n",
				[
					1575282446099,"fon-req",null,null,
					[
						41238905,null,null,null,-1000,null,null,null,null,null,
						null,null,null,null,0.002,2,null,null,null,null,null
					],
					null,"SUCCESS","Submitting funding bid of 1000.0 USD at 0.2000 for 2 days."
				]
			]`),
			expected: &notification.Notification{
				MTS:       1575282446099,
				Type:      "fon-req",
				MessageID: 0,
				NotifyInfo: fundingoffer.New{
					ID:         41238905,
					Symbol:     "",
					MTSCreated: 0,
					MTSUpdated: 0,
					Amount:     -1000,
					AmountOrig: 0,
					Type:       "",
					Flags:      map[string]interface{}(nil),
					Status:     "",
					Rate:       0.002,
					Period:     2,
					Notify:     false,
					Hidden:     false,
					Insure:     false,
					Renew:      false,
					RateReal:   0,
				},
				Code:   0,
				Status: "SUCCESS",
				Text:   "Submitting funding bid of 1000.0 USD at 0.2000 for 2 days.",
			},
			err: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			var raw []interface{}
			json.Unmarshal(v.pld, &raw)
			pldRaw := raw[len(raw)-1]
			pld := pldRaw.([]interface{})

			got, err := notification.FromRaw(pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
