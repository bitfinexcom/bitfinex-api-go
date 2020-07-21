package fundingoffer

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

type Offer struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	MTSUpdated int64
	Amount     float64
	AmountOrig float64
	Type       string
	Flags      interface{}
	Status     string
	Rate       float64
	Period     int64
	Notify     bool
	Hidden     bool
	Insure     bool
	Renew      bool
	RateReal   float64
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

// MarshalJSON converts the offer cancel object into the format required by the
// bitfinex websocket service.
func (cr *CancelRequest) MarshalJSON() ([]byte, error) {
	b, err := cr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"foc\", null, %s]", string(b))), nil
}

type SubmitRequest struct {
	Type   string
	Symbol string
	Amount float64
	Rate   float64
	Period int64
	Hidden bool
}

func (sr *SubmitRequest) ToJSON() ([]byte, error) {
	aux := struct {
		Type   string  `json:"type"`
		Symbol string  `json:"symbol"`
		Amount float64 `json:"amount,string"`
		Rate   float64 `json:"rate,string"`
		Period int64   `json:"period"`
		Flags  int     `json:"flags,omitempty"`
	}{
		Type:   sr.Type,
		Symbol: sr.Symbol,
		Amount: sr.Amount,
		Rate:   sr.Rate,
		Period: sr.Period,
	}
	if sr.Hidden {
		aux.Flags = aux.Flags + common.OrderFlagHidden
	}
	return json.Marshal(aux)
}

// MarshalJSON converts the offer submit object into the format required by the
// bitfinex websocket service.
func (sr *SubmitRequest) MarshalJSON() ([]byte, error) {
	aux, err := sr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"fon\", null, %s]", string(aux))), nil
}

func FromRaw(raw []interface{}) (o *Offer, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short for offer: %#v", raw)
	}

	o = &Offer{
		ID:         convert.I64ValOrZero(raw[0]),
		Symbol:     convert.SValOrEmpty(raw[1]),
		MTSCreated: convert.I64ValOrZero(raw[2]),
		MTSUpdated: convert.I64ValOrZero(raw[3]),
		Amount:     convert.F64ValOrZero(raw[4]),
		AmountOrig: convert.F64ValOrZero(raw[5]),
		Type:       convert.SValOrEmpty(raw[6]),
		Flags:      raw[9],
		Status:     convert.SValOrEmpty(raw[10]),
		Rate:       convert.F64ValOrZero(raw[14]),
		Period:     convert.I64ValOrZero(raw[15]),
		Notify:     convert.BValOrFalse(raw[16]),
		Hidden:     convert.BValOrFalse(raw[17]),
		Insure:     convert.BValOrFalse(raw[18]),
		Renew:      convert.BValOrFalse(raw[19]),
		RateReal:   convert.F64ValOrZero(raw[20]),
	}

	return
}

type New Offer
type Update Offer
type Cancel Offer
type Snapshot struct {
	Snapshot []*Offer
}

func SnapshotFromRaw(raw []interface{}) (snap *Snapshot, err error) {
	if len(raw) == 0 {
		return snap, fmt.Errorf("data slice too short for offer: %#v", raw)
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
