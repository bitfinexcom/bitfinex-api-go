package bitfinex

type AccountService struct {
	client *Client
}

type AccountPairFee struct {
	Pair      string
	MakerFees float64 `json:"maker_fees,string"`
	TakerFees float64 `json:"taker_fees,string"`
}

type AccountInfo struct {
	MakerFees float64 `json:"maker_fees,string"`
	TakerFees float64 `json:"taker_fees,string"`
	Fees      []AccountPairFee
}

// GET account_infos
func (a *AccountService) Info() (AccountInfo, error) {
	req, err := a.client.newAuthenticatedRequest("GET", "account_infos", nil)

	if err != nil {
		return AccountInfo{}, err
	}

	var v []AccountInfo
	_, err = a.client.do(req, &v)

	if err != nil {
		return AccountInfo{}, err
	}

	return v[0], nil
}

type KeyPerm struct {
	Read  bool
	Write bool
}

type Permissions struct {
	Account   KeyPerm
	History   KeyPerm
	Orders    KeyPerm
	Positions KeyPerm
	Funding   KeyPerm
	Wallets   KeyPerm
	Withdraw  KeyPerm
}

func (a *AccountService) KeyPermission() (Permissions, error) {
	req, err := a.client.newAuthenticatedRequest("GET", "key_info", nil)

	if err != nil {
		return Permissions{}, err
	}

	var v Permissions
	_, err = a.client.do(req, &v)
	if err != nil {
		return Permissions{}, err
	}
	return v, nil
}

type SummaryVolume struct {
	Currency string `json:"curr"`
	Volume   string `json:"vol"`
}
type SummaryProfit struct {
	Currency string `json:"curr"`
	Volume   string `json:"amount"`
}
type Summary struct {
	TradeVolume   SummaryVolume `json:"trade_vol_30d"`
	FundingProfit SummaryProfit `json:"funding_profit_30d"`
	MakerFee      string        `json:"maker_fee"`
	TakerFee      string        `json:"taker_fee"`
}

func (a *AccountService) Summary() (Summary, error) {
	req, err := a.client.newAuthenticatedRequest("GET", "summary", nil)

	if err != nil {
		return Summary{}, err
	}

	var v Summary
	_, err = a.client.do(req, &v)
	if err != nil {
		return Summary{}, err
	}
	return v, nil
}
