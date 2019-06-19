package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func TestAuthentication(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test
	async.Publish(`{"event":"info","version":2}`)
	ev, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.InfoEvent{Version: 2}, ev)

	// assert outgoing auth request
	if err := async.waitForMessage(0); err != nil {
		t.Fatal(err.Error())
	}
	assert(t, &websocket.SubscriptionRequest{SubID: "nonce1", Event: "auth", APIKey: "apiKeyABC"}, async.Sent[0].(*websocket.SubscriptionRequest))

	// auth ack
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)

	// assert incoming auth ack
	av, err := listener.nextAuthEvent()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &websocket.AuthEvent{Status: "OK", SubID: "nonce1", ChanID: 0}, av)
}

func TestWalletBalanceUpdates(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	//ws.SetReadTimeout(time.Second * 2)
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test--authentication assertions in TestAuthentication
	async.Publish(`{"event":"info","version":2}`)
	// eat event
	_, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}

	// auth ack
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)

	// eat event
	_, err = listener.nextAuthEvent()
	if err != nil {
		t.Fatal(err)
	}

	// publish account info post auth ack
	async.Publish(`[0,"wu",["exchange","BTC",30,0,30]]`)
	async.Publish(`[0,"wu",["exchange","USD",80000,0,80000]]`)
	async.Publish(`[0,"wu",["exchange","ETH",100,0,100]]`)
	async.Publish(`[0,"wu",["margin","BTC",10,0,10]]`)
	async.Publish(`[0,"wu",["funding","BTC",10,0,10]]`)
	async.Publish(`[0,"wu",["funding","USD",10000,0,10000]]`)
	async.Publish(`[0,"wu",["margin","USD",10000,0,10000]]`)
	async.Publish(`[0,"bu",[147260,147260]]`)

	// assert incoming wallet & balance updates
	wu, err := listener.nextWalletUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.WalletUpdate{Type: "exchange", Currency: "BTC", Balance: 30, BalanceAvailable: 30}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "exchange", Currency: "USD", Balance: 80000, BalanceAvailable: 80000}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "exchange", Currency: "ETH", Balance: 100, BalanceAvailable: 100}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "BTC", Balance: 10, BalanceAvailable: 10}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "funding", Currency: "BTC", Balance: 10, BalanceAvailable: 10}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "funding", Currency: "USD", Balance: 10000, BalanceAvailable: 10000}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "USD", Balance: 10000, BalanceAvailable: 10000}, wu)
	bu, err := listener.nextBalanceUpdate()
	if err != nil {
		t.Fatal(err)
	}
	// total aum, net aum
	assert(t, &bitfinex.BalanceUpdate{TotalAUM: 147260, NetAUM: 147260}, bu)
}

func TestNewOrder(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	//ws.SetReadTimeout(time.Second * 2)
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test
	async.Publish(`{"event":"info","version":2}`)
	_, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}

	// initial logon info--Authentication & WalletUpdate assertions in prior tests
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	async.Publish(`[0,"wu",["exchange","BTC",30,0,30]]`)
	async.Publish(`[0,"wu",["exchange","USD",80000,0,80000]]`)
	async.Publish(`[0,"wu",["exchange","ETH",100,0,100]]`)
	async.Publish(`[0,"wu",["margin","BTC",10,0,10]]`)
	async.Publish(`[0,"wu",["funding","BTC",10,0,10]]`)
	async.Publish(`[0,"wu",["funding","USD",10000,0,10000]]`)
	async.Publish(`[0,"wu",["margin","USD",10000,0,10000]]`)
	async.Publish(`[0,"bu",[147260,147260]]`)

	// submit order
	err = ws.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
		Symbol: "tBTCUSD",
		CID:    123,
		Amount: -0.456,
	})
	if err != nil {
		t.Fatal(err)
	}

	// assert outgoing order request
	if len(async.Sent) <= 1 {
		t.Fatalf("expected >1 sent messages, got %d", len(async.Sent))
	}
	assert(t, &bitfinex.OrderNewRequest{Symbol: "tBTCUSD", CID: 123, Amount: -0.456}, async.Sent[1].(*bitfinex.OrderNewRequest))

	// order ack
	async.Publish(`[0,"n",[null,"on-req",null,null,[1234567,null,123,"tBTCUSD",null,null,1,1,"MARKET",null,null,null,null,null,null,null,915.5,null,null,null,null,null,null,0,null,null],null,"SUCCESS","Submitting market buy order for 1.0 BTC."]]`)

	// assert order ack notification
	not, err := listener.nextNotification()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.Notification{Type: "on-req", NotifyInfo: &bitfinex.OrderNew{ID: 1234567, CID: 123, Symbol: "tBTCUSD", Amount: 1, AmountOrig: 1, Type: "MARKET", Price: 915.5}}, not)
}

