package main

import (
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

func main() {
	c := rest.NewClient()

	pulseHist, err := c.Pulse.PublicPulseHistory("1", "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("got public pulse message: %+v\n", pulseHist[0])
}
