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
//
// Below you can see different variations of using Order Multi ops

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClient().
		Credentials(key, secret)

	cancelOrdersByID(c)
}

func cancelOrdersByID(c *rest.Client) {
	resp, err := c.Orders.CancelOrders(rest.OrderIDs{1189452506, 1189452507})
	if err != nil {
		log.Fatalf("cancelAllOrders error: %s", err)
	}

	spew.Dump(resp)
}
