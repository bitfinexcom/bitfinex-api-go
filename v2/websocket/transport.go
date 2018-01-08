package websocket

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

func newWs(baseURL string) *ws {
	return &ws{
		BaseURL:    baseURL,
		downstream: make(chan []byte),
		shutdown:   make(chan struct{}),
		finished:   make(chan error),
	}
}

type ws struct {
	ws            *websocket.Conn
	wsLock        sync.Mutex
	BaseURL       string
	TLSSkipVerify bool
	timeout       int64
	downstream    chan []byte
	userShutdown  bool

	shutdown chan struct{} // signal to kill looping goroutines
	finished chan error    // signal to parent with error, if applicable
}

func (w *ws) Connect() error {
	if w.ws != nil {
		return nil // no op
	}
	w.wsLock.Lock()
	defer w.wsLock.Unlock()
	w.userShutdown = false
	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: w.TLSSkipVerify}

	log.Printf("connecting ws to %s", w.BaseURL)
	ws, _, err := d.Dial(w.BaseURL, nil)
	if err != nil {
		return err
	}

	w.ws = ws
	go w.listenWs()
	return nil
}

// Send marshals the given interface and then sends it to the API. This method
// can block so specify a context with timeout if you don't want to wait for too
// long.
func (w *ws) Send(ctx context.Context, msg interface{}) error {
	if w.ws == nil {
		return ErrWSNotConnected
	}

	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.shutdown: // abrupt ws shutdown
		return fmt.Errorf("websocket connection closed")
	default:
	}

	w.wsLock.Lock()
	defer w.wsLock.Unlock()
	err = w.ws.WriteMessage(websocket.TextMessage, bs)
	if err != nil {
		w.cleanup(err)
		return err
	}

	return nil
}

func (w *ws) Done() <-chan error {
	return w.finished
}

// listen on ws & fwd to listen()
func (w *ws) listenWs() {
	for {
		if w.ws == nil {
			return
		}
		if atomic.LoadInt64(&w.timeout) != 0 {
			w.ws.SetReadDeadline(time.Now().Add(time.Duration(w.timeout)))
		}

		select {
		case <-w.shutdown: // external shutdown request
			return
		default:
		}

		_, msg, err := w.ws.ReadMessage()
		if err != nil {
			w.cleanup(err)
			return
		}
		log.Print(string(msg))
		w.downstream <- msg
	}
}

func (w *ws) Listen() <-chan []byte {
	return w.downstream
}

func (w *ws) cleanup(err error) {
	close(w.downstream) // shut down caller's listen channel
	close(w.shutdown)   // signal to kill goroutines
	if err != nil && !w.userShutdown {
		w.finished <- err
	}
	close(w.finished) // signal to parent listeners
}

// Close the websocket connection
func (w *ws) Close() {
	w.wsLock.Lock()
	w.userShutdown = true
	if w.ws != nil {
		if err := w.ws.Close(); err != nil { // will trigger cleanup()
			log.Printf("[INFO]: error closing websocket: %s", err)
		}
		w.ws = nil
	}
	w.wsLock.Unlock()
}
