package bitfinex

type BalancesService struct {
	client *Client
}

type WalletBalance struct {
	Type      string
	Currency  string
	Amount    string
	Available string
}

// GET balances
func (b *BalancesService) All() ([]WalletBalance, error) {
	req, err := b.client.newAuthenticatedRequest("GET", "balances", nil)
	if err != nil {
		return nil, err
	}

	balances := make([]WalletBalance, 3)
	_, err = b.client.do(req, &balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
