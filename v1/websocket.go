package bitfinex

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/utils"

	"github.com/gorilla/websocket"
)

const (
	// wsReadTimeout is the maximum time to wait for a message from the server.
	wsReadTimeout = 90 * time.Second
	// wsWriteTimeout is the maximum time to wait for a write to complete.
	wsWriteTimeout = 10 * time.Second
)

// Pairs available
const (
	// Pairs
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

	// Channels
	ChanBook   = "book"
	ChanTrade  = "trades"
	ChanTicker = "ticker"
)

// WebSocketService allow to connect and receive stream data
// from bitfinex.com ws service.
// nolint:megacheck,structcheck
type WebSocketService struct {
	// http client
	client *Client
	// websocket client
	ws *websocket.Conn
	// special web socket for private messages
	privateWs *websocket.Conn
	// map internal channels to websocket's
	chanMap    map[float64]chan []float64
	subscribes []subscribeToChannel
}

type subscribeMsg struct {
	Event   string  `json:"event"`
	Channel string  `json:"channel"`
	Pair    string  `json:"pair"`
	ChanID  float64 `json:"chanId,omitempty"`
}

type subscribeToChannel struct {
	Channel string
	Pair    string
	Chan    chan []float64
}

// NewWebSocketService returns a WebSocketService using the given client.
func NewWebSocketService(c *Client) *WebSocketService {
	return &WebSocketService{
		client:     c,
		chanMap:    make(map[float64]chan []float64),
		subscribes: make([]subscribeToChannel, 0),
	}
}

