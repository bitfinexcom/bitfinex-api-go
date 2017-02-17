package main

// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_APIKEY=YOUR_API_KEY
// export BFX_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

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
