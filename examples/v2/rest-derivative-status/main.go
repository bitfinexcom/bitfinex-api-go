package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"log"
)


func main() {
	c := rest.NewClient()
	pLStats, err := c.Status.DerivativeStatus("tBTCF0:USTF0")
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(pLStats)
}

