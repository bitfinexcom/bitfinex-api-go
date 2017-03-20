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

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
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
