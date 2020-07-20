package ticker

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Ticker struct {
	Symbol          string
	Frr             float64
	Bid             float64
	BidPeriod       int64
	BidSize         float64
	Ask             float64
	AskPeriod       int64
	AskSize         float64
	DailyChange     float64
	DailyChangePerc float64
	LastPrice       float64
	Volume          float64
	High            float64
	Low             float64
}

type Update Ticker
type Snapshot struct {
	Snapshot []*Ticker
}

func SnapshotFromRaw(symbol string, raw [][]interface{}) (*Snapshot, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("data slice too short for ticker snapshot: %#v", raw)
	}

	snap := make([]*Ticker, 0)
	for _, f := range raw {
		c, err := FromRaw(symbol, f)
		if err == nil {
			snap = append(snap, c)
		}
	}

	return &Snapshot{Snapshot: snap}, nil
}

func FromRaw(symbol string, raw []interface{}) (t *Ticker, err error) {
	if len(raw) < 10 {
		return t, fmt.Errorf("data slice too short for ticker, expected %d got %d: %#v", 10, len(raw), raw)
	}

	// funding currency ticker
	// ignore bid/ask period for now
	if len(raw) == 13 {
		t = &Ticker{
			Symbol:          symbol,
			Bid:             convert.F64ValOrZero(raw[1]),
			BidSize:         convert.F64ValOrZero(raw[2]),
			Ask:             convert.F64ValOrZero(raw[4]),
			AskSize:         convert.F64ValOrZero(raw[5]),
			DailyChange:     convert.F64ValOrZero(raw[7]),
			DailyChangePerc: convert.F64ValOrZero(raw[8]),
			LastPrice:       convert.F64ValOrZero(raw[9]),
			Volume:          convert.F64ValOrZero(raw[10]),
			High:            convert.F64ValOrZero(raw[11]),
			Low:             convert.F64ValOrZero(raw[12]),
		}
		return
	}

	if len(raw) == 16 {
		t = &Ticker{
			Symbol:          symbol,
			Frr:             convert.F64ValOrZero(raw[0]),
			Bid:             convert.F64ValOrZero(raw[1]),
			BidPeriod:       convert.I64ValOrZero(raw[2]),
			BidSize:         convert.F64ValOrZero(raw[3]),
			Ask:             convert.F64ValOrZero(raw[4]),
			AskPeriod:       convert.I64ValOrZero(raw[5]),
			AskSize:         convert.F64ValOrZero(raw[6]),
			DailyChange:     convert.F64ValOrZero(raw[7]),
			DailyChangePerc: convert.F64ValOrZero(raw[8]),
			LastPrice:       convert.F64ValOrZero(raw[9]),
			Volume:          convert.F64ValOrZero(raw[10]),
			High:            convert.F64ValOrZero(raw[11]),
			Low:             convert.F64ValOrZero(raw[12]),
		}
		return
	}

	// all other tickers
	// on trading pairs (ex. tBTCUSD)
	t = &Ticker{
		Symbol:          symbol,
		Bid:             convert.F64ValOrZero(raw[0]),
		BidSize:         convert.F64ValOrZero(raw[1]),
		Ask:             convert.F64ValOrZero(raw[2]),
		AskSize:         convert.F64ValOrZero(raw[3]),
		DailyChange:     convert.F64ValOrZero(raw[4]),
		DailyChangePerc: convert.F64ValOrZero(raw[5]),
		LastPrice:       convert.F64ValOrZero(raw[6]),
		Volume:          convert.F64ValOrZero(raw[7]),
		High:            convert.F64ValOrZero(raw[8]),
		Low:             convert.F64ValOrZero(raw[9]),
	}
	return
}

func FromRestRaw(raw []interface{}) (t *Ticker, err error) {
	if len(raw) == 0 {
		return t, fmt.Errorf("data slice too short for ticker")
	}

	return FromRaw(raw[0].(string), raw[1:])
}
