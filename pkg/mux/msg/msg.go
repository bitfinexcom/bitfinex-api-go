package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
)

type Msg struct {
	Data     []byte
	Err      error
	CID      int
	IsPublic bool
}

func (m Msg) IsEvent() bool {
	t := bytes.TrimLeftFunc(m.Data, unicode.IsSpace)
	return bytes.HasPrefix(t, []byte("{"))
}

func (m Msg) IsRaw() bool {
	t := bytes.TrimLeftFunc(m.Data, unicode.IsSpace)
	return bytes.HasPrefix(t, []byte("["))
}

func (m Msg) ProcessRaw(chanInfo map[int64]event.Info) (interface{}, error) {
	var raw []interface{}
	if err := json.Unmarshal(m.Data, &raw); err != nil {
		return nil, fmt.Errorf("parsing msg: %s, err: %s", m.Data, err)
	}
	// payload data is always last element of the slice
	pld := raw[len(raw)-1]
	// chanID is always 1st element of the slice
	chID := convert.I64ValOrZero(raw[0])
	// allocate channel name by id to know how to transform raw data
	inf, ok := chanInfo[chID]
	if !ok {
		return nil, fmt.Errorf("unrecognized chanId:%d", chID)
	}

	switch data := pld.(type) {
	case string:
		return event.Info{
			ChanID:    chID,
			Subscribe: event.Subscribe{Event: data},
		}, nil
	case []interface{}:
		switch inf.Channel {
		case "trades":
			return trades.FromWSRaw(inf.Symbol, raw, data)
		case "ticker":
			return ticker.FromWSRaw(inf.Symbol, data)
		case "book":
			return book.FromWSRaw(inf.Symbol, inf.Precision, data)
		case "candles":
			return candle.FromWSRaw(inf.Key, data)
		case "status":
			return status.FromWSRaw(inf.Key, data)
		}
	}

	return raw, nil
}

func (m Msg) ProcessPrivateRaw() (interface{}, error) {
	var raw []interface{}
	if err := json.Unmarshal(m.Data, &raw); err != nil {
		return nil, fmt.Errorf("parsing auth msg: %s, err: %s", m.Data, err)
	}
	// payload data is always last element of the slice
	pld := raw[len(raw)-1]
	// op name is 2nd element
	op := convert.SValOrEmpty(raw[1])

	switch data := pld.(type) {
	case string:
		return event.Info{
			ChanID:    convert.I64ValOrZero(raw[0]),
			Subscribe: event.Subscribe{Event: data},
		}, nil
	case []interface{}:
		switch op {
		case "os":
			return order.SnapshotFromRaw(data)
		case "on":
			return order.NewFromRaw(data)
		case "ou":
			return order.UpdateFromRaw(data)
		case "oc":
			return order.CancelFromRaw(data)
		case "ws":
			return wallet.SnapshotFromRaw(data, wallet.FromWsRaw)
		case "wu":
			return wallet.UpdateFromWsRaw(data)
		}
	}

	return raw, nil
}

func (m Msg) ProcessEvent() (i event.Info, err error) {
	if err = json.Unmarshal(m.Data, &i); err != nil {
		return i, fmt.Errorf("parsing msg: %s, err: %s", m.Data, err)
	}
	return
}
