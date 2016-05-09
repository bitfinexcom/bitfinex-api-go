package bitfinex

import (
    "encoding/json"
    "fmt"
    "log"
    "reflect"
    "time"

    "golang.org/x/net/websocket"
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
    // special web socket for private messages
    privateWs *websocket.Conn
    // map internal channels to websocket's
    chanMap    map[float64]chan []float64
    subscribes []subscribeToChannel
}

type SubscribeMsg struct {
    Event   string  `json:"event"`
    Channel string  `json:"channel"`
    Pair    string  `json:"pair"`
    ChanId  float64 `json:"chanId,omitempty"`
}

type subscribeToChannel struct {
    Channel string
    Pair    string
    Chan    chan []float64
}

func NewWebSocketService(c *Client) *WebSocketService {
    return &WebSocketService{
        client:     c,
        chanMap:    make(map[float64]chan []float64),
        subscribes: make([]subscribeToChannel, 0),
    }
}

// Connect create new bitfinex websocket connection
func (w *WebSocketService) Connect() error {
    ws, err := websocket.Dial(w.client.WebSocketURL, "", "http://localhost/")
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

// Watch allows to subsribe to channels and watch for new updates.
// This method supports next channels: book, trade, ticker.
func (w *WebSocketService) Subscribe() {
    // Subscribe to each channel
    for _, s := range w.subscribes {
        msg, _ := json.Marshal(SubscribeMsg{
            Event:   "subscribe",
            Channel: s.Channel,
            Pair:    s.Pair,
        })

        _, err := w.ws.Write(msg)
        if err != nil {
            // Can't send message to web socket.
            log.Fatal(err)
        }
    }

    var clientMessage string
    for {
        if err := websocket.Message.Receive(w.ws, &clientMessage); err != nil {
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
                    var fullPayload []interface{}
                    err = json.Unmarshal([]byte(clientMessage), &fullPayload)
                    if err != nil {
                        // log.Println("Error decoding fullPayload", err)
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
                // Received "subscribed" resposne. Link channels.
                for _, k := range w.subscribes {
                    if event.Event == "subscribed" && event.Pair == k.Pair && event.Channel == k.Channel {
                        fmt.Println("!!!", event, "r:", k.Channel, k.Pair)
                        w.chanMap[event.ChanId] = k.Chan
                    }
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
    ApiKey      string `json:"apiKey"`
    AuthSig     string `json:"authSig"`
    AuthPayload string `json:"authPayload"`
}

// Private channel auth response
type privateResponse struct {
    Event  string  `json:"event"`
    Status string  `json:"status"`
    ChanId float64 `json:"chanId,omitempty"`
    UserId float64 `json:"userId"`
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
    ws, err := websocket.Dial(w.client.WebSocketURL, "", "http://localhost/")
    if err != nil {
        ch <- TermData{
            Error: err.Error(),
        }
        return
    }

    payload := "AUTH" + fmt.Sprintf("%v", time.Now().Unix())
    connectMsg, _ := json.Marshal(&privateConnect{
        Event:       "auth",
        ApiKey:      w.client.ApiKey,
        AuthSig:     w.client.signPayload(payload),
        AuthPayload: payload,
    })

    // Send auth message
    _, err = ws.Write(connectMsg)
    if err != nil {
        ch <- TermData{
            Error: err.Error(),
        }
        ws.Close()
        return
    }

    var msg string
    for {
        if err = websocket.Message.Receive(ws, &msg); err != nil {
            ch <- TermData{
                Error: err.Error(),
            }
            ws.Close()
            return
        } else {
            event := &privateResponse{}
            err = json.Unmarshal([]byte(msg), &event)
            if err != nil {
                // received data update
                var data []interface{}
                err = json.Unmarshal([]byte(msg), &data)
                if err == nil {
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
                    ws.Close()
                }
            }
        }
    }
}
