package fundingcredit

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Credit struct {
	ID            int64
	Symbol        string
	Side          int
	MTSCreated    int64
	MTSUpdated    int64
	Amount        float64
	Flags         map[string]interface{}
	Status        string
	RateType      string
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
	PositionPair  string
}

type New Credit
type Update Credit
type Cancel Credit

type Snapshot struct {
	Snapshot []*Credit
}

func FromRaw(raw []interface{}) (c *Credit, err error) {
	if len(raw) < 22 {
		return c, fmt.Errorf("data slice too short for funding credit: %#v", raw)
	}

	c = &Credit{
		ID:            convert.I64ValOrZero(raw[0]),
		Symbol:        convert.SValOrEmpty(raw[1]),
		Side:          convert.ToInt(raw[2]),
		MTSCreated:    convert.I64ValOrZero(raw[3]),
		MTSUpdated:    convert.I64ValOrZero(raw[4]),
		Amount:        convert.F64ValOrZero(raw[5]),
		Status:        convert.SValOrEmpty(raw[7]),
		RateType:      convert.SValOrEmpty(raw[8]),
		Rate:          convert.F64ValOrZero(raw[11]),
		Period:        convert.I64ValOrZero(raw[12]),
		MTSOpened:     convert.I64ValOrZero(raw[13]),
		MTSLastPayout: convert.I64ValOrZero(raw[14]),
		Notify:        convert.BValOrFalse(raw[15]),
		Hidden:        convert.BValOrFalse(raw[16]),
		Insure:        convert.BValOrFalse(raw[17]),
		Renew:         convert.BValOrFalse(raw[18]),
		RateReal:      convert.F64ValOrZero(raw[19]),
		NoClose:       convert.BValOrFalse(raw[20]),
		PositionPair:  convert.SValOrEmpty(raw[21]),
	}

	if flags, ok := raw[6].(map[string]interface{}); ok {
		c.Flags = flags
	}

	return
}

func NewFromRaw(raw []interface{}) (New, error) {
	c, err := FromRaw(raw)
	if err != nil {
		return New{}, nil
	}
	return New(*c), nil
}

func UpdateFromRaw(raw []interface{}) (Update, error) {
	c, err := FromRaw(raw)
	if err != nil {
		return Update{}, nil
	}
	return Update(*c), nil
}

func CancelFromRaw(raw []interface{}) (Cancel, error) {
	c, err := FromRaw(raw)
	if err != nil {
		return Cancel{}, nil
	}
	return Cancel(*c), nil
}

func SnapshotFromRaw(raw []interface{}) (snap *Snapshot, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for funding credit: %#v", raw)
	}

	fcs := make([]*Credit, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := FromRaw(l)
				if err != nil {
					return snap, err
				}
				fcs = append(fcs, o)
			}
		}
	default:
		return snap, fmt.Errorf("not a funding credit snapshot")
	}

	snap = &Snapshot{Snapshot: fcs}
	return
}
