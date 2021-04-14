package status

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type LiquidationsSnapshot struct {
	Snapshot []*Liquidation
}

type Liquidation struct {
	Symbol        string
	PositionID    int64
	MTS           int64
	Amount        float64
	BasePrice     float64
	IsMatch       int
	IsMarketSold  int
	PriceAcquired float64
}

func LiqFromRaw(raw []interface{}) (*Liquidation, error) {
	if len(raw) < 12 {
		return nil, fmt.Errorf("data slice too short for liquidation status: %#v", raw)
	}

	return &Liquidation{
		PositionID:    convert.I64ValOrZero(raw[1]),
		MTS:           convert.I64ValOrZero(raw[2]),
		Symbol:        convert.SValOrEmpty(raw[4]),
		Amount:        convert.F64ValOrZero(raw[5]),
		BasePrice:     convert.F64ValOrZero(raw[6]),
		IsMatch:       convert.ToInt(raw[8]),
		IsMarketSold:  convert.ToInt(raw[9]),
		PriceAcquired: convert.F64ValOrZero(raw[11]),
	}, nil
}

func LiqSnapshotFromRaw(raw [][]interface{}) (*LiquidationsSnapshot, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty data slice")
	}

	snapshot := make([]*Liquidation, len(raw))
	for i, r := range raw {
		l, err := LiqFromRaw(r)
		if err != nil {
			return nil, err
		}
		snapshot[i] = l
	}
	return &LiquidationsSnapshot{Snapshot: snapshot}, nil
}
