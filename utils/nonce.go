package utils

import (
	"fmt"
	"time"
)

// v1 support
const multiplier = 10000

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func GetNonce() string {
	return fmt.Sprintf("%v", time.Now().Unix() * multiplier)
}
