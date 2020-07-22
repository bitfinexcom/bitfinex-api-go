package bitfinex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
)

// Prefixes for available pairs
const (
	FundingPrefix = "f"
	TradingPrefix = "t"
)

var (
	ErrNotFound = errors.New("not found")
)

// Candle resolutions
const (
	OneMinute      CandleResolution = "1m"
	FiveMinutes    CandleResolution = "5m"
	FifteenMinutes CandleResolution = "15m"
	ThirtyMinutes  CandleResolution = "30m"
	OneHour        CandleResolution = "1h"
	ThreeHours     CandleResolution = "3h"
	SixHours       CandleResolution = "6h"
	TwelveHours    CandleResolution = "12h"
	OneDay         CandleResolution = "1D"
	OneWeek        CandleResolution = "7D"
	TwoWeeks       CandleResolution = "14D"
	OneMonth       CandleResolution = "1M"
)

type Mts int64
type SortOrder int

const (
	OldestFirst SortOrder = 1
	NewestFirst SortOrder = -1
)

type QueryLimit int

const QueryLimitMax QueryLimit = 1000

func CandleResolutionFromString(str string) (CandleResolution, error) {
	switch str {
	case string(OneMinute):
		return OneMinute, nil
	case string(FiveMinutes):
		return FiveMinutes, nil
	case string(FifteenMinutes):
		return FifteenMinutes, nil
	case string(ThirtyMinutes):
		return ThirtyMinutes, nil
	case string(OneHour):
		return OneHour, nil
	case string(ThreeHours):
		return ThreeHours, nil
	case string(SixHours):
		return SixHours, nil
	case string(TwelveHours):
		return TwelveHours, nil
	case string(OneDay):
		return OneDay, nil
	case string(OneWeek):
		return OneWeek, nil
	case string(TwoWeeks):
		return TwoWeeks, nil
	case string(OneMonth):
		return OneMonth, nil
	}
	return OneMinute, fmt.Errorf("could not convert string to resolution: %s", str)
}

type PermissionType string

const (
	PermissionRead  = "r"
	PermissionWrite = "w"
)

// private type--cannot instantiate.
type candleResolution string

// CandleResolution provides a typed set of resolutions for candle subscriptions.
type CandleResolution candleResolution

// Order sides
const (
	Bid   common.OrderSide = 1
	Ask   common.OrderSide = 2
	Long  common.OrderSide = 1
	Short common.OrderSide = 2
)

// Settings flags

const (
	Dec_s     int = 9
	Time_s    int = 32
	Timestamp int = 32768
	Seq_all   int = 65536
	Checksum  int = 131072
)

// Book precision levels
const (
	// Aggregate precision levels
	Precision0 BookPrecision = "P0"
	Precision2 BookPrecision = "P2"
	Precision1 BookPrecision = "P1"
	Precision3 BookPrecision = "P3"
	// Raw precision
	PrecisionRawBook BookPrecision = "R0"
)

// private type
type bookPrecision string

// BookPrecision provides a typed book precision level.
type BookPrecision bookPrecision

const (
	// FrequencyRealtime book frequency gives updates as they occur in real-time.
	FrequencyRealtime BookFrequency = "F0"
	// FrequencyTwoPerSecond delivers two book updates per second.
	FrequencyTwoPerSecond BookFrequency = "F1"
	// PriceLevelDefault provides a constant default price level for book subscriptions.
	PriceLevelDefault int = 25
)

type bookFrequency string

// BookFrequency provides a typed book frequency.
type BookFrequency bookFrequency

// TradeExecution represents the first message receievd for a trade on the private data feed.
type TradeExecution struct {
	ID         int64
	Pair       string
	MTS        int64
	OrderID    int64
	Amount     float64
	Price      float64
	OrderType  string
	OrderPrice float64
	Maker      int
}

func NewTradeExecutionFromRaw(raw []interface{}) (o *TradeExecution, err error) {
	if len(raw) < 6 {
		log.Printf("[ERROR] not enough members (%d, need at least 6) for trade execution: %#v", len(raw), raw)
		return o, fmt.Errorf("data slice too short for trade execution: %#v", raw)
	}

	// trade executions sometimes omit order type, price, and maker flag
	o = &TradeExecution{
		ID:      convert.I64ValOrZero(raw[0]),
		Pair:    convert.SValOrEmpty(raw[1]),
		MTS:     convert.I64ValOrZero(raw[2]),
		OrderID: convert.I64ValOrZero(raw[3]),
		Amount:  convert.F64ValOrZero(raw[4]),
		Price:   convert.F64ValOrZero(raw[5]),
	}

	if len(raw) >= 9 {
		o.OrderType = convert.SValOrEmpty(raw[6])
		o.OrderPrice = convert.F64ValOrZero(raw[7])
		o.Maker = convert.IValOrZero(raw[8])
	}

	return
}

type BalanceInfo struct {
	TotalAUM float64
	NetAUM   float64
	/*WalletType string
	Currency   string*/
}

