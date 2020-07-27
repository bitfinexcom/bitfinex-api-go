package position

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Status string

type Position struct {
	Id                   int64
	Symbol               string
	Status               Status
	Amount               float64
	BasePrice            float64
	MarginFunding        float64
	MarginFundingType    int64
	ProfitLoss           float64
	ProfitLossPercentage float64
	LiquidationPrice     float64
	Leverage             float64
}

func FromRaw(raw []interface{}) (o *Position, err error) {
	if len(raw) < 6 {
		return o, fmt.Errorf("data slice too short for position: %#v", raw)
	}

	o = &Position{
		Symbol:            convert.SValOrEmpty(raw[0]),
		Status:            Status(convert.SValOrEmpty(raw[1])),
		Amount:            convert.F64ValOrZero(raw[2]),
		BasePrice:         convert.F64ValOrZero(raw[3]),
		MarginFunding:     convert.F64ValOrZero(raw[4]),
		MarginFundingType: convert.I64ValOrZero(raw[5]),
	}

	if len(raw) == 10 {
		o.ProfitLoss = convert.F64ValOrZero(raw[6])
		o.ProfitLossPercentage = convert.F64ValOrZero(raw[7])
		o.LiquidationPrice = convert.F64ValOrZero(raw[8])
		o.Leverage = convert.F64ValOrZero(raw[9])
		return
	}

	if len(raw) > 10 {
		o.ProfitLoss = convert.F64ValOrZero(raw[6])
		o.ProfitLossPercentage = convert.F64ValOrZero(raw[7])
		o.LiquidationPrice = convert.F64ValOrZero(raw[8])
		o.Leverage = convert.F64ValOrZero(raw[9])
		o.Id = int64(convert.F64ValOrZero(raw[11]))
	}

	return
}

type Snapshot struct {
	Snapshot []*Position
}
type New Position
type Update Position
type Cancel Position

func SnapshotFromRaw(raw []interface{}) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for position: %#v", raw)
	}

	ps := make([]*Position, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				p, err := FromRaw(l)
				if err != nil {
					return s, err
				}
				ps = append(ps, p)
			}
		}
	default:
		return s, fmt.Errorf("not a position snapshot")
	}
	s = &Snapshot{Snapshot: ps}

	return
}

type ClaimRequest struct {
	Id int64
}

func (o *ClaimRequest) ToJSON() ([]byte, error) {
	aux := struct {
		Id int64 `json:"id"`
	}{
		Id: o.Id,
	}
	return json.Marshal(aux)
}
