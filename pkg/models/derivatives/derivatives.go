package derivatives

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type DerivativeStatusSnapshot struct {
	Snapshot []*DerivativeStatus
}

type StatusType string

const (
	DerivativeStatusType StatusType = "deriv"
)

type DerivativeStatus struct {
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
}

func NewDerivativeStatusFromWsRaw(symbol string, raw []interface{}) (*DerivativeStatus, error) {
	if len(raw) == 18 {
		ds := &DerivativeStatus{
			Symbol: symbol,
			MTS:    convert.I64ValOrZero(raw[0]),
			// placeholder
			Price:     convert.F64ValOrZero(raw[2]),
			SpotPrice: convert.F64ValOrZero(raw[3]),
			// placeholder
			InsuranceFundBalance: convert.F64ValOrZero(raw[5]),
			// placeholder
			FundingEventMTS: convert.I64ValOrZero(raw[7]),
			FundingAccrued:  convert.F64ValOrZero(raw[8]),
			FundingStep:     convert.F64ValOrZero(raw[9]),
			// placeholder
			CurrentFunding: convert.F64ValOrZero(raw[11]),
			// placeholder
			// placeholder
			MarkPrice: convert.F64ValOrZero(raw[14]),
			// placeholder
			// placeholder
			OpenInterest: convert.F64ValOrZero(raw[17]),
		}
		return ds, nil
	}

	return nil, fmt.Errorf("unexpected data slice length for derivative status: %#v", raw)
}

func NewDerivativeStatusFromRaw(raw []interface{}) (*DerivativeStatus, error) {
	if len(raw) == 19 {
		ds := &DerivativeStatus{
			Symbol: convert.SValOrEmpty(raw[0]),
			MTS:    convert.I64ValOrZero(raw[1]),
			// placeholder
			Price:     convert.F64ValOrZero(raw[3]),
			SpotPrice: convert.F64ValOrZero(raw[4]),
			// placeholder
			InsuranceFundBalance: convert.F64ValOrZero(raw[6]),
			// placeholder
			FundingEventMTS: convert.I64ValOrZero(raw[8]),
			FundingAccrued:  convert.F64ValOrZero(raw[9]),
			FundingStep:     convert.F64ValOrZero(raw[10]),
			// placeholder
			CurrentFunding: convert.F64ValOrZero(raw[12]),
			// placeholder
			// placeholder
			MarkPrice: convert.F64ValOrZero(raw[15]),
			// placeholder
			// placeholder
			OpenInterest: convert.F64ValOrZero(raw[18]),
		}

		return ds, nil
	}

	return nil, fmt.Errorf("unexpected data slice length for derivative status: %#v", raw)
}

func NewDerivativeSnapshotFromRaw(raw [][]interface{}) (*DerivativeStatusSnapshot, error) {
	snapshot := make([]*DerivativeStatus, len(raw))
	for i, rStatus := range raw {
		pStatus, err := NewDerivativeStatusFromRaw(rStatus)
		if err != nil {
			return nil, err
		}
		snapshot[i] = pStatus
	}
	return &DerivativeStatusSnapshot{Snapshot: snapshot}, nil
}
