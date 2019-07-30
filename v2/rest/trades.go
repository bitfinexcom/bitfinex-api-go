package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"net/url"
	"path"
	"strconv"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
// left this in her
func (s *TradeService) allAccountWithSymbol(symbol string) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("trades", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawPrivateToSnapshot(raw)
}

func (s *TradeService) allAccount() (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("trades", "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawPrivateToSnapshot(raw)
}

func (s *TradeService) AccountAll() (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	return s.allAccount()
}

func (s *TradeService) AccountAllWithSymbol(symbol string) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	return s.allAccountWithSymbol(symbol)
}

// return account trades that fit the given conditions
func (s *TradeService) AccountHistoryWithQuery(
	symbol string,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
	) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("trades", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	req.Params = make(url.Values)
	req.Params.Add("end", strconv.FormatInt(int64(end), 10))
	req.Params.Add("start", strconv.FormatInt(int64(start), 10))
	req.Params.Add("limit", strconv.FormatInt(int64(limit), 10))
	req.Params.Add("sort", strconv.FormatInt(int64(sort), 10))
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawPrivateToSnapshot(raw)
}

// return publicly executed trades that fit the given query conditions
func (s *TradeService) PublicHistoryWithQuery(
	symbol string,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
	) (*bitfinex.TradeSnapshot, error) {
		req := NewRequestWithMethod(path.Join("trades", symbol, "hist"), "GET")
		req.Params = make(url.Values)
		req.Params.Add("end", strconv.FormatInt(int64(end), 10))
		req.Params.Add("start", strconv.FormatInt(int64(start), 10))
		req.Params.Add("limit", strconv.FormatInt(int64(limit), 10))
		req.Params.Add("sort", strconv.FormatInt(int64(sort), 10))
		raw, err := s.Request(req)
		if err != nil {
			return nil, err
		}
		return parseRawPublicToSnapshot(symbol, raw)
}

func parseRawPublicToSnapshot(symbol string, raw []interface{}) (*bitfinex.TradeSnapshot, error) {
	if len(raw) <= 0 {
		// return empty
		return &bitfinex.TradeSnapshot{Snapshot: make([]*bitfinex.Trade, 0)}, nil
	}
	// convert to array of floats
	dat := make([][]float64, 0)
	for _, r := range raw {
		t := []float64{}
		for _, r2 := range r.([]interface{}) {
			t = append(t, r2.(float64))
		}
		dat = append(dat, t)
	}
	return bitfinex.NewTradeSnapshotFromRaw(symbol, dat)
}

func parseRawPrivateToSnapshot(raw []interface{}) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	if len(raw) <= 0 {
		return &bitfinex.TradeExecutionUpdateSnapshot{Snapshot: make([]*bitfinex.TradeExecutionUpdate, 0)}, nil
	}
	tradeExecutions, err := bitfinex.NewTradeExecutionUpdateSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return tradeExecutions, nil
}
