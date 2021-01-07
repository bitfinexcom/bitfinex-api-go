package trades

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Trade struct {
	Pair   string
	ID     int64
	MTS    int64
	Amount float64
	Price  float64
}

type TradeUpdate Trade
type TradeExecute Trade

type TradeSnapshot struct {
	Snapshot []Trade
}

func TFromRaw(pair string, raw []interface{}) (t Trade, err error) {
	if len(raw) >= 4 {
		t = Trade{
			Pair:   pair,
			ID:     convert.I64ValOrZero(raw[0]),
			MTS:    convert.I64ValOrZero(raw[1]),
			Amount: convert.F64ValOrZero(raw[2]),
			Price:  convert.F64ValOrZero(raw[3]),
		}
		return
	}

	err = fmt.Errorf("data slice too short: %#v", raw)
	return
}

func TUFromRaw(pair string, raw []interface{}) (TradeUpdate, error) {
	t, err := TFromRaw(pair, raw)
	if err != nil {
		return TradeUpdate{}, err
	}

	return TradeUpdate(t), nil
}

func TEFromRaw(pair string, raw []interface{}) (TradeExecute, error) {
	t, err := TFromRaw(pair, raw)
	if err != nil {
		return TradeExecute{}, err
	}

	return TradeExecute(t), nil
}

func TSnapshotFromRaw(pair string, raw [][]interface{}) (TradeSnapshot, error) {
	if len(raw) == 0 {
		return TradeSnapshot{}, fmt.Errorf("trade snapshot data slice too short:%#v", raw)
	}

	snapshot := make([]Trade, 0)
	for _, v := range raw {
		t, err := TFromRaw(pair, v)
		if err != nil {
			return TradeSnapshot{}, err
		}
		snapshot = append(snapshot, t)
	}

	return TradeSnapshot{Snapshot: snapshot}, nil
}
