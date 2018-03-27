package rest

import (
	"fmt"
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// OrderService manages data flow for the Order API endpoint
type OrderService struct {
	Synchronous
	Authenticator
}

// All returns all orders for the authenticated account.
func (s *OrderService) All(symbol string) (*bitfinex.OrderSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	r, err := s.NewAuthenticatedPostRequest(path.Join("auth", "r", "orders", symbol), nil)
	if err != nil {
		return nil, err
	}

	raw, err := s.Request(r)
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewOrderSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}

// Status retrieves the given order from the API. This is just a wrapper around
// the All() method, since the API does not provide lookup for a single Order.
func (s *OrderService) Status(orderID int64) (o *bitfinex.Order, err error) {
	os, err := s.All("")

	if err != nil {
		return o, err
	}

	if len(os.Snapshot) == 0 {
		return o, bitfinex.ErrNotFound
	}

	for _, e := range os.Snapshot {
		if e.ID == orderID {
			return e, nil
		}
	}

	return o, bitfinex.ErrNotFound
}

// All returns all orders for the authenticated account.
func (s *OrderService) History(symbol string, params ...map[string]interface{}) (*bitfinex.OrderSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	p := ReadParams(params...)
	r, err := s.NewAuthenticatedPostRequest(path.Join("auth", "r", "orders", symbol, "hist"), p)
	if err != nil {
		return nil, err
	}

	raw, err := s.Request(r)
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewOrderSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
