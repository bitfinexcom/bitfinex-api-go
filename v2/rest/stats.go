package rest

import (
	"fmt"
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/stats"
)

// StatsService manages the Stats endpoint.
type StatsService struct {
	requestFactory
	Synchronous
}

func (ss *StatsService) get(symbol string, key common.StatKey, extra string, section string) ([]interface{}, error) {
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

func (ss *StatsService) getHistory(symbol string, key common.StatKey, extra string) ([]*stats.Stat, error) {
	raw, err := ss.get(symbol, key, extra, "hist")
	if err != nil {
		return nil, err
	}

	return stats.SnapshotFromRaw(raw)
}

func (ss *StatsService) getLast(symbol string, key common.StatKey, extra string) (*stats.Stat, error) {
	raw, err := ss.get(symbol, key, extra, "last")
	if err != nil {
		return nil, err
	}

	return stats.FromRaw(raw)
}

// Retrieves platform statistics for funding history
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) FundingHistory(symbol string) ([]*stats.Stat, error) {
	return ss.getHistory(symbol, common.FundingSizeKey, "")
}

// Retrieves platform statistics for funding last
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) FundingLast(symbol string) (*stats.Stat, error) {
	return ss.getLast(symbol, common.FundingSizeKey, "")
}

// Retrieves platform statistics for credit size history
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) CreditSizeHistory(symbol string, side common.OrderSide) ([]*stats.Stat, error) {
	return ss.getHistory(symbol, common.CreditSizeKey, "")
}

// Retrieves platform statistics for credit size last
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) CreditSizeLast(symbol string, side common.OrderSide) (*stats.Stat, error) {
	return ss.getLast(symbol, common.CreditSizeKey, "")
}

// Retrieves platform statistics for credit size history
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) SymbolCreditSizeHistory(fundingSymbol string, tradingSymbol string) ([]*stats.Stat, error) {
	return ss.getHistory(fundingSymbol, common.CreditSizeSymKey, tradingSymbol)
}

// Retrieves platform statistics for credit size last
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) SymbolCreditSizeLast(fundingSymbol string, tradingSymbol string) (*stats.Stat, error) {
	return ss.getLast(fundingSymbol, common.CreditSizeSymKey, tradingSymbol)
}

// Retrieves platform statistics for position history
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) PositionHistory(symbol string, side common.OrderSide) ([]*stats.Stat, error) {
	var strSide string
	if side == common.Long {
		strSide = "long"
	} else if side == common.Short {
		strSide = "short"
	} else {
		return nil, fmt.Errorf("Unrecognized side %v in PositionHistory", side)
	}
	return ss.getHistory(symbol, common.PositionSizeKey, strSide)
}

// Retrieves platform statistics for position last
// see https://docs.bitfinex.com/reference#rest-public-stats for more info
func (ss *StatsService) PositionLast(symbol string, side common.OrderSide) (*stats.Stat, error) {
	var strSide string
	if side == common.Long {
		strSide = "long"
	} else if side == common.Short {
		strSide = "short"
	} else {
		return nil, fmt.Errorf("Unrecognized side %v in PositionHistory", side)
	}
	return ss.getLast(symbol, common.PositionSizeKey, strSide)
}
