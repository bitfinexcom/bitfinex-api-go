package bitfinex

import (
	"fmt"
	"math"
	"strconv"
)

const (
	ORDER_TYPE_MARKET                 = "market"
	ORDER_TYPE_LIMIT                  = "limit"
	ORDER_TYPE_STOP                   = "stop"
	ORDER_TYPE_TRAILING_STOP          = "trailing-stop"
	ORDER_TYPE_FILL_OR_KILL           = "fill-or-kill"
	ORDER_TYPE_EXCHANGE_MARKET        = "exchange market"
	ORDER_TYPE_EXCHANGE_LIMIT         = "exchange limit"
	ORDER_TYPE_EXCHANGE_STOP          = "exchange stop"
	ORDER_TYPE_EXCHANGE_TRAILING_STOP = "exchange trailing-stop"
	ORDER_TYPE_EXCHANGE_FILL_OR_KILL  = "exchange fill-or-kill"
)

type OrderService struct {
	client *Client
}

type Order struct {
	Id                int
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

// get all active orders
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

// Cancel all active orders
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

// Create a new order
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
		"amount":   fmt.Sprintf("%f", amount),
		"price":    fmt.Sprintf("%f", price),
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

// Cancel the order with id `orderId`
func (s *OrderService) Cancel(orderId int) error {
	payload := map[string]interface{}{
		"order_id": orderId,
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

type SubmitOrder struct {
	Symbol string
	Amount float64
	Price  float64
	Type   string
}

type MultipleOrderResponse struct {
	Orders []Order `json:"order_ids"`
	Status string
}

// Create Multiple Orders
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

// Cancel multiple orders
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
func (s *OrderService) Replace(orderId int64, useRemaining bool, newOrder SubmitOrder) (Order, error) {

	var side string
	if newOrder.Amount < 0 {
		newOrder.Amount = math.Abs(newOrder.Amount)
		side = "sell"
	} else {
		side = "buy"
	}

	payload := map[string]interface{}{
		"order_id":      strconv.FormatInt(orderId, 10),
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

// Retrieve the status of an order
func (s *OrderService) Status(orderId int64) (Order, error) {

	payload := map[string]interface{}{
		"order_id": orderId,
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
