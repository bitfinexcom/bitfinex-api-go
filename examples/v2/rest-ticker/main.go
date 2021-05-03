package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	getSingle(c)
	getMulti(c)
	all(c)
}

func getSingle(c *rest.Client) {
	t, err := c.Tickers.Get("tBTCUSD")

	if err != nil {
		log.Fatalf("getSingle: %s", err)
	}

	spew.Dump("getSingle:> ", t)
}

func getMulti(c *rest.Client) {
	symbols := []string{"tBTCUSD", "fUSD"}
	tm, err := c.Tickers.GetMulti(symbols)

	if err != nil {
		log.Fatalf("getMulti: %s", err)
	}

	spew.Dump("getMulti:> ", tm)
}

func all(c *rest.Client) {
	t, err := c.Tickers.All()

	if err != nil {
		log.Fatalf("all: %s", err)
	}

	spew.Dump("all:> ", t)
}
