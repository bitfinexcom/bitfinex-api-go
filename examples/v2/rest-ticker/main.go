package main

import (
	"log"
	bfx "github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)


func main() {
	c := rest.NewClient()
	symbols := []string{bfx.TradingPrefix+bfx.BTCUSD, bfx.TradingPrefix+bfx.EOSBTC}
	tickers, err := c.Tickers.GetMulti(symbols)

	if err != nil {
		log.Fatalf("getting ticker: %s", err)
	}

	log.Print(tickers)
}
