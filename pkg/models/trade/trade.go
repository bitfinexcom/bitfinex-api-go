package trade

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// Trade represents a trade on the public data feed.
type Trade struct {
	Pair   string
	ID     int64
	MTS    int64
	Amount float64
	Price  float64
	Rate   float64
	Period int
}

type Snapshot struct {
	Snapshot []*Trade
}

func FromRaw(pair string, raw []interface{}) (t *Trade, err error) {
	if len(raw) < 4 {
		return t, fmt.Errorf("data slice too short for trade: %#v", raw)
	}

	t = &Trade{
		Pair:   pair,
		ID:     convert.I64ValOrZero(raw[0]),
		MTS:    convert.I64ValOrZero(raw[1]),
		Amount: convert.F64ValOrZero(raw[2]),
	}

	if len(raw) == 4 {
		t.Price = convert.F64ValOrZero(raw[3])
	}

	if len(raw) >= 5 {
		t.Rate = convert.F64ValOrZero(raw[3])
		t.Period = convert.ToInt(raw[4])
	}

	return
}

func SnapshotFromRaw(pair string, raw [][]interface{}) (*Snapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice is too short for trade snapshot: %#v", raw)
	}

	snapshot := make([]*Trade, 0)
	for _, v := range raw {
		t, err := FromRaw(pair, v)
		if err != nil {
			return nil, err
		}
		snapshot = append(snapshot, t)
	}

	return &Snapshot{Snapshot: snapshot}, nil
}
