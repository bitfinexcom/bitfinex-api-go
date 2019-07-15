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

// Get all active orders
func (s *OrderService) All() (*bitfinex.OrderSnapshot, error) {
	// use no symbol, this will get all orders
	return s.getActiveOrders("")
}

// Get all active orders with the given symbol
func (s *OrderService) GetBySymbol(symbol string) (*bitfinex.OrderSnapshot, error) {
	// use no symbol, this will get all orders
	return s.getActiveOrders(symbol)
}

// Get an active order using its order id
func (s *OrderService) GetByOrderId(orderID int64) (o *bitfinex.Order, err error) {
	os, err := s.All()
	if err != nil {
		return nil, err
	}
	for _, order := range os.Snapshot {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, bitfinex.ErrNotFound
}

// Get all historical orders
func (s *OrderService) AllHistory() (*bitfinex.OrderSnapshot, error) {
	// use no symbol, this will get all orders
	return s.getHistoricalOrders("")
}

// Get all historical orders with the given symbol
func (s *OrderService) GetHistoryBySymbol(symbol string) (*bitfinex.OrderSnapshot, error) {
	// use no symbol, this will get all orders
	return s.getHistoricalOrders(symbol)
}

// Get a historical order using its order id
func (s *OrderService) GetHistoryByOrderId(orderID int64) (o *bitfinex.Order, err error) {
	os, err := s.AllHistory()
	if err != nil {
		return nil, err
	}
	for _, order := range os.Snapshot {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, bitfinex.ErrNotFound
}

// OrderTrades returns a set of executed trades related to an order.
func (s *OrderService) OrderTrades(symbol string, orderID int64) (*bitfinex.TradeExecutionUpdateSnapshot, error) {
	key := fmt.Sprintf("%s:%d", symbol, orderID)
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("order", key, "trades"))
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewTradeExecutionUpdateSnapshotFromRaw(raw)
}

func (s *OrderService) getActiveOrders(symbol string) (*bitfinex.OrderSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("orders", symbol))
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
	if os == nil {
		return &bitfinex.OrderSnapshot{}, nil
	}
	return os, nil
}

func (s *OrderService) getHistoricalOrders(symbol string) (*bitfinex.OrderSnapshot, error) {
	req, err := s.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("orders", symbol, "hist"))
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
	if os == nil {
		return &bitfinex.OrderSnapshot{}, nil
	}
	return os, nil
}
