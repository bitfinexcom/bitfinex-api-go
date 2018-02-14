package tests

import (
	"fmt"
)

// IncrementingNonceGenerator starts at nonce1 and increments each by +1: nonce1, nonce2, ..., nonceN
type IncrementingNonceGenerator struct {
	nonce int
}

// GetNonce returns an incrementing nonce value.
func (m *IncrementingNonceGenerator) GetNonce() string {
	m.nonce++
	return fmt.Sprintf("nonce%d", m.nonce)
}