func TestFills(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	//ws.SetReadTimeout(time.Second * 2)
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test
	async.Publish(`{"event":"info","version":2}`)
	_, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}

	// initial logon info
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	async.Publish(`[0,"ps",[["tBTCUSD","ACTIVE",7,916.52002351,0,0,null,null,null,null]]]`)
	async.Publish(`[0,"ws",[["exchange","BTC",30,0,null],["exchange","USD",80000,0,null],["exchange","ETH",100,0,null],["margin","BTC",10,0,null],["margin","USD",9987.16871968,0,null],["funding","BTC",10,0,null],["funding","USD",10000,0,null]]]`)
	// consume & assert snapshots
	ps, err := listener.nextPositionSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	eps := make([]*bitfinex.Position, 1)
	eps[0] = &bitfinex.Position{
		Symbol:    "tBTCUSD",
		Status:    "ACTIVE",
		Amount:    7,
		BasePrice: 916.52002351,
	}
	snap := &bitfinex.PositionSnapshot{
		Snapshot: eps,
	}
	assertSlice(t, snap, ps)
	w, err := listener.nextWalletSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	ews := make([]*bitfinex.Wallet, 7)
	ews[0] = &bitfinex.Wallet{Type: "exchange", Currency: "BTC", Balance: 30}
	ews[1] = &bitfinex.Wallet{Type: "exchange", Currency: "USD", Balance: 80000}
	ews[2] = &bitfinex.Wallet{Type: "exchange", Currency: "ETH", Balance: 100}
	ews[3] = &bitfinex.Wallet{Type: "margin", Currency: "BTC", Balance: 10}
	ews[4] = &bitfinex.Wallet{Type: "margin", Currency: "USD", Balance: 9987.16871968}
	ews[5] = &bitfinex.Wallet{Type: "funding", Currency: "BTC", Balance: 10}
	ews[6] = &bitfinex.Wallet{Type: "funding", Currency: "USD", Balance: 10000}
	wsnap := &bitfinex.WalletSnapshot{
		Snapshot: ews,
	}
	assertSlice(t, wsnap, w)

	// submit order
	err = ws.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
		Symbol: "tBTCUSD",
		CID:    123,
		Amount: -0.456,
	})
	if err != nil {
		t.Fatal(err)
	}

	// order ack
	async.Publish(`[0,"n",[null,"on-req",null,null,[1234567,null,123,"tBTCUSD",null,null,1,1,"MARKET",null,null,null,null,null,null,null,915.5,null,null,null,null,null,null,0,null,null],null,"SUCCESS","Submitting market buy order for 1.0 BTC."]]`)

	// assert order ack notification--Authentication, WalletUpdate, order acknowledgement assertions in prior tests
	_, err = listener.nextNotification()
	if err != nil {
		t.Fatal(err)
	}

	// <..crossing orders generates a fill..>
	// partial fills--position updates
	async.Publish(`[0,"pu",["tBTCUSD","ACTIVE",0.21679716,915.9,0,0,null,null,null,null]]`)
	async.Publish(`[0,"pu",["tBTCUSD","ACTIVE",1,916.13496085,0,0,null,null,null,null]]`)
	pu, err := listener.nextPositionUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.PositionUpdate{Symbol: "tBTCUSD", Status: "ACTIVE", Amount: 0.21679716, BasePrice: 915.9}, pu)
	pu, _ = listener.nextPositionUpdate()
	assert(t, &bitfinex.PositionUpdate{Symbol: "tBTCUSD", Status: "ACTIVE", Amount: 1, BasePrice: 916.13496085}, pu)

	// full fill--order terminal state message
	async.Publish(`[0,"oc",[1234567,0,123,"tBTCUSD",1514909325236,1514909325631,0,1,"MARKET",null,null,null,0,"EXECUTED @ 916.2(0.78): was PARTIALLY FILLED @ 915.9(0.22)",null,null,915.5,916.13496085,null,null,null,null,null,0,0,0]]`)
	oc, err := listener.nextOrderCancel()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.OrderCancel{ID: 1234567, CID: 123, Symbol: "tBTCUSD", MTSCreated: 1514909325236, MTSUpdated: 1514909325631, Amount: 0, AmountOrig: 1, Type: "MARKET", Status: "EXECUTED @ 916.2(0.78): was PARTIALLY FILLED @ 915.9(0.22)", Price: 915.5, PriceAvg: 916.13496085}, oc)

	// fills--trade executions
	async.Publish(`[0,"te",[1,"tBTCUSD",1514909325593,1234567,0.21679716,915.9,null,null,-1]]`)
	async.Publish(`[0,"te",[2,"tBTCUSD",1514909325597,1234567,0.78320284,916.2,null,null,-1]]`)
	te, err := listener.nextTradeExecution()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.TradeExecution{ID: 1, MTS: 1514909325593, Amount: 0.21679716, Price: 915.9}, te)
	te, err = listener.nextTradeExecution()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.TradeExecution{ID: 2, MTS: 1514909325597, Amount: 0.78320284, Price: 916.2}, te)

	// fills--trade updates
	async.Publish(`[0,"tu",[1,"tBTCUSD",1514909325593,1234567,0.21679716,915.9,"MARKET",915.5,-1,-0.39712904,"USD"]]`)
	async.Publish(`[0,"tu",[2,"tBTCUSD",1514909325597,1234567,0.78320284,916.2,"MARKET",915.5,-1,-1.43514088,"USD"]]`)
	tu, err := listener.nextTradeUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.TradeExecutionUpdate{ID: 1, Pair: "tBTCUSD", MTS: 1514909325593, ExecAmount: 0.21679716, ExecPrice: 915.9, OrderType: "MARKET", OrderPrice: 915.5, OrderID: 1234567, Maker: -1, Fee: -0.39712904, FeeCurrency: "USD"}, tu)
	tu, _ = listener.nextTradeUpdate()
	assert(t, &bitfinex.TradeExecutionUpdate{ID: 2, Pair: "tBTCUSD", MTS: 1514909325597, ExecAmount: 0.78320284, ExecPrice: 916.2, OrderType: "MARKET", OrderPrice: 915.5, OrderID: 1234567, Maker: -1, Fee: -1.43514088, FeeCurrency: "USD"}, tu)

	// fills--wallet updates from fee deduction
	async.Publish(`[0,"wu",["margin","USD",9999.60287096,0,null]]`)
	async.Publish(`[0,"wu",["margin","USD",9998.16773008,0,null]]`)
	wu, err := listener.nextWalletUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "USD", Balance: 9999.60287096}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "USD", Balance: 9998.16773008}, wu)

	// margin info update for executed trades
	async.Publish(`[0,"miu",["base",[-2.76536085,0,19150.16773008,19147.40236923]]]`)
	async.Publish(`[0,"miu",["sym","tBTCUSD",[60162.93960325,61088.2924336,60162.93960325,60162.93960325,null,null,null,null]]]`)
	mb, err := listener.nextMarginInfoBase()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.MarginInfoBase{}, mb)
	mu, err := listener.nextMarginInfoUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.MarginInfoUpdate{}, mu)

	// position update for executed trades
	async.Publish(`[0,"pu",["tBTCUSD","ACTIVE",1,916.13496085,0,0,-2.76536085,-0.30185082,0,43.7962]]`)
	pu, _ = listener.nextPositionUpdate()
	assert(t, &bitfinex.PositionUpdate{Symbol: "tBTCUSD", Status: "ACTIVE", Amount: 1, BasePrice: 916.13496085, ProfitLoss: -2.76536085, ProfitLossPercentage: -0.30185082, Leverage: 43.7962}, pu)

	// wallet margin update for executed trades
	async.Publish(`[0,"wu",["margin","BTC",10,0,10]]`)
	async.Publish(`[0,"wu",["margin","USD",9998.16773008,0,9998.16773008]]`)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "BTC", Balance: 10, BalanceAvailable: 10}, wu)
	wu, _ = listener.nextWalletUpdate()
	assert(t, &bitfinex.WalletUpdate{Type: "margin", Currency: "USD", Balance: 9998.16773008, BalanceAvailable: 9998.16773008}, wu)

	// funding update for executed trades
	async.Publish(`[0,"fiu",["sym","ftBTCUSD",[0,0,0,0]]]`)
	fi, err := listener.nextFundingInfo()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.FundingInfo{Symbol: "ftBTCUSD"}, fi)
}

