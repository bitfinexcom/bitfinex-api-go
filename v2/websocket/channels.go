package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

func (c *Client) handleChannel(msg []byte) error {
	var raw []interface{}
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		return err
	} else if len(raw) < 2 {
		return nil
	}

	chID, ok := raw[0].(float64)
	if !ok {
		return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
	}

	chanID := int64(chID)
	sub, err := c.subscriptions.lookupByChannelID(chanID)
	if err != nil {
		// no subscribed channel for message
		return err
	}
	c.subscriptions.heartbeat(chanID)
	if sub.Public {
		switch data := raw[1].(type) {
		case string:
			switch data {
			case "hb":
				// no-op, already updated heartbeat timeout from this event
				return nil
			default:
				body := raw[2].([]interface{})
				return c.handlePublicChannel(chanID, sub.Request.Channel, data, body)
			}
		case []interface{}:
			return c.handlePublicChannel(chanID, sub.Request.Channel, "", data)
		}
	} else {
		return c.handlePrivateChannel(raw)
	}
	return nil
}

func (c *Client) handlePublicChannel(chanID int64, channel, objType string, data []interface{}) error {
	// unauthenticated data slice
	// returns interface{} (which is really [][]float64)
	obj, err := c.processDataSlice(data)
	if err != nil {
		return err
	}
	// public data is returned as raw interface arrays, use a factory to convert to raw type & publish
	if factory, ok := c.factories[channel]; ok {
		flt := obj.([][]float64)
		if len(flt) == 1 {
			// single item
			arr := make([]interface{}, len(flt[0]))
			for i, ft := range flt[0] {
				arr[i] = ft
			}
			msg, err := factory.Build(chanID, objType, arr)
			if err != nil {
				return err
			}
			if msg != nil {
				c.listener <- msg
			}
		} else if len(flt) > 1 {
			msg, err := factory.BuildSnapshot(chanID, flt)
			if err != nil {
				return err
			}
			if msg != nil {
				c.listener <- msg
			}
		}
	} else {
		// factory lookup error
		return fmt.Errorf("could not find public factory for %s channel", channel)
	}
	return nil
}

func (c *Client) handlePrivateChannel(raw []interface{}) error {
	// authenticated data slice, or a heartbeat
	if raw[1].(string) == "hb" {
		chanID, ok := raw[0].(float64)
		if !ok {
			log.Printf("could not find chanID: %#v", raw)
			return nil
		}
		c.handleHeartbeat(int64(chanID))
	} else {
		// raw[2] is data slice
		// authenticated snapshots?
		if len(raw) > 2 {
			if arr, ok := raw[2].([]interface{}); ok {
				obj, err := c.handlePrivateDataMessage(raw[1].(string), arr)
				if err != nil {
					return err
				}
				// private data is returned as strongly typed data, publish directly
				if obj != nil {
					c.listener <- obj
				}
			}
		}
	}
	return nil
}

func (c *Client) handleHeartbeat(chanID int64) {
	c.subscriptions.heartbeat(chanID)
}

type unsubscribeMsg struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}

func (c *Client) processDataSlice(data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("unexpected data slice: %v", data)
	}

	var items [][]float64
	switch data[0].(type) {
	case []interface{}: // [][]float64
		for _, e := range data {
			if s, ok := e.([]interface{}); ok {
				item, err := bitfinex.F64Slice(s)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			} else {
				return nil, fmt.Errorf("expected slice of float64 slices but got: %v", data)
			}
		}
	case float64: // []float64
		item, err := bitfinex.F64Slice(data)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	default:
		return nil, fmt.Errorf("unexpected data slice: %v", data)
	}

	return items, nil
}

// public msg: [ChanID, [Data]]
// hb (both): [ChanID, "hb"]
// private update msg: [ChanID, "type", [Data]]
// private snapshot msg: [ChanID, "type", [[Data]]]
func (c *Client) handlePrivateDataMessage(term string, data []interface{}) (ms interface{}, err error) {
	if len(data) == 0 {
		// empty data msg
		return nil, nil
	}

	if term == "hb" { // Heartbeat
		// TODO: Consider adding a switch to enable/disable passing these along.
		return &bitfinex.Heartbeat{}, nil
	}
	/*
		list, ok := data[2].([]interface{})
		if !ok {
			return ms, fmt.Errorf("expected data list in third position but got %#v in %#v", data[2], data)
		}
	*/
	ms = c.convertRaw(term, data)

	return
}

