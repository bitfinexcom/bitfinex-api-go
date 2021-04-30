package rest

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tickerhist"
)

// TickerHistoryService manages the Tickers History endpoint.
type TickerHistoryService struct {
	requestFactory
	Synchronous
}

type GetTickerHistPayload struct {
	Symbols []string
	Start   int64
	End     int64
	Limit   uint32
}

// Get - retrieves the ticker history for the given symbol
// see https://docs.bitfinex.com/reference#tickers-history for more info
func (s *TickerHistoryService) Get(pld GetTickerHistPayload) ([]tickerhist.TickerHist, error) {
	if len(pld.Symbols) == 0 {
		return nil, fmt.Errorf("missing mandatory parameters: []Symbols")
	}

	req := NewRequestWithMethod("tickers/hist", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", strings.Join(pld.Symbols, ","))

	if pld.Start != 0 {
		req.Params.Add("start", fmt.Sprintf("%d", pld.Start))
	}

	if pld.End != 0 {
		req.Params.Add("end", fmt.Sprintf("%d", pld.End))
	}

	if pld.Limit != 0 {
		req.Params.Add("limit", fmt.Sprintf("%d", pld.Limit))
	}

	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	tickers := tickerhist.SnapshotFromRaw(convert.ToInterfaceArray(raw))
	return tickers.Snapshot, nil
}
