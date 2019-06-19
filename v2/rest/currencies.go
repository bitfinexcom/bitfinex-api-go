package rest

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
	"strings"
)

// TradeService manages the Trade endpoint.
type CurrenciesService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
func (cs *CurrenciesService) Conf(label, symbol, unit, explorer, pairs bool) ([]bitfinex.CurrencyConf, error) {
	segments := make([]string, 0)
	if label {
		segments = append(segments, string(bitfinex.CurrencyLabelMap))
	}
	if symbol {
		segments = append(segments, string(bitfinex.CurrencySymbolMap))
	}
	if unit {
		segments = append(segments, string(bitfinex.CurrencyUnitMap))
	}
	if explorer {
		segments = append(segments, string(bitfinex.CurrencyExplorerMap))
	}
	if pairs {
		segments = append(segments, string(bitfinex.CurrencyExchangeMap))
	}
	req := NewRequestWithMethod(path.Join("conf", strings.Join(segments,",")), "GET")
	raw, err := cs.Request(req)
	if err != nil {
		return nil, err
	}
	// add mapping to raw data
	parsedRaw := make([]bitfinex.RawCurrencyConf, len(raw))
	for index, d := range raw {
		parsedRaw = append(parsedRaw, bitfinex.RawCurrencyConf{Mapping: segments[index], Data: d})
	}
	// parse to config object
	configs, err := bitfinex.NewCurrencyConfFromRaw(parsedRaw)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

