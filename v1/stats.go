package bitfinex

import (
	"net/url"
	"strings"
)

type StatsService struct {
	client *Client
}

type Stats struct {
	Period int64
	Volume float64 `json:"volume,string"`
}

// All(pair) - Volume stats for specified pair
func (s *StatsService) All(pair string, period, volume string) ([]Stats, error) {
	pair = strings.ToUpper(pair)

	params := url.Values{}
	if period != "" {
		params.Add("period", period)
	}
	if volume != "" {
		params.Add("volume", volume)
	}
	req, err := s.client.newRequest("GET", "stats/"+strings.ToUpper(pair), params)

	if err != nil {
		return nil, err
	}

	var stats []Stats
	_, err = s.client.do(req, &stats)

	if err != nil {
		return nil, err
	}

	return stats, nil
}
