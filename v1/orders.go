package bitfinex

import (
	"math"
	"strconv"
)

// Order types that the API can return.
const (
	OrderTypeMarket               = "market"
	OrderTypeLimit                = "limit"
	OrderTypeStop                 = "stop"
	OrderTypeTrailingStop         = "trailing-stop"
	OrderTypeFillOrKill           = "fill-or-kill"
	OrderTypeExchangeMarket       = "exchange market"
	OrderTypeExchangeLimit        = "exchange limit"
	OrderTypeExchangeStop         = "exchange stop"
	OrderTypeExchangeTrailingStop = "exchange trailing-stop"
	OrderTypeExchangeFillOrKill   = "exchange fill-or-kill"
)

// OrderService manages the Order endpoint.
type OrderService struct {
	client *Client
}

// Order represents one order on the bitfinex platform.
type Order struct {
	ID                int64
	Symbol            string
	Exchange          string
	Price             string
	AvgExecutionPrice string `json:"avg_execution_price"`
	Side              string
	Type              string
	Timestamp         string
	IsLive            bool   `json:"is_live"`
	IsCanceled        bool   `json:"is_cancelled"`
	IsHidden          bool   `json:"is_hidden"`
	WasForced         bool   `json:"was_forced"`
	OriginalAmount    string `json:"original_amount"`
	RemainingAmount   string `json:"remaining_amount"`
	ExecutedAmount    string `json:"executed_amount"`
}

// All returns all orders for the authenticated account.
func (s *OrderService) All() ([]Order, error) {
	req, err := s.client.newAuthenticatedRequest("GET", "orders", nil)
	if err != nil {
		return nil, err
	}

	v := []Order{}
	_, err = s.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// CancelAll active orders for the authenticated account.
func (s *OrderService) CancelAll() error {
	req, err := s.client.newAuthenticatedRequest("POST", "order/cancel/all", nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// Create a new order.
func (s *OrderService) Create(symbol string, amount float64, price float64, orderType string) (*Order, error) {
	var side string
	if amount < 0 {
		amount = math.Abs(amount)
		side = "sell"
	} else {
		side = "buy"
	}

	payload := map[string]interface{}{
		"symbol":   symbol,
		"amount":   strconv.FormatFloat(amount, 'f', -1, 32),
		"price":    strconv.FormatFloat(price, 'f', -1, 32),
		"side":     side,
		"type":     orderType,
		"exchange": "bitfinex",
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/new", payload)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	_, err = s.client.do(req, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Cancel the order with id `orderID`.
func (s *OrderService) Cancel(orderID int64) error {
	payload := map[string]interface{}{
		"order_id": orderID,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/cancel", payload)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// SubmitOrder is an order to be created on the bitfinex platform.
type SubmitOrder struct {
	Symbol string
	Amount float64
	Price  float64
	Type   string
}

// MultipleOrderResponse bundles orders returned by the CreateMulti method.
type MultipleOrderResponse struct {
	Orders []Order `json:"order_ids"`
	Status string
}

// CreateMulti allows batch creation of orders.
func (s *OrderService) CreateMulti(orders []SubmitOrder) (MultipleOrderResponse, error) {
	ordersMap := make([]interface{}, 0)
	for _, order := range orders {
		var side string
		if order.Amount < 0 {
			order.Amount = math.Abs(order.Amount)
			side = "sell"
		} else {
			side = "buy"
		}
		ordersMap = append(ordersMap, map[string]interface{}{
			"symbol":   order.Symbol,
			"amount":   strconv.FormatFloat(order.Amount, 'f', -1, 32),
			"price":    strconv.FormatFloat(order.Price, 'f', -1, 32),
			"exchange": "bitfinex",
			"side":     side,
			"type":     order.Type,
		})
	}

	payload := map[string]interface{}{
		"orders": ordersMap,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/new/multi", payload)
	if err != nil {
		return MultipleOrderResponse{}, err
	}

	response := new(MultipleOrderResponse)
	_, err = s.client.do(req, response)

	return *response, err

}

// CancelMulti allows batch cancellation of orders.
func (s *OrderService) CancelMulti(orderIDS []int64) (string, error) {
	payload := map[string]interface{}{
		"order_ids": orderIDS,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/cancel/multi", payload)

	if err != nil {
		return "", err
	}

	response := make(map[string]string, 0)
	_, err = s.client.do(req, &response)

	return response["result"], err
}

// Replace an Order
func (s *OrderService) Replace(orderID int64, useRemaining bool, newOrder SubmitOrder) (Order, error) {
	var side string
	if newOrder.Amount < 0 {
		newOrder.Amount = math.Abs(newOrder.Amount)
		side = "sell"
	} else {
		side = "buy"
	}

	payload := map[string]interface{}{
		"order_id":      strconv.FormatInt(orderID, 10),
		"symbol":        newOrder.Symbol,
		"amount":        strconv.FormatFloat(newOrder.Amount, 'f', -1, 32),
		"price":         strconv.FormatFloat(newOrder.Price, 'f', -1, 32),
		"exchange":      "bitfinex",
		"side":          side,
		"type":          newOrder.Type,
		"use_remaining": useRemaining,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/cancel/replace", payload)
	if err != nil {
		return Order{}, err
	}

	order := new(Order)
	_, err = s.client.do(req, order)
	if err != nil {
		return *order, err
	}

	return *order, nil
}

// Status retrieves the given order from the API.
func (s *OrderService) Status(orderID int64) (Order, error) {

	payload := map[string]interface{}{
		"order_id": orderID,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "order/status", payload)

	if err != nil {
		return Order{}, err
	}

	order := new(Order)
	_, err = s.client.do(req, order)
	if err != nil {
		return *order, err
	}

	return *order, nil
}
