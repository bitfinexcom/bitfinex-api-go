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

    req, err := s.client.NewAuthenticatedRequest("POST", "history", payload)

    if err != nil {
        return nil, err
    }

    var v []Balance

    _, err = s.client.Do(req, &v)

    if err != nil {
        return nil, err
    }

    return v, nil
}
