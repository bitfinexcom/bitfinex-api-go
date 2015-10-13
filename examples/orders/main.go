package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go"
)

func main() {
	client := bitfinex.NewClient().Auth("api-key", "api-secret")

	// Sell 0.01 at price 260.99
	data, err := client.Orders.Create(bitfinex.BTCUSD, -0.01, 260.99, bitfinex.ORDER_TYPE_EXCHANGE_LIMIT)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", data)
	}
}
