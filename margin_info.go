package bitfinex

import (
	"fmt"
)

type MarginInfoServive struct {
	client *Client
}

type MarginInfo struct {
	MarginBalance     string        `json:"margin_balance"`
	TradableBalance   string        `json:"tradable_balance"`
	UnrealizedPl      float64       `json:"unrealized_pl"`
	UnrealizedSwap    float64       `json:"unrealized_swap"`
	NetValue          string        `json:"net_value"`
	RequiredMargin    float64       `json:"required_margin"`
	Leverage          string        `json:"leverage"`
	MarginRequirement string        `json:"margin_requirement"`
	MarginLimits      []MarginLimit `json:"margin_limits"`
	Message           string        `json:"message"`
}

type MarginLimit struct {
	OnPair            string `json:"on_pair"`
	InitialMargin     string `json:"initial_margin"`
	MarginRequirement string `json:"margin_requirement"`
	TradableBalance   string `json:"tradable_balance"`
}

// GET /margin_infos
func (s *MarginInfoServive) All() ([]MarginInfo, error) {
	req, err := s.client.NewAuthenticatedRequest("GET", "margin_infos", nil)
	if err != nil {
		return make([]MarginInfo, 0), err
	}

	var v []MarginInfo
	resp, err := s.client.Do(req, &v)
	if err != nil {
		return make([]MarginInfo, 0), err
	}

	fmt.Println(v)
	println(resp)

	return make([]MarginInfo, 0), nil
}
