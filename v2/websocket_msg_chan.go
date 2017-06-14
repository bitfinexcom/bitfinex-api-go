package bitfinex

import (
	"sync"
)

type msgChan struct {
	C    chan<- []byte
	mu   sync.Mutex
	c    chan []byte
	err  error
	done chan struct{}
}

func newMsgChan() *msgChan {
	m := &msgChan{
		c:    make(chan []byte),
		done: make(chan struct{}),
	}
	m.C = m.c

	return m
}

func (m *msgChan) Send(msg []byte) { m.c <- msg }

func (m *msgChan) Discard() {
	go func() {
		for _ = range m.c {
			// noop
		}
	}()
}

func (m *msgChan) Err() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.err
}

func (m *msgChan) Done() <-chan struct{} { return m.done }

func (m *msgChan) Receive() <-chan []byte { return m.c }

func (m *msgChan) Close(err error) {
	select { // Do nothing if we're already closed.
	default:
	case <-m.done:
		return
	}

	m.mu.Lock()
	m.err = err
	m.mu.Unlock()
	close(m.done)
	close(m.c)
}
