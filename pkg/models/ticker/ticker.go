package ticker

import (
	"errors"
	"fmt"
	"strings"

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
	// PLACEHOLDER,
	// PLACEHOLDER,
	FrrAmountAvailable float64
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
		if err != nil {
			return nil, err
		}
		snap = append(snap, c)
	}

	return &Snapshot{Snapshot: snap}, nil
}

func FromRaw(symbol string, raw []interface{}) (t *Ticker, err error) {
	// trading pair update / snapshot
	if strings.HasPrefix(symbol, "t") && len(raw) >= 10 {
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

	// funding pair update
	if strings.HasPrefix(symbol, "f") {
		if len(raw) >= 13 {
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
		}

		// funding pair snapshot
		if len(raw) >= 16 {
			t.FrrAmountAvailable = convert.F64ValOrZero(raw[15])
		}
		return
	}

	err = fmt.Errorf("unrecognized data slice format for pair:%s, data:%#v", symbol, raw)
	return
}

func FromRestRaw(raw []interface{}) (t *Ticker, err error) {
	if len(raw) == 0 {
		return t, fmt.Errorf("data slice too short for ticker")
	}

	return FromRaw(raw[0].(string), raw[1:])
}

// FromWSRaw - based on condition will return snapshot of tickers or single tick
func FromWSRaw(symbol string, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice")
	}

	_, isSnapshot := data[0].([]interface{})
	if isSnapshot {
		return SnapshotFromRaw(symbol, convert.ToInterfaceArray(data))
	}

	return FromRaw(symbol, data)
}
