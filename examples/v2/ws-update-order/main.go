package main

import (
	"log"
	"os"
	"time"
	"context"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func SubmitTestOrder(c *websocket.Client) {
	log.Printf("Submitting new order")
	err := c.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
		Symbol: "tBTCUSD",
		CID:    123,
		Amount: 0.02,
		Type: 	"EXCHANGE LIMIT",
		Price:  5000,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateTestOrder(orderId int64, c *websocket.Client) {
	log.Printf("Updating order")
	err := c.SubmitUpdateOrder(context.Background(), &bitfinex.OrderUpdateRequest{
		ID: orderId,
		Amount: 0.04,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	key := os.Getenv("BFX_KEY")
	secret := os.Getenv("BFX_SECRET")
	p := websocket.NewDefaultParameters()
	p.URL = "wss://test.bitfinex.com/ws/2"
	c := websocket.NewWithParams(p).Credentials(key, secret)

	err := c.Connect()
	if err != nil {
		log.Fatalf("connecting authenticated websocket: %s", err)
	}
	defer c.Close()

	// Begin listening to incoming messages

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
			break
		case *websocket.AuthEvent:
			// on authorize create new order
			SubmitTestOrder(c)
		case *bitfinex.OrderNew:
			// new order received so update it
			id := obj.(*bitfinex.OrderNew).ID
			UpdateTestOrder(id, c)
		default:
			log.Printf("MSG RECV: %#v", obj)
		}
	}

	time.Sleep(time.Second * 10)
}
