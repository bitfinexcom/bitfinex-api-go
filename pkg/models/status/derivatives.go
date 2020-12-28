package status

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type DerivativesSnapshot struct {
	Snapshot []*Derivative
}

type Derivative struct {
	Symbol               string
	MTS                  int64
	Price                float64
	SpotPrice            float64
	InsuranceFundBalance float64
	FundingEventMTS      int64
	FundingAccrued       float64
	FundingStep          float64
	CurrentFunding       float64
	MarkPrice            float64
	OpenInterest         float64
	ClampMIN             float64
	ClampMAX             float64
}

func DerivFromRaw(symbol string, raw []interface{}) (*Derivative, error) {
	if len(raw) < 22 {
		return nil, fmt.Errorf("data slice too short for derivative status: %#v", raw)
	}

	return &Derivative{
		Symbol:               symbol,
		MTS:                  convert.I64ValOrZero(raw[0]),
		Price:                convert.F64ValOrZero(raw[2]),
		SpotPrice:            convert.F64ValOrZero(raw[3]),
		InsuranceFundBalance: convert.F64ValOrZero(raw[5]),
		FundingEventMTS:      convert.I64ValOrZero(raw[7]),
		FundingAccrued:       convert.F64ValOrZero(raw[8]),
		FundingStep:          convert.F64ValOrZero(raw[9]),
		CurrentFunding:       convert.F64ValOrZero(raw[11]),
		MarkPrice:            convert.F64ValOrZero(raw[14]),
		OpenInterest:         convert.F64ValOrZero(raw[17]),
		ClampMIN:             convert.F64ValOrZero(raw[21]),
		ClampMAX:             convert.F64ValOrZero(raw[22]),
	}, nil
}

func DerivSnapshotFromRaw(symbol string, raw [][]interface{}) (*DerivativesSnapshot, error) {
	snapshot := make([]*Derivative, len(raw))
	for i, r := range raw {
		d, err := DerivFromRaw(symbol, r)
		if err != nil {
			return nil, err
		}
		snapshot[i] = d
	}
	return &DerivativesSnapshot{Snapshot: snapshot}, nil
}

func DerivFromRestRaw(raw []interface{}) (t *Derivative, err error) {
	if len(raw) < 2 {
		return t, fmt.Errorf("data slice too short for derivatives: %#v", raw)
	}

	return DerivFromRaw(raw[0].(string), raw[1:])
}
