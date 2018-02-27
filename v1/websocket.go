package bitfinex

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/bitfinexcom/bitfinex-api-go/utils"

	"github.com/gorilla/websocket"
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

// BfChanData :
type BfChanData struct {
	ChanID          int
	Channel, Symbol string
	Datas           []interface{}
}

// WebSocketService allow to connect and receive stream data
// from bitfinex.com ws service.
type WebSocketService struct {
	lock     sync.Mutex
	runtimes int
	// http client
	client *Client
	// websocket client
	ws *websocket.Conn
	// special web socket for private messages
	privateWs *websocket.Conn
	// map from channel ID to subscribe info
	chanIDMap   map[int]*subscribeToChannel
	subscribes  []*subscribeToChannel
	defaultChan chan BfChanData
}

type subscribeMsg struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Pair    string `json:"pair"`
	ChanID  int    `json:"chanId,omitempty"`
}

type subscribeToChannel struct {
	Channel string
	Pair    string
	Chan    chan BfChanData
}

// NewWebSocketService returns a WebSocketService using the given client.
func NewWebSocketService(c *Client) *WebSocketService {
	return &WebSocketService{
		client:      c,
		chanIDMap:   make(map[int]*subscribeToChannel, 0),
		subscribes:  make([]*subscribeToChannel, 0),
		defaultChan: make(chan BfChanData, 1024),
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

func (w *WebSocketService) GetDefaultChan() chan BfChanData {
	return w.defaultChan
}

func (w *WebSocketService) Subscribe(channel string, pair string, cn chan BfChanData) error {
	actChan := cn
	if actChan == nil {
		actChan = w.defaultChan
	}

	pair = strings.ToUpper(pair)
	w.lock.Lock()
	defer w.lock.Unlock()
	for idx := range w.subscribes {
		if w.subscribes[idx].Channel == channel && w.subscribes[idx].Pair == pair {
			w.subscribes[idx].Chan = actChan
			return w.sendSubscribeMessages(channel, pair)
		}
	}

	w.subscribes = append(w.subscribes, &subscribeToChannel{Channel: channel, Pair: pair, Chan: actChan})
	return w.sendSubscribeMessages(channel, pair)
}

// Unsubscribe : unsubscribe symbol's channel data
func (w *WebSocketService) Unsubscribe(channel, pair string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	for idx := range w.subscribes {
		if w.subscribes[idx].Channel == channel && w.subscribes[idx].Pair == pair {
			w.subscribes = append(w.subscribes[0:idx], w.subscribes[idx+1:]...)
			break
		}
	}
	for chanID, subptr := range w.chanIDMap {
		if subptr.Channel == channel && subptr.Pair == pair {
			delete(w.chanIDMap, chanID)
			break
		}
	}
}

func (w *WebSocketService) ClearSubscriptions() {
	w.lock.Lock()
	w.subscribes = make([]*subscribeToChannel, 0)
	w.chanIDMap = make(map[int]*subscribeToChannel, 0)
	w.lock.Unlock()
}

func (w *WebSocketService) sendSubscribeMessages(channel, pair string) error {
	msg, _ := json.Marshal(subscribeMsg{
		Event:   "subscribe",
		Channel: channel,
		Pair:    pair,
	})

	if w.ws == nil {
		log.Printf("websocket disconnected!")
		return fmt.Errorf("websocket disconnected")
	}

	err := w.ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

func (w *WebSocketService) Stop() {
	w.runtimes++
	if w.ws != nil {
		w.ws.Close()
	}
}

func (w *WebSocketService) Run() error {
	lastws := w.ws
	defer func() {
		log.Printf("close websocket")
		lastws.Close()
	}()
	w.runtimes++
	tmpTimes := w.runtimes
	for tmpTimes == w.runtimes {
		_, p, err := w.ws.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage failed : %v", err)
			return err
		}
		//log.Printf(string(p))
		if bytes.Contains(p, []byte("event")) {
			w.handleEventMessage(p)
		} else {
			w.handleDataMessage(p)
		}
	}
	log.Printf("WebSocketService.Run exit @ times %d/%d", tmpTimes, w.runtimes)
	return nil
}

func (w *WebSocketService) handleEventMessage(msg []byte) {
	// Check for first message(event:subscribed)
	event := &subscribeMsg{}
	err := json.Unmarshal(msg, event)

	// Received "subscribed" resposne. Link channels.
	if err != nil {
		log.Printf("unmarshal failed, error %v", err)
		return
	}
	if event.Event != "subscribed" {
		return
	}
	log.Printf("got event %s", string(msg))
	w.lock.Lock()
	defer w.lock.Unlock()
	for _, k := range w.subscribes {
		if strings.ToUpper(event.Pair) == strings.ToUpper(k.Pair) && event.Channel == k.Channel {
			w.chanIDMap[event.ChanID] = k
			return // should no duplicated, so return without more loop
		}
	}
}

func (w *WebSocketService) handleDataMessage(msg []byte) {
	// Received payload or data update
	var datas []interface{}
	err := json.Unmarshal(msg, &datas)
	if nil != err {
		log.Printf("Unmarshal failed, error : %v, msg : %s", err, string(msg))
		return
	}

	cnval, ok := datas[0].(float64)
	if !ok {
		log.Printf("datas[0]:%v is not a valid number", datas[0])
		return
	}
	chanID := int(cnval)
	subptr, ok := w.chanIDMap[chanID]
	if !ok {
		log.Printf("chanID '%d' not subscribe yet, %s\n", chanID, string(msg))
		return
	}
	// len(BfChanData.datas) == 1 means "snapshot" slice or "heartbeat", else means "update" slice,
	// receiver can idenfity them easily.
	subptr.Chan <- BfChanData{chanID, subptr.Channel, subptr.Pair, datas[1:]}
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
	connectMsg, _ := json.Marshal(&privateConnect{
		Event:       "auth",
		APIKey:      w.client.APIKey,
		AuthSig:     w.client.signPayload(payload),
		AuthPayload: payload,
	})

	// Send auth message
	err = ws.WriteMessage(websocket.TextMessage, connectMsg)
	if err != nil {
		ch <- TermData{
			Error: err.Error(),
		}
		log.Printf("ws.WriteMessage failed : %v", err)
		ws.Close()
		return
	}

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			ch <- TermData{
				Error: err.Error(),
			}
			log.Printf("ws.ReadMessage failed : %v", err)
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
					ch <- TermData{Term: data[1].(string)}
					return
				}

				dataTerm := data[1].(string)
				dataList := data[2].([]interface{})

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
				log.Printf("Auth failed, close websocket")
				ws.Close()
			}
		}
	}
}
