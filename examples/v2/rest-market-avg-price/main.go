package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	args := rest.AveragePriceArgs{
		Symbol: "fUSD",
		Amount: "100",
		Period: 2,
	}

	avgPrice, err := c.Market.AveragePrice(args)
	if err != nil {
		log.Fatalf("AveragePrice: %s", err)
	}

	spew.Dump(avgPrice)
}
