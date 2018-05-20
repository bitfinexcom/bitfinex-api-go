package rest

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// OrderService manages data flow for the Order API endpoint
type OrderService struct {
	requestFactory
	Synchronous
}

// All returns all orders for the authenticated account.
func (s *OrderService) All(symbol string) (*bitfinex.OrderSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(path.Join("orders", symbol))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
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
func (s *OrderService) History(symbol string) (*bitfinex.OrderSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}
	req, err := s.requestFactory.NewAuthenticatedRequest(path.Join("orders", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}

	os, err := bitfinex.NewOrderSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return os, nil
}

// OrderTrades returns a set of executed trades related to an order.
func (s *OrderService) OrderTrades(symbol string, orderID int64) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}
	key := fmt.Sprintf("%s:%d", symbol, orderID)
	req, err := s.requestFactory.NewAuthenticatedRequest(path.Join("order", key, "trades"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewTradeExecutionUpdateSnapshotFromRaw(raw)
}
