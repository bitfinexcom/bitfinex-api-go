package rest

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
)

type TickerService struct {
	Synchronous
}

func (t *TickerService) All() (*bitfinex.TickerSnapshot, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", "ALL")
	raw, err := t.Request(req)

	if err != nil {
		return nil, err
	}

	tickers := make([]*bitfinex.Ticker, len(raw))
	for i, ifacearr := range raw {
		arr, ok := ifacearr.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expecting array, got %T", ifacearr)
		}
		symbol, ok := arr[0].(string)
		if !ok {
			return nil, fmt.Errorf("expecting string, got %T", arr[0])
		}
		if len(symbol) <= 1 || (symbol[0] != 't' && symbol[0] != 'f') {
			return nil, errors.New("invalid symbol")
		}
		if (symbol[0] == 't' && len(arr) < 11) || (symbol[0] == 'f' && len(arr) < 14) {
			return nil, errors.New("invalid length of ticker")
		}
		sub := make([]float64, len(arr)-1)
		for j, iface := range arr[1:] {
			if iface == nil {
				sub[j] = 0
				continue
			}
			flt, ok := iface.(float64)
			if !ok {
				return nil, fmt.Errorf("expecting float64, got %T", iface)
			}
			sub[j] = flt
		}
		var entry *bitfinex.Ticker
		switch symbol[0] {
		case 't':
			entry = &bitfinex.Ticker{
				Symbol:          strings.ToLower(symbol[1:]),
				Bid:             sub[0],
				BidSize:         sub[1],
				Ask:             sub[2],
				AskSize:         sub[3],
				DailyChange:     sub[4],
				DailyChangePerc: sub[5],
				LastPrice:       sub[6],
				Volume:          sub[7],
				High:            sub[8],
				Low:             sub[9],
			}
		case 'f':
			entry = &bitfinex.Ticker{
				Symbol:          strings.ToLower(symbol[1:]),
				FRR:             sub[0],
				Bid:             sub[1],
				BidSize:         sub[2],
				BidPeriod:       int64(sub[3]),
				Ask:             sub[4],
				AskSize:         sub[5],
				AskPeriod:       int64(sub[6]),
				DailyChange:     sub[7],
				DailyChangePerc: sub[8],
				LastPrice:       sub[9],
				Volume:          sub[10],
				High:            sub[11],
				Low:             sub[12],
			}
		}
		tickers[i] = entry
	}
	return &bitfinex.TickerSnapshot{Snapshot: tickers}, nil
}
