package tests

import (
	"context"
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"testing"
	"time"
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

func TestReconnect(t *testing.T) {
	// create transport & nonce mocks
	wsPort := 4001
	wsService := NewTestWsService(wsPort)
	err := wsService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer wsService.Stop()

	// create client
	params := websocket.NewDefaultParameters().SetAutoReconnect(true).SetReconnectAttempts(5).SetReconnectInterval(time.Millisecond * 250).SetURL(fmt.Sprintf("ws://localhost:%d", wsPort))
	factory := websocket.NewWebsocketAsynchronousFactory(params)
	nonce := &MockNonceGenerator{}
	apiClient := websocket.NewWithParamsAsyncFactoryNonce(params, factory, nonce)

	// setup listener
	listener := newListener()
	listener.run(apiClient.Listen())

	// set ws options
	apiClient.Connect()
	defer apiClient.Close()

	// begin test
	nonce.Next("nonce1") // auth nonce
	wsService.Broadcast(`{"event":"info","version":2}`)
	msg, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	infoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &infoEv, msg)

	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
	// abrupt disconnect
	wsService.Stop()
	nonce.Next("nonce2")

	now := time.Now()
	// wait for client disconnect to start reconnect looping
	err = assertDisconnect(time.Second*20, apiClient)
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
	// ERROR client not reconnecting
	wsService.Start()
	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
	wsService.Broadcast(`{"event":"info","version":2}`)
	msg, err = listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &infoEv, msg)

	// API client thinks it's connected
	if !apiClient.IsConnected() {
		t.Fatal("not reconnected to websocket")
	}

	// done
}

