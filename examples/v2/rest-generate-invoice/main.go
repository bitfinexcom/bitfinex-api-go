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
		NewClient().
		Credentials(key, secret)

	args := rest.DepositInvoiceRequest{
		Currency: "LNX",
		Wallet:   "exchange",
		Amount:   "0.002",
	}

	resp, err := c.Invoice.GenerateInvoice(args)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
