package main

import (
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
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

	cancelOrdersMultiOp(c)
	cancelOrderMultiOp(c)
	orderMultiOp(c)
}

func cancelOrdersMultiOp(c *rest.Client) {
	resp, err := c.Orders.CancelOrdersMultiOp(rest.OrderIDs{1189452506, 1189452507})
	if err != nil {
		log.Fatalf("CancelOrdersMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func cancelOrderMultiOp(c *rest.Client) {
	resp, err := c.Orders.CancelOrderMultiOp(1189502586)
	if err != nil {
		log.Fatalf("CancelOrderMultiOp error: %s", err)
	}

	spew.Dump(resp)
}

func orderMultiOp(c *rest.Client) {
	pld := rest.OrderMultiArgs{
		Ops: [][]interface{}{
			{
				"on",
				bitfinex.OrderNewRequest{
					CID:    987,
					GID:    876,
					Type:   "EXCHANGE LIMIT",
					Symbol: "tBTCUSD",
					Price:  13,
					Amount: 0.001,
				},
			},
			{
				"oc",
				map[string]int{"id": 1189502430},
			},
			{
				"oc_multi",
				map[string][]int{"id": rest.OrderIDs{1189502431, 1189502432}},
			},
			{
				"ou",
				bitfinex.OrderUpdateRequest{
					ID:     1189502433,
					Price:  15,
					Amount: 0.002,
				},
			},
		},
	}

	resp, err := c.Orders.OrderMultiOp(pld)
	if err != nil {
		log.Fatalf("OrderMultiOp error: %s", err)
	}

	spew.Dump(resp)
}
