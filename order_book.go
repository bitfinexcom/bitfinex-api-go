package bitfinex

type OrderBookServive struct {
	client *Client
}

type OrderBook struct {
	Bids []struct {
		Price     string
		Amount    string
		Timestamp string
	}
	Asks []struct {
		Price     string
		Amount    string
		Timestamp string
	}
}

// TODO Convert price, amount to float64
// GET /book
func (s *OrderBookServive) Get(currency string) (OrderBook, error) {
	req, err := s.client.NewAuthenticatedRequest("GET", "book/"+currency, nil)
	if err != nil {
		return OrderBook{}, nil
	}

	var v OrderBook
	_, err = s.client.Do(req, &v)
	if err != nil {
		return OrderBook{}, err
	}

	return v, nil
}
