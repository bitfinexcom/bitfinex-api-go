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

	pulseHist, err := c.Pulse.PulseHistory(true)
	if err != nil {
		log.Fatalf("PulseHistory: %s", err)
	}

	spew.Dump(pulseHist)
}
