package fundingtrade

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type FundingTrade struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	OfferID    int64
	Amount     float64
	Rate       float64
	Period     int64
	Maker      int64
}

type Execution FundingTrade
type Update FundingTrade
type Snapshot struct {
	Snapshot []*FundingTrade
}
type HistoricalSnapshot Snapshot

func FromRaw(raw []interface{}) (ft *FundingTrade, err error) {
	if len(raw) < 8 {
		return ft, fmt.Errorf("data slice too short for funding trade: %#v", raw)
	}

	ft = &FundingTrade{
		ID:         convert.I64ValOrZero(raw[0]),
		Symbol:     convert.SValOrEmpty(raw[1]),
		MTSCreated: convert.I64ValOrZero(raw[2]),
		OfferID:    convert.I64ValOrZero(raw[3]),
		Amount:     convert.F64ValOrZero(raw[4]),
		Rate:       convert.F64ValOrZero(raw[5]),
		Period:     convert.I64ValOrZero(raw[6]),
		Maker:      convert.I64ValOrZero(raw[7]),
	}

	return
}

func SnapshotFromRaw(raw []interface{}) (snap *Snapshot, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for funding trade")
	}

	fts := make([]*FundingTrade, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := FromRaw(l)
				if err != nil {
					return snap, err
				}
				fts = append(fts, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding trade snapshot")
	}
	snap = &Snapshot{
		Snapshot: fts,
	}

	return
}

func HistoricalSnapshotFromRaw(raw []interface{}) (HistoricalSnapshot, error) {
	s, err := SnapshotFromRaw(raw)
	if err != nil {
		return HistoricalSnapshot{}, err
	}

	return HistoricalSnapshot(*s), nil
}
