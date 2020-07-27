package stats

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Stat struct {
	Period int64
	Volume float64
}

func FromRaw(raw []interface{}) (*Stat, error) {
	if len(raw) < 2 {
		return nil, fmt.Errorf("data slice too short (len=%d) for Stat: %#v", len(raw), raw)
	}

	return &Stat{
		Period: convert.I64ValOrZero(raw[0]),
		Volume: convert.F64ValOrZero(raw[1]),
	}, nil
}

func SnapshotFromRaw(raw []interface{}) (snap []*Stat, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for stats: %#v", raw)
	}

	stats := make([]*Stat, 0)
	for _, v := range raw {
		if v, ok := v.([]interface{}); ok {
			s, err := FromRaw(v)
			if err != nil {
				return snap, err
			}
			stats = append(stats, s)
		} else {
			return nil, fmt.Errorf("Invalid stats snapshot")
		}
	}

	return stats, nil
}
