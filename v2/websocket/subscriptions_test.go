package websocket

import (
	"testing"
	"time"

	"github.com/op/go-logging"
)

func TestSubscriptions_ControlGoroutineStopsOnClose(t *testing.T) {
	log := logging.MustGetLogger("test")
	subs := newSubscriptions(5*time.Second, log)

	// Close should not hang — the control goroutine should exit
	done := make(chan struct{})
	go func() {
		defer close(done)
		subs.Close()
	}()

	select {
	case <-done:
		// ok, goroutine exited
	case <-time.After(3 * time.Second):
		t.Fatal("Close() timed out waiting for control goroutine to exit")
	}
}

func TestSubscriptions_AddAndLookup(t *testing.T) {
	log := logging.MustGetLogger("test")
	subs := newSubscriptions(5*time.Second, log)
	defer subs.Close()

	req := &SubscriptionRequest{
		SubID:   "test-sub-1",
		Event:   "subscribe",
		Channel: ChanTicker,
		Symbol:  "tBTCUSD",
	}

	sub := subs.add(SocketId(0), req)
	if sub == nil {
		t.Fatal("expected subscription, got nil")
	}
	if !sub.Pending() {
		t.Error("new subscription should be pending")
	}

	// Lookup by subscription ID
	found, err := subs.lookupBySubscriptionID("test-sub-1")
	if err != nil {
		t.Fatalf("lookup by subscription ID: %v", err)
	}
	if found.SubID() != "test-sub-1" {
		t.Errorf("got sub ID %s, want test-sub-1", found.SubID())
	}

	// Activate
	err = subs.activate("test-sub-1", 12345)
	if err != nil {
		t.Fatalf("activate: %v", err)
	}

	// Lookup by channel ID after activation
	found2, err := subs.lookupBySocketChannelID(12345, SocketId(0))
	if err != nil {
		t.Fatalf("lookup by socket channel ID: %v", err)
	}
	if found2.Pending() {
		t.Error("activated subscription should not be pending")
	}
}

func TestSubscriptions_ResetSocketSubscriptions(t *testing.T) {
	log := logging.MustGetLogger("test")
	subs := newSubscriptions(5*time.Second, log)
	defer subs.Close()

	req1 := &SubscriptionRequest{SubID: "s1", Event: "subscribe", Channel: ChanTicker, Symbol: "tBTCUSD"}
	req2 := &SubscriptionRequest{SubID: "s2", Event: "subscribe", Channel: ChanBook, Symbol: "tBTCUSD"}

	subs.add(SocketId(0), req1)
	subs.add(SocketId(0), req2)

	reset := subs.ResetSocketSubscriptions(SocketId(0))
	if len(reset) != 2 {
		t.Errorf("expected 2 reset subscriptions, got %d", len(reset))
	}

	// Should not find them anymore
	_, err := subs.lookupBySubscriptionID("s1")
	if err == nil {
		t.Error("expected error looking up reset subscription s1")
	}
	_, err = subs.lookupBySubscriptionID("s2")
	if err == nil {
		t.Error("expected error looking up reset subscription s2")
	}
}
