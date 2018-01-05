package main

import (
	"log"
	"os"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
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
	c := websocket.NewClient().Credentials(key, secret)

	err := c.Connect()
	if err != nil {
		log.Fatalf("connecting authenticated websocket: %s", err)
	}
	go func() {
		for msg := range c.Listen() {
			if msg == nil {
				break
			}
			log.Printf("MSG RECV: %#v", msg)
		}
	}()

	//ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	//c.Authenticate(ctx)

	time.Sleep(time.Second * 10)
}
