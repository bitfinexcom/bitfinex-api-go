package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
	"strconv"
	"strings"
	"net/url"
	"fmt"
)

// CandleService manages the Candles endpoint.
type CandleService struct {
	Synchronous
}

// Return Candles for the public account.
func (c *CandleService) Last(symbol string, resolution bitfinex.CandleResolution) (*bitfinex.Candle, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	segments := []string{ "trade", string(resolution), symbol }

	req := NewRequestWithMethod(path.Join("candles", strings.Join(segments,":"), "LAST"), "GET")
	req.Params = make(url.Values)
	raw, err := c.Request(req)

	if err != nil {
		return nil, err
	}

	cs, err := bitfinex.NewCandleFromRaw(symbol, resolution, raw)

	if err != nil {
		return nil, err
	}

	return cs, nil
}

// Return Candles for the public account.
func (c *CandleService) History(symbol string, resolution bitfinex.CandleResolution) (*bitfinex.CandleSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	segments := []string{ "trade", string(resolution), symbol }

	req := NewRequestWithMethod(path.Join("candles", strings.Join(segments,":"), "HIST"), "GET")

	raw, err := c.Request(req)

	if err != nil {
		return nil, err
	}

	data := make([][]float64, 0, len(raw))
	for _, ifacearr := range raw {
		if arr, ok := ifacearr.([]interface{}); ok {
			sub := make([]float64, 0, len(arr))
			for _, iface := range arr {
				if flt, ok := iface.(float64); ok {
					sub = append(sub, flt)
				}
			}
			data = append(data, sub)
		}
	}

	cs, err := bitfinex.NewCandleSnapshotFromRaw(symbol, resolution, data)

	if err != nil {
		return nil, err
	}

	return cs, nil
}


// Return Candles for the public account.
func (c *CandleService) HistoryWithQuery(
	symbol string,
	resolution bitfinex.CandleResolution,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
	) (*bitfinex.CandleSnapshot, error) {
		if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	segments := []string{ "trade", string(resolution), symbol }

	req := NewRequestWithMethod(path.Join("candles", strings.Join(segments,":"), "HIST"), "GET")
	req.Params = make(url.Values)

	req.Params.Add("end", strconv.FormatInt(int64(end), 10))
	req.Params.Add("start", strconv.FormatInt(int64(start), 10))
	req.Params.Add("limit", strconv.FormatInt(int64(limit), 10))
	req.Params.Add("sort", strconv.FormatInt(int64(sort), 10))

	raw, err := c.Request(req)

	if err != nil {
		return nil, err
	}

	data := make([][]float64, 0, len(raw))
	for _, ifacearr := range raw {
		if arr, ok := ifacearr.([]interface{}); ok {
			sub := make([]float64, 0, len(arr))
			for _, iface := range arr {
				if flt, ok := iface.(float64); ok {
					sub = append(sub, flt)
				}
			}
			data = append(data, sub)
		}
	}

	cs, err := bitfinex.NewCandleSnapshotFromRaw(symbol, resolution, data)

	if err != nil {
		return nil, err
	}

	return cs, nil
}
