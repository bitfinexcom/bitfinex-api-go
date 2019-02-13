package tests

import (
	"context"
	"testing"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
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
	async.Publish(`[81757,[[0.000089,1,185.25362178],[0.00008898,1,1000],[0.00008896,2,4727.04575171],[0.00008895,2,21458.69508413],[0.00008893,1,4376.42464394],[0.00008892,2,3368.61853124],[0.00008891,1,19883.51942917],[0.0000889,3,9435.62761455],[0.00008889,3,2028.31355429],[0.00008887,2,1648.2231714],[0.00008885,1,2759.9202752],[0.00008881,1,54317.999],[0.00008879,1,20000],[0.00008875,1,2500],[0.00008873,1,623.90816043],[0.00008867,2,1623.72590867],[0.00008863,1,1001],[0.00008862,1,81476.999],[0.00008861,2,31514.63403132],[0.0000886,1,200000],[0.00008855,1,14000],[0.00008854,1,812.50832839],[0.00008851,1,1000],[0.0000885,1,100],[0.00008847,2,74338],[0.00008901,1,-1067.2702412],[0.00008905,1,-18296.32986369],[0.00008906,1,-15678.6],[0.00008907,1,-2696.32625547],[0.0000891,1,-2247.09217145],[0.00008911,1,-28169.38978256],[0.00008912,2,-7862.51772819],[0.00008913,1,-4491.60631071],[0.00008917,1,-2246.34063345],[0.00008918,3,-2502.36695768],[0.00008919,1,-2759.8023267],[0.0000892,1,-1716.543113],[0.00008923,3,-9841.07336657],[0.00008925,1,-2500],[0.00008929,1,-54317.999],[0.0000893,1,-20000],[0.00008931,2,-1779.82438355],[0.00008932,1,-3000],[0.00008934,1,-1626.72409133],[0.00008935,1,-7600],[0.00008937,1,-1655.81250531],[0.0000894,1,-14000],[0.00008943,1,-24069.65852907],[0.00008944,2,-8600],[0.00008946,1,-1768.039406]]]`)

	// publish new trade update
	async.Publish(`[81757,[0.00008918,2,-1379.90652441]]`)

	// publish new checksum
	pre := async.SentCount()
	async.Publish(`[81757,"cs",1217733465]`)

	// check that we did not send an unsubscribe message
	// because that woul mean the checksum was incorrect
	if err_unsub := async.waitForMessage(pre); err_unsub != nil {
		// no message sent
		return
	} else {
		t.Fatal("A new unsubscribe message was sent")
	}
}
