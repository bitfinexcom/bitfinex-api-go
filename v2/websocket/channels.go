package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingloan"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingtrade"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecutionupdate"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

type Heartbeat struct {
	//ChannelIDs []int64
}

func (c *Client) handleChannel(socketId SocketId, msg []byte) error {
	if c.terminal {
		return fmt.Errorf("received a message after close")
	}

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
	sub, err := c.subscriptions.lookupBySocketChannelID(chanID, socketId)
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
			case "cs":
				if checksum, ok := raw[2].(float64); ok {
					return c.handleChecksumChannel(sub, int(checksum))
				} else {
					c.log.Error("Unable to parse checksum")
				}
			default:
				body := raw[2].([]interface{})
				return c.handlePublicChannel(sub, sub.Request.Channel, data, body, msg)
			}
		case []interface{}:
			return c.handlePublicChannel(sub, sub.Request.Channel, "", data, msg)
		}
	} else {
		return c.handlePrivateChannel(raw)
	}
	return nil
}

func (c *Client) handleChecksumChannel(sub *subscription, checksum int) error {
	symbol := sub.Request.Symbol
	// force to signed integer
	bChecksum := uint32(checksum)
	var orderbook *Orderbook
	c.mtx.Lock()
	if ob, ok := c.orderbooks[symbol]; ok {
		orderbook = ob
	}
	c.mtx.Unlock()
	if orderbook != nil {
		oChecksum := orderbook.Checksum()
		// compare bitfinex checksum with local checksum
		if bChecksum == oChecksum {
			c.log.Debugf("Orderbook '%s' checksum verification successful.", symbol)
		} else {
			c.log.Warningf("Orderbook '%s' checksum is invalid got %d bot got %d. Data Out of sync, reconnecting.",
				symbol, bChecksum, oChecksum)
			err := c.sendUnsubscribeMessage(context.Background(), sub)
			if err != nil {
				return err
			}
			newSub := &SubscriptionRequest{
				SubID:   c.nonce.GetNonce(), // generate new subID
				Event:   sub.Request.Event,
				Channel: sub.Request.Channel,
				Symbol:  sub.Request.Symbol,
			}
			_, err_sub := c.Subscribe(context.Background(), newSub)
			if err_sub != nil {
				c.log.Warningf("could not resubscribe: %s", err_sub.Error())
				return err_sub
			}
		}
	}
	return nil
}

func (c *Client) handlePublicChannel(sub *subscription, channel, objType string, data []interface{}, raw_msg []byte) error {
	// unauthenticated data slice
	// public data is returned as raw interface arrays, use a factory to convert to raw type & publish
	if factory, ok := c.factories[channel]; ok {
		// convert to type array of interfaces
		if len(data) > 0 {
			if _, ok := data[0].([]interface{}); ok {
				interfaceArray := convert.ToInterfaceArray(data)
				// snapshot item
				c.mtx.Lock()
				// lock mutex since its mutates client struct
				msg, err := factory.BuildSnapshot(sub, interfaceArray, raw_msg)
				c.mtx.Unlock()
				if err != nil {
					return err
				}
				if msg != nil {
					c.listener <- msg
				}
			} else {
				// single item
				msg, err := factory.Build(sub, objType, data, raw_msg)
				if err != nil {
					return err
				}
				if msg != nil {
					c.listener <- msg
				}
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
	if val, ok := raw[1].(string); ok && val == "hb" {
		chanID, ok := raw[0].(float64)
		if !ok {
			c.log.Warningf("could not find chanID: %#v", raw)
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
		return &Heartbeat{}, nil
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
		o, err := position.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "pn":
		o, err := position.FromRaw(raw)
		if err != nil {
			return err
		}
		pn := position.New(*o)
		return &pn
	case "pu":
		o, err := position.FromRaw(raw)
		if err != nil {
			return err
		}
		pu := position.Update(*o)
		return &pu
	case "pc":
		o, err := position.FromRaw(raw)
		if err != nil {
			return err
		}
		pc := position.Cancel(*o)
		return &pc
	case "ws":
		o, err := wallet.SnapshotFromRaw(raw, wallet.FromWsRaw)
		if err != nil {
			return err
		}
		return o
	case "wu":
		o, err := wallet.FromWsRaw(raw)
		if err != nil {
			return err
		}
		wu := wallet.Update(*o)
		return &wu
	case "os":
		o, err := order.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "on":
		o, err := order.FromRaw(raw)
		if err != nil {
			return err
		}
		on := order.New(*o)
		return &on
	case "ou":
		o, err := order.FromRaw(raw)
		if err != nil {
			return err
		}
		ou := order.Update(*o)
		return &ou
	case "oc":
		o, err := order.FromRaw(raw)
		if err != nil {
			return err
		}
		oc := order.Cancel(*o)
		return &oc
	case "hts":
		tu, err := tradeexecutionupdate.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		hts := tradeexecutionupdate.HistoricalTradeSnapshot(*tu)
		return &hts
	case "te":
		o, err := bitfinex.NewTradeExecutionFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "tu":
		tu, err := tradeexecutionupdate.FromRaw(raw)
		if err != nil {
			return err
		}
		return tu
	case "fte":
		o, err := fundingtrade.FromRaw(raw)
		if err != nil {
			return err
		}
		fte := fundingtrade.Execution(*o)
		return &fte
	case "ftu":
		o, err := fundingtrade.FromRaw(raw)
		if err != nil {
			return err
		}
		ftu := fundingtrade.Update(*o)
		return &ftu
	case "hfts":
		fts, err := fundingtrade.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		nfts := fundingtrade.HistoricalSnapshot(*fts)
		return &nfts
	case "n":
		o, err := bitfinex.NewNotificationFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fos":
		o, err := fundingoffer.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fon":
		o, err := fundingoffer.FromRaw(raw)
		if err != nil {
			return err
		}
		fon := fundingoffer.New(*o)
		return &fon
	case "fou":
		o, err := fundingoffer.FromRaw(raw)
		if err != nil {
			return err
		}
		fou := fundingoffer.Update(*o)
		return &fou
	case "foc":
		o, err := fundingoffer.FromRaw(raw)
		if err != nil {
			return err
		}
		foc := fundingoffer.Cancel(*o)
		return &foc
	case "fiu":
		o, err := bitfinex.NewFundingInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcs":
		o, err := fundingcredit.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcn":
		o, err := fundingcredit.FromRaw(raw)
		if err != nil {
			return err
		}
		fcn := fundingcredit.New(*o)
		return &fcn
	case "fcu":
		o, err := fundingcredit.FromRaw(raw)
		if err != nil {
			return err
		}
		fcu := fundingcredit.Update(*o)
		return &fcu
	case "fcc":
		o, err := fundingcredit.FromRaw(raw)
		if err != nil {
			return err
		}
		fcc := fundingcredit.Cancel(*o)
		return &fcc
	case "fls":
		o, err := fundingloan.SnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fln":
		o, err := fundingloan.FromRaw(raw)
		if err != nil {
			return err
		}
		fln := fundingloan.New(*o)
		return &fln
	case "flu":
		o, err := fundingloan.FromRaw(raw)
		if err != nil {
			return err
		}
		flu := fundingloan.Update(*o)
		return &flu
	case "flc":
		o, err := fundingloan.FromRaw(raw)
		if err != nil {
			return err
		}
		flc := fundingloan.Cancel(*o)
		return &flc
	//case "uac":
	case "hb":
		return &Heartbeat{}
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
		c.log.Warningf("unhandled channel data, term: %s", term)
	}

	return fmt.Errorf("term %q not recognized", term)
}
