package main

import (
	"log"
	"os"
	"time"

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

	now := time.Now()
	millis := now.UnixNano() / 1000000
	prior := now.Add(time.Duration(-240) * time.Hour)
	millisStart := prior.UnixNano() / 1000000

	l, err := c.Ledgers.Ledgers("BTC", millisStart, millis, 1000)
	if err != nil {
		log.Fatalf("Ledgers: %s", err)
	}

	spew.Dump(l)
}
