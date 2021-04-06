package order

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Order struct {
	ID            int64
	GID           int64
	CID           int64
	Symbol        string
	MTSCreated    int64
	MTSUpdated    int64
	Amount        float64
	AmountOrig    float64
	Type          string
	TypePrev      string
	MTSTif        int64
	Flags         int64
	Status        string
	Price         float64
	PriceAvg      float64
	PriceTrailing float64
	PriceAuxLimit float64
	Notify        bool
	Hidden        bool
	PlacedID      int64
	Routing       string
	Meta          map[string]interface{}
}

// Snapshot is a collection of Orders that would usually be sent on
// inital connection.
type Snapshot struct {
	Snapshot []*Order
}

// Update is an Order that gets sent out after every change to an order.
type Update Order

// New gets sent out after an Order was created successfully.
type New Order

// Cancel gets sent out after an Order was cancelled successfully.
type Cancel Order

// FromRaw takes the raw list of values as returned from the websocket
// service and tries to convert it into an Order.
func FromRaw(raw []interface{}) (o *Order, err error) {
	if len(raw) < 32 {
		return o, fmt.Errorf("data slice too short for order: %#v", raw)
	}

	o = &Order{
		ID:            convert.I64ValOrZero(raw[0]),
		GID:           convert.I64ValOrZero(raw[1]),
		CID:           convert.I64ValOrZero(raw[2]),
		Symbol:        convert.SValOrEmpty(raw[3]),
		MTSCreated:    convert.I64ValOrZero(raw[4]),
		MTSUpdated:    convert.I64ValOrZero(raw[5]),
		Amount:        convert.F64ValOrZero(raw[6]),
		AmountOrig:    convert.F64ValOrZero(raw[7]),
		Type:          convert.SValOrEmpty(raw[8]),
		TypePrev:      convert.SValOrEmpty(raw[9]),
		MTSTif:        convert.I64ValOrZero(raw[10]),
		Flags:         convert.I64ValOrZero(raw[12]),
		Status:        convert.SValOrEmpty(raw[13]),
		Price:         convert.F64ValOrZero(raw[16]),
		PriceAvg:      convert.F64ValOrZero(raw[17]),
		PriceTrailing: convert.F64ValOrZero(raw[18]),
		PriceAuxLimit: convert.F64ValOrZero(raw[19]),
		Notify:        convert.BValOrFalse(raw[23]),
		Hidden:        convert.BValOrFalse(raw[24]),
		PlacedID:      convert.I64ValOrZero(raw[25]),
		Routing:       convert.SValOrEmpty(raw[28]),
	}

	if meta, ok := raw[31].(map[string]interface{}); ok {
		o.Meta = meta
	}

	return
}

// NewFromRaw reds "on" type message from data stream and
// maps it to order.New data structure
func NewFromRaw(raw []interface{}) (New, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return New{}, err
	}
	return New(*o), nil
}

// UpdateFromRaw reds "ou" type message from data stream and
// maps it to order.Update data structure
func UpdateFromRaw(raw []interface{}) (Update, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return Update{}, err
	}
	return Update(*o), nil
}

// CancelFromRaw reds "oc" type message from data stream and
// maps it to order.Cancel data structure
func CancelFromRaw(raw []interface{}) (Cancel, error) {
	o, err := FromRaw(raw)
	if err != nil {
		return Cancel{}, err
	}
	return Cancel(*o), nil
}

// SnapshotFromRaw takes a raw list of values as returned from the websocket
// service and tries to convert it into an Snapshot.
func SnapshotFromRaw(raw []interface{}) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for order: %#v", raw)
	}

	os := make([]*Order, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := FromRaw(l)
				if err != nil {
					return s, err
				}
				os = append(os, o)
			}
		}
	default:
		return s, fmt.Errorf("not an order snapshot")
	}
	s = &Snapshot{Snapshot: os}

	return
}
