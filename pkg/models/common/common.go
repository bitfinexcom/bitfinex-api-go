package common

import "errors"

const (
	OrderFlagHidden               int       = 64
	OrderFlagClose                int       = 512
	OrderFlagPostOnly             int       = 4096
	OrderFlagOCO                  int       = 16384
	OrderTypeMarket                         = "MARKET"
	OrderTypeExchangeMarket                 = "EXCHANGE MARKET"
	OrderTypeLimit                          = "LIMIT"
	OrderTypeExchangeLimit                  = "EXCHANGE LIMIT"
	OrderTypeStop                           = "STOP"
	OrderTypeExchangeStop                   = "EXCHANGE STOP"
	OrderTypeTrailingStop                   = "TRAILING STOP"
	OrderTypeExchangeTrailingStop           = "EXCHANGE TRAILING STOP"
	OrderTypeFOK                            = "FOK"
	OrderTypeExchangeFOK                    = "EXCHANGE FOK"
	OrderTypeStopLimit                      = "STOP LIMIT"
	OrderTypeExchangeStopLimit              = "EXCHANGE STOP LIMIT"
	Bid                           OrderSide = 1
	Ask                           OrderSide = 2
	Long                          OrderSide = 1
	Short                         OrderSide = 2
	FundingPrefix                           = "f"
	TradingPrefix                           = "t"
)

var (
	ErrNotFound = errors.New("not found")
)

// OrderSide provides a typed set of order sides.
type OrderSide byte
