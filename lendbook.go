package bitfinex

type LendbookService struct {
	client *Client
}

type Lendbook struct {
	Bids []struct {
		Rate      string
		Amount    string
		Period    float64
		Timestamp string
		Frr       string
	}
	Asks []struct {
		Rate      string
		Amount    string
		Period    float64
		Timestamp string
		Frr       string
	}
}

// TODO: Convert Rate and Amount to float64
// currency: BTC LTC DRK USD
// GET /lendbook/:currency
func (s *LendbookService) Get(currency string) (Lendbook, error) {
	req, err := s.client.NewRequest("GET", "lendbook/"+currency)
	if err != nil {
		return Lendbook{}, err
	}

	var v Lendbook
	_, err = s.client.Do(req, &v)
	if err != nil {
		return Lendbook{}, err
	}

	return v, nil
}

type Lends struct {
	Rate       string
	AmountLent string `json:"amount_lent"`
	AmountUsed string `json:"amount_used"`
	Timestamp  float64
}

// currency: BTC LTC DRK USD
// GET /lends/:currency
func (s *LendbookService) Lends(currency string) ([]Lends, error) {
	req, err := s.client.NewRequest("GET", "lends/"+currency)
	if err != nil {
		return nil, err
	}

	var v []Lends
	_, err = s.client.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
