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
	AutoReconnect     bool
	ReconnectInterval time.Duration
	ReconnectAttempts int
	reconnectTry      int
	ShutdownTimeout   time.Duration

	ResubscribeOnReconnect bool

	HeartbeatTimeout time.Duration

	URL string
}

func NewDefaultParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		ReconnectInterval:      time.Second,
		reconnectTry:           0,
		ReconnectAttempts:      5,
		URL:                    productionBaseURL,
		ShutdownTimeout:        time.Millisecond * 2500,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Millisecond * 3750, // HB ~ 2.5s, timeout = 3/2*HB
	}
}
