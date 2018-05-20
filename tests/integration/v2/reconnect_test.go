package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

var (
	wsPort    = 4001
	wsService *TestWsService
	apiRecv   *listener
	apiClient *websocket.Client
)

func assertDisconnect(maxWait time.Duration, client *websocket.Client) error {
	loops := 5
	delay := maxWait / time.Duration(loops)
	for i := 0; i < loops; i++ {
		if !client.IsConnected() {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("peer did not disconnect in %s", maxWait.String())
}

// convienence
func newTestParams(wsPort int) *websocket.Parameters {
	p := websocket.NewDefaultParameters()
	p.HeartbeatTimeout = time.Second * 4
	p.ShutdownTimeout = time.Second * 4
	p.URL = fmt.Sprintf("ws://localhost:%d", wsPort)
	p.AutoReconnect = true
	p.ReconnectInterval = time.Millisecond * 250 // first reconnect is instant, won't need to wait on this
	return p
}

func setup(t *testing.T, hbTimeout time.Duration, autoReconnect, auth bool) {
	if wsService != nil {
		wsService.Stop()
	}
	if apiClient != nil {
		apiClient.Close()
	}

	time.Sleep(time.Millisecond * 250)

	wsService = NewTestWsService(wsPort)
	wsService.PublishOnConnect(`{"event":"info","version":2}`)
	err := wsService.Start()
	if err != nil {
		t.Fatal(err)
	}

	// create client
	params := newTestParams(wsPort)
	params.HeartbeatTimeout = hbTimeout
	params.AutoReconnect = autoReconnect
	factory := websocket.NewWebsocketAsynchronousFactory(params)
	nonce := &IncrementingNonceGenerator{}
	apiClient = websocket.NewWithParamsAsyncFactoryNonce(params, factory, nonce)
	if auth {
		apiClient.Credentials("apiKey1", "apiSecret1")
	}

	// setup listener
	// listener closes when apiClient is closed
	apiRecv = newListener()
	apiRecv.run(apiClient.Listen())

	// set ws options
	apiClient.Connect()

	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
}

func TestReconnectResubscribeWithAuth(t *testing.T) {
	// create transport & nonce mocks
	setup(t, time.Second*10, true, true)

	// begin test
	infoEv, err := apiRecv.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	expInfoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &expInfoEv, infoEv)

	msg, err := wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce1","event":"auth","apiKey":"apiKey1","authSig":"6e7e3ab737bac9d36b6c3170356c9483edb0079cb65a2afa81efa9a6b906e0c3aeb16b574a44073dff4c0f604adbdd7d","authPayload":"AUTHnonce1","authNonce":"nonce1"}` != msg {
		t.Fatalf("[1] did not expect to receive msg: %s", msg)
	}
	wsService.Broadcast(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	authEv, err := apiRecv.nextAuthEvent()
	if err != nil {
		t.Fatal(err)
	}
	expAuthEv := websocket.AuthEvent{
		Event:  "auth",
		Status: "OK",
		ChanID: 0,
		UserID: 1,
		SubID:  "nonce1",
		AuthID: "valid-auth-guid",
	}
	assert(t, &expAuthEv, authEv)

	// subscriptions
	// trade sub
	_, err = apiClient.SubscribeTrades(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	msg, err = wsService.WaitForMessage(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce2","event":"subscribe","channel":"trades","symbol":"tBTCUSD"}` != msg {
		t.Fatalf("[2] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"trades","chanId":5,"symbol":"tBTCUSD","subId":"nonce2","pair":"BTCUSD"}`)
	tradeSub, err := apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTradeSub := websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce2",
		Channel: "trades",
	}
	assert(t, &expTradeSub, tradeSub)

	// book sub
	_, err = apiClient.SubscribeBook(context.Background(), "tBTCUSD", bitfinex.Precision0, bitfinex.FrequencyRealtime, 25)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = wsService.WaitForMessage(0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce3","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0","len":"25"}` != msg {
		t.Fatalf("[3] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce3","pair":"BTCUSD","prec":"P0","freq":"F0","len":"25"}`)
	bookSub, err := apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub := websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce3",
		Channel:   "book",
		Frequency: string(bitfinex.FrequencyRealtime),
		Precision: string(bitfinex.Precision0),
	}
	assert(t, &expBookSub, bookSub)

	// abrupt disconnect
	wsService.Stop()

	now := time.Now()
	// wait for client disconnect to start reconnect looping
	err = assertDisconnect(time.Second*10, apiClient)
	if err != nil {
		t.Fatal(err)
	}
	diff := time.Now().Sub(now)
	t.Logf("client disconnect detected in %s", diff.String())

	// recreate service
	wsService = NewTestWsService(wsPort)
	// fresh service, no clients
	if wsService.TotalClientCount() != 0 {
		t.Fatalf("total client count %d, expected non-zero", wsService.TotalClientCount())
	}
	wsService.Start()
	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
	wsService.Broadcast(`{"event":"info","version":2}`)
	infoEv, err = apiRecv.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &expInfoEv, infoEv)

	// assert authentication again
	msg, err = wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce4","event":"auth","apiKey":"apiKey1","authSig":"3e424670c0fa4dcb293eea38b9fe62cca49cacc595da01a493d6b9328517a5c940b22141fecf16f653c2662b298238f4","authPayload":"AUTHnonce4","authNonce":"nonce4"}` != msg {
		t.Fatalf("[4] did not expect to receive msg: %s", msg)
	}
	wsService.Broadcast(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce4","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	authEv, err = apiRecv.nextAuthEvent()
	if err != nil {
		t.Fatal(err)
	}
	expAuthEv = websocket.AuthEvent{
		Event:  "auth",
		Status: "OK",
		ChanID: 0,
		UserID: 1,
		SubID:  "nonce4",
		AuthID: "valid-auth-guid",
	}
	assert(t, &expAuthEv, authEv)

	// ensure client automatically resubscribes
	msg, err = wsService.WaitForMessage(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce5","event":"subscribe","channel":"trades","symbol":"tBTCUSD"}` != msg {
		t.Fatalf("[6] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"trades","chanId":5,"symbol":"tBTCUSD","subId":"nonce5","pair":"BTCUSD"}`)
	tradeSub, err = apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTradeSub = websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce5",
		Channel: "trades",
	}
	assert(t, &expTradeSub, tradeSub)
	msg, err = wsService.WaitForMessage(0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce6","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0","len":"25"}` != msg {
		t.Fatalf("[5] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce6","pair":"BTCUSD","prec":"P0","freq":"F0","len":"25"}`)
	bookSub, err = apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub = websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce6",
		Channel:   "book",
		Frequency: string(bitfinex.FrequencyRealtime),
		Precision: string(bitfinex.Precision0),
		Len:       "25",
	}
	assert(t, &expBookSub, bookSub)

	// API client thinks it's connected
	if !apiClient.IsConnected() {
		t.Fatal("not reconnected to websocket")
	}
}

func TestHeartbeatTimeoutNoReconnect(t *testing.T) {
	// create transport & nonce mocks
	setup(t, time.Second, false, false)

	// begin test
	msg, err := apiRecv.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	infoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &infoEv, msg)

	_, err = apiClient.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce1","pair":"BTCUSD"}`)

	// expect timeout channel heartbeat
	time.Sleep(time.Second * 2)

	if apiClient.IsConnected() {
		t.Fatal("API client still connected, expected heartbeat disconnect")
	}
}

// also tests resubscribes
func TestHeartbeatTimeoutReconnect(t *testing.T) {
	// create transport & nonce mocks
	setup(t, time.Second, true, false)

	// begin test
	// info msg automatically sends
	msg, err := apiRecv.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	infoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &infoEv, msg)

	// use ticker sub to check for reconnect
	_, err = apiClient.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	m, err := wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce1","event":"subscribe","channel":"ticker","symbol":"tBTCUSD"}` != m {
		t.Fatalf("[1] did not expect to receive: %s", m)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce1","pair":"BTCUSD"}`)
	tickerSub, err := apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTickerSub := websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce1",
		Channel: "ticker",
	}
	assert(t, &expTickerSub, tickerSub)

	// expect timeout channel heartbeat
	time.Sleep(time.Second * 2)

	// check reconnect subscriptions
	m, err = wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce2","event":"subscribe","channel":"ticker","symbol":"tBTCUSD"}` != m {
		t.Fatalf("[2] did not expect to receive: %s", m)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce2","pair":"BTCUSD"}`)
	tickerSub, err = apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTickerSub = websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce2",
		Channel: "ticker",
	}
	assert(t, &expTickerSub, tickerSub)
}

func TestHeartbeatNoTimeoutData(t *testing.T) {
	// create transport & nonce mocks
	setup(t, time.Second, true, false)

	// begin test
	// info msg automatically sends
	msg, err := apiRecv.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	infoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &infoEv, msg)

	// use ticker sub to check for reconnect
	_, err = apiClient.SubscribeTicker(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	m, err := wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce1","event":"subscribe","channel":"ticker","symbol":"tBTCUSD"}` != m {
		t.Fatalf("[1] did not expect to receive: %s", m)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"ticker","chanId":5,"symbol":"tBTCUSD","subId":"nonce1","pair":"BTCUSD"}`)
	tickerSub, err := apiRecv.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTickerSub := websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce1",
		Channel: "ticker",
	}
	assert(t, &expTickerSub, tickerSub)

	// would normally timeout here, but we can publish data to prevent
	// the 1 second timeout
	for i := 0; i < 8; i++ {
		wsService.Broadcast(`[5,[14957,68.17328796,14958,55.29588132,-659,-0.0422,14971,53723.08813995,16494,14454]]`)
		time.Sleep(time.Millisecond * 250)
	}

	tick, err := apiRecv.nextTick()
	if err != nil {
		log.Fatal(err)
	}
	expTicker := bitfinex.Ticker{
		Symbol:          "tBTCUSD",
		Bid:             14957,
		BidSize:         68.17328796,
		Ask:             14958,
		AskSize:         55.29588132,
		DailyChange:     -659,
		DailyChangePerc: -0.0422,
		LastPrice:       14971,
		Volume:          53723.08813995,
		High:            16494,
		Low:             14454,
	}
	assert(t, &expTicker, tick)

	if !apiClient.IsConnected() {
		t.Fatal("expected client connected, client has disconnected")
	}
}
