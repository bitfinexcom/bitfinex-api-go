package msg_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
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
			expected: &trade.Snapshot{
				Snapshot: []*trade.Trade{
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
			expected: &trade.Trade{
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
