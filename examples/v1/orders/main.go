package main

import (
	"fmt"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/v1"
)

// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

// WARNING: IF YOU RUN THIS EXAMPLE WITH A VALID KEY ON PRODUCTION
//          IT WILL SUBMIT AN ORDER !

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	client := bitfinex.NewClient().Auth(key, secret)

	// Sell 0.01BTC at $12.000
	data, err := client.Orders.Create(bitfinex.BTCUSD, -0.01, 12000, bitfinex.ORDER_TYPE_EXCHANGE_LIMIT)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", data)
	}
}
