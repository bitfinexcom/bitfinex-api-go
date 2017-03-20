package bitfinex

import (
	"strconv"
	"strings"
	"time"
)

type TickerService struct {
	client *Client
}

type Tick struct {
	Mid       string
	Bid       string
	Ask       string
	LastPrice string `json:"last_price"`
	Low       string
	High      string
	Volume    string
	Timestamp string
}

// ParseTime - return Timestamp in time.Time format
func (el *Tick) ParseTime() (*time.Time, error) {
	i, err := strconv.ParseFloat(el.Timestamp, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(i), 0)
	return &t, nil
}

// Get(pair) - return last Tick for specified pair
func (s *TickerService) Get(pair string) (Tick, error) {
	pair = strings.ToUpper(pair)
	req, err := s.client.newRequest("GET", "pubticker/"+pair, nil)

	if err != nil {
		return Tick{}, err
	}

	var v = &Tick{}
	_, err = s.client.do(req, v)

	if err != nil {
		return Tick{}, err
	}

	return *v, nil
}
