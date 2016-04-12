package bitfinex

import (
	"strconv"
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

func (el *Tick) ParseTime() (*time.Time, error) {
	i, err := strconv.ParseFloat(el.Timestamp, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(i), 0)
	return &t, nil
}

func (s *TickerService) Get(pair string) (Tick, error) {
	req, err := s.client.NewRequest("GET", "pubticker/"+pair)

	if err != nil {
		return Tick{}, err
	}

	var v = &Tick{}
	_, err = s.client.Do(req, v)

	if err != nil {
		return Tick{}, err
	}

	return *v, nil
}
