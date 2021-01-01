package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux"
)

func main() {
	m := mux.New().
		TransformRaw().
		Start()

	pairs := []string{}
	dat, err := ioutil.ReadFile("./testpairs.json")
	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(dat, &pairs); err != nil {
		log.Panic(err)
	}

	for _, pair := range pairs {
		tradePld := event.Subscribe{
			Event:   "subscribe",
			Channel: "trades",
			Symbol:  "t" + pair,
		}

		tickPld := event.Subscribe{
			Event:   "subscribe",
			Channel: "ticker",
			Symbol:  "t" + pair,
		}

		candlesPld := event.Subscribe{
			Event:   "subscribe",
			Channel: "candles",
			Key:     "trade:1m:t" + pair,
		}

		bookPld := event.Subscribe{
			Event:     "subscribe",
			Channel:   "book",
			Precision: "R0",
			Symbol:    "t" + pair,
		}

		m.Subscribe(tradePld)
		m.Subscribe(tickPld)
		m.Subscribe(candlesPld)
		m.Subscribe(bookPld)
	}

	derivStatusPld := event.Subscribe{
		Event:   "subscribe",
		Channel: "status",
		Key:     "deriv:tBTCF0:USTF0",
	}

	liqStatusPld := event.Subscribe{
		Event:   "subscribe",
		Channel: "status",
		Key:     "liq:global",
	}

	m.Subscribe(derivStatusPld)
	m.Subscribe(liqStatusPld)

	crash := make(chan error)

	go func() {
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("non crucial error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				log.Printf("%T: %+v\n", v, v)
			case *trade.Trade:
				log.Printf("%T: %+v\n", v, v)
			case *trade.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			case *ticker.Ticker:
				log.Printf("%T: %+v\n", v, v)
			case *ticker.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			case *book.Book:
				log.Printf("%T: %+v\n", v, v)
			case *book.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			case *candle.Candle:
				log.Printf("%T: %+v\n", v, v)
			case *candle.Snapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			case *status.Derivative:
				log.Printf("%T: %+v\n", v, v)
			case *status.DerivativesSnapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			case *status.Liquidation:
				log.Printf("%T: %+v\n", v, v)
			case *status.LiquidationsSnapshot:
				log.Printf("%T: %+v\n", v, v)
				for _, ss := range v.Snapshot {
					log.Printf("%T snapshot: %+v\n", ss, ss)
				}
			default:
				log.Printf("unrecognized msg: %T: %s\n", v, v)
			}
		})
	}()

	log.Fatal(<-crash)
}