func TestCancel(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	//ws.SetReadTimeout(time.Second * 2)
	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test
	async.Publish(`{"event":"info","version":2}`)
	_, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}

	// initial logon info--Authentication & WalletUpdate assertions in prior tests
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	async.Publish(`[0,"ps",[["tBTCUSD","ACTIVE",7,916.52002351,0,0,null,null,null,null]]]`)
	async.Publish(`[0,"ws",[["exchange","BTC",30,0,null],["exchange","USD",80000,0,null],["exchange","ETH",100,0,null],["margin","BTC",10,0,null],["margin","USD",9987.16871968,0,null],["funding","BTC",10,0,null],["funding","USD",10000,0,null]]]`)
	// consume & assert snapshots
	_, err_ps := listener.nextPositionSnapshot()
	if err_ps != nil {
		t.Fatal(err_ps)
	}
	_, err_was := listener.nextWalletSnapshot()
	if err_was != nil {
		t.Fatal(err_was)
	}

	// submit order
	err = ws.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
		Symbol: "tBTCUSD",
		CID:    123,
		Amount: -0.456,
		Type:   "LIMIT",
		Price:  900.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	// assert outgoing order request
	if len(async.Sent) <= 1 {
		t.Fatalf("expected >1 sent messages, got %d", len(async.Sent))
	}
	assert(t, &bitfinex.OrderNewRequest{Symbol: "tBTCUSD", CID: 123, Amount: -0.456, Type: "LIMIT", Price: 900.0}, async.Sent[1].(*bitfinex.OrderNewRequest))

	// order pending new
	async.Publish(`[0,"n",[null,"on-req",null,null,[1234567,null,123,"tBTCUSD",null,null,1,1,"LIMIT",null,null,null,null,null,null,null,900,null,null,null,null,null,null,0,null,null],null,"SUCCESS","Submitting limit buy order for 1.0 BTC."]]`)
	// order working--limit order
	async.Publish(`[0,"on",[1234567,0,123,"tBTCUSD",1515179518260,1515179518315,1,1,"LIMIT",null,null,null,0,"ACTIVE",null,null,900,0,null,null,null,null,null,0,0,0]]`)

	// eat order ack notification
	_, err_n := listener.nextNotification()
	if err_n != nil {
		t.Fatal(err_n)
	}

	on, err := listener.nextOrderNew()
	if err != nil {
		t.Fatal(err)
	}

	// assert order new update
	assert(t, &bitfinex.OrderNew{ID: 1234567, CID: 123, Symbol: "tBTCUSD", MTSCreated: 1515179518260, MTSUpdated: 1515179518315, Type: "LIMIT", Amount: 1, AmountOrig: 1, Status: "ACTIVE", Price: 900.0}, on)

	// publish cancel request
	req := &bitfinex.OrderCancelRequest{ID: on.ID}
	pre := async.SentCount()
	err = ws.SubmitCancel(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if err := async.waitForMessage(pre); err != nil {
		t.Fatal(err.Error())
	}
	// assert sent message
	assert(t, req, async.Sent[pre].(*bitfinex.OrderCancelRequest))

	// cancel ack notify
	async.Publish(`[0,"n",[null,"oc-req",null,null,[1149686139,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,0,null,null],null,"SUCCESS","Submitted for cancellation; waiting for confirmation (ID: 1149686139)."]]`)

	// cancel confirm
	async.Publish(`[0,"oc",[1234567,0,123,"tBTCUSD",1515179518260,1515179520203,1,1,"LIMIT",null,null,null,0,"CANCELED",null,null,900,0,null,null,null,null,null,0,0,0]]`)

	// assert cancel ack
	oc, err := listener.nextOrderCancel()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.OrderCancel{ID: 1234567, CID: 123, Symbol: "tBTCUSD", MTSCreated: 1515179518260, MTSUpdated: 1515179520203, Type: "LIMIT", Status: "CANCELED", Price: 900.0, Amount: 1, AmountOrig: 1}, oc)
}

