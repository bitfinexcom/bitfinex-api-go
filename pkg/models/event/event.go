package event

type Subscribe struct {
	Event     string `json:"event,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Len       string `json:"len,omitempty"`
	Key       string `json:"key,omitempty"`
	// auth related part
	APIKEY      string `json:"apiKey,omitempty"`
	AuthSig     string `json:"authSig,omitempty"`
	AuthPayload string `json:"authPayload,omitempty"`
	AuthNonce   string `json:"authNonce,omitempty"`
}

type Info struct {
	Subscribe
	ChanID   int64        `json:"chanId,omitempty"`
	Pair     string       `json:"pair,omitempty"`
	Code     int64        `json:"code,omitempty"`
	Version  int64        `json:"version,omitempty"`
	ServerID string       `json:"serverId,omitempty"`
	Status   string       `json:"status,omitempty"`
	UserID   int64        `json:"userId,omitempty"`
	SubID    string       `json:"subId,omitempty"`
	AuthID   string       `json:"auth_id,omitempty"`
	Message  string       `json:"msg,omitempty"`
	Caps     Capabilities `json:"caps,omitempty"`
	Platform struct {
		Status int `json:"status,omitempty"`
	} `json:"platform,omitempty"`
}

type Capability struct {
	Read  int `json:"read,omitempty"`
	Write int `json:"write,omitempty"`
}

type Capabilities struct {
	Orders    Capability `json:"orders,omitempty"`
	Account   Capability `json:"account,omitempty"`
	Funding   Capability `json:"funding,omitempty"`
	History   Capability `json:"history,omitempty"`
	Wallets   Capability `json:"wallets,omitempty"`
	Withdraw  Capability `json:"withdraw,omitempty"`
	Positions Capability `json:"positions,omitempty"`
}
