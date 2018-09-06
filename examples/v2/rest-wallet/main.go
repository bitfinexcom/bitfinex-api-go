package main

import (
"flag"
"log"
"os"
"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)


// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	flag.Parse()

	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	c := rest.NewClient().Credentials(key, secret)

	wallets, err := c.Wallet.Wallet()


	if err != nil {
		log.Fatalf("getting wallet %s", err)
	}

	spew.Dump(wallets)



}

