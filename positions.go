package bitfinex

// PositionsService structure
type PositionsService struct {
	client *Client
}

// Position structure
type Position struct {
	ID        int
	Symbol    string
	Amount    float64 `json:",string"`
	Status    string
	Base      float64 `json:",string"`
	Timestamp float64 `json:",string"`
	Swap      float64 `json:",string"`
	Pl        float64 `json:",string"`
}

// All - gets all positions
func (b *PositionsService) All() ([]Position, error) {
	req, err := b.client.NewAuthenticatedRequest("GET", "positions", nil)
	if err != nil {
		return nil, err
	}

	var positions []Position
	_, err = b.client.Do(req, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}
