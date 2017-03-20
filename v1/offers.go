package bitfinex

import "strconv"

type OffersService struct {
	client *Client
}

const (
	LEND = "lend"
	LOAN = "loan"
)

type Offer struct {
	Id              int64
	Currency        string
	Rate            string
	Period          int64
	Direction       string
	Timestamp       string
	IsLive          bool   `json:"is_live"`
	IsCancelled     bool   `json:"is_cancelled"`
	OriginalAmount  string `json:"original_amount:string"`
	RemainingAmount string `json:"remaining_amount:string"`
	ExecutedAmount  string `json:"executed_amount:string"`
	OfferId         int64  `json:"offer_id"`
}

// Create new offer for LEND or LOAN a currency, use LEND or LOAN constants as direction
func (s *OffersService) New(currency string, amount, rate float64, period int64, direction string) (Offer, error) {

	payload := map[string]interface{}{
		"currency":  currency,
		"amount":    strconv.FormatFloat(amount, 'f', -1, 32),
		"rate":      strconv.FormatFloat(rate, 'f', -1, 32),
		"period":    strconv.FormatInt(period, 10),
		"direction": direction,
	}

	req, err := s.client.newAuthenticatedRequest("POST", "offers/new", payload)

	if err != nil {
		return Offer{}, err
	}

	var offer = &Offer{}
	_, err = s.client.do(req, offer)

	if err != nil {
		return Offer{}, err
	}

	return *offer, nil

}

func (s *OffersService) Cancel(offerId int64) (Offer, error) {

	payload := map[string]interface{}{
		"offer_id": strconv.FormatInt(offerId, 10),
	}

	req, err := s.client.newAuthenticatedRequest("POST", "offers/cancel", payload)

	if err != nil {
		return Offer{}, err
	}

	var offer = &Offer{}

	_, err = s.client.do(req, offer)

	if err != nil {
		return Offer{}, err
	}

	return *offer, nil
}

func (s *OffersService) Status(offerId int64) (Offer, error) {

	payload := map[string]interface{}{
		"offer_id": strconv.FormatInt(offerId, 10),
	}

	req, err := s.client.newAuthenticatedRequest("POST", "offers/status", payload)

	if err != nil {
		return Offer{}, err
	}

	var offer = &Offer{}

	_, err = s.client.do(req, offer)

	if err != nil {
		return Offer{}, err
	}

	return *offer, nil
}
