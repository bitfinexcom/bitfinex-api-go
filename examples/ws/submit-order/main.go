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
	m := mux.New().
		TransformRaw().
		WithAPIKEY("URJ0mNaT0oM7cmwHoaWRqwF433urbujGlsVnFnSrYuG").
		WithAPISEC("Ovvs6zU5L2ZXjEXhtKJV3ITrY8KsS6tSLmvSV6Wn9j0").
		Start()

	crash := make(chan error)
	auth := make(chan bool)

	go func() {
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				if v.Event == "auth" && v.Status == "OK" {
					auth <- true
				}
			case *order.New:
				log.Printf("%T: %+v\n", v, v)
			case *order.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T item: %+v\n", ss, ss)
				}
			default:
				log.Printf("raw/unhandled: %T: %+v\n", v, v)
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
