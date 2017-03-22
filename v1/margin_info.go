package bitfinex

type MarginInfoService struct {
	client *Client
}

type MarginInfo struct {
	MarginBalance     float64       `json:"margin_balance,string"`
	TradableBalance   float64       `json:"tradable_balance,string"`
	UnrealizedPl      float64       `json:"unrealized_pl,string"`
	UnrealizedSwap    float64       `json:"unrealized_swap,string"`
	NetValue          float64       `json:"net_value,string"`
	RequiredMargin    float64       `json:"required_margin,string"`
	Leverage          float64       `json:"leverage,string"`
	MarginRequirement float64       `json:"margin_requirement,string"`
	MarginLimits      []MarginLimit `json:"margin_limits,string"`
	Message           string        `json:"message"`
}

type MarginLimit struct {
	OnPair            string  `json:"on_pair"`
	InitialMargin     float64 `json:"initial_margin,string"`
	MarginRequirement float64 `json:"margin_requirement,string"`
	TradableBalance   float64 `json:"tradable_balance,string"`
}

// GET /margin_infos
func (s *MarginInfoService) All() ([]MarginInfo, error) {
	req, err := s.client.newAuthenticatedRequest("GET", "margin_infos", nil)
	if err != nil {
		return nil, err
	}

	var v []MarginInfo
	_, err = s.client.do(req, &v)

	return v, err
}
