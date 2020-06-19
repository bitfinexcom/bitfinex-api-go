package rest

import (
	"net/url"
	"path"
	"strconv"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type MarketService struct {
	requestFactory
	Synchronous
}

// AveragePriceArgs data structure constructing average price query params
type AveragePriceArgs struct {
	Symbol    string
	Amount    string
	RateLimit string
	Period    int
}

// AveragePrice Calculate the average execution price for Trading or rate for Margin funding.
// See: https://docs.bitfinex.com/reference#rest-public-calc-market-average-price
func (ms *MarketService) AveragePrice(payload AveragePriceArgs) ([]float64, error) {
	req := NewRequestWithMethod(path.Join("calc", "trade", "avg"), "POST")
	req.Params = make(url.Values)
	req.Params.Add("symbol", payload.Symbol)
	req.Params.Add("amount", payload.Amount)
	req.Params.Add("rate_limit", payload.RateLimit)
	req.Params.Add("period", strconv.Itoa(payload.Period))

	raw, err := ms.Request(req)
	if err != nil {
		return nil, err
	}

	resp, err := convert.F64Slice(raw)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
