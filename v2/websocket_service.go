package bitfinex

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/bitfinexcom/bitfinex-api-go/utils"

	"github.com/gorilla/websocket"
)

// Available channels
const (
	ChanBook    = "book"
	ChanTrades  = "trades"
	ChanTicker  = "ticker"
	ChanCandles = "candles"
)

// Prefixes for available pairs
const (
	FundingPrefix = "f"
	TradingPrefix = "t"
)

var (
	ErrWSNotConnected     = fmt.Errorf("websocket connection not established")
	ErrWSAlreadyConnected = fmt.Errorf("websocket connection already established")
)

// bfxWebsocket is a wrapper around a simple websocket connection, that let's us
// manage callbacks and share a single websocket in a thread safe manner.
// It provides a single channel to write message to.
type bfxWebsocket struct {
	// Bitfinex client
	client *Client

	wsMu         sync.Mutex
	ws           *websocket.Conn
	timeout      int64
	webSocketURL string

	// TLSSkipVerify toggles if certificate verification should be skipped or not.
	TLSSkipVerify bool

	// The bitfinex API sends us untyped arrays as data, so we have to keep track
	// of which one belongs where.
	subMu       sync.Mutex
	pubSubIDs   map[string]publicSubInfo
	pubChanIDs  map[int64]PublicSubscriptionRequest // ChannelID -> SubscriptionRequest map
	privSubIDs  map[string]struct{}
	privChanIDs map[int64]struct{}

	eventHandler handlerT

	privateHandler  handlerT
	isAuthenticated bool

	handlersMu     sync.Mutex
	publicHandlers map[int64]handlerT

	done  chan struct{}
	errMu sync.Mutex
	err   error
}

type handlerT func(interface{})

type publicSubInfo struct {
	req PublicSubscriptionRequest
	h   handlerT
}

func newBfxWebsocket(c *Client, wsURL string) *bfxWebsocket {
	b := &bfxWebsocket{
		client:       c,
		webSocketURL: wsURL,
	}
	b.init()

	return b
}

func (b *bfxWebsocket) Connect() error {
	if b.ws != nil {
		return nil // We're already connected.
	}

	b.init()
	return b.connect()
}

func (b *bfxWebsocket) connect() error {
	b.wsMu.Lock()
	defer b.wsMu.Unlock()
	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: b.TLSSkipVerify}

	ws, _, err := d.Dial(b.webSocketURL, nil)
	if err != nil {
		return err
	}

	b.ws = ws

	go b.receiver()

	return nil
}

func (b *bfxWebsocket) init() {
	b.privSubIDs = map[string]struct{}{}
	b.pubSubIDs = map[string]publicSubInfo{}
	b.pubChanIDs = map[int64]PublicSubscriptionRequest{}
	b.publicHandlers = map[int64]handlerT{}
	b.privChanIDs = map[int64]struct{}{}
	b.done = make(chan struct{})
}

func (b *bfxWebsocket) receiver() {
	for {
		if b.ws == nil {
			return
		}
		if atomic.LoadInt64(&b.timeout) != 0 {
			b.ws.SetReadDeadline(time.Now().Add(time.Duration(b.timeout)))
		}

		select {
		case <-b.Done():
			return
		default:
		}

		_, msg, err := b.ws.ReadMessage()
		if err != nil {
			b.close(err)
			return
		}

		// Errors here should be non critical so we just log them.
		err = b.handleMessage(msg)
		if err != nil {
			log.Printf("[WARN]: %s\n", err)
		}
	}
}

// Done returns a channel that will be closed if the underlying websocket
// connection gets closed.
func (b *bfxWebsocket) Done() <-chan struct{} { return b.done }

// Err returns an error if the done channel was closed due to an error.
func (b *bfxWebsocket) Err() error {
	b.errMu.Lock()
	defer b.errMu.Unlock()
	return b.err
}

func (b *bfxWebsocket) close(e error) {
	b.wsMu.Lock()
	if b.ws != nil {
		if err := b.ws.Close(); err != nil {
			log.Printf("[INFO]: error closing websocket: %s", err)
		}
		b.ws = nil
	}
	b.wsMu.Unlock()

	b.errMu.Lock()
	b.err = e
	b.errMu.Unlock()

	select { // Do nothing if we're already closed.
	default:
	case <-b.done:
		return
	}

	close(b.done)
}

