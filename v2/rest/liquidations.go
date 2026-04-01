package rest

// LiquidationsService manages the Liquidations endpoint.
type LiquidationsService struct {
	requestFactory
	Synchronous
}

// History - Get historical liquidation data
// see https://docs.bitfinex.com/reference/rest-public-liquidations
func (s *LiquidationsService) History() ([]interface{}, error) {
	req := NewRequestWithMethod("liquidations/hist", "GET")
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
