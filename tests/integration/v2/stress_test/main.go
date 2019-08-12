package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
	"github.com/op/go-logging"
	"time"
)

var symbols = []string{
	"tBTCUSD", "tETHUSD", "tEOSUSD", "tVETUSD", "tDGBUSD", "tXRPBTC", "tTRXUSD",
	"tLEOUSD", "tLEOBTC", "tLEOUST", "tBTCUST", "tETHBTC", "tETHUST", "tXRPUSD",
}

func main() {
	useChannels := flag.Bool("channels", false, "subscribes to a lot of channels")
	useMultiplexor := flag.Bool("multiplexer", false, "subscribes/un-subscribes forcing ws connection re-shuffling")
	runTime := flag.Int64("time", 5, "runtime of the stress test in minutes")
	logMode := flag.String("log-level", "INFO", "level of logging. Can be INFO, DEBUG, WARN or ERROR.")
	flag.Parse()

	p := websocket.NewDefaultParameters()
	logger := logging.MustGetLogger("bfx-websocket")
	p.ManageOrderbook = true
	p.Logger = logger
	logLevel, err := logging.LogLevel(*logMode)
	if err != nil {
		panic(err)
	}
	logging.SetLevel(logLevel, "bfx-websocket")
	c := websocket.NewWithParams(p)
	err = c.Connect()
	if err != nil {
		panic(err)
	}

	quit := make(chan interface{})
	if (*useMultiplexor) {
		go runMultiplexerStressTest(c, logger, quit)
	} else if (*useChannels) {
		go runChannelsStressTest(c, logger, quit)
	} else {
		fmt.Println("No flags set, use:")
		flag.PrintDefaults()
	}

	timeout := time.After(time.Minute * time.Duration(*runTime))
	for {
		select {
		case msg := <- c.Listen():
			switch msg.(type) {
			case error:
				logger.Error(msg)
			default:
				logger.Debugf("MSG RECV: %#v", msg)
			}
		case <- timeout:
			logger.Warningf("Test timeout of %d mins reached. Killing test.", *runTime)
			panic("time reached")
		case <-quit:
			return
		}
	}
}

func runChannelsStressTest(client *websocket.Client, log *logging.Logger, quit chan interface{}) {
	fmt.Println("Starting channel stress test...")
	for _, ticker := range symbols {
		log.Infof("Subscribing to trades (%s) (socketCount=%d)", ticker, client.ConnectionCount())
		_, err := client.SubscribeTrades(context.Background(), ticker)
		if err != nil {
			panic(fmt.Sprintf("could not subscribe to trades: %s", err.Error()))
		}
		_, err = client.SubscribeCandles(context.Background(), ticker, bitfinex.FifteenMinutes)
		if err != nil {
			panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.FifteenMinutes))
		}
		_, err = client.SubscribeCandles(context.Background(), ticker, bitfinex.ThirtyMinutes)
		if err != nil {
			panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.ThirtyMinutes))
		}
		_, err = client.SubscribeCandles(context.Background(), ticker, bitfinex.OneHour)
		if err != nil {
			panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.OneHour))
		}
		_, err = client.SubscribeCandles(context.Background(), ticker, bitfinex.OneMinute)
		if err != nil {
			panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.OneMinute))
		}
	}
}

func runMultiplexerStressTest(client *websocket.Client, log *logging.Logger, quit chan interface{}) {
	fmt.Println("Starting multiplexer stress test...")
	for {
		subIds := make([]string, 0)
		for _, ticker := range symbols {
			log.Infof("Subscribing to trades (%s) (socketCount=%d)", ticker, client.ConnectionCount())
			subId1, err := client.SubscribeTrades(context.Background(), ticker)
			if err != nil {
				panic(fmt.Sprintf("could not subscribe to trades: %s", err.Error()))
			}
			subIds = append(subIds, subId1)
			subId2, err := client.SubscribeCandles(context.Background(), ticker, bitfinex.FifteenMinutes)
			if err != nil {
				panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.FifteenMinutes))
			}
			subIds = append(subIds, subId2)
			subId3, err := client.SubscribeBook(context.Background(), ticker, bitfinex.Precision0, bitfinex.FrequencyRealtime, 25)
			if err != nil {
				panic(fmt.Sprintf("could not subscribe to candles %s: %s", err.Error(), bitfinex.FifteenMinutes))
			}
			subIds = append(subIds, subId3)
		}
		// wait for a set amount of time before un-subscribing
		time.Sleep(time.Second * 10)

		// un-subscribe from all channels
		for _, subId := range subIds {
			log.Infof("Un-subscribing from (%s) (socketCount=%d)", subId, client.ConnectionCount())
			err := client.Unsubscribe(context.Background(), subId)
			if err != nil {
				panic(fmt.Sprintf("Could not un-subscribe from channel %s", subId))
			}
		}

		// wait again
		time.Sleep(time.Second * 10)
	}
}
