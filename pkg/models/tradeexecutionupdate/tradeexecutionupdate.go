package tradeexecutionupdate

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// TradeExecutionUpdate represents a full update to a trade on the private data feed.  Following a TradeExecution,
// TradeExecutionUpdates include additional details, e.g. the trade's execution ID (TradeID).
type TradeExecutionUpdate struct {
	ID          int64
	Pair        string
	MTS         int64
	OrderID     int64
	ExecAmount  float64
	ExecPrice   float64
	OrderType   string
	OrderPrice  float64
	Maker       int
	Fee         float64
	FeeCurrency string
}

type Snapshot struct {
	Snapshot []*TradeExecutionUpdate
}

type HistoricalTradeSnapshot Snapshot

// public trade update just looks like a trade
func FromRaw(raw []interface{}) (tu *TradeExecutionUpdate, err error) {
	if len(raw) == 4 {
		tu = &TradeExecutionUpdate{
			ID:         convert.I64ValOrZero(raw[0]),
			MTS:        convert.I64ValOrZero(raw[1]),
			ExecAmount: convert.F64ValOrZero(raw[2]),
			ExecPrice:  convert.F64ValOrZero(raw[3]),
		}
		return
	}
	if len(raw) > 10 {
		tu = &TradeExecutionUpdate{
			ID:          convert.I64ValOrZero(raw[0]),
			Pair:        convert.SValOrEmpty(raw[1]),
			MTS:         convert.I64ValOrZero(raw[2]),
			OrderID:     convert.I64ValOrZero(raw[3]),
			ExecAmount:  convert.F64ValOrZero(raw[4]),
			ExecPrice:   convert.F64ValOrZero(raw[5]),
			OrderType:   convert.SValOrEmpty(raw[6]),
			OrderPrice:  convert.F64ValOrZero(raw[7]),
			Maker:       convert.ToInt(raw[8]),
			Fee:         convert.F64ValOrZero(raw[9]),
			FeeCurrency: convert.SValOrEmpty(raw[10]),
		}
		return
	}
	return tu, fmt.Errorf("data slice too short for trade update: %#v", raw)
}

func SnapshotFromRaw(raw []interface{}) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("data slice is too short for trade execution update snapshot: %#v", raw)
	}

	ts := make([]*TradeExecutionUpdate, 0)
	for _, v := range raw {
		if l, ok := v.([]interface{}); ok {
			t, err := FromRaw(l)
			if err != nil {
				return s, err
			}
			ts = append(ts, t)
		}
	}

	s = &Snapshot{Snapshot: ts}
	return
}
