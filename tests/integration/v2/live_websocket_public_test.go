package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func wait(wg *sync.WaitGroup, bc <-chan error, to time.Duration) error {
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
	c := websocket.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()
	c.SetReadTimeout(time.Second * 2)

	errch := make(chan error)
	go func() {
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch msg.(type) {
				case error:
					errch <- msg.(error)
				default:
					t.Logf("test recv: %#v", msg)
				}
				wg.Done()
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*2)
	defer cxl()
	id, err := c.SubscribeTicker(ctx, bitfinex.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Fatalf("failed to receive first message from websocket: %s", err)
	}

	// here?
	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive second message from websocket: %s", err)
	}
}

func TestPublicTrades(t *testing.T) {
	c := websocket.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. 3 x data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()
	c.SetReadTimeout(time.Second * 2)

	errch := make(chan error)
	go func() {
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch msg.(type) {
				case error:
					errch <- msg.(error)
				default:
					t.Logf("test recv: %#v", msg)
				}
				wg.Done()
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*2)
	defer cxl()
	id, err := c.SubscribeTrades(ctx, bitfinex.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive first message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive second message from websocket: %s", err)
	}
}

func TestPublicBooks(t *testing.T) {
	c := websocket.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()
	c.SetReadTimeout(time.Second * 2)

	errch := make(chan error)
	go func() {
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch msg.(type) {
				case error:
					errch <- msg.(error)
				default:
					t.Logf("test recv: %#v", msg)
				}
				wg.Done()
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*2)
	defer cxl()
	id, err := c.SubscribeBook(ctx, bitfinex.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Fatalf("failed to receive first message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive second message from websocket: %s", err)
	}
}

func TestPublicCandles(t *testing.T) {
	c := websocket.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(3) // 1. Info with version, 2. Subscription event, 3. data message

	err := c.Connect()
	if err != nil {
		t.Fatal("Error connecting to web socket : ", err)
	}
	defer c.Close()
	c.SetReadTimeout(time.Second * 2)

	errch := make(chan error)
	go func() {
		for {
			select {
			case msg := <-c.Listen():
				if msg == nil {
					return
				}
				log.Printf("recv msg: %#v", msg)
				switch msg.(type) {
				case error:
					errch <- msg.(error)
				default:
					t.Logf("test recv: %#v", msg)
				}
				wg.Done()
			}
		}
	}()

	ctx, cxl := context.WithTimeout(context.Background(), time.Second*2)
	defer cxl()
	id, err := c.SubscribeCandles(ctx, bitfinex.TradingPrefix+bitfinex.BTCUSD, bitfinex.OneMonth)
	if err != nil {
		t.Fatal(err)
	}

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive first message from websocket: %s", err)
	}

	err = c.Unsubscribe(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)

	if err := wait(&wg, errch, 2*time.Second); err != nil {
		t.Errorf("failed to receive second message from websocket: %s", err)
	}
}
