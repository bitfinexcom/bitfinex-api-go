package tests

import (
	"testing"
)

func TestOrderBook(t *testing.T) {
	order_book, err := client.OrderBook.Get("BTCUSD", 10, 10, true)

	if err != nil {
		t.Fatalf("OrderBook.Get() returned error: %v", err)
	}

	if len(order_book.Bids) == 0 {
		t.Fatal("Order book should contain Bids")
	}
	if len(order_book.Asks) == 0 {
		t.Fatal("Order book should contain Asks")
	}
}
