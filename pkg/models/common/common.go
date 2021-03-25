package common

import (
	"errors"
	"fmt"
)

const (
	OrderFlagHidden               int              = 64
	OrderFlagClose                int              = 512
	OrderFlagPostOnly             int              = 4096
	OrderFlagOCO                  int              = 16384
	Checksum                      int              = 131072
	OrderStatusActive                              = "ACTIVE"
	OrderStatusExecuted                            = "EXECUTED"
	OrderStatusPartiallyFilled                     = "PARTIALLY FILLED"
	OrderStatusCanceled                            = "CANCELED"
	OrderTypeExchangeLimit                         = "EXCHANGE LIMIT"
	OrderTypeMarket                                = "MARKET"
	OrderTypeExchangeMarket                        = "EXCHANGE MARKET"
	OrderTypeLimit                                 = "LIMIT"
	OrderTypeStop                                  = "STOP"
	OrderTypeExchangeStop                          = "EXCHANGE STOP"
	OrderTypeTrailingStop                          = "TRAILING STOP"
	OrderTypeExchangeTrailingStop                  = "EXCHANGE TRAILING STOP"
	OrderTypeFOK                                   = "FOK"
	OrderTypeExchangeFOK                           = "EXCHANGE FOK"
	OrderTypeStopLimit                             = "STOP LIMIT"
	OrderTypeExchangeStopLimit                     = "EXCHANGE STOP LIMIT"
	PermissionRead                                 = "r"
	PermissionWrite                                = "w"
	FundingPrefix                                  = "f"
	TradingPrefix                                  = "t"
	FundingSizeKey                StatKey          = "funding.size"
	CreditSizeKey                 StatKey          = "credits.size"
	CreditSizeSymKey              StatKey          = "credits.size.sym"
	PositionSizeKey               StatKey          = "pos.size"
	Bid                           OrderSide        = 1
	Ask                           OrderSide        = 2
	Long                          OrderSide        = 1
	Short                         OrderSide        = 2
	OldestFirst                   SortOrder        = 1
	NewestFirst                   SortOrder        = -1
	OneMinute                     CandleResolution = "1m"
	FiveMinutes                   CandleResolution = "5m"
	FifteenMinutes                CandleResolution = "15m"
	ThirtyMinutes                 CandleResolution = "30m"
	OneHour                       CandleResolution = "1h"
	ThreeHours                    CandleResolution = "3h"
	SixHours                      CandleResolution = "6h"
	TwelveHours                   CandleResolution = "12h"
	OneDay                        CandleResolution = "1D"
	OneWeek                       CandleResolution = "7D"
	TwoWeeks                      CandleResolution = "14D"
	OneMonth                      CandleResolution = "1M"
	Precision0                    BookPrecision    = "P0" // Aggregate precision levels
	Precision1                    BookPrecision    = "P1" // Aggregate precision levels
	Precision2                    BookPrecision    = "P2" // Aggregate precision levels
	Precision3                    BookPrecision    = "P3" // Aggregate precision levels
	PrecisionRawBook              BookPrecision    = "R0" // Raw precision
	// FrequencyRealtime book frequency gives updates as they occur in real-time.
	FrequencyRealtime BookFrequency = "F0"
	// FrequencyTwoPerSecond delivers two book updates per second.
	FrequencyTwoPerSecond BookFrequency = "F1"
	// PriceLevelDefault provides a constant default price level for book subscriptions.
	PriceLevelDefault int = 25
)

var (
	ErrNotFound = errors.New("not found")
)

// OrderSide provides a typed set of order sides.
type OrderSide byte

// CandleResolution provides a typed set of resolutions for candle subscriptions.
type CandleResolution string

// BookPrecision provides a typed book precision level.
type BookPrecision string

// BookFrequency provides a typed book frequency.
type BookFrequency string

type SortOrder int

type QueryLimit int

type PermissionType string

type Mts int64

type StatKey string

type StatusType string

type OrderType string

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
