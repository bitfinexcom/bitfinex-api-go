package bitfinex

import "time"

type HistoryService struct {
	client *Client
}

type Balance struct {
	Currency    string
	Amount      string
	Balance     string
	Description string
	Timestamp   string
}

func (s *HistoryService) Balance(currency, wallet string, since, until time.Time, limit int) ([]Balance, error) {

	payload := map[string]interface{}{"currency": currency}

	if !since.IsZero() {
		payload["since"] = since.Unix()
	}
	if !until.IsZero() {
		payload["until"] = until.Unix()
	}
	if limit != 0 {
		payload["limit"] = limit
	}

	req, err := s.client.newAuthenticatedRequest("POST", "history", payload)

	if err != nil {
		return nil, err
	}

	var v []Balance

	_, err = s.client.do(req, &v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

type Movement struct {
	ID          int64 `json:",int"`
	Currency    string
	Method      string
	Type        string
	Amount      string
	Description string
	Status      string
	Timestamp   string
}

func (s *HistoryService) Movements(currency, method string, since, until time.Time, limit int) ([]Movement, error) {

	payload := map[string]interface{}{"currency": currency, "method": method}

	if !since.IsZero() {
		payload["since"] = since.Unix()
	}
	if !until.IsZero() {
		payload["until"] = until.Unix()
	}
	if limit != 0 {
		payload["limit"] = limit
	}

	req, err := s.client.newAuthenticatedRequest("POST", "history/movements", payload)

	if err != nil {
		return nil, err
	}

	var v []Movement

	_, err = s.client.do(req, &v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

type PastTrade struct {
	Price       string
	Amount      string
	Timestamp   string
	Exchange    string
	Type        string
	FeeCurrency string `json:"fee_currency"`
	FeeAmount   string `json:"fee_amount"`
	TID         int64
	OrderId     int64 `json:"order_id,int"`
}

func (s *HistoryService) Trades(pair string, since, until time.Time, limit int, reverse bool) ([]PastTrade, error) {
	payload := map[string]interface{}{"symbol": pair}

	if !since.IsZero() {
		payload["timestamp"] = since.Unix()
	}
	if !until.IsZero() {
		payload["until"] = until.Unix()
	}
	if limit != 0 {
		payload["limit_trades"] = limit
	}
	if reverse {
		payload["reverse"] = 1
	}

	req, err := s.client.newAuthenticatedRequest("POST", "mytrades", payload)

	if err != nil {
		return nil, err
	}

	var v []PastTrade

	_, err = s.client.do(req, &v)

	if err != nil {
		return nil, err
	}

	return v, nil
}
