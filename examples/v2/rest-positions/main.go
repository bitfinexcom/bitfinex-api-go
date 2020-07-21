package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
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
	c := rest.NewClient().Credentials(key, secret)

	all(c)
	claim(c)
}

func all(c *rest.Client) {
	p, err := c.Positions.All()
	if err != nil {
		log.Fatalf("All: %s", err)
	}

	for _, ps := range p.Snapshot {
		fmt.Println(ps)
	}
}

func claim(c *rest.Client) {
	pc, err := c.Positions.Claim(&position.ClaimRequest{
		Id: 36228736,
	})
	if err != nil {
		log.Fatalf("Claim: %s", err)
	}

	spew.Dump(pc)
}
