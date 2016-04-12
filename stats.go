package bitfinex

import "strings"

type StatsService struct {
	client *Client
}

type Stats struct {
	Period int64
	Volume float64 `json:"volume,string"`
}

// All(pair) - Volume stats for specified pair
func (s *StatsService) All(pair string) ([]Stats, error) {
	pair = strings.ToUpper(pair)
	req, err := s.client.NewRequest("GET", "stats/"+pair)

	if err != nil {
		return nil, err
	}

	var stats []Stats
	_, err = s.client.Do(req, &stats)

	if err != nil {
		return nil, err
	}

	return stats, nil
}
