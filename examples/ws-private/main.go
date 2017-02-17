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

    dataChan := make(chan bitfinex.TermData)
    go client.WebSocket.ConnectPrivate(dataChan)

    for {
        select {
        case data := <-dataChan:
            if data.HasError() {
                // Data has error - websocket channel will be closed.
                fmt.Println("Error:", data.Error)
                return
            } else {
                fmt.Println("Data:", data)
            }
        }
    }
}
