package rest

import (
	"path"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/currency"
)

// CurrenciesService manages the conf endpoint.
type CurrenciesService struct {
	requestFactory
	Synchronous
}

// Conf - retreive currency and symbol service configuration data
// see https://docs.bitfinex.com/reference#rest-public-conf for more info
func (cs *CurrenciesService) Conf(label, symbol, unit, explorer, pairs bool) ([]currency.Conf, error) {
	segments := make([]string, 0)
	if label {
		segments = append(segments, string(currency.LabelMap))
	}
	if symbol {
		segments = append(segments, string(currency.SymbolMap))
	}
	if unit {
		segments = append(segments, string(currency.UnitMap))
	}
	if explorer {
		segments = append(segments, string(currency.ExplorerMap))
	}
	if pairs {
		segments = append(segments, string(currency.ExchangeMap))
	}

	req := NewRequestWithMethod(path.Join("conf", strings.Join(segments, ",")), "GET")
	raw, err := cs.Request(req)
	if err != nil {
		return nil, err
	}

	// add mapping to raw data
	parsedRaw := make([]currency.RawConf, len(raw))
	for index, d := range raw {
		parsedRaw = append(parsedRaw, currency.RawConf{Mapping: segments[index], Data: d})
	}

	// parse to config object
	configs, err := currency.FromRaw(parsedRaw)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
