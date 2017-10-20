package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
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
	c := bitfinex.NewClient().Credentials(key, secret)

	err := c.Websocket.Connect()
	if err != nil {
		log.Fatalf("connecting authenticated websocket: %s", err)
	}

	c.Websocket.AttachEventHandler(func(ev interface{}) {
		log.Printf("EVENT: %#v", ev)
	})

	c.Websocket.AttachPrivateHandler(func(msg interface{}) {
		log.Printf("PRIV MSG: %#v", msg)
	})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	c.Websocket.Authenticate(ctx)

	time.Sleep(time.Second * 10)
}
