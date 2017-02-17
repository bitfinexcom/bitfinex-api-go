package main

import (
    "fmt"
    "os"

    "github.com/bitfinexcom/bitfinex-api-go"
)

func main() {
    key := os.Getenv("BFX_APIKEY")
    secret := os.Getenv("BFX_SECRET")
    client := bitfinex.NewClient().Auth(key, secret)

    // Sell 0.01BTC at $12.000
    data, err := client.Orders.Create(bitfinex.BTCUSD, -0.01, 12000, bitfinex.ORDER_TYPE_EXCHANGE_LIMIT)

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Response:", data)
    }
}
