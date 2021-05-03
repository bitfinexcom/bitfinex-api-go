package tickerhist

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type TickerHist struct {
	Symbol string
	Bid    float64
	// PLACEHOLDER,
	Ask float64
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	// PLACEHOLDER,
	MTS int64
}

var tickerHistFields = map[string]int{
	"Symbol": 0,
	"Bid":    1,
	"Ask":    3,
	"Mts":    12,
}

type Snapshot struct {
	Snapshot []TickerHist
}

func SnapshotFromRaw(raw [][]interface{}) (ss Snapshot) {
	if len(raw) == 0 {
		return
	}

	snap := make([]TickerHist, 0)
	for _, r := range raw {
		th, err := FromRaw(r)
		if err != nil {
			continue
		}
		snap = append(snap, th)
	}

	return Snapshot{Snapshot: snap}
}

func FromRaw(raw []interface{}) (t TickerHist, err error) {
	// to avoid index out of range issue
	if len(raw) < 13 {
		err = fmt.Errorf("data slice too short for ticker history, data:%#v", raw)
		return
	}

	t = TickerHist{
		Symbol: convert.SValOrEmpty(raw[tickerHistFields["Symbol"]]),
		Bid:    convert.F64ValOrZero(raw[tickerHistFields["Bid"]]),
		Ask:    convert.F64ValOrZero(raw[tickerHistFields["Ask"]]),
		MTS:    convert.I64ValOrZero(raw[tickerHistFields["Mts"]]),
	}
	return
}
