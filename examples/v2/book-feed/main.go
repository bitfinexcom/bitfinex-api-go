package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func main() {
	client := websocket.New()
	err := client.Connect()
	if err != nil {
		log.Printf("could not connect: %s", err.Error())
		return
	}
	go func() {
		for msg := range client.Listen() {
			log.Printf("recv: %#v", msg)
			if _, ok := msg.(*websocket.InfoEvent); ok {
				_, err := client.SubscribeBook(context.Background(), "BTCUSD", common.Precision0, common.FrequencyRealtime, 1)
				if err != nil {
					log.Printf("could not subscribe to book: %s", err.Error())
				}
			}
		}
	}()
	done := make(chan bool, 1)
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	go func() {
		<-interrupt
		client.Close()
		done <- true
		os.Exit(0)
	}()
	<-done
}
