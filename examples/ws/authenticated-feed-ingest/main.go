package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux"
)

func main() {
	m := mux.New().
		TransformRaw().
		WithAPIKEY("YOUR_API_KEY").
		WithAPISEC("YOUR_API_SECRET").
		Start()

	crash := make(chan error)

	go func() {
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				log.Printf("%T: %+v\n", v, v)
			case *order.New:
				log.Printf("%T: %+v\n", v, v)
			case *order.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T item: %+v\n", ss, ss)
				}
			case *order.Update:
				log.Printf("%T: %+v\n", v, v)
			case *order.Cancel:
				log.Printf("%T: %+v\n", v, v)
			case *wallet.Update:
				log.Printf("%T: %+v\n", v, v)
			case *wallet.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T item: %+v\n", ss, ss)
				}
			case balanceinfo.Update:
				log.Printf("%T: %+v\n", v, v)
			default:
				log.Printf("raw/unhandled: %T: %+v\n", v, v)
			}
		})
	}()

	log.Fatal(<-crash)
}
