package bitfinex

import (
	"fmt"
	"path"
)

// OrderService manages the Order endpoint.
type OrderService struct {
	client *Client
}

// All returns all orders for the authenticated account.
func (s *OrderService) All(symbol string) (OrderSnapshot, error) {
	req, err := s.client.newAuthenticatedRequest("POST", path.Join("orders", symbol), nil)
	if err != nil {
		return nil, err
	}

	var raw []interface{}
	_, err = s.client.do(req, &raw)
	if err != nil {
		return nil, err
	}

	os, err := orderSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}

// Status retrieves the given order from the API. This is just a wrapper around
// the All() method, since the API does not provide lookup for a single Order.
func (s *OrderService) Status(orderID int64) (o Order, err error) {
	os, err := s.All("")
	if err != nil {
		return o, err
	}
	if len(os) == 0 {
		return o, ErrNotFound
	}

	for _, e := range os {
		if e.ID == orderID {
			return e, nil
		}
	}

	return o, ErrNotFound
}

// All returns all orders for the authenticated account.
func (s *OrderService) History(symbol string) (OrderSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	req, err := s.client.newAuthenticatedRequest("POST", path.Join("orders", symbol, "hist"), nil)
	if err != nil {
		return nil, err
	}

	var raw []interface{}
	_, err = s.client.do(req, &raw)
	if err != nil {
		return nil, err
	}
	fmt.Printf("raw: %#v\n", raw)

	os, err := orderSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}
