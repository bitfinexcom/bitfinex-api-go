package tests

import (
	"context"
	"testing"
	"time"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

// method of testing with mocked endpoints
func TestTicker(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &MockNonceGenerator{}

	// create client
	ws := websocket.NewClientWithAsyncNonce(async, nonce)

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	ws.SetReadTimeout(time.Second * 2)
	ws.Connect()
	defer ws.Close()

	// info welcome msg
	async.Publish(`{"event":"info","version":2}`)
	nonce.Next("1514401173001")
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.InfoEvent{Version: 2}, ev)

	// subscribe
	id, err := ws.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	// subscribe ack
	async.Publish(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"1514401173001","pair":"BTCUSD"}`)
	sub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.SubscribeEvent{
		SubID:   "1514401173001",
		Channel: "ticker",
		ChanID:  5,
		Symbol:  "tBTCUSD",
		Pair:    "BTCUSD",
	}, sub)

	// tick data
	async.Publish(`[5,[14957,68.17328796,14958,55.29588132,-659,-0.0422,14971,53723.08813995,16494,14454]]`)
	tick, err := listener.nextTick()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.Ticker{
		Symbol:          "tBTCUSD",
		Bid:             14957,
		Ask:             14958,
		BidSize:         68.17328796,
		AskSize:         55.29588132,
		DailyChange:     -659,
		DailyChangePerc: -0.0422,
		LastPrice:       14971,
		Volume:          53723.08813995,
		High:            16494,
		Low:             14454,
	}, tick)

	// unsubscribe
	ws.Unsubscribe(context.Background(), id)
	async.Publish(`{"event":"unsubscribed","chanId":5,"status":"OK"}`)
	unsub, err := listener.nextUnsubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.UnsubscribeEvent{ChanID: 5, Status: "OK"}, unsub)
}
