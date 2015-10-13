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

// GET orders
func (s *OrderService) All() ([]Order, error) {
	req, err := s.client.NewAuthenticatedRequest("GET", "orders", nil)
	if err != nil {
		return nil, err
	}

	v := []Order{}
	_, err = s.client.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// POST order/cancel/all
func (s *OrderService) CancelAll() error {
	req, err := s.client.NewAuthenticatedRequest("POST", "order/cancel/all", nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// POST order/new
// symbol   # The name of the symbol (see `/symbols`).
// price    # Price to buy or sell at. May omit if a market order.
// amount   # Order size: how much to buy or sell. Use negative amount to create sell order.
// side     # Either "buy" or "sell".
// type     # Either "market" / "limit" / "stop" / "trailing-stop" / "fill-or-kill" /
//                   "exchange market" / "exchange limit" / "exchange stop" /
//                   "exchange trailing-stop" / "exchange fill-or-kill"
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
		"type":     "exchange limit",
		"exchange": "bitfinex",
	}

	req, err := s.client.NewAuthenticatedRequest("POST", "order/new", payload)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	_, err = s.client.Do(req, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
