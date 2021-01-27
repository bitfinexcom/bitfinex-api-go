package trades

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// AuthExecutionUpdate used for mapping authenticated trade execution update raw messages
type AuthExecutionUpdate struct {
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

func AuthExecutionUpdateFromRaw(raw []interface{}) (eu AuthExecutionUpdate, err error) {
	if len(raw) < 11 {
		return AuthExecutionUpdate{}, fmt.Errorf("data slice too short for auth trade execution update: %#v", raw)
	}

	eu = AuthExecutionUpdate{
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