func (b *bfxWebsocket) Close() {
	b.close(nil)
}

// Send marshals the given interface and then sends it to the API. This method
// can block so specify a context with timeout if you don't want to wait for too
// long.
func (b *bfxWebsocket) Send(ctx context.Context, msg interface{}) error {
	if b.ws == nil {
		return ErrWSNotConnected
	}

	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.Done():
		return fmt.Errorf("websocket closed: ", b.Err())
	default:
	}

	b.wsMu.Lock()
	defer b.wsMu.Unlock()
	err = b.ws.WriteMessage(websocket.TextMessage, bs)
	if err != nil { // If WriteMessage returns an error, it's permanent.
		b.close(err)
		return err
	}

	return nil
}

func (b *bfxWebsocket) handleMessage(msg []byte) error {
	t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
	if bytes.HasPrefix(t, []byte("[")) { // Data messages are just arrays of values.
		var raw []interface{}
		err := json.Unmarshal(msg, &raw)
		if err != nil {
			return err
		} else if len(raw) < 2 {
			return nil
		}

		chanID, ok := raw[0].(float64)
		if !ok {
			return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
		}

		if _, has := b.privChanIDs[int64(chanID)]; has {
			td, err := b.handlePrivateDataMessage(raw)
			if err != nil {
				return err
			} else if td == nil {
				return nil
			}
			if b.privateHandler != nil {
				go b.privateHandler(td)
				return nil
			}
		} else if _, has := b.pubChanIDs[int64(chanID)]; has {
			td, err := b.handlePublicDataMessage(raw)
			if err != nil {
				return err
			} else if td == nil {
				return nil
			}
			if h, has := b.publicHandlers[int64(chanID)]; has {
				go h(td)
				return nil
			}
		} else {
			// TODO: log unhandled message?
		}
	} else if bytes.HasPrefix(t, []byte("{")) { // Events are encoded as objects.
		ev, err := b.onEvent(msg)
		if err != nil {
			return err
		}
		if b.eventHandler != nil {
			go b.eventHandler(ev)
		}
	} else {
		return fmt.Errorf("unexpected message: %s", msg)
	}

	return nil
}

type subscriptionRequest struct {
	Event       string   `json:"event"`
	APIKey      string   `json:"apiKey"`
	AuthSig     string   `json:"authSig"`
	AuthPayload string   `json:"authPayload"`
	AuthNonce   string   `json:"authNonce"`
	Filter      []string `json:"filter"`
	SubID       string   `json:"subId"`
}

// Authenticate creates the payload for the authentication request and sends it
// to the API. The filters will be applied to the authenticated channel, i.e.
// only subscribe to the filtered messages.
func (b *bfxWebsocket) Authenticate(ctx context.Context, filter ...string) error {
	nonce := utils.GetNonce()
	payload := "AUTH" + nonce
	s := &subscriptionRequest{
		Event:       "auth",
		APIKey:      b.client.APIKey,
		AuthSig:     b.client.sign(payload),
		AuthPayload: payload,
		AuthNonce:   nonce,
		Filter:      filter,
		SubID:       nonce,
	}

	b.subMu.Lock()
	b.privSubIDs[nonce] = struct{}{}
	b.subMu.Unlock()

	if err := b.Send(ctx, s); err != nil {
		return err
	}
	b.isAuthenticated = true

	return nil
}

func (b *bfxWebsocket) AttachEventHandler(f handlerT) error {
	b.eventHandler = f
	return nil
}

func (b *bfxWebsocket) AttachPrivateHandler(f handlerT) error {
	b.privateHandler = f
	return nil
}

func (b *bfxWebsocket) RemoveEventHandler() error {
	b.eventHandler = nil
	return nil
}

func (b *bfxWebsocket) RemovePrivateHandler() error {
	b.privateHandler = nil
	return nil
}

// SetReadTimeout sets the read timeout for the underlying websocket connections.
func (b *bfxWebsocket) SetReadTimeout(t time.Duration) {
	atomic.StoreInt64(&b.timeout, t.Nanoseconds())
}
