package bitfinex

import (
	"net/url"
	"strconv"
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

func (s *TradesService) All(pair string, timestamp time.Time, limitTrades int) ([]Trade, error) {
	pair = strings.ToUpper(pair)

	params := url.Values{}
	if !time.Time.IsZero(timestamp) {
		params.Add("timestamp", strconv.FormatInt(timestamp.Unix(), 10))
	}
	if limitTrades != 0 {
		params.Add("limit_trades", strconv.Itoa(limitTrades))
	}
	req, err := s.client.newRequest("GET", "trades/"+pair, params)
	if err != nil {
		return nil, err
	}

	var v []Trade

	_, err = s.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
