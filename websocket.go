package bitfinex

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
)

// Pairs available
const (
	// Pairs
	BTCUSD = "BTCUSD"
	LTCUSD = "LTCUSD"
	LTCBTC = "LTCBTC"

	// Channels
	CHAN_BOOK   = "book"
	CHAN_TRADE  = "trade"
	CHAN_TICKER = "ticker"
)

// WebSocketService allow to connect and receive stream data
// from bitfinex.com ws service.
type WebSocketService struct {
	// http client
	client *Client
	// websocket client
	ws *websocket.Conn
	// map internal channles to websocket's
	chanMap map[float64]chan []float64
}

func NewWebSocketService(c *Client) *WebSocketService {
	return &WebSocketService{
		client:  c,
		chanMap: make(map[float64]chan []float64),
	}
}

type SubscribeMsg struct {
	Event   string `json:"Event"`
	Channel string `json:"Channel"`
	Pair    string `json:"Pair"`
	// in response
	ChanId float64 `json:"ChanId,omitempty"`
}

// Connect create new bitfinex websocket connection
func (w *WebSocketService) Connect() error {
	ws, err := websocket.Dial(WebSocketURL, "", "http://localhost/")
	if err != nil {
		return err
	}
	w.ws = ws
	return nil
}

// Close web socket connection
func (w *WebSocketService) Close() {
	w.ws.Close()
}

// Watch allows to subsribe to channels and watch for new updates.
// This method supports next channels: book, trade, ticker.
func (w *WebSocketService) Watch(channel string, pair string, c chan []float64) {
	msg, _ := json.Marshal(SubscribeMsg{
		Event:   "subscribe",
		Channel: channel,
		Pair:    pair,
	})

	_, err := w.ws.Write(msg)
	if err != nil {
		// Can't send message to web socket.
		log.Fatal(err)
	}

	var clientMessage string
	for {
		if err = websocket.Message.Receive(w.ws, &clientMessage); err != nil {
			log.Fatal("Error reading message: ", err)
		} else {
			// Check for first message(event:subscribed)
			event := &SubscribeMsg{}
			err = json.Unmarshal([]byte(clientMessage), &event)
			if err != nil {
				// Received payload or data update
				var dataUpdate []float64
				err = json.Unmarshal([]byte(clientMessage), &dataUpdate)
				if err == nil {
					chanId := dataUpdate[0]
					// Remove chanId from data update
					// and send message to internal chan
					w.chanMap[chanId] <- dataUpdate[1:]
				} else {
					// Payload received
					// TODO: Refactor this!
					var fullPayload []interface{}
					err = json.Unmarshal([]byte(clientMessage), &fullPayload)
					if err != nil {
						log.Println("Error decoding fullPayload", err)
					} else {
						itemsSlice := fullPayload[1]
						i, _ := json.Marshal(itemsSlice)
						var items [][]float64
						err = json.Unmarshal(i, &items)
						if err == nil {
							chanId := fullPayload[0].(float64)
							for _, v := range items {
								w.chanMap[chanId] <- v
							}
						}
					}
				}
			} else {
				// Received "subscribed" resposne. Lets link channles.
				if event.Event == "subscribed" {
					w.chanMap[event.ChanId] = c
				}
			}
		}
	}
}
