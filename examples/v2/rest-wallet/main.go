package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)


// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_KEY")
	secret := os.Getenv("BFX_SECRET")
	c := rest.NewClientWithURL("https://test.bitfinex.com/v2/").Credentials(key, secret)

	wallets, err := c.Wallet.Wallet()
	if err != nil {
		log.Fatalf("getting wallet %s", err)
	}
	spew.Dump(wallets)

	notfication, err := c.Wallet.Transfer("exchange", "margin", "BTC", "BTC", 0.1)
	if err != nil {
		panic(err)
	}
	fmt.Println(notfication)

	notfication2, err := c.Wallet.DepositAddress("exchange", "bitcoin")
	if err != nil {
		panic(err)
	}
	fmt.Println(notfication2)
}
