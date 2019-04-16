package main

import (
	"log"
	bfx "github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"time"
)


func main() {
	c := rest.NewClient()


	log.Printf("1) Query Last Candle")
	candle, err := c.Candles.Last(bfx.TradingPrefix+bfx.BTCUSD, bfx.FiveMinutes)

	if err != nil {
		log.Fatalf("getting candle: %s", err)
	}

	log.Printf("last candle: %#v\n", candle)

	now := time.Now()
	millis := now.UnixNano() / 1000000

	prior := now.Add(time.Duration(-24) * time.Hour)
	millisStart := prior.UnixNano() / 1000000


	log.Printf("2) Query Candle History with no params")
	candles, err := c.Candles.History(bfx.TradingPrefix+bfx.BTCUSD, bfx.FiveMinutes)

	if err != nil {
		log.Fatalf("getting candles: %s", err)
	}

	log.Printf("length of candles is: %v", len(candles.Snapshot))

	log.Printf("first candle is: %#v\n", candles.Snapshot[0])
	log.Printf("last candle is: %#v\n", candles.Snapshot[len(candles.Snapshot)-1])

	start := bfx.Mts(millisStart)
	end := bfx.Mts(millis)

	log.Printf("3) Query Candle History with params")
	candlesMore, err := c.Candles.HistoryWithQuery(
		bfx.TradingPrefix+bfx.BTCUSD,
		bfx.FiveMinutes,
		start,
		end,
		200,
		bfx.OldestFirst,
		)

	if err != nil {
		log.Fatalf("getting candles: %s", err)
	}

	log.Printf("length of candles is: %v", len(candlesMore.Snapshot))
	log.Printf("first candle is: %#v\n", candlesMore.Snapshot[0])
	log.Printf("last candle is: %#v\n", candlesMore.Snapshot[len(candlesMore.Snapshot)-1])



}

