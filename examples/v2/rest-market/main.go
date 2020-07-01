package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	averagePrice(c)
	foreignExchangeRate(c)
}

func averagePrice(c *rest.Client) {
	args := rest.AveragePriceRequest{
		Symbol: "fUSD",
		Amount: "100",
		Period: 2,
	}

	avgPrice, err := c.Market.AveragePrice(args)
	if err != nil {
		log.Fatalf("AveragePrice err: %s", err)
	}

	spew.Dump(avgPrice)
}

func foreignExchangeRate(c *rest.Client) {
	args := rest.ForeignExchangeRateRequest{
		FirstCurrency:  "BTC",
		SecondCurrency: "USD",
	}

	fxRate, err := c.Market.ForeignExchangeRate(args)
	if err != nil {
		log.Fatalf("ForeignExchangeRate err: %s", err)
	}

	spew.Dump(fxRate)
}
