package tests

import (
	"context"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"testing"
	"time"
)

// method of testing with mocked endpoints
func TestTicker(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce)

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	ws.Connect()
	defer ws.Close()

	// info welcome msg
	async.Publish(`{"event":"info","version":2}`)
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
	async.Publish(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce1","pair":"BTCUSD"}`)
	sub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.SubscribeEvent{
		SubID:   "nonce1",
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

func TestMockTickerCallback(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce)

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	ws.Connect()
	defer ws.Close()

	// register callbacks for messages
	callbackInfo := make(chan interface{}, 10)
	callbackTick1 := make(chan interface{}, 10)
	callbackTick2 := make(chan interface{}, 10)
	ws.RegisterCallback(websocket.InfoEvent{}, func (info interface{}) { callbackInfo <- info })
	ws.RegisterCallback(bitfinex.Ticker{}, func (tick interface{}) { callbackTick1 <- tick })
	ws.RegisterCallback(bitfinex.Ticker{}, func (tick interface{}) { callbackTick2 <- tick })

	// info welcome msg
	expectedInfoEvent := websocket.InfoEvent{Version: 2}
	async.Publish(`{"event":"info","version":2}`)
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &expectedInfoEvent, ev)
	select {
	case val := <-callbackInfo:
		assert(t, &expectedInfoEvent, val)
	case <-time.After(time.Second):
		t.Fatal("Did not receive an info callback")
	}


	// subscribe
	id, err := ws.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}

	// subscribe ack
	async.Publish(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce1","pair":"BTCUSD"}`)
	sub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.SubscribeEvent{
		SubID:   "nonce1",
		Channel: "ticker",
		ChanID:  5,
		Symbol:  "tBTCUSD",
		Pair:    "BTCUSD",
	}, sub)

	// publish the tick data
	expectedTick := bitfinex.Ticker{
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
	}
	async.Publish(`[5,[14957,68.17328796,14958,55.29588132,-659,-0.0422,14971,53723.08813995,16494,14454]]`)
	tick, err := listener.nextTick()
	if err != nil {
		t.Fatal(err)
	}
	// ensure delivery of the tick matches between serial methods and the multiple callbacks
	assert(t, &expectedTick, tick)
	select {
	case val := <-callbackTick1:
		assert(t, &expectedTick, val)
	case <-time.After(time.Second):
		t.Fatal("Did not receive an tick callback 1/2")
	}
	select {
	case val := <-callbackTick2:
		assert(t, &expectedTick, val)
	case <-time.After(time.Second):
		t.Fatal("Did not receive an tick callback 2/2")
	}

	// unsubscribe
	ws.Unsubscribe(context.Background(), id)
	async.Publish(`{"event":"unsubscribed","chanId":5,"status":"OK"}`)
	unsub, err := listener.nextUnsubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.UnsubscribeEvent{ChanID: 5, Status: "OK"}, unsub)
}

