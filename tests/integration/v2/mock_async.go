package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

// does not work for reconnect tests
type TestAsyncFactory struct {
	Count int
	Async websocket.Asynchronous
}

func (t *TestAsyncFactory) Create() websocket.Asynchronous {
	t.Count += 1
	// if first creation then send given async
	if t.Count == 1 {
		return t.Async
	}
	// otherwise create a new async for each new creation
	return newTestAsync()
}

func newTestAsyncFactory(async websocket.Asynchronous) websocket.AsynchronousFactory {
	return &TestAsyncFactory{Async: async, Count: 0}
}

type TestAsync struct {
	done      chan error
	bridge    chan []byte
	connected bool
	Sent      []interface{}
	mutex     sync.Mutex
}

func (t *TestAsync) SentCount() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return len(t.Sent)
}

func (t *TestAsync) waitForMessage(num int) error {
	seconds := 4
	loops := 20
	delay := time.Duration(float64(time.Second) * float64(seconds) / float64(loops))
	for i := 0; i < loops; i++ {
		t.mutex.Lock()
		len := len(t.Sent)
		t.mutex.Unlock()
		if num+1 <= len {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("did not send a message in pos %d", num)
}

func (t *TestAsync) Connect() error {
	t.connected = true
	return nil
}

func (t *TestAsync) Send(ctx context.Context, msg interface{}) error {
	if !t.connected {
		return errors.New("must connect before sending")
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Sent = append(t.Sent, msg)
	return nil
}

func (t *TestAsync) DumpSentMessages() {
	for i, msg := range t.Sent {
		log.Printf("%2d: %#v", i, msg)
	}
}

func (t *TestAsync) Listen() <-chan []byte {
	return t.bridge
}

func (t *TestAsync) Publish(raw string) {
	t.bridge <- []byte(raw)
}

func (t *TestAsync) Close() {
	close(t.bridge)
	close(t.done)
}

func (t *TestAsync) Done() <-chan error {
	return t.done
}

func newTestAsync() *TestAsync {
	return &TestAsync{
		bridge:    make(chan []byte),
		connected: false,
		Sent:      make([]interface{}, 0),
		done:      make(chan error),
	}
}
