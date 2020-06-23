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
	cancelByGroupIDs(c)
	cancelByClientOrderID(c)
	cancelOrderMix(c)
}

func cancelAllOrders(c *rest.Client) {
	args := rest.CancelOrderMultiArgs{All: 1}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelAllOrders error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByOrderIDs(c *rest.Client) {
	args := rest.CancelOrderMultiArgs{ID: []int{1189413241, 1189413242}}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByOrderIDs error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByGroupIDs(c *rest.Client) {
	args := rest.CancelOrderMultiArgs{GID: []int{888, 777}}
	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByGroupIDs error: %s", err)
	}

	spew.Dump(resp)
}

func cancelByClientOrderID(c *rest.Client) {
	args := rest.CancelOrderMultiArgs{
		CID: [][]interface{}{
			{123, "2020-06-23"},
			{321, "2020-06-23"},
		},
	}

	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelByClientOrderID error: %s", err)
	}

	spew.Dump(resp)
}

func cancelOrderMix(c *rest.Client) {
	args := rest.CancelOrderMultiArgs{
		ID:  []int{1189412946},
		GID: []int{888},
		CID: [][]interface{}{
			{321, "2020-06-23"},
		},
	}

	resp, err := c.Orders.CancelOrderMulti(args)
	if err != nil {
		log.Fatalf("cancelOrderMix error: %s", err)
	}

	spew.Dump(resp)
}
