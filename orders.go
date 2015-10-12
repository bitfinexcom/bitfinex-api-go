package bitfinex

type OrderService struct {
	client *Client
}

type Order struct {
	Id                string
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
	ExecutedAmount    string `json:executed_amount`
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
func (s *OrderService) Create(pair string, amount float64, order_type string) error {
	req, err := s.client.NewAuthenticatedRequest("POST", "order/new", nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