// convertRaw takes a term and the raw data attached to it to try and convert that
// untyped list into a proper type.
func (c *Client) convertRaw(term string, raw []interface{}) interface{} {
	// The things you do to get proper types.
	switch term {
	case "bu":
		o, err := bitfinex.NewBalanceInfoFromRaw(raw)
		if err != nil {
			return err
		}
		bu := bitfinex.BalanceUpdate(*o)
		return &bu
	case "ps":
		o, err := bitfinex.NewPositionSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "pn":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pn := bitfinex.PositionNew(*o)
		return &pn
	case "pu":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pu := bitfinex.PositionUpdate(*o)
		return &pu
	case "pc":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pc := bitfinex.PositionCancel(*o)
		return &pc
	case "ws":
		o, err := bitfinex.NewWalletSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "wu":
		o, err := bitfinex.NewWalletFromRaw(raw)
		if err != nil {
			return err
		}
		wu := bitfinex.WalletUpdate(*o)
		return &wu
	case "os":
		o, err := bitfinex.NewOrderSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "on":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		on := bitfinex.OrderNew(*o)
		return &on
	case "ou":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		ou := bitfinex.OrderUpdate(*o)
		return &ou
	case "oc":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		oc := bitfinex.OrderCancel(*o)
		return &oc
	case "hts":
		o, err := bitfinex.NewTradeExecutionUpdateSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		hts := bitfinex.HistoricalTradeSnapshot(*o)
		return &hts
	case "te":
		o, err := bitfinex.NewTradeExecutionFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "tu":
		tu, err := bitfinex.NewTradeExecutionUpdateFromRaw(raw)
		if err != nil {
			return err
		}
		return tu
	case "fte":
		o, err := bitfinex.NewFundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		fte := bitfinex.FundingTradeExecution(*o)
		return &fte
	case "ftu":
		o, err := bitfinex.NewFundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		ftu := bitfinex.FundingTradeUpdate(*o)
		return &ftu
	case "hfts":
		o, err := bitfinex.NewFundingTradeSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		nfts := bitfinex.HistoricalFundingTradeSnapshot(*o)
		return &nfts
	case "n":
		o, err := bitfinex.NewNotificationFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fos":
		o, err := bitfinex.NewFundingOfferSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fon":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		fon := bitfinex.FundingOfferNew(*o)
		return &fon
	case "fou":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		fou := bitfinex.FundingOfferUpdate(*o)
		return &fou
	case "foc":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		foc := bitfinex.FundingOfferCancel(*o)
		return &foc
	case "fiu":
		o, err := bitfinex.NewFundingInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcs":
		o, err := bitfinex.NewFundingCreditSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcn":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcn := bitfinex.FundingCreditNew(*o)
		return &fcn
	case "fcu":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcu := bitfinex.FundingCreditUpdate(*o)
		return &fcu
	case "fcc":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcc := bitfinex.FundingCreditCancel(*o)
		return &fcc
	case "fls":
		o, err := bitfinex.NewFundingLoanSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fln":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		fln := bitfinex.FundingLoanNew(*o)
		return &fln
	case "flu":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		flu := bitfinex.FundingLoanUpdate(*o)
		return &flu
	case "flc":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		flc := bitfinex.FundingLoanCancel(*o)
		return &flc
		//case "uac":
	case "hb":
		return &bitfinex.Heartbeat{}
	case "ats":
		// TODO: Is not in documentation, so figure out what it is.
		return nil
	case "oc-req":
		// TODO
		return nil
	case "on-req":
		// TODO
		return nil
	case "mis": // Should not be sent anymore as of 2017-04-01
		return nil
	case "miu":
		o, err := bitfinex.NewMarginInfoFromRaw(raw)
		if err != nil {
			return err
		}
		// return a strongly typed reference, rather than dereference a generic interface
		// too bad golang doesn't inherit an interface's underlying type when creating a reference to the interface
		if base, ok := o.(*bitfinex.MarginInfoBase); ok {
			return base
		}
		if update, ok := o.(*bitfinex.MarginInfoUpdate); ok {
			return update
		}
		return o // better than nothing
	default:
		log.Printf("unhandled channel data, term: %s", term)
	}

	return fmt.Errorf("term %q not recognized", term)
}
