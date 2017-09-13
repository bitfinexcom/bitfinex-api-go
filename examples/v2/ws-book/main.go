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

	c.Websocket.AttachEventHandler(func(_ context.Context, ev interface{}) {
		log.Printf("EVENT: %#v", ev)
	})

	c.Websocket.AttachPublicHandler(func(_ context.Context, ev interface{}) {
		log.Printf("PUBLIC MSG: %#v", ev)
	})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	msg := &bitfinex.PublicSubscriptionRequest{
		Event:   "subscribe",
		Channel: bitfinex.ChanTicker,
		Symbol:  bitfinex.TradingPrefix + bitfinex.BTCUSD,
	}
	err = c.Websocket.Subscribe(ctx, msg)
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
