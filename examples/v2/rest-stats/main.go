package main

import (
	"fmt"
	"log"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
)


func main() {
	c := rest.NewClient()
	pLStats, err := c.Stats.PositionLast("tBTCUSD", bitfinex.Long)
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(pLStats)

	pHStats, err := c.Stats.PositionHistory("tBTCUSD", bitfinex.Long)
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(pHStats)

	scsStats, err := c.Stats.SymbolCreditSizeLast("fUSD", "tBTCUSD")
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(scsStats)

	scsHistStats, err := c.Stats.SymbolCreditSizeHistory("fUSD", "tBTCUSD")
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(scsHistStats)

	fStats, err := c.Stats.FundingLast("fUSD")
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(fStats)

	fhStats, err := c.Stats.FundingHistory("fUSD")
	if err != nil {
		log.Fatalf("getting getting last position stats: %s", err)
	}
	fmt.Println(fhStats)
}
