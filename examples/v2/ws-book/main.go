package main

import (
	"context"
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

func main() {
	c := bitfinex.NewClient()

	err := c.Websocket.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}
	c.Websocket.SetReadTimeout(time.Second * 2)

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		log.Printf("EVENT: %#v", ev)
	})

	h := func(ev interface{}) {
		log.Printf("PUBLIC MSG BTCUSD: %#v", ev)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanTicker,
		Symbol:  bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h)
	if err != nil {
		log.Fatal(err)
	}

	h2 := func(ev interface{}) {
		log.Printf("PUBLIC MSG IOTUSD: %#v", ev)
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Second*1)
	msg = &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanTrades,
		Symbol:  bitfinex.TradingPrefix + bitfinex.IOTUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg, h2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-c.Websocket.Done():
			log.Printf("channel closed: %s", c.Websocket.Err())
			return
		}
	}
}