func TestReconnectResubscribeNoAuth(t *testing.T) {
	// create transport & nonce mocks
	wsPort := 4001
	wsService := NewTestWsService(wsPort)
	err := wsService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer wsService.Stop()

	// create client
	params := websocket.NewDefaultParameters().SetAutoReconnect(true).SetReconnectAttempts(5).SetReconnectInterval(time.Millisecond * 250).SetURL(fmt.Sprintf("ws://localhost:%d", wsPort))
	factory := websocket.NewWebsocketAsynchronousFactory(params)
	nonce := &MockNonceGenerator{}
	apiClient := websocket.NewWithParamsAsyncFactoryNonce(params, factory, nonce)

	// setup listener
	listener := newListener()
	listener.run(apiClient.Listen())

	// set ws options
	apiClient.Connect()
	defer apiClient.Close()

	// begin test
	nonce.Next("nonce1") // auth nonce
	wsService.Broadcast(`{"event":"info","version":2}`)
	infoEv, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	expInfoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &expInfoEv, infoEv)

	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}

	// subscriptions
	nonce.Next("nonce2")
	_, err = apiClient.SubscribeTrades(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	msg, err := wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce2","event":"subscribe","channel":"trades","symbol":"tBTCUSD"}` != msg {
		t.Fatalf("[1] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"trades","chanId":5,"symbol":"tBTCUSD","subId":"nonce2","pair":"BTCUSD"}`)
	tradeSub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTradeSub := websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce2",
		Channel: "trades",
	}
	assert(t, &expTradeSub, tradeSub)

	nonce.Next("nonce3")
	_, err = apiClient.SubscribeBook(context.Background(), "tBTCUSD", websocket.Precision0, websocket.FrequencyRealtime)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = wsService.WaitForMessage(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce3","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0"}` != msg {
		t.Fatalf("[2] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce3","pair":"BTCUSD","prec":"P0","freq":"F0"}`)
	bookSub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub := websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce3",
		Channel:   "book",
		Frequency: string(websocket.FrequencyRealtime),
		Precision: string(websocket.Precision0),
	}
	assert(t, &expBookSub, bookSub)

	// abrupt disconnect
	wsService.Stop()
	nonce.Next("nonce2")

	now := time.Now()
	// wait for client disconnect to start reconnect looping
	err = assertDisconnect(time.Second*20, apiClient)
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
	// ERROR client not reconnecting
	wsService.Start()
	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
	wsService.Broadcast(`{"event":"info","version":2}`)
	infoEv, err = listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &expInfoEv, infoEv)

	// ensure client automatically resubscribes
	msg, err = wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce2","event":"subscribe","channel":"trades","symbol":"tBTCUSD"}` != msg {
		t.Fatalf("[3] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"trades","chanId":5,"symbol":"tBTCUSD","subId":"nonce2","pair":"BTCUSD"}`)
	tradeSub, err = listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTradeSub = websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce2",
		Channel: "trades",
	}
	assert(t, &expTradeSub, tradeSub)
	msg, err = wsService.WaitForMessage(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce3","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0"}` != msg {
		t.Fatalf("[4] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce3","pair":"BTCUSD","prec":"P0","freq":"F0"}`)
	bookSub, err = listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub = websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce3",
		Channel:   "book",
		Frequency: string(websocket.FrequencyRealtime),
		Precision: string(websocket.Precision0),
	}
	assert(t, &expBookSub, bookSub)

	// API client thinks it's connected
	if !apiClient.IsConnected() {
		t.Fatal("not reconnected to websocket")
	}

	// done
}

func TestReconnectResubscribeWithAuth(t *testing.T) {
	// create transport & nonce mocks
	wsPort := 4001
	wsService := NewTestWsService(wsPort)
	err := wsService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer wsService.Stop()

	// create client
	params := websocket.NewDefaultParameters().SetAutoReconnect(true).SetReconnectAttempts(5).SetReconnectInterval(time.Millisecond * 250).SetURL(fmt.Sprintf("ws://localhost:%d", wsPort))
	factory := websocket.NewWebsocketAsynchronousFactory(params)
	nonce := &MockNonceGenerator{}
	apiClient := websocket.NewWithParamsAsyncFactoryNonce(params, factory, nonce).Credentials("apiKey1", "apiSecret1")

	// setup listener
	listener := newListener()
	listener.run(apiClient.Listen())

	// set ws options
	apiClient.Connect()
	defer apiClient.Close()

	// begin test
	nonce.Next("nonce1") // auth nonce
	wsService.Broadcast(`{"event":"info","version":2}`)
	infoEv, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	expInfoEv := websocket.InfoEvent{
		Version: 2,
	}
	assert(t, &expInfoEv, infoEv)

	if err := wsService.WaitForClientCount(1); err != nil {
		t.Fatal(err)
	}
	msg, err := wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce1","event":"auth","apiKey":"apiKey1","authSig":"6e7e3ab737bac9d36b6c3170356c9483edb0079cb65a2afa81efa9a6b906e0c3aeb16b574a44073dff4c0f604adbdd7d","authPayload":"AUTHnonce1","authNonce":"nonce1"}` != msg {
		t.Fatalf("[1] did not expect to receive msg: %s", msg)
	}
	wsService.Broadcast(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	authEv, err := listener.nextAuthEvent()
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
	nonce.Next("nonce2")
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
	tradeSub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expTradeSub := websocket.SubscribeEvent{
		Symbol:  "tBTCUSD",
		SubID:   "nonce2",
		Channel: "trades",
	}
	assert(t, &expTradeSub, tradeSub)

	nonce.Next("nonce3")
	_, err = apiClient.SubscribeBook(context.Background(), "tBTCUSD", websocket.Precision0, websocket.FrequencyRealtime)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = wsService.WaitForMessage(0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if `{"subId":"nonce3","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0"}` != msg {
		t.Fatalf("[3] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce3","pair":"BTCUSD","prec":"P0","freq":"F0"}`)
	bookSub, err := listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub := websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce3",
		Channel:   "book",
		Frequency: string(websocket.FrequencyRealtime),
		Precision: string(websocket.Precision0),
	}
	assert(t, &expBookSub, bookSub)

	// abrupt disconnect
	wsService.Stop()
	nonce.Next("nonce4")

	now := time.Now()
	// wait for client disconnect to start reconnect looping
	err = assertDisconnect(time.Second*20, apiClient)
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
	infoEv, err = listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &expInfoEv, infoEv)

	// assert authentication again
	msg, err = wsService.WaitForMessage(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	nonce.Next("nonce5") // for next subs
	if `{"subId":"nonce4","event":"auth","apiKey":"apiKey1","authSig":"3e424670c0fa4dcb293eea38b9fe62cca49cacc595da01a493d6b9328517a5c940b22141fecf16f653c2662b298238f4","authPayload":"AUTHnonce4","authNonce":"nonce4"}` != msg {
		t.Fatalf("[4] did not expect to receive msg: %s", msg)
	}
	wsService.Broadcast(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce4","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	authEv, err = listener.nextAuthEvent()
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
	tradeSub, err = listener.nextSubscriptionEvent()
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
	if `{"subId":"nonce5","event":"subscribe","channel":"book","symbol":"tBTCUSD","prec":"P0","freq":"F0"}` != msg {
		t.Fatalf("[5] did not expect to receive: %s", msg)
	}
	wsService.Broadcast(`{"event":"subscribed","channel":"book","chanId":8,"symbol":"tBTCUSD","subId":"nonce5","pair":"BTCUSD","prec":"P0","freq":"F0"}`)
	bookSub, err = listener.nextSubscriptionEvent()
	if err != nil {
		t.Fatal(err)
	}
	expBookSub = websocket.SubscribeEvent{
		Symbol:    "tBTCUSD",
		SubID:     "nonce5",
		Channel:   "book",
		Frequency: string(websocket.FrequencyRealtime),
		Precision: string(websocket.Precision0),
	}
	assert(t, &expBookSub, bookSub)

	// API client thinks it's connected
	if !apiClient.IsConnected() {
		t.Fatal("not reconnected to websocket")
	}

	// done
	// TODO extend this test with private-only subscription requests
}
