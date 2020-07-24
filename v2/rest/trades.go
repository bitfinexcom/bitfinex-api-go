package rest

import (
	"net/url"
	"path"
	"strconv"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecutionupdate"
)

// TradeService manages the Trade endpoint.
type TradeService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
// left this in her
func (s *TradeService) allAccountWithSymbol(symbol string) (*tradeexecutionupdate.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("trades", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawPrivateToSnapshot(raw)
}

func (s *TradeService) allAccount() (*tradeexecutionupdate.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("trades", "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return parseRawPrivateToSnapshot(raw)
}

// Retrieves all matched trades for the account
// see https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info
func (s *TradeService) AccountAll() (*tradeexecutionupdate.Snapshot, error) {
	return s.allAccount()
}

// Retrieves all matched trades with the given symbol for the account
// see https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info
func (s *TradeService) AccountAllWithSymbol(symbol string) (*tradeexecutionupdate.Snapshot, error) {
	return s.allAccountWithSymbol(symbol)
}

// Queries all matched trades with group of optional parameters
// see https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info
func (s *TradeService) AccountHistoryWithQuery(
	symbol string,
	start common.Mts,
	end common.Mts,
	limit common.QueryLimit,
	sort common.SortOrder,
) (*tradeexecutionupdate.Snapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(common.PermissionRead, path.Join("trades", symbol, "hist"))
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

// Queries all public trades with a group of optional paramters
// see https://docs.bitfinex.com/reference#rest-public-trades for more info
func (s *TradeService) PublicHistoryWithQuery(
	symbol string,
	start common.Mts,
	end common.Mts,
	limit common.QueryLimit,
	sort common.SortOrder,
) (*trade.Snapshot, error) {
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

	if len(raw) <= 0 {
		return &trade.Snapshot{Snapshot: make([]*trade.Trade, 0)}, nil
	}

	return trade.SnapshotFromRaw(symbol, convert.ToInterfaceArray(raw))
}

func parseRawPrivateToSnapshot(raw []interface{}) (*tradeexecutionupdate.Snapshot, error) {
	if len(raw) <= 0 {
		return &tradeexecutionupdate.Snapshot{Snapshot: make([]*tradeexecutionupdate.TradeExecutionUpdate, 0)}, nil
	}
	tradeExecutions, err := tradeexecutionupdate.SnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return tradeExecutions, nil
}
