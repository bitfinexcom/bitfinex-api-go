package tests

/*
import (
	"context"
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"testing"
	"time"
)
*/
/*
// Socket read timeouts should terminate connections, invalidating this test
func TestSocketReadTimeout(t *testing.T) {
	port := 4001
	wsService := NewTestWsService(port)
	nonce := &MockNonceGenerator{}
	wsService.Start()
	nonce.Next("nonce1")

	// create client
	ws := websocket.NewClientWithURLNonce(fmt.Sprintf("ws://localhost:%d", port), nonce)

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	ws.SetReadTimeout(time.Second * 2)
	ws.Connect()
	defer ws.Close()

	// wait for test harness goroutines to start.. turn this into a signal
	// if this sleep is too unreliable & causes failures
	time.Sleep(time.Millisecond * 250)

	wsService.Broadcast(`{"event":"info","version":2}`)

	// info welcome msg
	nonce.Next("nonce2")
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.InfoEvent{Version: 2}, ev)

	// subscribe
	_, err = ws.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}

	// subscribe ack
	wsService.Broadcast(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce2","pair":"BTCUSD"}`)
	sub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.SubscribeEvent{
		SubID:   "nonce2",
		Channel: "ticker",
		ChanID:  5,
		Symbol:  "tBTCUSD",
		Pair:    "BTCUSD",
	}, sub)

	// tick data
	wsService.Broadcast(`[5,[14957,68.17328796,14958,55.29588132,-659,-0.0422,14971,53723.08813995,16494,14454]]`)
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

	// trigger a socket read timeout, do not disconnect
	time.Sleep(time.Second * 3)

	if !ws.IsConnected() {
		t.Fatal("socket not connected after read timeout")
	}
}
*/
