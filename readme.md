# Bitfinex Trading API for Golang. Bitcoin, Ether and Litecoin trading
* Official implementation
* REST API
* WebSockets API 

## Installation

``` bash
go get github.com/bitfinexcom/bitfinex-api-go
```

## Usage

### Basic requests

``` go
package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
)

func main() {
	client := bitfinex.NewClient().Auth("api-key", "api-secret")
	info, err := client.Account.Info()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}
```

### Authentication

``` go
func main() {
	client := bitfinex.NewClient().Auth("api-key", "api-secret")
}
```

### Order create

``` go
order, err := client.Orders.Create(bitfinex.BTCUSD, -0.01, 260.99, bitfinex.ORDER_TYPE_EXCHANGE_LIMIT)

if err != nil {
    return err
} else {
    return order
}
```

See [examples](https://github.com/bitfinexcom/bitfinex-api-go/tree/master/examples) and [doc.go](https://github.com/bitfinexcom/bitfinex-api-go/blob/master/doc.go) for more examples.

## Testing

All integration tests are stored in `tests/integration` directory. Because these tests are running using live data, there is a much higher probability of false positives in test failures due to network issues, test data having been changed, etc.

Run tests using:
``` bash
export BFX_API_KEY="api-key"
export BFX_API_SECRET="api-secret"
go test -v ./tests/integration
```

## Contributing

1. Fork it (https://github.com/bitfinexcom/bitfinex-api-go/fork)
2. Create your feature branch (`git checkout -b my-new-feature)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
