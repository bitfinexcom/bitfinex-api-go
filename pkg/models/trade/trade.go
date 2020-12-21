package trade

import (
	"errors"
	"fmt"
	"strings"

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
	if strings.HasPrefix(pair, "t") && len(raw) >= 4 {
		t = &Trade{
			Pair:   pair,
			ID:     convert.I64ValOrZero(raw[0]),
			MTS:    convert.I64ValOrZero(raw[1]),
			Amount: convert.F64ValOrZero(raw[2]),
			Price:  convert.F64ValOrZero(raw[3]),
		}
		return
	}

	if strings.HasPrefix(pair, "f") && len(raw) >= 5 {
		t = &Trade{
			Pair:   pair,
			ID:     convert.I64ValOrZero(raw[0]),
			MTS:    convert.I64ValOrZero(raw[1]),
			Amount: convert.F64ValOrZero(raw[2]),
			Rate:   convert.F64ValOrZero(raw[3]),
			Period: convert.ToInt(raw[4]),
		}
		return
	}

	err = fmt.Errorf("data slice too short for %s pair: %#v", pair, raw)
	return
}

func SnapshotFromRaw(pair string, raw [][]interface{}) (*Snapshot, error) {
	if len(raw) == 0 {
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

// FromWSRaw - based on condition will return snapshot of trades or single trade
func FromWSRaw(pair string, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice")
	}

	_, isSnapshot := data[0].([]interface{})
	if isSnapshot {
		return SnapshotFromRaw(pair, convert.ToInterfaceArray(data))
	}

	return FromRaw(pair, data)
}
