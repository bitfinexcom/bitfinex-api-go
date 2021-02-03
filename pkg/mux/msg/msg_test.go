package msg_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/margin"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/notification"
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
		"ticker trading pair snapshot": {
			pld: []byte(`[
				111,
				[[
					7616.5,31.89055171,7617.5,43.358118629999986,
					-550.8,-0.0674,7617.1,8314.71200815,8257.8,7500
				]]
			]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "ticker",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: &ticker.Snapshot{
				Snapshot: []*ticker.Ticker{
					{
						Symbol:          "tBTCUST",
						Bid:             7616.5,
						BidSize:         31.89055171,
						Ask:             7617.5,
						AskSize:         43.358118629999986,
						DailyChange:     -550.8,
						DailyChangePerc: -0.0674,
						LastPrice:       7617.1,
						Volume:          8314.71200815,
						High:            8257.8,
						Low:             7500,
					},
				},
			},
		},
		"ticker trading pair update": {
			pld: []byte(`[
				111,
				[
					7617,52.98726298,7617.1,53.601795929999994,
					-550.9,-0.0674,7617,8318.92961981,8257.8,7500
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
				Frr:             0,
				Bid:             7617,
				BidPeriod:       0,
				BidSize:         52.98726298,
				Ask:             7617.1,
				AskPeriod:       0,
				AskSize:         53.601795929999994,
				DailyChange:     -550.9,
				DailyChangePerc: -0.0674,
				LastPrice:       7617,
				Volume:          8318.92961981,
				High:            8257.8,
				Low:             7500,
			},
		},
		"ticker funding pair snapshot": {
			pld: []byte(`[
				111,
				[[
					0.0003193369863013699,0.0002401,30,3939629.6177260396,0.00019012,2,
					307776.1592138799,-0.00005823,-0.2344,0.00019016,122156333.45260866,
					0.00027397,6.8e-7,null,null,3441851.73330503
				]]
			]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "ticker",
						Symbol:  "fUSD",
					},
				},
			},
			expected: &ticker.Snapshot{
				Snapshot: []*ticker.Ticker{
					{
						Symbol:             "fUSD",
						Frr:                0.0003193369863013699,
						Bid:                0.0002401,
						BidPeriod:          30,
						BidSize:            3.9396296177260396e+06,
						Ask:                0.00019012,
						AskPeriod:          2,
						AskSize:            307776.1592138799,
						DailyChange:        -5.823e-05,
						DailyChangePerc:    -0.2344,
						LastPrice:          0.00019016,
						Volume:             1.2215633345260866e+08,
						High:               0.00027397,
						Low:                6.8e-07,
						FrrAmountAvailable: 3.44185173330503e+06,
					},
				},
			},
		},
		"ticker funding trading pair update": {
			pld: []byte(`[
				111,
				[
					0.0003193315068493151,0.0002401,30,4037829.0804227195,0.000189,4,
					384507.7314462898,-0.00005939,-0.2391,0.000189,122159083.98991197,
					0.00027397,6.8e-7,null,null,3441851.73330503
				]
			]`),
			inf: map[int64]event.Info{
				111: {
					Subscribe: event.Subscribe{
						Channel: "ticker",
						Symbol:  "fUSD",
					},
				},
			},
			expected: &ticker.Ticker{
				Symbol:             "fUSD",
				Frr:                0.0003193315068493151,
				Bid:                0.0002401,
				BidPeriod:          30,
				BidSize:            4.0378290804227195e+06,
				Ask:                0.000189,
				AskPeriod:          4,
				AskSize:            384507.7314462898,
				DailyChange:        -5.939e-05,
				DailyChangePerc:    -0.2391,
				LastPrice:          0.000189,
				Volume:             1.2215908398991197e+08,
				High:               0.00027397,
				Low:                6.8e-07,
				FrrAmountAvailable: 3.44185173330503e+06,
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
		"trade execution": {
			pld: []byte(`[17470,"te",[401597395,1574694478808,0.005,7245.3]]`),
			inf: map[int64]event.Info{
				17470: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: trades.TradeExecuted{
				Pair:   "tBTCUST",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
		},
		"trade execution update": {
			pld: []byte(`[17470,"tu",[401597395,1574694478808,0.005,7245.3]]`),
			inf: map[int64]event.Info{
				17470: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "tBTCUST",
					},
				},
			},
			expected: trades.TradeExecutionUpdate{
				Pair:   "tBTCUST",
				ID:     401597395,
				MTS:    1574694478808,
				Amount: 0.005,
				Price:  7245.3,
			},
		},
		"funding trade execution": {
			pld: []byte(`[337371,"fte",[133323543,1574694605000,-59.84,0.00023647,2]]`),
			inf: map[int64]event.Info{
				337371: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "fUSD",
					},
				},
			},
			expected: trades.FundingTradeExecuted{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
		},
		"funding trade execution update": {
			pld: []byte(`[337371,"ftu",[133323543,1574694605000,-59.84,0.00023647,2]]`),
			inf: map[int64]event.Info{
				337371: {
					Subscribe: event.Subscribe{
						Channel: "trades",
						Symbol:  "fUSD",
					},
				},
			},
			expected: trades.FundingTradeExecutionUpdate{
				Symbol: "fUSD",
				ID:     133323543,
				MTS:    1574694605000,
				Amount: -59.84,
				Rate:   0.00023647,
				Period: 2,
			},
		},
		"book snapshot trading pair bid entry": {
			pld: []byte(`[17082,[[7254.7,3,3.3]]]`),
			inf: map[int64]event.Info{
				17082: {
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
						Count:       3,
						Price:       7254.7,
						Amount:      3.3,
						PriceJsNum:  "7254.7",
						AmountJsNum: "3.3",
						Side:        1,
						Action:      0,
					},
				},
			},
		},
		"book snapshot trading pair ask entry": {
			pld: []byte(`[17082,[[7254.7,3,-3.3]]]`),
			inf: map[int64]event.Info{
				17082: {
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
						Count:       3,
						Price:       7254.7,
						Amount:      3.3,
						PriceJsNum:  "7254.7",
						AmountJsNum: "-3.3",
						Side:        2,
						Action:      0,
					},
				},
			},
		},
		"book snapshot trading pair exit": {
			pld: []byte(`[17082,[[7254.7,0,3.3]]]`),
			inf: map[int64]event.Info{
				17082: {
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
						Count:       0,
						Price:       7254.7,
						Amount:      3.3,
						PriceJsNum:  "7254.7",
						AmountJsNum: "3.3",
						Side:        1,
						Action:      1,
					},
				},
			},
		},
		"book snapshot funding pair": {
			pld: []byte(`[431549,[[0.00023112,30,1,-15190.7005375]]]`),
			inf: map[int64]event.Info{
				431549: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "fUSD",
						Precision: "P0",
						Frequency: "F0",
					},
				},
			},
			expected: &book.Snapshot{
				Snapshot: []*book.Book{
					{
						Symbol:      "fUSD",
						Count:       1,
						Period:      30,
						Amount:      -15190.7005375,
						Rate:        0.00023112,
						AmountJsNum: "-15190.7005375",
					},
				},
			},
		},
		"book trading pair update": {
			pld: []byte(`[17082,[7254.7,3,3.3]]`),
			inf: map[int64]event.Info{
				17082: {
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
				Count:       3,
				Price:       7254.7,
				Amount:      3.3,
				PriceJsNum:  "7254.7",
				AmountJsNum: "3.3",
				Side:        1,
				Action:      0,
			},
		},
		"book funding pair update": {
			pld: []byte(`[348748,[0.00023157,2,1,66.35007188]]`),
			inf: map[int64]event.Info{
				348748: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "fUSD",
						Precision: "P0",
						Frequency: "F0",
					},
				},
			},
			expected: &book.Book{
				Symbol:      "fUSD",
				Count:       1,
				Period:      2,
				Amount:      66.35007188,
				Rate:        0.00023157,
				AmountJsNum: "66.35007188",
			},
		},
		"raw trading pair book snapshot bid entry": {
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
		"raw trading pair book snapshot ask entry": {
			pld: []byte(`[869944,[[55804480297,33766,-2]]]`),
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
						PriceJsNum:  "33766",
						AmountJsNum: "-2",
						Side:        2,
						Action:      0,
						Price:       33766,
						Amount:      2,
					},
				},
			},
		},
		"raw trading pair book snapshot remove entry": {
			pld: []byte(`[869944,[[55804480297,0,2]]]`),
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
						Price:       0,
						Amount:      2,
						PriceJsNum:  "0",
						AmountJsNum: "2",
						Side:        1,
						Action:      1,
					},
				},
			},
		},
		"raw funding pair book snapshot": {
			pld: []byte(`[472778,[[658282397,30,0.000233,-530]]]`),
			inf: map[int64]event.Info{
				472778: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "fUSD",
						Precision: "R0",
					},
				},
			},
			expected: &book.Snapshot{
				Snapshot: []*book.Book{
					{
						Symbol:      "fUSD",
						ID:          658282397,
						Period:      30,
						Amount:      -530,
						Rate:        0.000233,
						AmountJsNum: "-530",
					},
				},
			},
		},
		"raw trading pair book update": {
			pld: []byte(`[433290,[34753006045,0,-1]]`),
			inf: map[int64]event.Info{
				433290: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "tBTCUSD",
						Precision: "R0",
					},
				},
			},
			expected: &book.Book{
				Symbol:      "tBTCUSD",
				ID:          34753006045,
				PriceJsNum:  "0",
				AmountJsNum: "-1",
				Amount:      1,
				Side:        2,
				Action:      1,
			},
		},
		"raw funding pair book update": {
			pld: []byte(`[472778,[658286906,2,0,1]]`),
			inf: map[int64]event.Info{
				472778: {
					Subscribe: event.Subscribe{
						Channel:   "book",
						Symbol:    "fUSD",
						Precision: "R0",
					},
				},
			},
			expected: &book.Book{
				Symbol:      "fUSD",
				ID:          658286906,
				Period:      2,
				Amount:      1,
				Rate:        0,
				AmountJsNum: "1",
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
		"funding credits snapshot": {
			pld: []byte(`[
				0,
				"fcs",
				[[
					26223578,"fUST",1,1575052261000,1575296187000,350,0,"ACTIVE",null,null,
					null,0,30,1575052261000,1575293487000,0,0,null,0,null,0,"tBTCUST"
				]]
			]`),
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
						Period:        30,
						MTSOpened:     1575052261000,
						MTSLastPayout: 1575293487000,
						PositionPair:  "tBTCUST",
					},
				},
			},
		},
		"funding credit new": {
			pld: []byte(`[
				0,
				"fcn",
				[
					26223578,"fUST",1,1575052261000,1575296787000,350,0,"ACTIVE",null,null,
					null,0,30,1575052261000,1575293487000,0,0,null,0,null,0,"tBTCUST"
				]
			]`),
			expected: fundingcredit.New{
				ID:            26223578,
				Symbol:        "fUST",
				Side:          1,
				MTSCreated:    1575052261000,
				MTSUpdated:    1575296787000,
				Amount:        350,
				Flags:         map[string]interface{}(nil),
				Status:        "ACTIVE",
				RateType:      "",
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
		"funding credit update": {
			pld: []byte(`[
				0,
				"fcu",
				[
					26223578,"fUST",1,1575052261000,1575296787000,350,0,"ACTIVE",null,null,
					null,0,30,1575052261000,1575293487000,0,0,null,0,null,0,"tBTCUST"
				]
			]`),
			expected: fundingcredit.Update{
				ID:            26223578,
				Symbol:        "fUST",
				Side:          1,
				MTSCreated:    1575052261000,
				MTSUpdated:    1575296787000,
				Amount:        350,
				Flags:         map[string]interface{}(nil),
				Status:        "ACTIVE",
				RateType:      "",
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
		"funding credit close": {
			pld: []byte(`[
				0,
				"fcc",
				[
					26223578,"fUST",1,1575052261000,1575296787000,350,0,"ACTIVE",null,null,
					null,0,30,1575052261000,1575293487000,0,0,null,0,null,0,"tBTCUST"
				]
			]`),
			expected: fundingcredit.Cancel{
				ID:            26223578,
				Symbol:        "fUST",
				Side:          1,
				MTSCreated:    1575052261000,
				MTSUpdated:    1575296787000,
				Amount:        350,
				Flags:         map[string]interface{}(nil),
				Status:        "ACTIVE",
				RateType:      "",
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
		"funding loan snapshot": {
			pld: []byte(`[
				0,
				"fls",
				[[
					2995442,"fUSD",-1,1575291961000,1575295850000,820,0,"ACTIVE",null,
					null,null,0.002,7,1575282446000,1575295850000,0,0,null,0,null,0
				]]
			]`),
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
					},
				},
			},
		},
		"funding loan new": {
			pld: []byte(`[
				0,
				"fln",
				[
					2995444,"fUSD",-1,1575298742000,1575298742000,1000,0,"ACTIVE",null,
					null,null,0.002,7,1575298742000,1575298742000,0,0,null,0,null,0
				]
			]`),
			expected: fundingloan.New{
				ID:            2995444,
				Symbol:        "fUSD",
				Side:          -1,
				MTSCreated:    1575298742000,
				MTSUpdated:    1575298742000,
				Amount:        1000,
				Status:        "ACTIVE",
				Rate:          0.002,
				Period:        7,
				MTSOpened:     1575298742000,
				MTSLastPayout: 1575298742000,
			},
		},
		"funding loan update": {
			pld: []byte(`[
				0,
				"flu",
				[
					2995444,"fUSD",-1,1575298742000,1575298742000,1000,0,"ACTIVE",null,
					null,null,0.002,7,1575298742000,1575298742000,0,0,null,0,null,0
				]
			]`),
			expected: fundingloan.Update{
				ID:            2995444,
				Symbol:        "fUSD",
				Side:          -1,
				MTSCreated:    1575298742000,
				MTSUpdated:    1575298742000,
				Amount:        1000,
				Status:        "ACTIVE",
				Rate:          0.002,
				Period:        7,
				MTSOpened:     1575298742000,
				MTSLastPayout: 1575298742000,
			},
		},
		"funding loan close": {
			pld: []byte(`[
				0,
				"flc",
				[
					2995444,"fUSD",-1,1575298742000,1575298742000,1000,0,"ACTIVE",null,
					null,null,0.002,7,1575298742000,1575298742000,0,0,null,0,null,0
				]
			]`),
			expected: fundingloan.Cancel{
				ID:            2995444,
				Symbol:        "fUSD",
				Side:          -1,
				MTSCreated:    1575298742000,
				MTSUpdated:    1575298742000,
				Amount:        1000,
				Status:        "ACTIVE",
				Rate:          0.002,
				Period:        7,
				MTSOpened:     1575298742000,
				MTSLastPayout: 1575298742000,
			},
		},
		"margin info update symbol calc": {
			pld: []byte(`[
				0,
				"miu",
				[
					"sym","tETHUSD",
					[149361.09689202666,149639.26293509,830.0182168075556,895.0658432466332,null,null,null,null]
				]
			]`),
			expected: &margin.InfoUpdate{
				Symbol:          "tETHUSD",
				TradableBalance: 149361.09689202666,
				GrossBalance:    149639.26293509,
				Buy:             830.0182168075556,
				Sell:            895.0658432466332,
			},
		},
		"margin info update base calc": {
			pld: []byte(`[
				0,
				"miu",
				["base",[-13.014640000000007,0,49331.70267297,49318.68803297,27]]
			]`),
			expected: &margin.InfoBase{
				UserProfitLoss: -13.014640000000007,
				UserSwaps:      0,
				MarginBalance:  49331.70267297,
				MarginNet:      49318.68803297,
				MarginRequired: 27,
			},
		},
		"order new notification": {
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
		},
		"order udate notification": {
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
		},
		"funding offer new notification": {
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
