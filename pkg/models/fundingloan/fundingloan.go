package fundingloan

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

const (
	LoanStatusActive          LoanStatus = "ACTIVE"
	LoanStatusExecuted        LoanStatus = "EXECUTED"
	LoanStatusPartiallyFilled LoanStatus = "PARTIALLY FILLED"
	LoanStatusCanceled        LoanStatus = "CANCELED"
)

type LoanStatus string
type HistoricalLoan Loan
type New Loan
type Update Loan
type Cancel Loan

type Snapshot struct {
	Snapshot []*Loan
}

type Loan struct {
	ID            int64
	Symbol        string
	Side          string
	MTSCreated    int64
	MTSUpdated    int64
	Amount        float64
	Flags         interface{}
	Status        LoanStatus
	Rate          float64
	Period        int64
	MTSOpened     int64
	MTSLastPayout int64
	Notify        bool
	Hidden        bool
	Insure        bool
	Renew         bool
	RateReal      float64
	NoClose       bool
}

type CancelRequest struct {
	ID int64
}

func (cr *CancelRequest) ToJSON() ([]byte, error) {
	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: cr.ID,
	}
	return json.Marshal(resp)
}

// MarshalJSON converts the funding loan cancel request into the format required by the
// bitfinex websocket service.
func (cr *CancelRequest) MarshalJSON() ([]byte, error) {
	b, err := cr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"flc\", null, %s]", string(b))), nil
}

func FromRaw(raw []interface{}) (l *Loan, err error) {
	if len(raw) < 21 {
		return l, fmt.Errorf("data slice too short (len=%d) for loan: %#v", len(raw), raw)
	}

	l = &Loan{
		ID:            convert.I64ValOrZero(raw[0]),
		Symbol:        convert.SValOrEmpty(raw[1]),
		Side:          convert.SValOrEmpty(raw[2]),
		MTSCreated:    convert.I64ValOrZero(raw[3]),
		MTSUpdated:    convert.I64ValOrZero(raw[4]),
		Amount:        convert.F64ValOrZero(raw[5]),
		Flags:         raw[6],
		Status:        LoanStatus(convert.SValOrEmpty(raw[7])),
		Rate:          convert.F64ValOrZero(raw[11]),
		Period:        convert.I64ValOrZero(raw[12]),
		MTSOpened:     convert.I64ValOrZero(raw[13]),
		MTSLastPayout: convert.I64ValOrZero(raw[14]),
		Notify:        convert.BValOrFalse(raw[15]),
		Hidden:        convert.BValOrFalse(raw[16]),
		Insure:        convert.BValOrFalse(raw[17]), // DS: marked as _PLACEHOLDER in docs WS and REST
		Renew:         convert.BValOrFalse(raw[18]),
		RateReal:      convert.F64ValOrZero(raw[19]),
		NoClose:       convert.BValOrFalse(raw[20]),
	}

	return
}

func SnapshotFromRaw(raw []interface{}) (snap *Snapshot, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for loan: %#v", raw)
	}

	loans := make([]*Loan, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := FromRaw(l)
				if err != nil {
					return snap, err
				}
				loans = append(loans, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding loan snapshot")
	}

	snap = &Snapshot{Snapshot: loans}

	return
}
