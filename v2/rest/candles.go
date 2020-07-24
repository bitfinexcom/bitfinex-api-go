package rest

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// CandleService manages the Candles endpoint.
type CandleService struct {
	Synchronous
}

func getPathSegments(symbol string, resolution common.CandleResolution) (s string, err error) {
	if len(symbol) == 0 {
		err = fmt.Errorf("symbol cannot be empty")
		return
	}

	segments := []string{"trade", string(resolution), symbol}
	s = strings.Join(segments, ":")
	return
}

// Last - retrieve the last candle for the given symbol with the given resolution
// See https://docs.bitfinex.com/reference#rest-public-candles for more info
func (c *CandleService) Last(symbol string, resolution common.CandleResolution) (*candle.Candle, error) {
	segments, err := getPathSegments(symbol, resolution)
	if err != nil {
		return nil, err
	}

	req := NewRequestWithMethod(path.Join("candles", segments, "LAST"), "GET")
	raw, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	cs, err := candle.FromRaw(symbol, resolution, raw)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// History - retrieves all candles (Max=1000) with the given symbol and the given candle resolution
// See https://docs.bitfinex.com/reference#rest-public-candles for more info
func (c *CandleService) History(symbol string, resolution common.CandleResolution) (*candle.Snapshot, error) {
	segments, err := getPathSegments(symbol, resolution)
	if err != nil {
		return nil, err
	}

	req := NewRequestWithMethod(path.Join("candles", segments, "HIST"), "GET")
	raw, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	cs, err := candle.SnapshotFromRaw(symbol, resolution, convert.ToInterfaceArray(raw))
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// HistoryWithQuery - retrieves all candles (Max=1000) that fit the given query criteria
// See https://docs.bitfinex.com/reference#rest-public-candles for more info
func (c *CandleService) HistoryWithQuery(
	symbol string,
	resolution common.CandleResolution,
	start common.Mts,
	end common.Mts,
	limit common.QueryLimit,
	sort common.SortOrder,
) (*candle.Snapshot, error) {
	segments, err := getPathSegments(symbol, resolution)
	if err != nil {
		return nil, err
	}

	req := NewRequestWithMethod(path.Join("candles", segments, "HIST"), "GET")
	req.Params = make(url.Values)
	req.Params.Add("end", strconv.FormatInt(int64(end), 10))
	req.Params.Add("start", strconv.FormatInt(int64(start), 10))
	req.Params.Add("limit", strconv.FormatInt(int64(limit), 10))
	req.Params.Add("sort", strconv.FormatInt(int64(sort), 10))

	raw, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	cs, err := candle.SnapshotFromRaw(symbol, resolution, convert.ToInterfaceArray(raw))
	if err != nil {
		return nil, err
	}

	return cs, nil
}
