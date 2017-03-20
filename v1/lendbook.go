package bitfinex

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

type LendbookService struct {
	client *Client
}

type Lend struct {
	Rate      string
	Amount    string
	Period    int
	Timestamp string
	Frr       string
}

func (el *Lend) ParseTime() (*time.Time, error) {
	i, err := strconv.ParseFloat(el.Timestamp, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(i), 0)
	return &t, nil
}

type Lendbook struct {
	Bids []Lend
	Asks []Lend
}

// GET /lendbook/:currency
func (s *LendbookService) Get(currency string, limitBids, limitAsks int) (Lendbook, error) {
	currency = strings.ToUpper(currency)

	params := url.Values{}
	if limitBids != 0 {
		params.Add("limit_bids", strconv.Itoa(limitBids))
	}
	if limitAsks != 0 {
		params.Add("limit_asks", strconv.Itoa(limitAsks))
	}

	req, err := s.client.newRequest("GET", "lendbook/"+currency, params)
	if err != nil {
		return Lendbook{}, err
	}

	var v Lendbook
	_, err = s.client.do(req, &v)
	if err != nil {
		return Lendbook{}, err
	}

	return v, nil
}

type Lends struct {
	Rate       string
	AmountLent string `json:"amount_lent"`
	AmountUsed string `json:"amount_used"`
	Timestamp  int64
}

func (el *Lends) Time() *time.Time {
	t := time.Unix(el.Timestamp, 0)
	return &t
}

// GET /lends/:currency
func (s *LendbookService) Lends(currency string) ([]Lends, error) {
	currency = strings.ToUpper(currency)
	req, err := s.client.newRequest("GET", "lends/"+currency, nil)
	if err != nil {
		return nil, err
	}

	var v []Lends
	_, err = s.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