func NewBalanceInfoFromRaw(raw []interface{}) (o *BalanceInfo, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for balance info: %#v", raw)
	}

	o = &BalanceInfo{
		TotalAUM: convert.F64ValOrZero(raw[0]),
		NetAUM:   convert.F64ValOrZero(raw[1]),
		/*WalletType: convert.SValOrEmpty(raw[2]),
		Currency:   convert.SValOrEmpty(raw[3]),*/
	}

	return
}

type BalanceUpdate BalanceInfo

// marginInfoFromRaw returns either a MarginInfoBase or MarginInfoUpdate, since
// the Margin Info is split up into a base and per symbol parts.
func NewMarginInfoFromRaw(raw []interface{}) (o interface{}, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	typ, ok := raw[0].(string)
	if !ok {
		return o, fmt.Errorf("expected margin info type in first position for margin info but got %#v", raw)
	}

	if len(raw) == 2 && typ == "base" { // This should be ["base", [...]]
		data, ok := raw[1].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in second position for margin info but got %#v", raw)
		}

		return NewMarginInfoBaseFromRaw(data)
	} else if len(raw) == 3 && typ == "sym" { // This should be ["sym", SYMBOL, [...]]
		symbol, ok := raw[1].(string)
		if !ok {
			return o, fmt.Errorf("expected margin info symbol in second position for margin info update but got %#v", raw)
		}

		data, ok := raw[2].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in third position for margin info update but got %#v", raw)
		}

		return NewMarginInfoUpdateFromRaw(symbol, data)
	}

	return nil, fmt.Errorf("invalid margin info type in %#v", raw)
}

type MarginInfoUpdate struct {
	Symbol          string
	TradableBalance float64
}

func NewMarginInfoUpdateFromRaw(symbol string, raw []interface{}) (o *MarginInfoUpdate, err error) {
	if len(raw) < 1 {
		return o, fmt.Errorf("data slice too short for margin info update: %#v", raw)
	}

	o = &MarginInfoUpdate{
		Symbol:          symbol,
		TradableBalance: convert.F64ValOrZero(raw[0]),
	}

	return
}

type MarginInfoBase struct {
	UserProfitLoss float64
	UserSwaps      float64
	MarginBalance  float64
	MarginNet      float64
}

func NewMarginInfoBaseFromRaw(raw []interface{}) (o *MarginInfoBase, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	o = &MarginInfoBase{
		UserProfitLoss: convert.F64ValOrZero(raw[0]),
		UserSwaps:      convert.F64ValOrZero(raw[1]),
		MarginBalance:  convert.F64ValOrZero(raw[2]),
		MarginNet:      convert.F64ValOrZero(raw[3]),
	}

	return
}

type FundingInfo struct {
	Symbol       string
	YieldLoan    float64
	YieldLend    float64
	DurationLoan float64
	DurationLend float64
}

func NewFundingInfoFromRaw(raw []interface{}) (o *FundingInfo, err error) {
	if len(raw) < 3 { // "sym", symbol, data
		return o, fmt.Errorf("data slice too short for funding info: %#v", raw)
	}

	sym, ok := raw[1].(string)
	if !ok {
		return o, fmt.Errorf("expected symbol in second position of funding info: %v", raw)
	}

	data, ok := raw[2].([]interface{})
	if !ok {
		return o, fmt.Errorf("expected list in third position of funding info: %v", raw)
	}

	if len(data) < 4 {
		return o, fmt.Errorf("data too short: %#v", data)
	}

	o = &FundingInfo{
		Symbol:       sym,
		YieldLoan:    convert.F64ValOrZero(data[0]),
		YieldLend:    convert.F64ValOrZero(data[1]),
		DurationLoan: convert.F64ValOrZero(data[2]),
		DurationLend: convert.F64ValOrZero(data[3]),
	}

	return
}

type Notification struct {
	MTS        int64
	Type       string
	MessageID  int64
	NotifyInfo interface{}
	Code       int64
	Status     string
	Text       string
}

