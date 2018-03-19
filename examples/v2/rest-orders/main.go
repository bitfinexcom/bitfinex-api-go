package main

import (
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
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
	c := rest.NewClient().Credentials(key, secret)

	// request without params
	os, err := c.Orders.History(bitfinex.TradingPrefix + bitfinex.BTCUSD)
	if err != nil {
		log.Fatalf("getting orders: %s", err)
	}

	log.Printf("orders: %#v\n", os)

	// request orders with params
	params := make(map[string]interface{})
	params["limit"] = 2

	os, err = c.Orders.History(bitfinex.TradingPrefix+bitfinex.BTCUSD, params)
	if err != nil {
		log.Fatalf("getting orders: %s", err)
	}

	log.Printf("orders with params: %#v\n", os)
}
