package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func TestWebsocketOrder(t *testing.T) {
	if !auth {
		t.Skip("no credentials, skipping order creation")
	}

	wg := sync.WaitGroup{}
	wg.Add(1) // 1. Authentication event

	c := websocket.NewClient().Credentials(key, secret)

	errch := make(chan error)

	err := c.Connect()
	if err != nil {
		t.Fatalf("connecting to websocket service: %s", err)
	}
	defer c.Close()
	c.SetReadTimeout(time.Second * 2)

	go func() {
		for ev := range c.Listen() {
			switch e := ev.(type) {
			case *bitfinex.Notification:
				if e.Status == "ERROR" && e.Type == "on-req" {
					t.Errorf("failed to create order: %s", e.Text)
				}
			case *bitfinex.OrderNew:
				wg.Done()
			case *bitfinex.OrderCancel:
				wg.Done()
			case error:
				t.Logf("Listen() error: %s", ev)
				errch <- ev.(error)
				wg.Done()
			}
		}
	}()
	/*
		err = c.Authenticate(context.Background())
		if err != nil {
			t.Fatalf("authenticating with websocket service: %s", err)
		}
		if err := wait(&wg, errch, 2*time.Second); err != nil {
			t.Fatalf("failed to authenticate with websocket service: %s", err)
		}
	*/
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

	ctx, cxl1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cxl1()
	err = c.SubmitOrder(ctx, o)
	if err != nil {
		t.Fatalf("failed to send OrderNewRequest: %s", err)
	}
	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Fatalf("failed to create order: %s", err)
	}

	oc := &bitfinex.OrderCancelRequest{
		CID:     cid,
		CIDDate: cidDate,
	}

	wg.Add(1)
	ctx, cxl2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cxl2()
	err = c.SubmitCancel(ctx, oc)
	if err != nil {
		t.Fatalf("failed to send OrderCancelRequest: %s", err)
	}
	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Fatalf("failed to cancel order: %s", err)
	}
}
