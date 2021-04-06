package position

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Position struct {
	Id                   int64
	Symbol               string
	Status               string
	Amount               float64
	BasePrice            float64
	MarginFunding        float64
	MarginFundingType    int64
	ProfitLoss           float64
	ProfitLossPercentage float64
	LiquidationPrice     float64
	Leverage             float64
	Flag                 interface{}
	MtsCreate            int64
	MtsUpdate            int64
	Type                 string
	Collateral           float64
	CollateralMin        float64
	Meta                 map[string]interface{}
}

type New Position
type Update Position
type Cancel Position

type Snapshot struct {
	Snapshot []*Position
}

func FromRaw(raw []interface{}) (p *Position, err error) {
	if len(raw) < 20 {
		return p, fmt.Errorf("data slice too short for position: %#v", raw)
	}

	p = &Position{
		Symbol:               convert.SValOrEmpty(raw[0]),
		Status:               convert.SValOrEmpty(raw[1]),
		Amount:               convert.F64ValOrZero(raw[2]),
		BasePrice:            convert.F64ValOrZero(raw[3]),
		MarginFunding:        convert.F64ValOrZero(raw[4]),
		MarginFundingType:    convert.I64ValOrZero(raw[5]),
		ProfitLoss:           convert.F64ValOrZero(raw[6]),
		ProfitLossPercentage: convert.F64ValOrZero(raw[7]),
		LiquidationPrice:     convert.F64ValOrZero(raw[8]),
		Leverage:             convert.F64ValOrZero(raw[9]),
		Id:                   convert.I64ValOrZero(raw[11]),
		MtsCreate:            convert.I64ValOrZero(raw[12]),
		MtsUpdate:            convert.I64ValOrZero(raw[13]),
		Type:                 convert.SValOrEmpty(raw[15]),
		Collateral:           convert.F64ValOrZero(raw[17]),
		CollateralMin:        convert.F64ValOrZero(raw[18]),
	}

	if meta, ok := raw[19].(map[string]interface{}); ok {
		p.Meta = meta
	}

	return
}

func NewFromRaw(raw []interface{}) (New, error) {
	p, err := FromRaw(raw)
	if err != nil {
		return New{}, err
	}
	p.Type = "pn"
	return New(*p), nil
}

func UpdateFromRaw(raw []interface{}) (Update, error) {
	p, err := FromRaw(raw)
	if err != nil {
		return Update{}, err
	}
	p.Type = "pu"
	return Update(*p), nil
}

func CancelFromRaw(raw []interface{}) (Cancel, error) {
	p, err := FromRaw(raw)
	if err != nil {
		return Cancel{}, err
	}
	p.Type = "pc"
	return Cancel(*p), nil
}

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
