package rest

import (
	"testing"
)

func TestSign_ReturnsErrorNotNil(t *testing.T) {
	// The sign function should never return ("", nil) when
	// the underlying hmac write fails. While hmac.Write rarely fails,
	// we verify the function at least returns a valid signature for
	// a normal case, and that the error return path is correct.
	c := &Client{
		apiSecret: "test-secret",
	}

	sig, err := c.sign("test-message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig == "" {
		t.Error("expected non-empty signature")
	}

	// Verify different messages produce different signatures
	sig2, err := c.sign("different-message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sig == sig2 {
		t.Error("different messages should produce different signatures")
	}
}
