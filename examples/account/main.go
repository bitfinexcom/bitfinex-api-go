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
    info, err := client.Account.Info()

    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(info)
    }
}
