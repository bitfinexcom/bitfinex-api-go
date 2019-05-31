package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

var (
	orderid = flag.String("id", "", "lookup trades for an order ID")
	api     = flag.String("api", "https://api-pub.bitfinex.com/v2/", "v2 REST API URL")
)

// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	flag.Parse()

	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	c := rest.NewClientWithURL(*api).Credentials(key, secret)

	available, err := c.Platform.Status()
	if err != nil {
		log.Fatalf("getting status: %s", err)
	}

	if !available {
		log.Fatalf("API not available")
	}

	if *orderid != "" {
		ordid, err := strconv.ParseInt(*orderid, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		os, err := c.Orders.OrderTrades(bitfinex.TradingPrefix+bitfinex.BTCUSD, ordid)
		if err != nil {
			log.Fatalf("getting order trades: %s", err)
		}

		log.Printf("order trades: %#v\n", os)
	} else {
		os, err := c.Orders.AllHistory()
		if err != nil {
			log.Fatalf("getting orders: %s", err)
		}

		log.Printf("orders: %#v\n", os)
	}
}
