package bitfinex

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
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

// CandleResolution provides a typed set of resolutions for candle subscriptions.
type CandleResolution string

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

// BookFrequency provides a typed book frequency.
type BookFrequency string

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
