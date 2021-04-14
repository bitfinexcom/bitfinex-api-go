package trades

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// Trade data structure for mapping trading pair currency raw
// data with "t" prefix in SYMBOL from public feed
type Trade struct {
	Pair   string
	ID     int64
	MTS    int64
	Amount float64
	Price  float64
}

type TradeExecutionUpdate Trade
type TradeExecuted Trade

type TradeSnapshot struct {
	Snapshot []Trade
}

// TFromRaw maps raw data slice to instance of Trade
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

// TEUFromRaw maps raw data slice to instance of TradeExecutionUpdate
func TEUFromRaw(pair string, raw []interface{}) (TradeExecutionUpdate, error) {
	t, err := TFromRaw(pair, raw)
	if err != nil {
		return TradeExecutionUpdate{}, err
	}

	return TradeExecutionUpdate(t), nil
}

// TEFromRaw maps raw data slice to instance of TradeExecuted
func TEFromRaw(pair string, raw []interface{}) (TradeExecuted, error) {
	t, err := TFromRaw(pair, raw)
	if err != nil {
		return TradeExecuted{}, err
	}

	return TradeExecuted(t), nil
}

// TSnapshotFromRaw maps raw data slice to trading data structures
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
