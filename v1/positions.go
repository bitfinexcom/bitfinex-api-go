package bitfinex

import (
	"strconv"
	"time"
)

// PositionsService structure
type PositionsService struct {
	client *Client
}

// Position structure
type Position struct {
	ID        int
	Symbol    string
	Amount    string
	Status    string
	Base      string
	Timestamp string
	Swap      string
	Pl        string
}

func (p *Position) ParseTime() (*time.Time, error) {
	i, err := strconv.ParseFloat(p.Timestamp, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(i), 0)
	return &t, nil
}

// All - gets all positions
func (b *PositionsService) All() ([]Position, error) {
	req, err := b.client.newAuthenticatedRequest("POST", "positions", nil)
	if err != nil {
		return nil, err
	}

	var positions []Position
	_, err = b.client.do(req, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

// Claim a position
func (b *PositionsService) Claim(positionId int, amount string) (Position, error) {

	request := map[string]interface{}{
		"position_id": positionId,
		"amount":      amount,
	}

	req, err := b.client.newAuthenticatedRequest("POST", "position/claim", request)

	if err != nil {
		return Position{}, err
	}

	var position = &Position{}

	_, err = b.client.do(req, position)

	if err != nil {
		return Position{}, err
	}

	return *position, nil
}
