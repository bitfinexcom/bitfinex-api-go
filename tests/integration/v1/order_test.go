package tests

// Test api client return error messages.
/*
// This test always fails unless the caller provides a X-BFX-APIKEY.  Commented in favor of mocked tests.
func TestFailCreateOrder(t *testing.T) {
	_, err := client.Orders.Create("BTCUSD", 1, 299.0, bitfinex.OrderTypeExchangeLimit)

	if err.Error() != "POST https://api.bitfinex.com/v1/order/new: 400 Invalid order: not enough exchange balance for 1.0 BTCUSD at 299.0" {
		t.Fatalf("OrderBook.Get() returned error: %v", err)
	}
}
*/
