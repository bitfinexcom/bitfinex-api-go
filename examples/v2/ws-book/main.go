package main

import (
	"context"
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func main() {
	c := websocket.New()

	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD ticker
	ctx, cxl1 := context.WithTimeout(context.Background(), time.Second*1)
	defer cxl1()
	_, err = c.SubscribeTicker(ctx, bitfinex.TradingPrefix+bitfinex.BTCUSD)
	if err != nil {
		log.Fatal(err)
	}

	// subscribe to IOTUSD trades
	ctx, cxl2 := context.WithTimeout(context.Background(), time.Second*1)
	defer cxl2()
	_, err = c.SubscribeTrades(ctx, bitfinex.TradingPrefix+bitfinex.IOTUSD)
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
			break
		default:
		}
		log.Printf("MSG RECV: %#v", obj)
	}
}
