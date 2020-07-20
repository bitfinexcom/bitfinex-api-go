package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	get(c)
	getMulti(c)
	all(c)
}

func get(c *rest.Client) {
	t, err := c.Tickers.Get("tBTCUSD")

	if err != nil {
		log.Fatalf("get: %s", err)
	}

	spew.Dump(t)
}

func getMulti(c *rest.Client) {
	symbols := []string{"tBTCUSD", "tEOSBTC"}
	tm, err := c.Tickers.GetMulti(symbols)

	if err != nil {
		log.Fatalf("getMulti: %s", err)
	}

	spew.Dump(tm)
}

func all(c *rest.Client) {
	t, err := c.Tickers.All()

	if err != nil {
		log.Fatalf("all: %s", err)
	}

	spew.Dump(t)
}
