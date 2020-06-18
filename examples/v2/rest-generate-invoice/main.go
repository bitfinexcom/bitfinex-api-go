package main

import (
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

// Set BFX_API_KEY and BFX_API_SECRET:
//
// export BFX_API_KEY=<your-api-key>
// export BFX_API_SECRET=<your-api-secret>
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClientWithURL("https://api.staging.bitfinex.com/v2/").
		Credentials(key, secret)

	resp, err := c.Invoice.GenerateInvoice("LNX", "exchange", "0.002")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
