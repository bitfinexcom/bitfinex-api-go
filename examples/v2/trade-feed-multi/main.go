package main

import (
	"context"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"log"
	_ "net/http/pprof"
)

var tickers = []string{"tBTCUSD", "tETHUSD", "tBTCUSD", "tVETUSD", "tDGBUSD", "tEOSUSD", "tTRXUSD"}

func main() {
	client := websocket.New()
	err := client.Connect()
	if err != nil {
		log.Printf("could not connect: %s", err.Error())
		return
	}

	for obj := range client.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
			return
		case *bitfinex.Trade:
			log.Printf("New trade: %s", obj)
		case *websocket.InfoEvent:
			// Info event confirms connection to the bfx websocket
			for _, ticker := range tickers {
				_, err := client.SubscribeTrades(context.Background(), ticker)
				if err != nil {
					log.Printf("could not subscribe to trades: %s", err.Error())
				}
			}
		default:
			log.Printf("MSG RECV: %#v", obj)
		}
	}
}
