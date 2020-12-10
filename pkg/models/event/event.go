package event

type Event struct {
	Event     string       `json:"event,omitempty"`
	Channel   string       `json:"channel,omitempty"`
	ChanID    int64        `json:"chanId,omitempty"`
	Symbol    string       `json:"symbol,omitempty"`
	Precision string       `json:"prec,omitempty"`
	Frequency string       `json:"freq,omitempty"`
	Key       string       `json:"key,omitempty"`
	Len       string       `json:"len,omitempty"`
	Pair      string       `json:"pair,omitempty"`
	Code      int64        `json:"code,omitempty"`
	Version   int64        `json:"version,omitempty"`
	ServerID  string       `json:"serverId,omitempty"`
	Status    string       `json:"status"`
	UserID    int64        `json:"userId,omitempty"`
	SubID     string       `json:"subId"`
	AuthID    string       `json:"auth_id,omitempty"`
	Message   string       `json:"msg,omitempty"`
	Caps      Capabilities `json:"caps"`
	Platform  struct {
		Status int `json:"status,omitempty"`
	} `json:"platform"`
}

type Capability struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type Capabilities struct {
	Orders    Capability `json:"orders"`
	Account   Capability `json:"account"`
	Funding   Capability `json:"funding"`
	History   Capability `json:"history"`
	Wallets   Capability `json:"wallets"`
	Withdraw  Capability `json:"withdraw"`
	Positions Capability `json:"positions"`
}
