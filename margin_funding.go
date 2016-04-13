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

    req, err := s.client.NewAuthenticatedRequest("POST", "offer/new", payload)

    if err != nil {
        return MarginOffer{}, err
    }

    var v MarginOffer
    _, err = s.client.Do(req, &v)

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
