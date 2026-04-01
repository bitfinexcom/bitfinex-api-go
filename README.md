# Bitfinex Trading Library for Go

A Go reference implementation of the Bitfinex API for both REST and websocket interaction.

### Features
* REST V1/V2 and Websocket V2
* Connection multiplexing with automatic load balancing
* Types for all data schemas
* Safe concurrent close and reconnection handling
* Go 1.21+

## Installation

```bash
go get github.com/bitfinexcom/bitfinex-api-go
```

## Quickstart

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
    "github.com/bitfinexcom/bitfinex-api-go/v2/rest"
)

func main() {
    key := os.Getenv("BFX_API_KEY")
    secret := os.Getenv("BFX_API_SECRET")
    client := rest.NewClient().Credentials(key, secret)

    // Check platform status
    available, err := client.Platform.Status()
    if err != nil || !available {
        log.Fatal("platform not available")
    }

    // Submit order
    response, err := client.Orders.SubmitOrder(&order.NewRequest{
        Symbol: "tBTCUSD",
        CID:    time.Now().Unix() / 1000,
        Amount: 0.02,
        Type:   "EXCHANGE LIMIT",
        Price:  5000,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response)
}
```

## Available REST V2 Services

| Service | Description |
|---------|-------------|
| `Platform` | Platform operational status |
| `Tickers` | Ticker data (single, multi, all) |
| `TickersHistory` | Historical ticker data |
| `Book` | Order book with configurable precision |
| `Candles` | OHLCV candlestick data |
| `Trades` | Public and authenticated trade history |
| `Stats` | Funding, credit, position, volume statistics |
| `Status` | Derivative status (current and historical) |
| `Currencies` | Platform configuration and currency info |
| `Market` | Average execution price, FX rates |
| `Orders` | Active orders, history, submit/update/cancel, multi-ops |
| `Positions` | Active positions, history, audit, snapshots, claim, increase |
| `Wallet` | Wallets, transfer, deposit address, withdraw |
| `Ledgers` | Ledger history |
| `Funding` | Offers, loans, credits, trades, auto-renew, keep, close |
| `Derivatives` | Collateral set/limits |
| `Invoice` | Lightning Network deposit invoices |
| `Account` | User info, summary, margin info, available balance, login/audit history |
| `Movements` | Deposit/withdrawal history |
| `Alerts` | Price alert list, set, delete |
| `Settings` | User settings read/write/delete |
| `Liquidations` | Historical liquidation data |
| `Rankings` | Leaderboard standings |

## Docs

* **[V1](docs/v1.md)** - Legacy API documentation (all methods prefixed with `V1`)
* **[V2 Rest](docs/rest_v2.md)** - REST V2 documentation
* **[V2 Websocket](docs/ws_v2.md)** - Websocket V2 documentation

## Examples

#### Authentication

```go
client := rest.NewClient().Credentials("API_KEY", "API_SEC")
```

#### Subscribe to Trades (Websocket)

```go
// using github.com/bitfinexcom/bitfinex-api-go/v2/websocket
_, err := client.SubscribeTrades(context.Background(), "tBTCUSD")
if err != nil {
    log.Printf("could not subscribe to trades: %s", err.Error())
}
```

#### Get Order History (REST)

```go
// using github.com/bitfinexcom/bitfinex-api-go/v2/rest
os, err := client.Orders.AllHistory()
if err != nil {
    log.Fatalf("getting orders: %s", err)
}
```

#### Get Account Info (REST)

```go
info, err := client.Account.UserInfo()
if err != nil {
    log.Fatalf("getting account info: %s", err)
}
```

See the **[examples](https://github.com/bitfinexcom/bitfinex-api-go/tree/master/examples)** directory for more, including:

- [Creating/updating an order](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/ws-update-order/main.go)
- [Subscribing to orderbook updates](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/book-feed/main.go)
- [Integrating a custom logger](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/ws-custom-logger/main.go)
- [Submitting funding offers](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/rest-funding/main.go)
- [Retrieving active positions](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/examples/v2/rest-positions/main.go)

## V1 API Note

The V1 API is fully functional but considered legacy. All V1 service methods are prefixed with `V1` for explicit version identification:

```go
import "github.com/bitfinexcom/bitfinex-api-go/v1"

client := bitfinex.NewClient().Auth(key, secret)

// V1 methods use V1 prefix
ticker, err := client.Ticker.V1Get("btcusd")
orders, err := client.Orders.V1All()
```

## FAQ

### Is there any rate limiting?

For a Websocket connection there is no limit to the number of requests sent down the connection (unlimited order operations) however an account can only create 15 new connections every 5 mins and each connection is only able to subscribe to 30 inbound data channels. This library handles all of the load balancing/multiplexing for channels and will automatically create/destroy new connections when needed.

For REST the base limit per-user is 1,000 orders per 5 minute interval, and is shared between all account API connections. It increases proportionally to your trade volume based on the following formula:

```
1000 + (TOTAL_PAIRS_PLATFORM * 60 * 5) / (250000000 / USER_VOL_LAST_30d)
```

### Will I always receive an `on` packet?

No; if your order fills immediately, the first packet referencing the order will be an `oc` signaling the order has closed. If the order fills partially immediately after creation, an `on` packet will arrive with a status of `PARTIALLY FILLED...`

### nonce too small

If you make multiple parallel requests and receive a nonce error, note that nonces guard against replay attacks. When multiple HTTP requests arrive at the API with the wrong nonce (e.g. due to async timing), the API will reject the request. For parallel requests, use multiple API keys.

### How do `te` and `tu` messages differ?

A `te` packet is sent first immediately after a trade has been matched & executed, followed by a `tu` message once processing completes. During high load, `tu` may be delayed -- use `te` for realtime feeds.

### What is the difference between R* and P* order books?

Order books with precision `R0` are 'raw' and contain entries for each individual order, whereas `P*` books contain entries for each price level (aggregating orders).

## Contributing

1. Fork it (https://github.com/bitfinexcom/bitfinex-api-go/fork)
2. Create your feature branch (`git checkout -b feat/my-new-feature`)
3. Commit your changes (`git commit -m 'feat: add some feature'`)
4. Push to the branch (`git push origin feat/my-new-feature`)
5. Create a new Pull Request