func TestUpdateOrder(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	nonce := &IncrementingNonceGenerator{}

	// create client
	ws := websocket.NewWithAsyncFactoryNonce(newTestAsyncFactory(async), nonce).Credentials("apiKeyABC", "apiSecretXYZ")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	err_ws := ws.Connect()
	if err_ws != nil {
		t.Fatal(err_ws)
	}
	defer ws.Close()

	// begin test
	async.Publish(`{"event":"info","version":2}`)
	_, err := listener.nextInfoEvent()
	if err != nil {
		t.Fatal(err)
	}

	// initial logon info--Authentication & WalletUpdate assertions in prior tests
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	async.Publish(`[0,"ps",[["tBTCUSD","ACTIVE",7,916.52002351,0,0,null,null,null,null]]]`)
	async.Publish(`[0,"ws",[["exchange","BTC",30,0,null],["exchange","USD",80000,0,null],["exchange","ETH",100,0,null],["margin","BTC",10,0,null],["margin","USD",9987.16871968,0,null],["funding","BTC",10,0,null],["funding","USD",10000,0,null]]]`)
	// consume & assert snapshots
	_, errps := listener.nextPositionSnapshot()
	if errps != nil {
		t.Fatal(errps)
	}
	_, errws := listener.nextWalletSnapshot()
	if errws != nil {
		t.Fatal(errws)
	}

	// submit order
	err = ws.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
		Symbol: "tBTCUSD",
		CID:    123,
		Amount: -0.456,
		Type:   "LIMIT",
		Price:  900.0,
	})
	if err != nil {
		t.Fatal(err)
	}

	// assert outgoing order request
	if len(async.Sent) <= 1 {
		t.Fatalf("expected >1 sent messages, got %d", len(async.Sent))
	}
	assert(t, &bitfinex.OrderNewRequest{Symbol: "tBTCUSD", CID: 123, Amount: -0.456, Type: "LIMIT", Price: 900.0}, async.Sent[1].(*bitfinex.OrderNewRequest))

	// order pending new
	async.Publish(`[0,"n",[null,"on-req",null,null,[1234567,null,123,"tBTCUSD",null,null,1,1,"LIMIT",null,null,null,null,null,null,null,900,null,null,null,null,null,null,0,null,null],null,"SUCCESS","Submitting limit buy order for 1.0 BTC."]]`)
	// order working--limit order
	async.Publish(`[0,"on",[1234567,0,123,"tBTCUSD",1515179518260,1515179518315,1,1,"LIMIT",null,null,null,0,"ACTIVE",null,null,900,0,null,null,null,null,null,0,0,0]]`)

	// eat order ack notification
	_, errn := listener.nextNotification()
	if errn != nil {
		t.Fatal(errn)
	}

	on, err := listener.nextOrderNew()
	if err != nil {
		t.Fatal(err)
	}

	// assert order new update
	assert(t, &bitfinex.OrderNew{ID: 1234567, CID: 123, Symbol: "tBTCUSD", MTSCreated: 1515179518260, MTSUpdated: 1515179518315, Type: "LIMIT", Amount: 1, AmountOrig: 1, Status: "ACTIVE", Price: 900.0}, on)

	// publish update request
	req := &bitfinex.OrderUpdateRequest{
		ID: on.ID,
		Amount: 0.04,
		Price: 1200,
	}
	pre := async.SentCount()
	err = ws.SubmitUpdateOrder(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if err := async.waitForMessage(pre); err != nil {
		t.Fatal(err.Error())
	}
	// assert sent message
	assert(t, req, async.Sent[pre].(*bitfinex.OrderUpdateRequest))

	// cancel ack notify
	async.Publish(`[0,"n",[1547469854094,"ou-req",null,null,[1234567,0,123,"tBTCUSD",1547469854025,1547469854042,0.04,0.04,"LIMIT",null,null,null,0,"ACTIVE",null,null,1200,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null],null,"SUCCESS","Submitting update to exchange limit buy order for 0.04 BTC."]]`)
	// cancel confirm
	async.Publish(`[0,"ou",[1234567,0,123,"tBTCUSD",1547469854025,1547469854121,0.04,0.04,"LIMIT",null,null,null,0,"ACTIVE",null,null,1200,0,0,0,null,null,null,0,0,null,null,null,"API>BFX",null,null,null]]`)

	// assert cancel ack
	ou, err := listener.nextOrderUpdate()
	if err != nil {
		t.Fatal(err)
	}
	assert(t, &bitfinex.OrderUpdate{ID:1234567, GID:0, CID:123, Symbol:"tBTCUSD", MTSCreated:1547469854025, MTSUpdated:1547469854121, Amount:0.04, AmountOrig:0.04, Type:"LIMIT", TypePrev:"", Flags:0, Status:"ACTIVE", Price:1200, PriceAvg:0, PriceTrailing:0, PriceAuxLimit:0, Notify:false, Hidden:false, PlacedID:0}, ou)
}

