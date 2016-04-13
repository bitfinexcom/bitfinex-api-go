package bitfinex

type AccountService struct {
    client *Client
}

// TODO return struct
// GET account_infos
func (a *AccountService) Info() (string, error) {
    req, err := a.client.NewAuthenticatedRequest("GET", "account_infos", nil)
    if err != nil {
        return "", err
    }

    resp, err := a.client.Do(req, nil)
    if err != nil {
        return "", err
    }

    return resp.String(), nil
}
