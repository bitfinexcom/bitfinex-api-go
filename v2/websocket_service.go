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

// Available pairs
const (
	BTCUSD = "BTCUSD"
	LTCUSD = "LTCUSD"
	LTCBTC = "LTCBTC"
	ETHUSD = "ETHUSD"
	ETHBTC = "ETHBTC"
	ETCUSD = "ETCUSD"
	ETCBTC = "ETCBTC"
	BFXUSD = "BFXUSD"
	BFXBTC = "BFXBTC"
	ZECUSD = "ZECUSD"
	ZECBTC = "ZECBTC"
	XMRUSD = "XMRUSD"
	XMRBTC = "XMRBTC"
	RRTUSD = "RRTUSD"
	RRTBTC = "RRTBTC"
	XRPUSD = "XRPUSD"
	XRPBTC = "XRPBTC"
	EOSETH = "EOSETH"
	EOSUSD = "EOSUSD"
	EOSBTC = "EOSBTC"
	IOTUSD = "IOTUSD"
	IOTBTC = "IOTBTC"
	IOTETH = "IOTETH"
	BCCBTC = "BCCBTC"
	BCUBTC = "BCUBTC"
	BCCUSD = "BCCUSD"
	BCUUSD = "BCUUSD"
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

var ErrWSNotConnected = fmt.Errorf("websocket connection not established")

// bfxWebsocket is a wrapper around a simple websocket connection, that let's us
// manage callbacks and share a single websocket in a thread safe manner.
// It provides a single channel to write message to.
type bfxWebsocket struct {
	// Bitfinex client
	client *Client

	wsMu          sync.Mutex
	ws            *websocket.Conn
	timeout       int64
	webSocketURL  string
	TLSSkipVerify bool

	// The bitfinex API sends us untyped arrays as data, so we have to keep track
	// of which one belongs where.
	mu          sync.Mutex
	pubSubIDs   map[string]PublicSubscriptionRequest
	pubChanIDs  map[int64]PublicSubscriptionRequest // ChannelID -> SubscriptionRequest map
	privSubIDs  map[string]struct{}
	privChanIDs map[int64]struct{}

	eventHandler   handlerT
	privateHandler handlerT
	publicHandler  handlerT

	mc   *msgChan
	stop chan struct{}
}

type handlerT func(interface{})

func newBfxWebsocket(c *Client, wsURL string) *bfxWebsocket {
	b := &bfxWebsocket{
		client:       c,
		privSubIDs:   map[string]struct{}{},
		pubSubIDs:    map[string]PublicSubscriptionRequest{},
		pubChanIDs:   map[int64]PublicSubscriptionRequest{},
		privChanIDs:  map[int64]struct{}{},
		webSocketURL: wsURL,
		mc:           newMsgChan(),
		stop:         make(chan struct{}),
	}

	return b
}

func (b *bfxWebsocket) Connect() error {
	b.wsMu.Lock()
	defer b.wsMu.Unlock()
	if b.ws != nil {
		return nil // We're already connected.
	}

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

	b.privSubIDs = map[string]struct{}{}
	b.pubSubIDs = map[string]PublicSubscriptionRequest{}
	b.pubChanIDs = map[int64]PublicSubscriptionRequest{}
	b.privChanIDs = map[int64]struct{}{}
	b.mc = newMsgChan()
	b.stop = make(chan struct{})

	go b.sender()
	go b.receiver()

	return nil
}

func (b *bfxWebsocket) sender() {
	for {
		select {
		default:
		case <-b.mc.Done():
			return
		case msg := <-b.mc.Receive():
			err := b.ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil { // If WriteMessage returns an error, it's permanent.
				b.mc.Close(err)
				return
			}
		}
	}
}

// Done returns a channel that will be closed if the underlying websocket
// connection gets closed.
func (b *bfxWebsocket) Done() <-chan struct{} { return b.mc.Done() }

// Err returns an error if the done channel was closed due to an error.
func (b *bfxWebsocket) Err() error { return b.mc.Err() }

//
func (b *bfxWebsocket) Close() {
	b.wsMu.Lock()
	defer b.wsMu.Unlock()
	b.mc.Close(nil)
	b.ws.Close()
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
	case b.mc.C <- bs:
		return nil
	case <-b.Done():
		return fmt.Errorf("websocket closed: ", b.mc.Err())
	}

	return nil
}

func (b *bfxWebsocket) receiver() {
	for {
		if atomic.LoadInt64(&b.timeout) != 0 {
			b.ws.SetReadDeadline(time.Now().Add(time.Duration(b.timeout)))
		}
		_, msg, err := b.ws.ReadMessage()
		if err != nil {
			b.mc.Close(err)
			return
		}

		t := bytes.TrimLeftFunc(msg, unicode.IsSpace)
		if bytes.HasPrefix(t, []byte("[")) {
			var raw []interface{}
			err := json.Unmarshal(msg, &raw)
			if err != nil {
				//return err
				continue
			} else if len(raw) < 2 {
				//return nil
				continue
			}

			chanID, ok := raw[0].(float64)
			if !ok {
				//return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
				continue
			}

			if _, has := b.privChanIDs[int64(chanID)]; has {
				td, err := b.handlePrivateDataMessage(raw)
				if err != nil {
					log.Printf("[WARN]: %s\n", err)
					continue
				} else if td == nil {
					continue
				}
				if b.privateHandler != nil {
					go b.privateHandler(td)
				}
			} else if _, has := b.pubChanIDs[int64(chanID)]; has {
				td, err := b.handlePublicDataMessage(raw)
				if err != nil {
					log.Printf("[WARN]: %s\n", err)
					continue
				} else if td == nil {
					continue
				}
				if b.publicHandler != nil {
					go b.publicHandler(td)
				}
			} else {
				// TODO: log unhandled message?
			}
		} else if bytes.HasPrefix(t, []byte("{")) { // Events are encoded as objects.
			ev, err := b.onEvent(msg)
			if err != nil {
				log.Printf("[WARN]: %s\n", err)
			}
			if b.eventHandler != nil {
				go b.eventHandler(ev)
			}
		} else {
			log.Printf("[WARN]: unexpected message: %s\n", msg)
		}
	}
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
		AuthSig:     b.client.signPayload(payload),
		AuthPayload: payload,
		AuthNonce:   nonce,
		Filter:      filter,
		SubID:       nonce,
	}

	b.mu.Lock()
	b.privSubIDs[nonce] = struct{}{}
	b.mu.Unlock()

	return b.Send(ctx, s)
}

func (b *bfxWebsocket) AttachEventHandler(f handlerT) error {
	b.eventHandler = f
	return nil
}

func (b *bfxWebsocket) AttachPrivateHandler(f handlerT) error {
	b.privateHandler = f
	return nil
}

func (b *bfxWebsocket) AttachPublicHandler(f handlerT) error {
	b.publicHandler = f
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

func (b *bfxWebsocket) RemovePublicHandler() error {
	b.publicHandler = nil
	return nil
}

// SetReadTimeout sets the read timeout for the underlying websocket connections.
func (b *bfxWebsocket) SetReadTimeout(t time.Duration) {
	atomic.StoreInt64(&b.timeout, t.Nanoseconds())
}