// Connect create new bitfinex websocket connection
func (w *WebSocketService) Connect() error {
	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	if w.client.WebSocketTLSSkipVerify {
		d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ws, _, err := d.Dial(w.client.WebSocketURL, nil)
	if err != nil {
		return err
	}
	w.ws = ws
	return nil
}

// Close web socket connection
func (w *WebSocketService) Close() {
	if w.ws != nil {
		w.ws.Close()
	}
}

func (w *WebSocketService) AddSubscribe(channel string, pair string, c chan []float64) {
	s := subscribeToChannel{
		Channel: channel,
		Pair:    pair,
		Chan:    c,
	}
	w.subscribes = append(w.subscribes, s)
}

func (w *WebSocketService) ClearSubscriptions() {
	w.subscribes = make([]subscribeToChannel, 0)
}

func (w *WebSocketService) sendSubscribeMessages() error {
	for _, s := range w.subscribes {
		msg, err := json.Marshal(subscribeMsg{
			Event:   "subscribe",
			Channel: s.Channel,
			Pair:    s.Pair,
		})
		if err != nil {
			return err
		}

		w.ws.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
		err = w.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe allows to subsribe to channels and watch for new updates.
// This method supports next channels: book, trade, ticker.
func (w *WebSocketService) Subscribe() error {
	// Subscribe to each channel
	if err := w.sendSubscribeMessages(); err != nil {
		return err
	}

	for {
		w.ws.SetReadDeadline(time.Now().Add(wsReadTimeout))
		_, p, err := w.ws.ReadMessage()
		if err != nil {
			return err
		}

		if bytes.Contains(p, []byte("event")) {
			w.handleEventMessage(p)
		} else {
			w.handleDataMessage(p)
		}
	}
	// nolint
	return nil
}

func (w *WebSocketService) handleEventMessage(msg []byte) {
	// Check for first message(event:subscribed)
	event := &subscribeMsg{}
	err := json.Unmarshal(msg, event)

	// Received "subscribed" resposne. Link channels.
	if err == nil {
		for _, k := range w.subscribes {
			if event.Event == "subscribed" && event.Pair == k.Pair && event.Channel == k.Channel {
				w.chanMap[event.ChanID] = k.Chan
			}
		}
	}
}

func (w *WebSocketService) handleDataMessage(msg []byte) {
	// Received payload or data update
	var dataUpdate []float64
	err := json.Unmarshal(msg, &dataUpdate)
	if err == nil && len(dataUpdate) > 1 {
		chanID := dataUpdate[0]
		// Remove chanID from data update
		// and send message to internal chan
		if ch, ok := w.chanMap[chanID]; ok {
			ch <- dataUpdate[1:]
		}
		return
	}

	// Payload received
	var fullPayload []interface{}
	err = json.Unmarshal(msg, &fullPayload)

	if err != nil {
		log.Println("Error decoding fullPayload", err)
	} else {
		if len(fullPayload) < 2 {
			return
		}
		chanID, ok := fullPayload[0].(float64)
		if !ok {
			return
		}
		ch, ok := w.chanMap[chanID]
		if !ok {
			return
		}
		if len(fullPayload) > 3 {
			itemsSlice := fullPayload[3:]
			i, _ := json.Marshal(itemsSlice)
			var item []float64
			err = json.Unmarshal(i, &item)
			if err == nil {
				ch <- item
			}
		} else {
			itemsSlice := fullPayload[1]
			i, _ := json.Marshal(itemsSlice)
			var items [][]float64
			err = json.Unmarshal(i, &items)
			if err == nil {
				for _, v := range items {
					ch <- v
				}
			}
		}
	}
}

/////////////////////////////
// Private websocket messages
/////////////////////////////

type privateConnect struct {
	Event       string `json:"event"`
	APIKey      string `json:"apiKey"`
	AuthSig     string `json:"authSig"`
	AuthPayload string `json:"authPayload"`
}

// Private channel auth response
type privateResponse struct {
	Event  string  `json:"event"`
	Status string  `json:"status"`
	ChanID float64 `json:"chanId,omitempty"`
	UserID float64 `json:"userId"`
}

type TermData struct {
	// Data term. E.g: ps, ws, ou, etc... See official documentation for more details.
	Term string
	// Data will contain different number of elements for each term.
	// Examples:
	// Term: ws, Data: ["exchange","BTC",0.01410829,0]
	// Term: oc, Data: [0,"BTCUSD",0,-0.01,"","CANCELED",270,0,"2015-10-15T11:26:13Z",0]
	Data  []interface{}
	Error string
}

func (c *TermData) HasError() bool {
	return len(c.Error) > 0
}

func (w *WebSocketService) ConnectPrivate(ch chan TermData) {

	var d = websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	if w.client.WebSocketTLSSkipVerify {
		d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ws, _, err := d.Dial(w.client.WebSocketURL, nil)
	if err != nil {
		ch <- TermData{
			Error: err.Error(),
		}
		return
	}

	nonce := utils.GetNonce()
	payload := "AUTH" + nonce
	sig, err_sig := w.client.signPayload(payload)
	if err_sig != nil {
		ch <- TermData{
			Error: err_sig.Error(),
		}
		ws.Close()
		return
	}
	connectMsg, err := json.Marshal(&privateConnect{
		Event:       "auth",
		APIKey:      w.client.APIKey,
		AuthSig:     sig,
		AuthPayload: payload,
	})
	if err != nil {
		ch <- TermData{
			Error: err.Error(),
		}
		ws.Close()
		return
	}

	// Send auth message
	ws.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
	err = ws.WriteMessage(websocket.TextMessage, connectMsg)
	if err != nil {
		ch <- TermData{
			Error: err.Error(),
		}
		ws.Close()
		return
	}

	for {
		ws.SetReadDeadline(time.Now().Add(wsReadTimeout))
		_, p, err := ws.ReadMessage()
		if err != nil {
			ch <- TermData{
				Error: err.Error(),
			}
			ws.Close()
			return
		}

		event := &privateResponse{}
		err = json.Unmarshal(p, &event)
		if err != nil {
			// received data update
			var data []interface{}
			err = json.Unmarshal(p, &data)
			if err == nil {
				if len(data) == 2 { // Heartbeat
					// XXX: Consider adding a switch to enable/disable passing these along.
					if term, ok := data[1].(string); ok {
						ch <- TermData{Term: term}
					}
					continue
				}

				if len(data) < 3 {
					continue
				}
				dataTerm, ok := data[1].(string)
				if !ok {
					continue
				}
				dataList, ok := data[2].([]interface{})
				if !ok {
					continue
				}

				// check for empty data
				if len(dataList) > 0 {
					if reflect.TypeOf(dataList[0]) == reflect.TypeOf([]interface{}{}) {
						// received list of lists
						for _, v := range dataList {
							ch <- TermData{
								Term: dataTerm,
								Data: v.([]interface{}),
							}
						}
					} else {
						// received flat list
						ch <- TermData{
							Term: dataTerm,
							Data: dataList,
						}
					}
				}
			}
		} else {
			// received auth response
			if event.Event == "auth" && event.Status != "OK" {
				ch <- TermData{
					Error: "Error connecting to private web socket channel.",
				}
				ws.Close()
			}
		}
	}
}
