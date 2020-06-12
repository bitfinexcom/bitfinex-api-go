package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	pulseHist, err := c.Pulse.PublicPulseHistory(0, 0)
	if err != nil {
		log.Fatalf("PublicPulseHistory: %s", err)
	}

	spew.Dump(pulseHist)
}
