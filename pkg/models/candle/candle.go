package candle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

type Candle struct {
	Symbol     string
	Resolution common.CandleResolution
	MTS        int64
	Open       float64
	Close      float64
	High       float64
	Low        float64
	Volume     float64
}

type Snapshot struct {
	Snapshot []*Candle
}

func SnapshotFromRaw(symbol string, resolution common.CandleResolution, raw [][]interface{}) (*Snapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for candle snapshot: %#v", raw)
	}

	snap := make([]*Candle, 0)
	for _, f := range raw {
		c, err := FromRaw(symbol, resolution, f)
		if err == nil {
			snap = append(snap, c)
		}
	}

	return &Snapshot{Snapshot: snap}, nil
}

func FromRaw(symbol string, resolution common.CandleResolution, raw []interface{}) (c *Candle, err error) {
	if len(raw) < 6 {
		return c, fmt.Errorf("data slice too short for candle, expected %d got %d: %#v", 6, len(raw), raw)
	}

	c = &Candle{
		Symbol:     symbol,
		Resolution: resolution,
		MTS:        convert.I64ValOrZero(raw[0]),
		Open:       convert.F64ValOrZero(raw[1]),
		Close:      convert.F64ValOrZero(raw[2]),
		High:       convert.F64ValOrZero(raw[3]),
		Low:        convert.F64ValOrZero(raw[4]),
		Volume:     convert.F64ValOrZero(raw[5]),
	}

	return
}

// FromWSRaw - based on condition will return snapshot of candles or single candle
func FromWSRaw(key string, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice")
	}

	_, isSnapshot := data[0].([]interface{})
	ss := strings.Split(key, ":")

	// basic signature is: "trade:1m:tBTCUSD"
	if len(ss) < 3 {
		return nil, fmt.Errorf("invalid key:%s", key)
	}

	res := common.CandleResolution(ss[1])
	symbol := ss[2]

	if isSnapshot {
		return SnapshotFromRaw(symbol, res, convert.ToInterfaceArray(data))
	}
	return FromRaw(symbol, res, data)
}
