package bitfinex

import (
	"strings"
	"time"
)

type TradesService struct {
	client *Client
}

type Trade struct {
	Price     string
	Amount    string
	Exchange  string
	Type      string
	Timestamp int64
	TradeId   int64 `json:"tid,int"`
}

func (el *Trade) Time() *time.Time {
	t := time.Unix(el.Timestamp, 0)
	return &t
}

func (s *TradesService) All(pair string) ([]Trade, error) {
	pair = strings.ToUpper(pair)
	req, err := s.client.NewRequest("GET", "trades/"+pair)
	if err != nil {
		return nil, err
	}

	var v []Trade

	_, err = s.client.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
