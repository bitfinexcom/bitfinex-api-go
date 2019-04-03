package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

func main() {

	key := "Key"
	secret := "Secret"
	c := rest.NewClient().Credentials(key, secret)

	ledgers, err := c.Ledgers.Ledgers("BTC",1539159010000,1554237570367,200)
	// specify ("CURRENCY", Epochmiliseconds start, Epochmiliseconds stop, limit of ledgers to return max 500)
	if err != nil {
		fmt.Print("getting ledgers %s", err)
	}
	fmt.Printf("%s\n", ledgers)
}
