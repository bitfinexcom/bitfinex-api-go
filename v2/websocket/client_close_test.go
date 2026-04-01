package websocket

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/op/go-logging"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/utils"
)

// mockAsync is a minimal Asynchronous implementation for testing.
type mockAsync struct {
	done     chan error
	quit     chan error
	listen   chan []byte
	mu       sync.Mutex
	closed   bool
}

func newMockAsync() *mockAsync {
	return &mockAsync{
		done:   make(chan error, 1),
		quit:   make(chan error, 1),
		listen: make(chan []byte, 10),
	}
}

func (m *mockAsync) Connect() error { return nil }
func (m *mockAsync) Send(ctx context.Context, msg interface{}) error { return nil }
func (m *mockAsync) Listen() <-chan []byte { return m.listen }
func (m *mockAsync) Done() <-chan error { return m.quit }
func (m *mockAsync) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		m.closed = true
		// Signal done to any listener waiting on Done()
		select {
		case m.quit <- nil:
		default:
		}
	}
}

type mockAsyncFactory struct{}

func (f *mockAsyncFactory) Create() Asynchronous {
	return newMockAsync()
}

func newTestClient() *Client {
	log := logging.MustGetLogger("test")
	params := NewDefaultParameters()
	params.Logger = log
	params.ShutdownTimeout = 1 * time.Second

	return &Client{
		asyncFactory:   &mockAsyncFactory{},
		Authentication: NoAuthentication,
		factories:      make(map[string]messageFactory),
		subscriptions:  newSubscriptions(params.HeartbeatTimeout, params.Logger),
		orderbooks:     make(map[string]*Orderbook),
		nonce:          utils.NewEpochNonceGenerator(),
		parameters:     params,
		listener:       make(chan interface{}),
		terminal:       false,
		shutdown:       nil,
		sockets:        make(map[SocketId]*Socket),
		mtx:            &sync.RWMutex{},
		log:            log,
	}
}

func TestClient_DoubleCloseNoPanic(t *testing.T) {
	c := newTestClient()
	// No connected sockets — test that double-close on listener channel
	// doesn't panic thanks to sync.Once

	// First close
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("first Close() panicked: %v", r)
			}
		}()
		c.Close()
	}()

	// Second close should not panic (sync.Once protects against double close)
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("second Close() panicked: %v", r)
			}
		}()
		c.Close()
	}()
}

func TestClient_CloseWithConnectedSocket(t *testing.T) {
	c := newTestClient()
	c.parameters.ShutdownTimeout = 2 * time.Second

	mock := newMockAsync()
	socket := &Socket{
		Id:           0,
		Asynchronous: mock,
		IsConnected:  true,
	}
	c.sockets[0] = socket

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.Close()
	}()

	select {
	case <-done:
		// ok
	case <-time.After(5 * time.Second):
		t.Fatal("Close() timed out")
	}

	if mock.mu.Lock(); !mock.closed {
		t.Error("expected mock async to be closed")
	}
	mock.mu.Unlock()
}

func TestClient_CloseWithNoSockets(t *testing.T) {
	c := newTestClient()

	// Should not panic with empty socket map
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Close() with no sockets panicked: %v", r)
			}
		}()
		c.Close()
	}()
}