func NewNotificationFromRaw(raw []interface{}) (o *Notification, err error) {
	if len(raw) < 8 {
		return o, fmt.Errorf("data slice too short for notification: %#v", raw)
	}

	o = &Notification{
		MTS:       convert.I64ValOrZero(raw[0]),
		Type:      convert.SValOrEmpty(raw[1]),
		MessageID: convert.I64ValOrZero(raw[2]),
		//NotifyInfo: raw[4],
		Code:   convert.I64ValOrZero(raw[5]),
		Status: convert.SValOrEmpty(raw[6]),
		Text:   convert.SValOrEmpty(raw[7]),
	}

	// raw[4] = notify info
	var nraw []interface{}
	if raw[4] != nil {
		nraw = raw[4].([]interface{})
		switch o.Type {
		case "on-req":
			if len(nraw) <= 0 {
				o.NotifyInfo = nil
				break
			}
			// will be a set of orders if created via rest
			// this is to accommodate OCO orders
			if _, ok := nraw[0].([]interface{}); ok {
				o.NotifyInfo, err = order.SnapshotFromRaw(nraw)
				if err != nil {
					return nil, err
				}
			} else {
				on, err := order.FromRaw(nraw)
				if err != nil {
					return nil, err
				}
				oNew := order.New(*on)
				o.NotifyInfo = &oNew
			}
		case "ou-req":
			on, err := order.FromRaw(nraw)
			if err != nil {
				return nil, err
			}
			ou := order.Update(*on)
			o.NotifyInfo = &ou
		case "oc-req":
			// if list of list then parse to order snapshot
			on, err := order.FromRaw(nraw)
			if err != nil {
				return o, err
			}
			oc := order.Cancel(*on)
			o.NotifyInfo = &oc
		case "fon-req":
			fon, err := fundingoffer.FromRaw(nraw)
			if err != nil {
				return o, err
			}
			fundingOffer := fundingoffer.New(*fon)
			o.NotifyInfo = &fundingOffer
		case "foc-req":
			foc, err := fundingoffer.FromRaw(nraw)
			if err != nil {
				return o, err
			}
			fundingOffer := fundingoffer.Cancel(*foc)
			o.NotifyInfo = &fundingOffer
		case "uca":
			o.NotifyInfo = raw[4]
		case "acc_tf":
			o.NotifyInfo = raw[4]
		case "pm-req":
			p, err := position.FromRaw(nraw)
			if err != nil {
				return o, err
			}
			cp := position.Cancel(*p)
			o.NotifyInfo = &cp
		default:
			o.NotifyInfo = raw[4]
		}
	}

	return
}

type bookAction byte

// BookAction represents a new/update or removal for a book entry.
type BookAction bookAction

const (
	//BookUpdateEntry represents a new or updated book entry.
	BookUpdateEntry BookAction = 0
	//BookRemoveEntry represents a removal of a book entry.
	BookRemoveEntry BookAction = 1
)

// BookUpdate represents an order book price update.
type BookUpdate struct {
	ID          int64            // the book update ID, optional
	Symbol      string           // book symbol
	Price       float64          // updated price
	PriceJsNum  json.Number      // update price as json.Number
	Count       int64            // updated count, optional
	Amount      float64          // updated amount
	AmountJsNum json.Number      // update amount as json.Number
	Side        common.OrderSide // side
	Action      BookAction       // action (add/remove)
}

type BookUpdateSnapshot struct {
	Snapshot []*BookUpdate
}

func NewBookUpdateSnapshotFromRaw(symbol, precision string, raw [][]float64, raw_numbers interface{}) (*BookUpdateSnapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for book snapshot: %#v", raw)
	}
	snap := make([]*BookUpdate, len(raw))
	for i, f := range raw {
		b, err := NewBookUpdateFromRaw(symbol, precision, convert.ToInterface(f), raw_numbers.([]interface{})[i])
		if err != nil {
			return nil, err
		}
		snap[i] = b
	}
	return &BookUpdateSnapshot{Snapshot: snap}, nil
}

func IsRawBook(precision string) bool {
	return precision == "R0"
}

// NewBookUpdateFromRaw creates a new book update object from raw data.  Precision determines how
// to interpret the side (baked into Count versus Amount)
// raw book updates [ID, price, qty], aggregated book updates [price, amount, count]
func NewBookUpdateFromRaw(symbol, precision string, data []interface{}, raw_numbers interface{}) (b *BookUpdate, err error) {
	if len(data) < 3 {
		return b, fmt.Errorf("data slice too short for book update, expected %d got %d: %#v", 5, len(data), data)
	}
	var px float64
	var px_num json.Number
	var id, cnt int64
	raw_num_array := raw_numbers.([]interface{})
	amt := convert.F64ValOrZero(data[2])
	amt_num := convert.FloatToJsonNumber(raw_num_array[2])

	var side common.OrderSide
	var actionCtrl float64
	if IsRawBook(precision) {
		// [ID, price, amount]
		id = convert.I64ValOrZero(data[0])
		px = convert.F64ValOrZero(data[1])
		px_num = convert.FloatToJsonNumber(raw_num_array[1])
		actionCtrl = px
	} else {
		// [price, amount, count]
		px = convert.F64ValOrZero(data[0])
		px_num = convert.FloatToJsonNumber(raw_num_array[0])
		cnt = convert.I64ValOrZero(data[1])
		actionCtrl = float64(cnt)
	}

	if amt > 0 {
		side = Bid
	} else {
		side = Ask
	}

	var action BookAction
	if actionCtrl <= 0 {
		action = BookRemoveEntry
	} else {
		action = BookUpdateEntry
	}

	b = &BookUpdate{
		Symbol:      symbol,
		Price:       math.Abs(px),
		PriceJsNum:  px_num,
		Count:       cnt,
		Amount:      math.Abs(amt),
		AmountJsNum: amt_num,
		Side:        side,
		Action:      action,
		ID:          id,
	}

	return
}

type StatKey string

const (
	FundingSizeKey   StatKey = "funding.size"
	CreditSizeKey    StatKey = "credits.size"
	CreditSizeSymKey StatKey = "credits.size.sym"
	PositionSizeKey  StatKey = "pos.size"
)

type Stat struct {
	Period int64
	Volume float64
}

type StatusType string
