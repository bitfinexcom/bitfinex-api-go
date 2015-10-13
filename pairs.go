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
