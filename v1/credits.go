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

// Returns an array of Credit
func (c *CreditsService) All() ([]Credit, error) {
	req, err := c.client.newAuthenticatedRequest("GET", "credits", nil)
	if err != nil {
		return nil, err
	}

	credits := make([]Credit, 0)
	_, err = c.client.do(req, &credits)
	if err != nil {
		return nil, err
	}

	return credits, nil
}
