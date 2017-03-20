package bitfinex

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

type OrderBookService struct {
	client *Client
}

type OrderBookEntry struct {
	Price     string
	Rate      string
	Amount    string
	Period    int
	Timestamp string
	Frr       string
}

type OrderBook struct {
	Bids []OrderBookEntry
	Asks []OrderBookEntry
}

func (el *OrderBookEntry) ParseTime() (*time.Time, error) {
	i, err := strconv.ParseFloat(el.Timestamp, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(i), 0)
	return &t, nil
}

// GET /book
func (s *OrderBookService) Get(pair string, limitBids, limitAsks int, noGroup bool) (OrderBook, error) {
	pair = strings.ToUpper(pair)

	params := url.Values{}
	if limitBids != 0 {
		params.Add("limit_bids", strconv.Itoa(limitBids))
	}
	if limitAsks != 0 {
		params.Add("limit_asks", strconv.Itoa(limitAsks))
	}
	if noGroup {
		params.Add("group", "0")
	}

	req, err := s.client.newRequest("GET", "book/"+pair, params)

	if err != nil {
		return OrderBook{}, err
	}

	var v OrderBook
	_, err = s.client.do(req, &v)

	if err != nil {
		return OrderBook{}, err
	}

	return v, nil
}
