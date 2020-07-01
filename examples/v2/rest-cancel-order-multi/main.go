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
// Below you can see different variations of using CancelOrderMulti function

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClient().
		Credentials(key, secret)

	cancelAllOrders(c)
	cancelByOrderIDs(c)
	cancelByGroupOrderIDs(c)
	cancelByClientOrderID(c)
	cancelOrderMix(c)
}

func cancelAllOrders(c *rest.Client) {
	args := rest.CancelOrderMultiRequest{All: 1}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelAllOrders error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByOrderIDs(c *rest.Client) {
	args := rest.CancelOrderMultiRequest{OrderIDs: rest.OrderIDs{1189452509, 1189452510, 1189452511}}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByOrderIDs error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByGroupOrderIDs(c *rest.Client) {
	args := rest.CancelOrderMultiRequest{GroupOrderIDs: rest.GroupOrderIDs{888, 777}}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByGroupIDs error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByClientOrderID(c *rest.Client) {
	args := rest.CancelOrderMultiRequest{
		ClientOrderIDs: rest.ClientOrderIDs{
			{123, "2020-06-24"},
			{321, "2020-06-24"},
		},
	}

	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByClientOrderID error: %s", err)
	}

	spew.Dump(resp)
}

func cancelOrderMix(c *rest.Client) {
	args := rest.CancelOrderMultiRequest{
		OrderIDs:      rest.OrderIDs{1189452432, 1189452435},
		GroupOrderIDs: rest.GroupOrderIDs{888},
		ClientOrderIDs: rest.ClientOrderIDs{
			{321, "2020-06-24"},
		},
	}

	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelOrderMix error: %s", err)
	}

	spew.Dump(resp)
}
