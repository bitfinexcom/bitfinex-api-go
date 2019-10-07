package rest

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"net/url"
	"path"
	"strings"
)

// TradeService manages the Trade endpoint.
type StatusService struct {
	requestFactory
	Synchronous
}

const (
	DERIV_TYPE = "deriv"
)

func (ss *StatusService) get(sType string, key string) (*bitfinex.DerivativeStatusSnapshot, error) {
	req := NewRequestWithMethod(path.Join("status", sType), "GET")
	req.Params = make(url.Values)
	req.Params.Add("keys", key)
	raw, err := ss.Request(req)
	if err != nil {
		return nil, err
	}
	trueRaw := make([][]interface{}, len(raw))
	for i, r := range raw {
		trueRaw[i] = r.([]interface{})
	}
	s, err := bitfinex.NewDerivativeSnapshotFromRaw(trueRaw)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Retrieves derivative status information for the given symbol from the platform
// see https://docs.bitfinex.com/reference#rest-public-status for more info
func (ss *StatusService) DerivativeStatus(symbol string) (*bitfinex.DerivativeStatus, error) {
	data, err := ss.get(DERIV_TYPE, symbol)
	if err != nil {
		return nil, err
	}
	if len(data.Snapshot) == 0 {
		return nil, fmt.Errorf("no status found for symbol %s", symbol)
	}
	return data.Snapshot[0], err
}

// Retrieves derivative status information for the given symbols from the platform
// see https://docs.bitfinex.com/reference#rest-public-status for more info
func (ss *StatusService) DerivativeStatusMulti(symbols []string) ([]*bitfinex.DerivativeStatus, error) {
	key := strings.Join(symbols, ",")
	data, err := ss.get(DERIV_TYPE, key)
	if err != nil {
		return nil, err
	}
	return data.Snapshot, err
}

// Retrieves derivative status information for all symbols from the platform
// see https://docs.bitfinex.com/reference#rest-public-status for more info
func (ss *StatusService) DerivativeStatusAll() ([]*bitfinex.DerivativeStatus, error) {
	data, err := ss.get(DERIV_TYPE, "ALL")
	if err != nil {
		return nil, err
	}
	return data.Snapshot, err
}

