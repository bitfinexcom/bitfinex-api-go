package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go"
)

func main() {
	client := bitfinex.NewClient().Auth("api-key", "api-secret")

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
