package rest

import (
	"net/url"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
)

// TickerService manages the Ticker endpoint.
type TickerService struct {
	requestFactory
	Synchronous
}

func (s *TickerService) getTickers(symbols []string) ([]*ticker.Ticker, error) {
	req := NewRequestWithMethod("tickers", "GET")
	req.Params = make(url.Values)
	req.Params.Add("symbols", strings.Join(symbols, ","))
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	tickers := make([]*ticker.Ticker, 0)
	for _, traw := range raw {
		t, err := ticker.FromRestRaw(traw.([]interface{}))
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, t)
	}

	return tickers, nil
}

// Get - retrieves the ticker for the given symbol
// see https://docs.bitfinex.com/reference#rest-public-tickers for more info
func (s *TickerService) Get(symbol string) (*ticker.Ticker, error) {
	t, err := s.getTickers([]string{symbol})
	if err != nil {
		return nil, err
	}

	return t[0], nil
}

// GetMulti - retrieves the tickers for the given symbols
// see https://docs.bitfinex.com/reference#rest-public-tickers for more info
func (s *TickerService) GetMulti(symbols []string) ([]*ticker.Ticker, error) {
	return s.getTickers(symbols)
}

// All - retrieves all tickers for all symbols
// see https://docs.bitfinex.com/reference#rest-public-tickers for more info
func (s *TickerService) All() ([]*ticker.Ticker, error) {
	return s.getTickers([]string{"ALL"})
}
