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
		t.Fatalf("did not expect to receive: %s", msg)
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

	_, err = apiClient.SubscribeBook(context.Background(), "tBTCUSD", websocket.Precision0, websocket.FrequencyRealtime)
	if err != nil {
		t.Fatal(err)
	}
	// TODO assert ws sub

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

	// API client thinks it's connected
	if !apiClient.IsConnected() {
		t.Fatal("not reconnected to websocket")
	}

	// done
}
