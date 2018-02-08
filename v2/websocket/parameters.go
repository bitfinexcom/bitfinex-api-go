package websocket

import (
	"time"
)

// Parameters defines adapter behavior.
type Parameters struct {
	// TODO implement failover hierarchy
	/*
		AutoResubscribe bool
		ResubscribeAttempts         int
		ResubscribeHeartbeatTimeout time.Duration
		resubscribeTry              int
	*/
	autoReconnect     bool
	reconnectInterval time.Duration
	reconnectAttempts int
	reconnectTry      int
	shutdownTimeout   time.Duration

	ResubscribeOnReconnect bool

	url string
}

func NewDefaultParameters() *Parameters {
	return &Parameters{
		autoReconnect:          true,
		reconnectInterval:      time.Second,
		reconnectTry:           0,
		reconnectAttempts:      5,
		url:                    productionBaseURL,
		shutdownTimeout:        time.Millisecond * 2500,
		ResubscribeOnReconnect: true,
	}
}

func (p *Parameters) SetAutoReconnect(auto bool) *Parameters {
	p.autoReconnect = auto
	return p
}

func (p *Parameters) SetReconnectInterval(t time.Duration) *Parameters {
	p.reconnectInterval = t
	return p
}

func (p *Parameters) SetReconnectAttempts(attempts int) *Parameters {
	p.reconnectAttempts = attempts
	return p
}

func (p *Parameters) SetURL(url string) *Parameters {
	p.url = url
	return p
}
