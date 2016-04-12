package bitfinex

type PairsService struct {
	client *Client
}

type Pairs []string

// GET /symbols
func (p *PairsService) All() (*Pairs, error) {
	req, err := p.client.NewRequest("GET", "symbols")
	if err != nil {
		return nil, err
	}

	v := &Pairs{}
	_, err = p.client.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type Pair struct {
	Pair             string
	PricePrecision   int     `json:"price_precision,int"`
	InitialMargin    float64 `json:"initial_margin,string"`
	MinimumMargin    float64 `json:"minimum_margin,string"`
	MaximumOrderSize float64 `json:"maximum_order_size,string"`
	MinimumOrderSize float64 `json:"minimum_order_size,string"`
	Espiration       string
}

// GET /symbols
func (p *PairsService) AllDetailed() ([]Pair, error) {
	req, err := p.client.NewRequest("GET", "symbols_details")
	if err != nil {
		return nil, err
	}

	var v []Pair
	_, err = p.client.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
