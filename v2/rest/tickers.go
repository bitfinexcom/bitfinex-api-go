package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"net/url"
	"strings"
)

// TradeService manages the Trade endpoint.
type TickerService struct {
	requestFactory
	Synchronous
}

// Retrieves the ticker for the given symbol
// see https://docs.bitfinex.com/reference#rest-public-ticker for more info
func (s *TickerService) Get(symbol string) (*bitfinex.Ticker, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", symbol)
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	ticker, err := bitfinex.NewTickerFromRestRaw(raw[0].([]interface{}))
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

// Retrieves the tickers for the given symbols
// see https://docs.bitfinex.com/reference#rest-public-ticker for more info
func (s *TickerService) GetMulti(symbols []string) (*[]bitfinex.Ticker, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", strings.Join(symbols, ","))
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	tickers := make([]bitfinex.Ticker, 0)
	for _, ticker := range raw {
		t, err := bitfinex.NewTickerFromRestRaw(ticker.([]interface{}))
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, *t)
	}
	return &tickers, nil
}

// Retrieves all tickers for all symbols
// see https://docs.bitfinex.com/reference#rest-public-ticker for more info
func (s *TickerService) All() (*[]bitfinex.Ticker, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", "ALL")
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	tickers := make([]bitfinex.Ticker, 0)
	for _, ticker := range raw {
		t, err := bitfinex.NewTickerFromRestRaw(ticker.([]interface{}))
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, *t)
	}
	return &tickers, nil
}
