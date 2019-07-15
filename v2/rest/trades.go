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
func (s *TradeService) All(symbol string) (*bitfinex.TradeSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(bitfinex.PermissionRead, path.Join("trades", symbol, "hist"), map[string]interface{}{"start": nil, "end": nil, "limit": nil})
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawToSnapshot(symbol, raw)
}

func (s *TradeService) AccountAll(symbol string) (*bitfinex.TradeSnapshot, error) {
	return s.All(symbol)
}

// return account trades that fit the given conditions
func (s *TradeService) AccountHistoryWithQuery(
	symbol string,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
	) (*bitfinex.TradeSnapshot, error) {
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
	return parseRawToSnapshot(symbol, raw)
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
		return parseRawToSnapshot(symbol, raw)
}

func parseRawToSnapshot(symbol string, raw []interface{}) (*bitfinex.TradeSnapshot, error) {
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
