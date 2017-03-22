package tests

import (
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"testing"
)

// Test api client return error messages.
func TestFailCreateOrder(t *testing.T) {
	_, err := client.Orders.Create("BTCUSD", 1, 299.0, bitfinex.ORDER_TYPE_EXCHANGE_LIMIT)

	if err.Error() != "POST https://api.bitfinex.com/v1/order/new: 400 Invalid order: not enough exchange balance for 1.0 BTCUSD at 299.0" {
		t.Fatalf("OrderBook.Get() returned error: %v", err)
	}
}
