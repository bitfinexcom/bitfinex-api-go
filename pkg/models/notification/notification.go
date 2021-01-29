package notification

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
)

type Notification struct {
	MTS        int64
	Type       string
	MessageID  int64
	NotifyInfo interface{}
	Code       int64
	Status     string
	Text       string
}

func FromRaw(raw []interface{}) (n *Notification, err error) {
	if len(raw) < 8 {
		return n, fmt.Errorf("data slice too short for notification: %#v", raw)
	}

	n = &Notification{
		MTS:       convert.I64ValOrZero(raw[0]),
		Type:      convert.SValOrEmpty(raw[1]),
		MessageID: convert.I64ValOrZero(raw[2]),
		Code:      convert.I64ValOrZero(raw[5]),
		Status:    convert.SValOrEmpty(raw[6]),
		Text:      convert.SValOrEmpty(raw[7]),
	}

	// raw[4] = notify info
	if raw[4] == nil {
		return
	}

	nraw := raw[4].([]interface{})
	if len(nraw) == 0 {
		return
	}

	switch n.Type {
	case "on-req":
		// will be a set of orders if created via rest
		// this is to accommodate OCO orders
		if _, isSnapshot := nraw[0].([]interface{}); isSnapshot {
			n.NotifyInfo, err = order.SnapshotFromRaw(nraw)
			return
		}

		n.NotifyInfo, err = order.NewFromRaw(nraw)
		return
	case "ou-req", "ou":
		n.NotifyInfo, err = order.UpdateFromRaw(nraw)
		return
	case "oc-req":
		n.NotifyInfo, err = order.CancelFromRaw(nraw)
		return
	case "fon-req":
		n.NotifyInfo, err = fundingoffer.NewFromRaw(nraw)
		return
	case "foc-req":
		n.NotifyInfo, err = fundingoffer.CancelFromRaw(nraw)
		return
	case "pm-req", "pc":
		n.NotifyInfo, err = position.CancelFromRaw(nraw)
		return
	default:
		n.NotifyInfo = raw[4]
	}

	return
}
