package rest

import (
	"fmt"
	"net/url"
)

// RankingsService manages the Rankings/Leaderboard endpoint.
type RankingsService struct {
	requestFactory
	Synchronous
}

// History - Get leaderboard standings
// key: "plu_diff" (unrealized P/L), "plu_perc" (% change), "vol" (volume), "plr" (realized P/L)
// timeFrame: "3h", "1w", "1M"
// symbol: e.g. "tBTCUSD"
// see https://docs.bitfinex.com/reference/rest-public-rankings
func (s *RankingsService) History(key, timeFrame, symbol string, limit int) ([]interface{}, error) {
	endpoint := fmt.Sprintf("rankings/%s:%s:%s/hist", key, timeFrame, symbol)
	req := NewRequestWithMethod(endpoint, "GET")
	if limit > 0 {
		req.Params = make(url.Values)
		req.Params.Add("limit", fmt.Sprintf("%d", limit))
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
