package bitfinex

import (
    "fmt"
    "math"
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
