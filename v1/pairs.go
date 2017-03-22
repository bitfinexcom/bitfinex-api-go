package bitfinex

type PairsService struct {
	client *Client
}

// Get all Pair names as array of strings
func (p *PairsService) All() ([]string, error) {
	req, err := p.client.newRequest("GET", "symbols", nil)
	if err != nil {
		return nil, err
	}

	var v []string

	_, err = p.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Detailed Pair
type Pair struct {
	Pair             string
	PricePrecision   int     `json:"price_precision,int"`
	InitialMargin    float64 `json:"initial_margin,string"`
	MinimumMargin    float64 `json:"minimum_margin,string"`
	MaximumOrderSize float64 `json:"maximum_order_size,string"`
	MinimumOrderSize float64 `json:"minimum_order_size,string"`
	Espiration       string
}

// Return a list of detailed pairs
func (p *PairsService) AllDetailed() ([]Pair, error) {
	req, err := p.client.newRequest("GET", "symbols_details", nil)
	if err != nil {
		return nil, err
	}

	var v []Pair
	_, err = p.client.do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
