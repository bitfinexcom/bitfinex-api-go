package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type MarketService struct {
	requestFactory
	Synchronous
}

// AveragePriceRequest data structure for constructing average price query params
type AveragePriceRequest struct {
	Symbol    string
	Amount    string
	RateLimit string
	Period    int
}

// ForeignExchangeRateRequest data structure for constructing foreign
// exchange rate request payload
type ForeignExchangeRateRequest struct {
	FirstCurrency  string `json:"ccy1"`
	SecondCurrency string `json:"ccy2"`
}

// AveragePrice Calculate the average execution price for Trading or rate for Margin funding.
// See: https://docs.bitfinex.com/reference#rest-public-calc-market-average-price
func (ms *MarketService) AveragePrice(pld AveragePriceRequest) ([]float64, error) {
	req := NewRequestWithMethod(path.Join("calc", "trade", "avg"), "POST")
	req.Params = make(url.Values)
	req.Params.Add("symbol", pld.Symbol)
	req.Params.Add("amount", pld.Amount)
	req.Params.Add("rate_limit", pld.RateLimit)
	req.Params.Add("period", strconv.Itoa(pld.Period))

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

// ForeignExchangeRate - Calculate the exchange rate between two currencies
// See: https://docs.bitfinex.com/reference#rest-public-calc-foreign-exchange-rate
func (ms *MarketService) ForeignExchangeRate(pld ForeignExchangeRateRequest) ([]float64, error) {
	if len(pld.FirstCurrency) == 0 || len(pld.SecondCurrency) == 0 {
		return nil, fmt.Errorf("FirstCurrency and SecondCurrency are required arguments")
	}

	bytes, err := json.Marshal(pld)
	if err != nil {
		return nil, err
	}

	req := NewRequestWithBytes(path.Join("calc", "fx"), bytes)
	req.Headers["Content-Type"] = "application/json"

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
