package ledger

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Ledger struct {
	ID       int64
	Currency string
	// placeholder
	MTS int64
	// placeholder
	Amount  float64
	Balance float64
	// placeholder
	Description string
}

type Snapshot struct {
	Snapshot []*Ledger
}

type transformerFn func(raw []interface{}) (w *Ledger, err error)

// FromRaw takes the raw list of values as returned from the websocket
// service and tries to convert it into an Ledger.
func FromRaw(raw []interface{}) (o *Ledger, err error) {
	if len(raw) < 9 {
		return o, fmt.Errorf("data slice too short for ledger: %#v", raw)
	}

	o = &Ledger{
		ID:          convert.I64ValOrZero(raw[0]),
		Currency:    convert.SValOrEmpty(raw[1]),
		MTS:         convert.I64ValOrZero(raw[3]),
		Amount:      convert.F64ValOrZero(raw[5]),
		Balance:     convert.F64ValOrZero(raw[6]),
		Description: convert.SValOrEmpty(raw[8]),
	}

	return
}

// SnapshotFromRaw takes a raw list of values as returned from the websocket
// service and tries to convert it into an Snapshot.
func SnapshotFromRaw(raw []interface{}, t transformerFn) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for ledgers: %#v", raw)
	}

	lss := make([]*Ledger, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := t(l)
				if err != nil {
					return s, err
				}
				lss = append(lss, o)
			}
		}
	default:
		return s, fmt.Errorf("not an ledger snapshot")
	}
	s = &Snapshot{Snapshot: lss}
	return
}
