package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

func wait(wg *sync.WaitGroup, bc <-chan struct{}, to time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-bc:
		return fmt.Errorf("websocket closed while waiting") // timed out
	case <-c:
		return nil // completed normally
	case <-time.After(to):
		return fmt.Errorf("timed out waiting") // timed out
	}
}

func TestPublicTicker(t *testing.T) {
	c := bitfinex.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Websocket.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Websocket.Close()
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		wg.Done()
	})

	h := func(ev interface{}) {
		wg.Done()
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanTicker,
		Symbol:  bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Fatalf("failed to receive message from websocket: %s", err)
	}

	err = c.Websocket.Unsubscribe(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}
}

func TestPublicTrades(t *testing.T) {
	c := bitfinex.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Websocket.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Websocket.Close()
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		wg.Done()
	})

	h := func(ev interface{}) {
		wg.Done()
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanTrades,
		Symbol:  bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}

	err = c.Websocket.Unsubscribe(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}
}

func TestPublicBooks(t *testing.T) {
	c := bitfinex.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Websocket.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Websocket.Close()
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		wg.Done()
	})

	h := func(ev interface{}) {
		wg.Done()
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanBook,
		Symbol:  bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Fatalf("failed to receive message from websocket: %s", err)
	}

	err = c.Websocket.Unsubscribe(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}
}

func TestPublicCandles(t *testing.T) {
	c := bitfinex.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Websocket.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Websocket.Close()
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		wg.Done()
	})

	h := func(ev interface{}) {
		wg.Done()
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanCandles,
		Key:     "trade:1M:" + bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}

	err = c.Websocket.Unsubscribe(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, c.Websocket.Done(), 2*time.Second); err != nil {
		t.Errorf("failed to receive message from websocket: %s", err)
	}
}