func TestUsesAuthenticatedSocket(t *testing.T) {
	// create transport & nonce mocks
	async := newTestAsync()
	// create client
	p := websocket.NewDefaultParameters()
	// lock the capacity to 3
	p.CapacityPerConnection = 3
	ws := websocket.NewWithParamsAsyncFactory(p, newTestAsyncFactory(async)).Credentials("apiKeyABC", "apiSecretXYZ")

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
	// auth ack
	async.Publish(`{"event":"auth","status":"OK","chanId":0,"userId":1,"subId":"nonce1","auth_id":"valid-auth-guid","caps":{"orders":{"read":1,"write":0},"account":{"read":1,"write":0},"funding":{"read":1,"write":0},"history":{"read":1,"write":0},"wallets":{"read":1,"write":0},"withdraw":{"read":0,"write":0},"positions":{"read":1,"write":0}}}`)
	// force websocket to create new connections
	tickers := []string{"tBTCUSD", "tETHUSD", "tBTCUSD", "tVETUSD", "tDGBUSD", "tEOSUSD", "tTRXUSD", "tEOSETH", "tBTCETH",
		"tBTCEOS", "tXRPUSD", "tXRPBTC", "tTRXETH", "tTRXBTC", "tLTCUSD", "tLTCBTC", "tLTCETH"}
	for i, ticker := range tickers {
		// subscribe to 15m candles
		id, err := ws.SubscribeCandles(context.Background(), ticker, bitfinex.FifteenMinutes)
		if err != nil {
			t.Fatal(err)
		}
		async.Publish(`{"event":"subscribed","channel":"candles","chanId":`+string(i)+`,"key":"trade:15m:`+ticker+`","subId":"`+id+`"}`)
	}
	authSocket, err := ws.GetAuthenticatedSocket()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*authSocket)
}
