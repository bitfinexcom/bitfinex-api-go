package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

func main() {
	c := rest.NewClient()

	// calculate start and end
	now := time.Now()
	millis := now.UnixNano() / 1000000
	prior := now.Add(time.Duration(-24) * time.Hour)
	millisStart := prior.UnixNano() / 1000000
	start := common.Mts(millisStart)
	end := common.Mts(millis)
	// send request
	trades, err := c.Trades.PublicHistoryWithQuery("tBTCUSD", start, end, 10, common.OldestFirst)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, trade := range trades.Snapshot {
		fmt.Println(trade)
	}
}
