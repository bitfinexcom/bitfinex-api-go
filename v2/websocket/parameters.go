package websocket

import (
	"time"
)

// Parameters defines adapter behavior.
type Parameters struct {
	// TODO implement failover hierarchy
	/*
		AutoResubscribe             bool
		ResubscribeAttempts         int
		ResubscribeHeartbeatTimeout time.Duration
		resubscribeTry              int
	*/
	AutoReconnect     bool
	ReconnectInterval time.Duration
	ReconnectAttempts int
	reconnectTry      int

	URL string
}

func NewDefaultParameters() *Parameters {
	return &Parameters{
		AutoReconnect:     true,
		ReconnectInterval: time.Second,
		reconnectTry:      0,
		ReconnectAttempts: 5,
		URL:               productionBaseURL,
	}
}

func (p *Parameters) SetAutoReconnect(auto bool) *Parameters {
	p.AutoReconnect = auto
	return p
}

func (p *Parameters) SetReconnectInterval(t time.Duration) *Parameters {
	p.ReconnectInterval = t
	return p
}

func (p *Parameters) SetReconnectAttempts(attempts int) *Parameters {
	p.ReconnectAttempts = attempts
	return p
}

func (p *Parameters) SetURL(url string) *Parameters {
	p.URL = url
	return p
}
