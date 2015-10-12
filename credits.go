package bitfinex

type CreditsService struct {
	client *Client
}

type Credit struct {
	Id        int
	Currency  string
	Status    string
	Rate      float64
	Period    float64
	Amount    float64
	Timestamp string
}

// GET /credits
func (c *CreditsService) All() ([]Credit, error) {
	req, err := c.client.NewAuthenticatedRequest("GET", "credits", nil)
	if err != nil {
		return nil, err
	}

	credits := make([]Credit, 0)
	_, err = c.client.Do(req, &credits)
	if err != nil {
		return nil, err
	}

	return credits, nil
}
