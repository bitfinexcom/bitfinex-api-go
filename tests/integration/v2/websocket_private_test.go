package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

func TestWebsocketOrder(t *testing.T) {
	if !auth {
		t.Skip("no credentials, skipping order creation")
	}

	wg := sync.WaitGroup{}
	wg.Add(1) // 1. Authentication event

	c := bitfinex.NewClient().Credentials(key, secret)

	err := c.Websocket.Connect()
	if err != nil {
		t.Fatalf("connecting to websocket service: %s", err)
	}
	defer c.Websocket.Close()
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		switch e := ev.(type) {
		case bitfinex.AuthEvent:
			if e.Status == "OK" {
				wg.Done()
			}
		case bitfinex.UnsubscribeEvent:
			wg.Done()
		}
	})

	c.Websocket.AttachPrivateHandler(func(ev interface{}) {
		switch e := ev.(type) {
		case bitfinex.Notification:
			if e.Status == "ERROR" && e.Type == "on-req" {
				t.Errorf("failed to create order: %s", e.Text)
			}
		case bitfinex.OrderNew:
			wg.Done()
		case bitfinex.OrderCancel:
			wg.Done()
		}
	})

	err = c.Websocket.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("authenticating with websocket service: %s", err)
	}
	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Fatalf("failed to authenticate with websocket service: %s", err)
	}

	wg.Add(1)
	n := time.Now()
	cid := n.Unix()
	cidDate := n.Format("2006-01-02")
	o := &bitfinex.OrderNewRequest{
		CID:    cid,
		Type:   bitfinex.OrderTypeExchangeLimit,
		Symbol: bitfinex.TradingPrefix + bitfinex.BTCUSD,
		Amount: 1.0,
		Price:  28.5,
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = c.Websocket.Send(ctx, o)
	if err != nil {
		t.Fatalf("failed to send OrderNewRequest: %s", err)
	}
	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Fatalf("failed to create order: %s", err)
	}

	oc := &bitfinex.OrderCancelRequest{
		CID:     &cid,
		CIDDate: &cidDate,
	}

	wg.Add(1)
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = c.Websocket.Send(ctx, oc)
	if err != nil {
		t.Fatalf("failed to send OrderCancelRequest: %s", err)
	}
	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Fatalf("failed to cancel order: %s", err)
	}
}
