package tests

/*
import (
	"context"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"testing"
	"time"
)

// the following test is used to run the API against

func TestAPI(t *testing.T) {
	// create transport & nonce mocks

	// create client
	params := websocket.NewDefaultParameters()
	params.URL = "wss://dev-prdn.bitfinex.com:2997/ws/2"
	ws := websocket.NewWithParams(params) //.Credentials("U83q9jkML2GVj1fVxFJOAXQeDGaXIzeZ6PwNPQLEXt4", "77SWIRggvw0rCOJUgk9GVcxbldjTxOJP5WLCjWBFIVc")

	// setup listener
	listener := newListener()
	listener.run(ws.Listen())

	// set ws options
	//ws.SetReadTimeout(time.Second * 2)
	ws.Connect()
	defer ws.Close()

		// begin test
		//ev, err := listener.nextAuthEvent()
		//if err != nil {
		//	t.Fatal(err)
		//}
		//assert(t, &websocket.AuthEvent{Event: "auth", Status: "OK"}, ev)


	tradeSubID, err := ws.SubscribeTrades(context.Background(), "tBTCUSD")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("trade sub ID: %s", tradeSubID)

		//bookSubID, err := ws.SubscribeBook(context.Background(), "tBTCUSD", websocket.Precision0, websocket.FrequencyRealtime)
		//t.Logf("book sub ID: %s", bookSubID)

	time.Sleep(time.Second * 15)
}
*/
