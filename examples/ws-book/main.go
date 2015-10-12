package main

import (
	"github.com/bitfinexcom/bitfinex-api-go"
	"log"
)

// This example shows how to work with WebSocket book api
func main() {
	api := bitfinex.NewClient()
	err := api.WebSocket.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket")
	}
	defer api.WebSocket.Close()

	book_btcusd_chan := make(chan []float64)
	book_ltcusd_chan := make(chan []float64)
	trades_chan := make(chan []float64)
	ticker_chan := make(chan []float64)

	go api.WebSocket.Watch(bitfinex.CHAN_BOOK, bitfinex.BTCUSD, book_btcusd_chan)
	go api.WebSocket.Watch(bitfinex.CHAN_BOOK, bitfinex.LTCUSD, book_ltcusd_chan)
	go api.WebSocket.Watch(bitfinex.CHAN_TRADE, bitfinex.BTCUSD, trades_chan)
	go api.WebSocket.Watch(bitfinex.CHAN_TICKER, bitfinex.BTCUSD, ticker_chan)

	// After api client successfully connect to remote web socket
	// channel will reveive current payload as separate messages.
	// Each channel will receive order book updates: [PRICE, COUNT, ±AMOUNT]
	for {
		select {
		case btcusd_msg := <-book_btcusd_chan:
			log.Println("BOOK BTCUSD:", btcusd_msg)
		case ltcusd_msg := <-book_ltcusd_chan:
			log.Println("BOOK LTCUSD:", ltcusd_msg)
		case trade_msg := <-trades_chan:
			log.Println("TRADES:", trade_msg)
		case ticker_msg := <-ticker_chan:
			log.Println("TICKER:", ticker_msg)
		}
	}
}
