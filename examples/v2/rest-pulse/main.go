package main

import (
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

func main() {
	c := rest.NewClient()

	profile, err := c.Pulse.PublicPulseProfile("Bitfinex")
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("got public pulse profile: %+v\n", profile)
}
