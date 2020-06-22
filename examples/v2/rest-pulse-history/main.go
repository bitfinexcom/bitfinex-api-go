package main

import (
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()

	now := time.Now()
	millis := now.UnixNano() / 1000000
	end := bitfinex.Mts(millis)

	pulseHist, err := c.Pulse.PublicPulseHistory(0, end)
	if err != nil {
		log.Fatalf("PublicPulseHistory: %s", err)
	}

	spew.Dump(pulseHist)
}
