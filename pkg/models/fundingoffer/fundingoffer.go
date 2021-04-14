package fundingoffer

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Offer struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	MTSUpdated int64
	Amount     float64
	AmountOrig float64
	Type       string
	Flags      map[string]interface{}
	Status     string
	Rate       float64
	Period     int64
	Notify     bool
	Hidden     bool
	Insure     bool
	Renew      bool
	RateReal   float64
}

type New Offer
type Update Offer
type Cancel Offer
type Snapshot struct {
	Snapshot []*Offer
}

func FromRaw(raw []interface{}) (o *Offer, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short for funding offer: %#v", raw)
	}

	o = &Offer{
		ID:         convert.I64ValOrZero(raw[0]),
		Symbol:     convert.SValOrEmpty(raw[1]),
		MTSCreated: convert.I64ValOrZero(raw[2]),
		MTSUpdated: convert.I64ValOrZero(raw[3]),
		Amount:     convert.F64ValOrZero(raw[4]),
		AmountOrig: convert.F64ValOrZero(raw[5]),
		Type:       convert.SValOrEmpty(raw[6]),
		Status:     convert.SValOrEmpty(raw[10]),
		Rate:       convert.F64ValOrZero(raw[14]),
		Period:     convert.I64ValOrZero(raw[15]),
		Notify:     convert.BValOrFalse(raw[16]),
		Hidden:     convert.BValOrFalse(raw[17]),
		Insure:     convert.BValOrFalse(raw[18]),
		Renew:      convert.BValOrFalse(raw[19]),
		RateReal:   convert.F64ValOrZero(raw[20]),
	}

	if flags, ok := raw[9].(map[string]interface{}); ok {
		o.Flags = flags
	}

	return
}

func CancelFromRaw(raw []interface{}) (Cancel, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return Cancel{}, err
	}
	return Cancel(*o), nil
}

func NewFromRaw(raw []interface{}) (New, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return New{}, err
	}
	return New(*o), nil
}

func UpdateFromRaw(raw []interface{}) (Update, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return Update{}, err
	}
	return Update(*o), nil
}

func SnapshotFromRaw(raw []interface{}) (snap *Snapshot, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for funding offer: %#v", raw)
	}

	fos := make([]*Offer, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := FromRaw(l)
				if err != nil {
					return snap, err
				}
				fos = append(fos, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding offer snapshot")
	}

	snap = &Snapshot{Snapshot: fos}
	return
}
