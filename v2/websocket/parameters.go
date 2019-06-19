package websocket

import (
	"github.com/op/go-logging"
	"time"
)

// Parameters defines adapter behavior.
type Parameters struct {
	AutoReconnect          bool
	ReconnectInterval      time.Duration
	ReconnectAttempts      int
	reconnectTry           int
	ShutdownTimeout        time.Duration
	CapacityPerConnection  int
	Logger                 *logging.Logger

	ResubscribeOnReconnect bool

	HeartbeatTimeout       time.Duration
	LogTransport           bool

	URL                    string
	ManageOrderbook        bool
}

func NewDefaultParameters() *Parameters {
	return &Parameters{
		AutoReconnect:          true,
		CapacityPerConnection:  25,
		ReconnectInterval:      time.Second * 3,
		reconnectTry:           0,
		ReconnectAttempts:      15,
		URL:                    productionBaseURL,
		ManageOrderbook:        false,
		ShutdownTimeout:        time.Second * 5,
		ResubscribeOnReconnect: true,
		HeartbeatTimeout:       time.Second * 30,
		LogTransport:           false,           // log transport send/recv
		Logger:                 logging.MustGetLogger("bitfinex-ws"),
	}
}
