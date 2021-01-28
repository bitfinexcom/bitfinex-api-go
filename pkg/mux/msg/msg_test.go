package msg_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/msg"
	"github.com/stretchr/testify/assert"
)

func TestIsEvent(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected bool
	}{
		"event type": {
			pld:      []byte(`{}`),
			expected: true,
		},
		"not event type": {
			pld:      []byte(`[]`),
			expected: false,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{
				Data: v.pld,
			}

			got := m.IsEvent()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestIsRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected bool
	}{
		"raw msg type": {
			pld:      []byte(`[]`),
			expected: true,
		},
		"raw info type": {
			pld:      []byte(`{}`),
			expected: false,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{
				Data: v.pld,
			}

			got := m.IsRaw()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	m := msg.Msg{
		Data: []byte(`{
			"event": "info",
			"version": 2,
			"serverId": "dbea77ee-4740-4a82-84f3-c6bc1b5abb9a",
			"platform": {
				"status":1
			}
		}`),
	}

	expected := event.Info{
		Subscribe: event.Subscribe{
			Event: "info",
		},
		Version:  2,
		ServerID: "dbea77ee-4740-4a82-84f3-c6bc1b5abb9a",
		Platform: event.Platform{Status: 1},
	}

	got, err := m.ProcessEvent()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestProcessRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected interface{}
		inf      map[int64]event.Info
	}{
		"info event": {
			pld: []byte(`[123, "hb"]`),
			inf: map[int64]event.Info{123: {}},
			expected: event.Info{
				ChanID:    123,
				Subscribe: event.Subscribe{Event: "hb"},
			},
		},
		"trades snapshot": {
			pld: []byte(`[111,[[559273857,1609665708633,-0.0048,34113]]]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: trades.TradeSnapshot{
				Snapshot: []trades.Trade{
					{
						Pair:   "tBTCUST",
						ID:     559273857,
						MTS:    1609665708633,
						Amount: -0.0048,
						Price:  34113,
					},
				},
			},
		},
		"trade": {
			pld: []byte(`[111,[559273857,1609665708633,-0.0048,34113]]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: trades.Trade{
				Pair:   "tBTCUST",
				ID:     559273857,
				MTS:    1609665708633,
				Amount: -0.0048,
				Price:  34113,
			},
		},
		"ticker": {
			pld: []byte(`[
				111,
				[
					34072,0.019999999999999997,34080,6.69793272,4350,
					0.1464,34062,4047.85335915,34758,29490
				]
			]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "ticker",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: &ticker.Ticker{
				Symbol:          "tBTCUST",
				Bid:             34072,
				BidSize:         0.019999999999999997,
				Ask:             34080,
				AskSize:         6.69793272,
				DailyChange:     4350,
				DailyChangePerc: 0.1464,
				LastPrice:       34062,
				Volume:          4047.85335915,
				High:            34758,
				Low:             29490,
			},
		},
		"candles snapshot": {
			pld: []byte(`[111,[[1609668540000,828.01,827.67,828.42,827.67,2.32080241]]]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "candles",
						Key:     "trade:1m:tETHUST",
					},
				},
			},
			expected: &candle.Snapshot{
				Snapshot: []*candle.Candle{
					{
						Symbol:     "tETHUST",
						Resolution: "1m",
						MTS:        1609668540000,
						Open:       828.01,
						Close:      827.67,
						High:       828.42,
						Low:        827.67,
						Volume:     2.32080241,
					},
				},
			},
		},
		"candle": {
			pld: []byte(`[111,[1609668540000,828.01,827.67,828.42,827.67,2.32080241]]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "candles",
						Key:     "trade:1m:tETHUST",
					},
				},
			},
			expected: &candle.Candle{
				Symbol:     "tETHUST",
				Resolution: "1m",
				MTS:        1609668540000,
				Open:       828.01,
				Close:      827.67,
				High:       828.42,
				Low:        827.67,
				Volume:     2.32080241,
			},
		},
		"raw book snapshot": {
			pld: []byte(`[869944,[[55804480297,33766,2]]]`),
			inf: map[int64]event.Info{
				869944: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "tBTCUSD",
						Precision: "R0",
					},
				},
			},
			expected: &book.Snapshot{
				Snapshot: []*book.Book{
					{
						Symbol:      "tBTCUSD",
						ID:          55804480297,
						Price:       33766,
						Amount:      2,
						PriceJsNum:  "33766",
						AmountJsNum: "2",
						Side:        1,
						Action:      0,
					},
				},
			},
		},
		"raw book": {
			pld: []byte(`[869944,[55804480297,33766,2]]`),
			inf: map[int64]event.Info{
				869944: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "tBTCUSD",
						Precision: "R0",
					},
				},
			},
			expected: &book.Book{
				Symbol:      "tBTCUSD",
				ID:          55804480297,
				Price:       33766,
				Amount:      2,
				PriceJsNum:  "33766",
				AmountJsNum: "2",
				Side:        1,
				Action:      0,
			},
		},
		"book snapshot": {
			pld: []byte(`[793767,[[676.3,1,5]]]`),
			inf: map[int64]event.Info{
				793767: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "tETHEUR",
						Precision: "P0",
						Frequency: "F0",
					},
				},
			},
			expected: &book.Snapshot{
				Snapshot: []*book.Book{
					{
						Symbol:      "tETHEUR",
						Count:       1,
						Period:      0,
						Price:       676.3,
						Amount:      5,
						Rate:        0,
						PriceJsNum:  "676.3",
						AmountJsNum: "5",
						Side:        1,
						Action:      0,
					},
				},
			},
		},
		"book": {
			pld: []byte(`[793767,[676.3,1,5]]`),
			inf: map[int64]event.Info{
				793767: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "tETHEUR",
						Precision: "P0",
						Frequency: "F0",
					},
				},
			},
			expected: &book.Book{
				Symbol:      "tETHEUR",
				Count:       1,
				Period:      0,
				Price:       676.3,
				Amount:      5,
				Rate:        0,
				PriceJsNum:  "676.3",
				AmountJsNum: "5",
				Side:        1,
				Action:      0,
			},
		},
		"derivatives status snapshot": {
			pld: []byte(`[
				799830,
				[[
					1609921474000,null,34568.786626655,34575.5,null,1856521.42387705,
					null,1609948800000,-0.00004348,481,null,0,null,null,34593.64333333333,
					null,null,11153.74635347,null,null,null,null,null
				]]
			]`),
			inf: map[int64]event.Info{
				799830: {
					Subscribe: event.Subscribe{
						Channel: "status",
						Key:     "deriv:tBTCF0:USTF0",
					},
				},
			},
			expected: &status.DerivativesSnapshot{
				Snapshot: []*status.Derivative{
					{
						Symbol:               "tBTCF0:USTF0",
						MTS:                  1609921474000,
						Price:                34568.786626655,
						SpotPrice:            34575.5,
						InsuranceFundBalance: 1.85652142387705e+06,
						FundingEventMTS:      1609948800000,
						FundingAccrued:       -4.348e-05,
						FundingStep:          481,
						MarkPrice:            34593.64333333333,
						OpenInterest:         11153.74635347,
					},
				},
			},
		},
		"derivatives status": {
			pld: []byte(`[
				799830,
				[
					1609921474000,null,34568.786626655,34575.5,null,1856521.42387705,
					null,1609948800000,-0.00004348,481,null,0,null,null,34593.64333333333,
					null,null,11153.74635347,null,null,null,null,null
				]
			]`),
			inf: map[int64]event.Info{
				799830: {
					Subscribe: event.Subscribe{
						Channel: "status",
						Key:     "deriv:tBTCF0:USTF0",
					},
				},
			},
			expected: &status.Derivative{
				Symbol:               "tBTCF0:USTF0",
				MTS:                  1609921474000,
				Price:                34568.786626655,
				SpotPrice:            34575.5,
				InsuranceFundBalance: 1.85652142387705e+06,
				FundingEventMTS:      1609948800000,
				FundingAccrued:       -4.348e-05,
				FundingStep:          481,
				MarkPrice:            34593.64333333333,
				OpenInterest:         11153.74635347,
			},
		},
		"liquidation status snapshot": {
			pld: []byte(`[
				521209,
				[[
					"pos",145511476,1609921778489,null,"tBTCF0:USTF0",
					0.12173,34618.82986269,null,1,1,null,34281
				]]
			]`),
			inf: map[int64]event.Info{
				521209: {
					Subscribe: event.Subscribe{
						Channel: "status",
						Key:     "liq:global",
					},
				},
			},
			expected: &status.LiquidationsSnapshot{
				Snapshot: []*status.Liquidation{
					{
						Symbol:        "tBTCF0:USTF0",
						PositionID:    145511476,
						MTS:           1609921778489,
						Amount:        0.12173,
						BasePrice:     34618.82986269,
						IsMatch:       1,
						IsMarketSold:  1,
						PriceAcquired: 34281,
					},
				},
			},
		},
		"liquidation status": {
			pld: []byte(`[
				521209,
				[
					"pos",145511476,1609921778489,null,"tBTCF0:USTF0",
					0.12173,34618.82986269,null,1,1,null,34281
				]
			]`),
			inf: map[int64]event.Info{
				521209: {
					Subscribe: event.Subscribe{
						Channel: "status",
						Key:     "liq:global",
					},
				},
			},
			expected: &status.Liquidation{
				Symbol:        "tBTCF0:USTF0",
				PositionID:    145511476,
				MTS:           1609921778489,
				Amount:        0.12173,
				BasePrice:     34618.82986269,
				IsMatch:       1,
				IsMarketSold:  1,
				PriceAcquired: 34281,
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{Data: v.pld}
			got, err := m.ProcessRaw(v.inf)
			assert.NoError(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestProcessPrivateRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected interface{}
		inf      map[int64]event.Info
	}{
		"info event": {
			pld: []byte(`[0, "hb"]`),
			expected: event.Info{
				ChanID:    0,
				Subscribe: event.Subscribe{Event: "hb"},
			},
		},
		"balance info update": {
			pld: []byte(`[0,"bu",[4131,4131.85]]`),
			expected: balanceinfo.Update{
				TotalAUM: 4131,
				NetAUM:   4131.85,
			},
		},
		"position snapshot": {
			pld: []byte(`[
				0,
				"ps",
				[[
					"tETHUST","ACTIVE",0.2,153.71,0,0,null,null,null,
					null,null,142420429,null,null,null,0,null,0,null,
					{
						"reason":"TRADE",
						"order_id":34934099168,
						"order_id_oppo":34934090814,
						"liq_stage":null,
						"trade_price":"153.71",
						"trade_amount":"0.2"
					}
				]]
			]`),
			expected: &position.Snapshot{
				Snapshot: []*position.Position{
					{
						Id:        142420429,
						Symbol:    "tETHUST",
						Status:    "ACTIVE",
						Amount:    0.2,
						BasePrice: 153.71,
						Type:      "ps",
						Meta: map[string]interface{}{
							"liq_stage":     nil,
							"order_id":      3.4934099168e+10,
							"order_id_oppo": 3.4934090814e+10,
							"reason":        "TRADE",
							"trade_amount":  "0.2",
							"trade_price":   "153.71",
						},
					},
				},
			},
		},
		"position new": {
			pld: []byte(`[
				0,
				"pn",
				[
					"tETHUST","ACTIVE",0.2,153.71,0,0,-0.07944800000000068,-0.05855181835925015,
					67.52755254906451, 1.409288545397275,null,142420429,null,null,null,0,null,0,0,
					{
						"reason":"TRADE",
						"order_id":34934099168,
						"order_id_oppo":34934090814,
						"liq_stage":null,
						"trade_price":"153.71",
						"trade_amount":"0.2"
					}
				]
			]`),
			expected: position.New{
				Id:                   142420429,
				Symbol:               "tETHUST",
				Status:               "ACTIVE",
				Amount:               0.2,
				BasePrice:            153.71,
				Type:                 "pn",
				ProfitLoss:           -0.07944800000000068,
				ProfitLossPercentage: -0.05855181835925015,
				LiquidationPrice:     67.52755254906451,
				Leverage:             1.409288545397275,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      3.4934099168e+10,
					"order_id_oppo": 3.4934090814e+10,
					"reason":        "TRADE",
					"trade_amount":  "0.2",
					"trade_price":   "153.71",
				},
			},
		},
		"position update": {
			pld: []byte(`[
				0,
				"pu",
				[
					"tETHUST","ACTIVE",0.2,153.71,0,0,-0.07944800000000068,-0.05855181835925015,
					67.52755254906451, 1.409288545397275,null,142420429,null,null,null,0,null,0,0,
					{
						"reason":"TRADE",
						"order_id":34934099168,
						"order_id_oppo":34934090814,
						"liq_stage":null,
						"trade_price":"153.71",
						"trade_amount":"0.2"
					}
				]
			]`),
			expected: position.Update{
				Id:                   142420429,
				Symbol:               "tETHUST",
				Status:               "ACTIVE",
				Amount:               0.2,
				BasePrice:            153.71,
				Type:                 "pu",
				ProfitLoss:           -0.07944800000000068,
				ProfitLossPercentage: -0.05855181835925015,
				LiquidationPrice:     67.52755254906451,
				Leverage:             1.409288545397275,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      3.4934099168e+10,
					"order_id_oppo": 3.4934090814e+10,
					"reason":        "TRADE",
					"trade_amount":  "0.2",
					"trade_price":   "153.71",
				},
			},
		},
		"position close": {
			pld: []byte(`[
				0,
				"pc",
				[
					"tETHUST","ACTIVE",0.2,153.71,0,0,-0.07944800000000068,-0.05855181835925015,
					67.52755254906451, 1.409288545397275,null,142420429,null,null,null,0,null,0,0,
					{
						"reason":"TRADE",
						"order_id":34934099168,
						"order_id_oppo":34934090814,
						"liq_stage":null,
						"trade_price":"153.71",
						"trade_amount":"0.2"
					}
				]
			]`),
			expected: position.Cancel{
				Id:                   142420429,
				Symbol:               "tETHUST",
				Status:               "ACTIVE",
				Amount:               0.2,
				BasePrice:            153.71,
				Type:                 "pc",
				ProfitLoss:           -0.07944800000000068,
				ProfitLossPercentage: -0.05855181835925015,
				LiquidationPrice:     67.52755254906451,
				Leverage:             1.409288545397275,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      3.4934099168e+10,
					"order_id_oppo": 3.4934090814e+10,
					"reason":        "TRADE",
					"trade_amount":  "0.2",
					"trade_price":   "153.71",
				},
			},
		},
		"wallet snapshot": {
			pld: []byte(`[0,"ws",[["exchange","SAN",19.76,0,null,null,null]]]`),
			expected: &wallet.Snapshot{
				Snapshot: []*wallet.Wallet{
					{
						Type:         "exchange",
						Currency:     "SAN",
						Balance:      19.76,
						TradeDetails: nil,
					},
				},
			},
		},
		"wallet update": {
			pld: []byte(`[
				0,
				"wu",
				[
					"exchange","BTC",1.61169184,0,null,"Exchange 0.01 BTC for USD @ 7804.6",
					{
						"reason":"TRADE",
						"order_id":34988418651,
						"order_id_oppo":34990541044,
						"trade_price":"7804.6",
						"trade_amount":"0.01"
					}
				]
			]`),
			expected: wallet.Update{
				Type:              "exchange",
				Currency:          "BTC",
				Balance:           1.61169184,
				UnsettledInterest: 0,
				BalanceAvailable:  0,
				LastChange:        "Exchange 0.01 BTC for USD @ 7804.6",
				TradeDetails: map[string]interface{}{
					"order_id":      3.4988418651e+10,
					"order_id_oppo": 3.4990541044e+10,
					"reason":        "TRADE",
					"trade_amount":  "0.01",
					"trade_price":   "7804.6",
				},
			},
		},
		"order snapshot": {
			pld: []byte(`[
				0,
				"os",
				[[
					34930659963,null,1574955083558,"tETHUSD",1574955083558,1574955083573,
					0.201104,0.201104,"EXCHANGE LIMIT",null,null,null,0,"ACTIVE",null,null,
					120,0,0,0,null,null,null,0,0,null,null,null,"BFX",null,null,null
				]]
			]`),
			expected: &order.Snapshot{
				Snapshot: []*order.Order{
					{
						ID:            34930659963,
						GID:           0,
						CID:           1574955083558,
						Symbol:        "tETHUSD",
						MTSCreated:    1574955083558,
						MTSUpdated:    1574955083573,
						Amount:        0.201104,
						AmountOrig:    0.201104,
						Type:          "EXCHANGE LIMIT",
						TypePrev:      "",
						MTSTif:        0,
						Flags:         0,
						Status:        "ACTIVE",
						Price:         120,
						PriceAvg:      0,
						PriceTrailing: 0,
						PriceAuxLimit: 0,
						Notify:        false,
						Hidden:        false,
						PlacedID:      0,
						Routing:       "BFX",
						Meta:          nil,
					},
				},
			},
		},
		"order new": {
			pld: []byte(`[
				0,
				"on",
				[
					34930659963,null,1574955083558,"tETHUSD",1574955083558,1574955354487,
					0.201104,0.201104,"EXCHANGE LIMIT",null,null,null,0,"CANCELED",null,
					null,120,0,0,0,null,null,null,0,0,null,null,null,"BFX",null,null,null
				]
			]`),
			expected: order.New{
				ID:            34930659963,
				GID:           0,
				CID:           1574955083558,
				Symbol:        "tETHUSD",
				MTSCreated:    1574955083558,
				MTSUpdated:    1574955354487,
				Amount:        0.201104,
				AmountOrig:    0.201104,
				Type:          "EXCHANGE LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "CANCELED",
				Price:         120,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        false,
				PlacedID:      0,
				Routing:       "BFX",
				Meta:          nil,
			},
		},
		"order update": {
			pld: []byte(`[
				0,
				"ou",
				[
					34930659963,null,1574955083558,"tETHUSD",1574955083558,1574955354487,
					0.201104,0.201104,"EXCHANGE LIMIT",null,null,null,0,"CANCELED",null,
					null,120,0,0,0,null,null,null,0,0,null,null,null,"BFX",null,null,null
				]
			]`),
			expected: order.Update{
				ID:            34930659963,
				GID:           0,
				CID:           1574955083558,
				Symbol:        "tETHUSD",
				MTSCreated:    1574955083558,
				MTSUpdated:    1574955354487,
				Amount:        0.201104,
				AmountOrig:    0.201104,
				Type:          "EXCHANGE LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "CANCELED",
				Price:         120,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        false,
				PlacedID:      0,
				Routing:       "BFX",
				Meta:          nil,
			},
		},
		"order cancel": {
			pld: []byte(`[
				0,
				"oc",
				[
					34930659963,null,1574955083558,"tETHUSD",1574955083558,1574955354487,
					0.201104,0.201104,"EXCHANGE LIMIT",null,null,null,0,"CANCELED",null,
					null,120,0,0,0,null,null,null,0,0,null,null,null,"BFX",null,null,null
				]
			]`),
			expected: order.Cancel{
				ID:            34930659963,
				GID:           0,
				CID:           1574955083558,
				Symbol:        "tETHUSD",
				MTSCreated:    1574955083558,
				MTSUpdated:    1574955354487,
				Amount:        0.201104,
				AmountOrig:    0.201104,
				Type:          "EXCHANGE LIMIT",
				TypePrev:      "",
				MTSTif:        0,
				Flags:         0,
				Status:        "CANCELED",
				Price:         120,
				PriceAvg:      0,
				PriceTrailing: 0,
				PriceAuxLimit: 0,
				Notify:        false,
				Hidden:        false,
				PlacedID:      0,
				Routing:       "BFX",
				Meta:          nil,
			},
		},
		"trade execution": {
			pld: []byte(`[
				0,
				"te",
				[402088407,"tETHUST",1574963975602,34938060782,-0.2,153.57,"MARKET",0,-1,null,null,0]
			]`),
			expected: trades.AuthTradeExecution{
				ID:            402088407,
				Pair:          "tETHUST",
				MTS:           1574963975602,
				OrderID:       34938060782,
				ExecAmount:    -0.2,
				ExecPrice:     153.57,
				OrderType:     "MARKET",
				OrderPrice:    0,
				Maker:         -1,
				ClientOrderID: 0,
			},
		},
		"trade execution update": {
			pld: []byte(`[
				0,
				"tu",
				[402088407,"tETHUST",1574963975602,34938060782,-0.2,153.57,"MARKET",0,-1,-0.061668,"USD"]
			]`),
			expected: trades.AuthTradeExecutionUpdate{
				ID:          402088407,
				Pair:        "tETHUST",
				MTS:         1574963975602,
				OrderID:     34938060782,
				ExecAmount:  -0.2,
				ExecPrice:   153.57,
				OrderType:   "MARKET",
				OrderPrice:  0,
				Maker:       -1,
				Fee:         -0.061668,
				FeeCurrency: "USD",
			},
		},
		"funding trade executed": {
			pld: []byte(`[0,"fte",[636854,"fUSD",1575282446000,41238905,-1000,0.002,7,null]]`),
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
		},
		"funding trade update": {
			pld: []byte(`[0,"ftu",[636854,"fUSD",1575282446000,41238905,-1000,0.002,7,null]]`),
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
		},
		"funding offer snapshot": {
			pld: []byte(`[
				0,
				"fos",
				[[
					41237920,"fETH",1573912039000,1573912039000,0.5,0.5,"LIMIT",
					null,null,0,"ACTIVE",null,null,null,0.0024,2,0,0,null,0,null
				]]
			]`),
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
					},
				},
			},
		},
		"funding offer new": {
			pld: []byte(`[
				0,
				"fon",
				[
					41238747,"fUST",1575026670000,1575026670000,5000,5000,"LIMIT",null,null,
					0,"ACTIVE",null,null,null,0.006000000000000001,30,0,0,null,0,null
				]
			]`),
			expected: fundingoffer.New{
				ID:         41238747,
				Symbol:     "fUST",
				MTSCreated: 1575026670000,
				MTSUpdated: 1575026670000,
				Amount:     5000,
				AmountOrig: 5000,
				Type:       "LIMIT",
				Flags:      nil,
				Status:     "ACTIVE",
				Rate:       0.006000000000000001,
				Period:     30,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
		},
		"funding offer update": {
			pld: []byte(`[
				0,
				"fou",
				[
					41238747,"fUST",1575026670000,1575026670000,5000,5000,"LIMIT",null,null,
					0,"ACTIVE",null,null,null,0.006000000000000001,30,0,0,null,0,null
				]
			]`),
			expected: fundingoffer.Update{
				ID:         41238747,
				Symbol:     "fUST",
				MTSCreated: 1575026670000,
				MTSUpdated: 1575026670000,
				Amount:     5000,
				AmountOrig: 5000,
				Type:       "LIMIT",
				Flags:      nil,
				Status:     "ACTIVE",
				Rate:       0.006000000000000001,
				Period:     30,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
		},
		"funding offer cancel": {
			pld: []byte(`[
				0,
				"foc",
				[
					41238747,"fUST",1575026670000,1575026670000,5000,5000,"LIMIT",null,null,
					0,"ACTIVE",null,null,null,0.006000000000000001,30,0,0,null,0,null
				]
			]`),
			expected: fundingoffer.Cancel{
				ID:         41238747,
				Symbol:     "fUST",
				MTSCreated: 1575026670000,
				MTSUpdated: 1575026670000,
				Amount:     5000,
				AmountOrig: 5000,
				Type:       "LIMIT",
				Flags:      nil,
				Status:     "ACTIVE",
				Rate:       0.006000000000000001,
				Period:     30,
				Notify:     false,
				Hidden:     false,
				Insure:     false,
				Renew:      false,
				RateReal:   0,
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{Data: v.pld}
			got, err := m.ProcessPrivateRaw()
			assert.NoError(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}
