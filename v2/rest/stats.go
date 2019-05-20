package rest

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// TradeService manages the Trade endpoint.
type StatsService struct {
	requestFactory
	Synchronous
}

func (ss *StatsService) get(symbol string, key bitfinex.StatKey, extra string, section string) ([]interface{}, error) {
	var params string
	if extra != "" {
		params = fmt.Sprintf("%s:1m:%s:%s", string(key), symbol, extra)
	} else {
		params = fmt.Sprintf("%s:1m:%s", string(key), symbol)
	}
	req := NewRequestWithMethod(path.Join("stats1", params, section), "GET")
	raw, err := ss.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (ss *StatsService) getHistory(symbol string, key bitfinex.StatKey, extra string) ([]bitfinex.Stat, error) {
	stats, err := ss.get(symbol, key, extra, "hist")
	if err != nil {
		return nil, err
	}
	res := make([]bitfinex.Stat, len(stats))
	for index, stat := range stats {
		arr := stat.([]interface{})
		period := arr[0].(float64)
		volume := arr[1].(float64)
		res[index] = bitfinex.Stat{Period: int64(period), Volume: volume}
	}
	return res, nil
}

func (ss *StatsService) getLast(symbol string, key bitfinex.StatKey, extra string) (*bitfinex.Stat, error) {
	stat, err := ss.get(symbol, key, extra, "last")
	if err != nil {
		return nil, err
	}
	if len(stat) == 0 {
		return nil, fmt.Errorf("Unable to get last stat for %s:%s", symbol, key)
	}
	period := stat[0].(float64)
	volume := stat[1].(float64)
	return &bitfinex.Stat{Period: int64(period), Volume: volume}, nil
}

func (ss *StatsService) FundingHistory(symbol string) ([]bitfinex.Stat, error) {
	return ss.getHistory(symbol, bitfinex.FundingSizeKey, "")
}

func (ss *StatsService) FundingLast(symbol string) (*bitfinex.Stat, error) {
	return ss.getLast(symbol, bitfinex.FundingSizeKey, "")
}

func (ss *StatsService) CreditSizeHistory(symbol string, side bitfinex.OrderSide) ([]bitfinex.Stat, error) {
	return ss.getHistory(symbol, bitfinex.CreditSizeKey, "")
}

func (ss *StatsService) CreditSizeLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error) {
	return ss.getLast(symbol, bitfinex.CreditSizeKey, "")
}

func (ss *StatsService) SymbolCreditSizeHistory(fundingSymbol string, tradingSymbol string) ([]bitfinex.Stat, error) {
	return ss.getHistory(fundingSymbol, bitfinex.CreditSizeSymKey, tradingSymbol)
}

func (ss *StatsService) SymbolCreditSizeLast(fundingSymbol string, tradingSymbol string) (*bitfinex.Stat, error) {
	return ss.getLast(fundingSymbol, bitfinex.CreditSizeSymKey, tradingSymbol)
}

func (ss *StatsService) PositionHistory(symbol string, side bitfinex.OrderSide) ([]bitfinex.Stat, error) {
	var strSide string
	if side == bitfinex.Long {
		strSide = "long"
	} else if side == bitfinex.Short {
		strSide = "short"
	} else {
		return nil, fmt.Errorf("Unrecognized side %v in PositionHistory", side)
	}
	return ss.getHistory(symbol, bitfinex.PositionSizeKey, strSide)
}

func (ss *StatsService) PositionLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error) {
	var strSide string
	if side == bitfinex.Long {
		strSide = "long"
	} else if side == bitfinex.Short {
		strSide = "short"
	} else {
		return nil, fmt.Errorf("Unrecognized side %v in PositionHistory", side)
	}
	return ss.getLast(symbol, bitfinex.PositionSizeKey, strSide)
}
