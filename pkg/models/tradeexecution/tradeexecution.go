package tradeexecution

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// TradeExecution represents the first message receievd for a trade on the private data feed.
type TradeExecution struct {
	ID         int64
	Pair       string
	MTS        int64
	OrderID    int64
	ExecAmount float64
	ExecPrice  float64
	OrderType  string
	OrderPrice float64
	Maker      int
}

func FromRaw(raw []interface{}) (te *TradeExecution, err error) {
	if len(raw) < 6 {
		return te, fmt.Errorf("data slice too short for trade execution: %#v", raw)
	}

	// trade executions sometimes omit order type, price, and maker flag
	te = &TradeExecution{
		ID:         convert.I64ValOrZero(raw[0]),
		Pair:       convert.SValOrEmpty(raw[1]),
		MTS:        convert.I64ValOrZero(raw[2]),
		OrderID:    convert.I64ValOrZero(raw[3]),
		ExecAmount: convert.F64ValOrZero(raw[4]),
		ExecPrice:  convert.F64ValOrZero(raw[5]),
	}

	if len(raw) >= 9 {
		te.OrderType = convert.SValOrEmpty(raw[6])
		te.OrderPrice = convert.F64ValOrZero(raw[7])
		te.Maker = convert.ToInt(raw[8])
	}

	return
}
