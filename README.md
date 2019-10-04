# Bitfinex Trading Library for GoLang  - Bitcoin, Ethereum, Ripple and more

![https://api.travis-ci.org/bitfinexcom/bitfinex-api-go.svg?branch=master](https://api.travis-ci.org/bitfinexcom/bitfinex-api-go.svg?branch=master)

A Golang reference implementation of the Bitfinex API for both REST and websocket interaction.

### Features
* Official implementation
* REST V1/V2 and Websocket
* Connection multiplexing
* Types for all data schemas

## Instillation


``` bash
go get github.com/bitfinexcom/bitfinex-api-go
```

Optional - run the 'trade-feed' example to begin receiving realtime trade updates via the websocket

```bash
cd $GOPATH/src/github.com/bitfinexcom/bitfinex-api-go
go run examples/v2/trade-feed/main.go
```

## Quickstart


``` go
package main

import (
	"log"
	"os"
	"time"
	"context"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

func listenToOutput(c *websocket.Client) {
	for obj := range c.Listen() {
        log.Printf("MSG RECV: %#v", obj)
    }
}

func main() {
	c := websocket.New().Credentials("BFX_KEY", "BFX_SECRET")

	err := c.Connect()
	if err != nil {
		log.Fatalf("Error websocket: %s", err)
	}
	defer c.Close()
    go listenToOutput(c)

	err := c.SubmitOrder(context.Background(), &bitfinex.OrderNewRequest{
        Symbol: "tBTCUSD",
        CID:    time.Now().Unix() / 1000,
        Amount: 0.02,
        Type: 	"EXCHANGE LIMIT",
        Price:  5000,
    })
    if err != nil {
        log.Fatalf("Error creating order: %s", err)
    }

	time.Sleep(time.Second * 20)
}

```

``` go
package main

import (
    "fmt"
    "github.com/bitfinexcom/bitfinex-api-go/v2"
)

func main() {
    client := bitfinex.NewClient().Credentials("API_KEY", "API_SEC")
	
    // create order
    response, err := c.Orders.SubmitOrder(&bitfinex.OrderNewRequest{
        Symbol: "tBTCUSD",
        CID:    time.Now().Unix() / 1000,
        Amount: 0.02,
        Type: 	"EXCHANGE LIMIT",
        Price:  5000,
    })
    if err != nil {
        panic(err)
    }
}
```

## Docs

* <b>[V1](docs/v1.md)</b> - Documentation (depreciated)
* <b>[V2 Rest](docs/rest_v2.md)</b> - Documentation
* <b>[V2 Websocket](docs/ws_v2.md)</b> - Documentation

## Examples

#### Authentication

``` go
func main() {
	client := bitfinex.NewClient().Credentials("API_KEY", "API_SEC")
}
```

#### Subscribe to Trades

``` go
// using github.com/bitfinexcom/bitfinex-api-go/v2/websocket as client
_, err := client.SubscribeTrades(context.Background(), "tBTCUSD")
if err != nil {
    log.Printf("Could not subscribe to trades: %s", err.Error())
}
```

#### Get candles via REST

```go
// using github.com/bitfinexcom/bitfinex-api-go/v2/rest as client
os, err := client.Orders.AllHistory()
if err != nil {
    log.Fatalf("getting orders: %s", err)
}
```


See the <b>[examples](https://github.com/bitfinexcom/bitfinex-api-go/tree/master/examples)</b> directory for more, like:

- [Creating/updating an order](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/ws-update-order/main.go)
- [Subcribing to orderbook updates](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/book-feed/main.go)
- [Integrating a custom logger](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/ws-custom-logger/main.go)
- [Submitting funding offers](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/rest-funding/main.go)
- [Retrieving active positions](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/rest-positions/main.go)

## FAQ

### Is there any rate limiting?

For a Websocket connection there is no limit to the number of requests sent down the connection (unlimited order operations) however an account can only create 15 new connections every 5 mins and each connection is only able to subscribe to 30 inbound data channels. Fortunately this library handles all of the load balancing/multiplexing for channels and will automatically create/destroy new connections when needed, however the user may still encounter the max connections rate limiting error.

For rest the base limit per-user is 1,000 orders per 5 minute interval, and is shared between all account API connections. It increases proportionally to your trade volume based on the following formula:

1000 + (TOTAL_PAIRS_PLATFORM * 60 * 5) / (250000000 / USER_VOL_LAST_30d)

Where TOTAL_PAIRS_PLATFORM is the number of pairs on the Bitfinex platform (currently ~101) and USER_VOL_LAST_30d is in USD.

### Will I always receive an `on` packet?

No; if your order fills immediately, the first packet referencing the order will be an `oc` signaling the order has closed. If the order fills partially immediately after creation, an `on` packet will arrive with a status of `PARTIALLY FILLED...`

For example, if you submit a `LIMIT` buy for 0.2 BTC and it is added to the order book, an `on` packet will arrive via ws2. After a partial fill of 0.1 BTC, an `ou` packet will arrive, followed by a final `oc` after the remaining 0.1 BTC fills.

On the other hand, if the order fills immediately for 0.2 BTC, you will only receive an `oc` packet.

### My websocket won't connect!

Did you call `client.Connect()`? :)

### nonce too small

I make multiple parallel request and I receive an error that the nonce is too small. What does it mean?

Nonces are used to guard against replay attacks. When multiple HTTP requests arrive at the API with the wrong nonce, e.g. because of an async timing issue, the API will reject the request.

If you need to go parallel, you have to use multiple API keys right now.

### How do `te` and `tu` messages differ?

A `te` packet is sent first to the client immediately after a trade has been matched & executed, followed by a `tu` message once it has completed processing. During times of high load, the `tu` message may be noticably delayed, and as such only the `te` message should be used for a realtime feed.

### What are the sequence numbers for?

If you enable sequencing on v2 of the WS API, each incoming packet will have a public sequence number at the end, along with an auth sequence number in the case of channel `0` packets. The public seq numbers increment on each packet, and the auth seq numbers increment on each authenticated action (new orders, etc). These values allow you to verify that no packets have been missed/dropped, since they always increase monotonically.

### What is the difference between R* and P* order books?

Order books with precision `R0` are considered 'raw' and contain entries for each order submitted to the book, whereas `P*` books contain entries for each price level (which aggregate orders).

## Contributing

1. Fork it (https://github.com/bitfinexcom/bitfinex-api-go/fork)
2. Create your feature branch (`git checkout -b my-new-feature)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
