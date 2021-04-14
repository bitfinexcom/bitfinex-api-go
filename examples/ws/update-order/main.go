package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux"
)

func main() {
	m := mux.
		New().
		TransformRaw().
		WithAPIKEY("YOUR_API_KEY").
		WithAPISEC("YOUR_API_SECRET").
		Start()

	crash := make(chan error)
	auth := make(chan bool)

	go func() {
		// if listener will fail, program will exit by passing error to crash channel
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				if v.Event == "auth" && v.Status == "OK" {
					// notify auth channel about successful login
					auth <- true
				}
			case order.New:
				// new order received, can update it now
				log.Printf("%T: %+v\n", v, v)
				m.Send(&order.UpdateRequest{
					ID:     v.ID,
					Amount: 0.002,
				})
			case order.Update:
				// order update performed, exiting
				log.Printf("%T: %+v\n", v, v)
				close(crash)
			}
		})
	}()

	for {
		select {
		case err := <-crash:
			fmt.Printf("err: %s\n", err)
			os.Exit(1)
		case <-auth:
			// authenticated, safe to submit orders etc
			if err := m.Send(&order.NewRequest{
				CID:    788,
				Type:   "EXCHANGE LIMIT",
				Symbol: "tBTCUSD",
				Price:  33,
				Amount: 0.001,
			}); err != nil {
				fmt.Printf("err submitting new order: %s\n", err)
			}
		}
	}
}
