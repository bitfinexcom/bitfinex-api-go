package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go"
)

func main() {
	client := bitfinex.NewClient().Auth("api-key", "api-secret")
	info, err := client.Account.Info()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}
