package tests

import (
	"testing"
)

func TestGetBalances(t *testing.T) {
	_, err := client.Balances.All()

	if err != nil {
		t.Fatalf("Balances.All() returned error: %v", err)
	}
}
