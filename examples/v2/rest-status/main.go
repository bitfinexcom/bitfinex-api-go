package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()
	derivativeStatus(c)
	derivativeStatusMulti(c)
	derivativeStatusAll(c)
}

func derivativeStatus(c *rest.Client) {
	ds, err := c.Status.DerivativeStatus("tBTCF0:USTF0")
	if err != nil {
		log.Fatalf("derivativeStatus: %s", err)
	}

	spew.Dump(ds)
}

func derivativeStatusMulti(c *rest.Client) {
	ds, err := c.Status.DerivativeStatusMulti([]string{"tBTCF0:USTF0", "tETHF0:USTF0"})
	if err != nil {
		log.Fatalf("derivativeStatusMulti: %s", err)
	}

	spew.Dump(ds)
}

func derivativeStatusAll(c *rest.Client) {
	ds, err := c.Status.DerivativeStatusAll()
	if err != nil {
		log.Fatalf("DerivativeStatusAll: %s", err)
	}

	spew.Dump(ds)
}
