package bitfinex

import "strconv"

type MarginFundingService struct {
	client *Client
}

type MarginOffer struct {
	ID              int64
	Currency        string
	Rate            string
	Period          int
	Direction       string
	Timestamp       string
	IsLive          bool   `json:"is_live"`
	IsCancelled     bool   `json:"is_cancelled"`
	OriginalAmount  string `json:"original_amount"`
	RemainingAmount string `json:"remaining_amount"`
	ExecutedAmount  string `json:"executed_amount"`
	OfferId         int
}

func (s *MarginFundingService) new(currency, direction string, amount, rate float64, period int) (MarginOffer, error) {
	payload := map[string]interface{}{
		"currency":  currency,
		"amount":    strconv.FormatFloat(amount, 'f', -1, 32),
		"rate":      strconv.FormatFloat(rate, 'f', -1, 32),
		"period":    period,
		"direction": direction,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "offer/new", payload)

	if err != nil {
		return MarginOffer{}, err
	}

	var v MarginOffer
	_, err = s.client.do(req, &v)

	if err != nil {
		return MarginOffer{}, err
	}

	return v, nil
}

func (s *MarginFundingService) NewLend(currency string, amount, rate float64, period int) (MarginOffer, error) {
	return s.new(currency, "lend", amount, rate, period)
}
func (s *MarginFundingService) NewLoan(currency string, amount, rate float64, period int) (MarginOffer, error) {
	return s.new(currency, "loan", amount, rate, period)
}

func (s *MarginFundingService) Cancel(offerId int64) (MarginOffer, error) {
	payload := map[string]interface{}{"offer_id": offerId}

	req, err := s.client.newAuthenticatedRequest("POST", "offer/cancel", payload)

	if err != nil {
		return MarginOffer{}, err
	}

	var v MarginOffer
	_, err = s.client.do(req, &v)

	if err != nil {
		return MarginOffer{}, err
	}
	return v, nil
}

func (s *MarginFundingService) Status(offerId int64) (MarginOffer, error) {
	payload := map[string]interface{}{"offer_id": offerId}

	req, err := s.client.newAuthenticatedRequest("POST", "offer/status", payload)

	if err != nil {
		return MarginOffer{}, err
	}

	var v MarginOffer
	_, err = s.client.do(req, &v)

	if err != nil {
		return MarginOffer{}, err
	}
	return v, nil
}

type ActiveOffer struct {
	ID              int64
	Currency        string
	Rate            string
	Period          int
	Direction       string
	Timestamp       string
	IsLive          bool   `json:"is_live"`
	IsCancelled     bool   `json:"is_cancelled"`
	OriginalAmount  string `json:"original_amount"`
	RemainingAmount string `json:"remaining_amount"`
	ExecutedAmount  string `json:"executed_amount"`
}

func (s *MarginFundingService) Credits() ([]ActiveOffer, error) {

	req, err := s.client.newAuthenticatedRequest("POST", "credits", nil)

	if err != nil {
		return nil, err
	}

	var v []ActiveOffer
	_, err = s.client.do(req, &v)

	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *MarginFundingService) Offers() ([]ActiveOffer, error) {

	req, err := s.client.newAuthenticatedRequest("POST", "offers", nil)

	if err != nil {
		return nil, err
	}

	var v []ActiveOffer
	_, err = s.client.do(req, &v)

	if err != nil {
		return nil, err
	}
	return v, nil
}
