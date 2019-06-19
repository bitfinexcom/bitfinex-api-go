package tests

import (
	"context"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
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
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
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
	err_unsub := ws.Unsubscribe(context.Background(), id)
	if err_unsub != nil {
		t.Fatal(err_unsub)
	}
	async.Publish(`{"event":"unsubscribed","chanId":5,"status":"OK"}`)
	unsub, err := listener.nextUnsubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.UnsubscribeEvent{ChanID: 5, Status: "OK"}, unsub)
}

func TestOrderbook(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()

	// create client
	p := websocket.NewDefaultParameters()
	p.ManageOrderbook = true
	ws := websocket.NewWithParamsAsyncFactory(p, newTestAsyncFactory(async))

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// info welcome msg
	async.Publish(`{"event":"info","version":2}`)
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.InfoEvent{Version: 2}, ev)

	// we will use XRPBTC since that uses reallllyy small numbers
	bId, err_st := ws.SubscribeBook(context.Background(), bitfinex.TradingPrefix+bitfinex.XRPBTC, bitfinex.Precision0, bitfinex.FrequencyRealtime, 25)
	if err_st != nil {
		t.Fatal(err_st)
	}
	// checksum enabled ack
	async.Publish(`{"event":"conf","status":"OK","flags":131072}`)
	// subscribe ack
	async.Publish(`{"event":"subscribed","channel":"book","chanId":81757,"symbol":"tXRPBTC","prec":"P0","freq":"F0","len":"25","subId":"` + bId + `","pair":"XRPBTC"}`)

	// publish a snapshot
	async.Publish(`[81757,[[0.0000011,13,271510.49],[0.00000109,4,500793.10790141],[0.00000108,5,776367.43],[0.00000107,1,23329.54842056],[0.00000106,3,116868.87735849],[0.00000105,3,205000],[0.00000103,3,227308.25386407],[0.00000102,2,105000],[0.00000101,1,2970],[0.000001,2,21000],[7e-7,1,10000],[6.6e-7,1,10000],[6e-7,1,100000],[4.9e-7,1,10000],[2.5e-7,1,2000],[6e-8,1,100000],[5e-8,1,200000],[1e-8,4,640000],[0.00000111,1,-4847.13],[0.00000112,7,-528102.69042633],[0.00000113,5,-302397.07],[0.00000114,3,-339088.93],[0.00000126,4,-245944.06],[0.00000127,1,-5000],[0.0000013,1,-5000],[0.00000134,1,-8249.18938656],[0.00000136,1,-13161.25184766],[0.00000145,1,-2914],[0.0000015,3,-54448.5],[0.00000152,2,-5538.54849594],[0.00000153,1,-62691.75475079],[0.00000159,1,-2914],[0.0000016,1,-52631.10296831],[0.00000164,1,-4000],[0.00000166,1,-3831.46784605],[0.00000171,1,-14575.17730379],[0.00000174,1,-3124.81815395],[0.0000018,1,-18000],[0.00000182,1,-16000],[0.00000186,1,-4000],[0.00000189,1,-10000.686624],[0.00000191,1,-14500],[0.00000193,1,-2422]]]`)

	// publish new trade update
	async.Publish(`[81757,[0.0000011,12,266122.94]]`)

	// test that we can retrieve the orderbook
	ob, err_ob := ws.GetOrderbook("tXRPBTC")
	if err_ob != nil {
		t.Fatal(err_ob)
	}

	// test that changing the orderbook values will not invalidate the checksum
	// since they have been dereferenced
	ob.Bids()[0].Amount = 9999999

	// publish new checksum
	pre := async.SentCount()
	async.Publish(`[81757,"cs",-1175357890]`)

	// test that the new trade has been added to the orderbook
	newTrade := ob.Bids()[0]
	// check that it has overwritten the original trade in the book at that price
	if newTrade.PriceJsNum.String() != "0.0000011" {
		t.Fatal("Newly submitted trade did not update into orderbook")
	}
	if newTrade.AmountJsNum.String() != "266122.94" {
		t.Fatal("Newly submitted trade did not update into orderbook")
	}
	// check that we did not send an unsubscribe message
	// because that would mean the checksum was incorrect
	if err_unsub := async.waitForMessage(pre); err_unsub != nil {
		// no message sent
		return
	} else {
		t.Fatal("A new unsubscribe message was sent")
	}
}

func TestCreateNewSocket(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()

	// create client
	p := websocket.NewDefaultParameters()
	// lock the capacity to 10
	p.CapacityPerConnection = 10
	p.ManageOrderbook = true
	ws := websocket.NewWithParamsAsyncFactory(p, newTestAsyncFactory(async))

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// info welcome msg
	async.Publish(`{"event":"info","version":2}`)
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.InfoEvent{Version: 2}, ev)

	tickers := []string{"tBTCUSD", "tETHUSD", "tBTCUSD", "tVETUSD", "tDGBUSD", "tEOSUSD", "tTRXUSD", "tEOSETH", "tBTCETH",
		"tBTCEOS", "tXRPUSD", "tXRPBTC", "tTRXETH", "tTRXBTC", "tLTCUSD", "tLTCBTC", "tLTCETH"}
	for i, ticker := range tickers {
		id := i*10
		// subscribe to 15m candles
		id1, err := ws.SubscribeCandles(context.Background(), ticker, bitfinex.FifteenMinutes)
		if err != nil {
			t.Fatal(err)
		}
		async.Publish(`{"event":"subscribed","channel":"candles","chanId":`+string(id)+`,"key":"trade:15m:`+ticker+`","subId":"`+id1+`"}`)
		// subscribe to 1hr candles
		id2, err := ws.SubscribeCandles(context.Background(), ticker, bitfinex.OneHour)
		if err != nil {
			t.Fatal(err)
		}
		// subscribe ack
		async.Publish(`{"event":"subscribed","channel":"candles","chanId":`+string(id+1)+`,"key":"trade:1hr:`+ticker+`","subId":"`+id2+`"}`)
		// subscribe to 30min candles
		id3, err := ws.SubscribeCandles(context.Background(), ticker, bitfinex.OneHour)
		if err != nil {
			t.Fatal(err)
		}
		// subscribe ack
		async.Publish(`{"event":"subscribed","channel":"candles","chanId":`+string(id+2)+`,"key":"trade:30m:`+ticker+`","subId":"`+id3+`"}`)
	}
	conCount := ws.ConnectionCount()
	if conCount != 6 {
		t.Fatal("Expected socket count to be 6 but got", conCount)
	}
}
