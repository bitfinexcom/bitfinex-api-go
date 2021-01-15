package main

import (
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/status"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trades"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux"
)

func main() {
	m := mux.New().
		TransformRaw().
		Start()

	pairs := []string{
		"BTCUSD",
		"BTCUST",
		"BTCXCH",
		"ETHBTC",
		"ETHEUR",
		"ETHGBP",
		"ETHJPY",
		"ETHUSD",
		"ETHUST",
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

		rawBookPld := event.Subscribe{
			Event:     "subscribe",
			Channel:   "book",
			Precision: "R0",
			Symbol:    "t" + pair,
		}

		bookPld := event.Subscribe{
			Event:     "subscribe",
			Channel:   "book",
			Precision: "P0",
			Frequency: "F0",
			Symbol:    "t" + pair,
		}

		m.Subscribe(tradePld)
		m.Subscribe(tickPld)
		m.Subscribe(candlesPld)
		m.Subscribe(rawBookPld)
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

	fundingPairTrade := event.Subscribe{
		Event:   "subscribe",
		Channel: "trades",
		Symbol:  "fUSD",
	}

	m.Subscribe(derivStatusPld)
	m.Subscribe(liqStatusPld)
	m.Subscribe(fundingPairTrade)

	crash := make(chan error)

	go func() {
		crash <- m.Listen(func(msg interface{}, err error) {
			if err != nil {
				log.Printf("non crucial error received: %s\n", err)
			}

			switch v := msg.(type) {
			case event.Info:
				log.Printf("%T: %+v\n", v, v)
			case trades.TradeSnapshot:
				log.Printf("%T: %+v\n", v, v)
			case trades.FundingTradeSnapshot:
				log.Printf("%T: %+v\n", v, v)
			case trades.TradeUpdate:
				log.Printf("%T: %+v\n", v, v)
			case trades.TradeExecuted:
				log.Printf("%T: %+v\n", v, v)
			case trades.FundingTradeUpdate:
				log.Printf("%T: %+v\n", v, v)
			case trades.FundingTradeExecuted:
				log.Printf("%T: %+v\n", v, v)
			case *ticker.Ticker:
				log.Printf("%T: %+v\n", v, v)
			case *ticker.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case *book.Book:
				log.Printf("%T: %+v\n", v, v)
			case *book.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case *candle.Candle:
				log.Printf("%T: %+v\n", v, v)
			case *candle.Snapshot:
				log.Printf("%T: %+v\n", v, v)
			case *status.Derivative:
				log.Printf("%T: %+v\n", v, v)
			case *status.DerivativesSnapshot:
				log.Printf("%T: %+v\n", v, v)
			case *status.Liquidation:
				log.Printf("%T: %+v\n", v, v)
			case *status.LiquidationsSnapshot:
				log.Printf("%T: %+v\n", v, v)
			default:
				log.Printf("raw/unrecognized msg: %T: %s\n", v, v)
			}
		})
	}()

	log.Fatal(<-crash)
}
