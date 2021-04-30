package main

import (
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	getSinglePair(c)
	getSinglePairWithLimit(c)
	getMultiPair(c)
	getForRange(c)
}

func getSinglePair(c *rest.Client) {
	args := rest.GetTickerHistPayload{
		Symbols: []string{"tLTCUSD"},
	}

	t, err := c.TickersHistory.Get(args)

	if err != nil {
		log.Fatalf("get single pair: %s", err)
	}

	spew.Dump("get single pair:> ", t)
}

func getSinglePairWithLimit(c *rest.Client) {
	args := rest.GetTickerHistPayload{
		Symbols: []string{"tLTCUSD"},
		Limit:   10,
	}

	t, err := c.TickersHistory.Get(args)

	if err != nil {
		log.Fatalf("get single pair with limit: %s", err)
	}

	spew.Dump("get single pair with limit:> ", t)
}

func getMultiPair(c *rest.Client) {
	args := rest.GetTickerHistPayload{
		Symbols: []string{"tLTCUSD", "tBTCUSD"},
	}

	t, err := c.TickersHistory.Get(args)

	if err != nil {
		log.Fatalf("get multi pair: %s", err)
	}

	spew.Dump("get multi pair:> ", t)
}

func getForRange(c *rest.Client) {
	now := time.Now()
	yesterday := now.Add(time.Duration(-24) * time.Hour)
	end := now.UnixNano() / 1000000
	start := yesterday.UnixNano() / 1000000

	args := rest.GetTickerHistPayload{
		Symbols: []string{"tBTCUSD"},
		Limit:   5,
		Start:   start,
		End:     end,
	}

	t, err := c.TickersHistory.Get(args)

	if err != nil {
		log.Fatalf("get for range: %s", err)
	}

	spew.Dump("get for range:> ", t)
}
