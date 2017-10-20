package main

import (
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	c := bitfinex.NewClient().Credentials(key, secret)

	available, err := c.Platform.Status()
	if err != nil {
		log.Fatalf("getting status: %s", err)
	}

	if !available {
		log.Fatalf("API not available")
	}

	os, err := c.Orders.History(bitfinex.TradingPrefix + bitfinex.IOTBTC)
	if err != nil {
		log.Fatalf("getting orders: %s", err)
	}

	log.Printf("orders: %#v\n", os)
}
