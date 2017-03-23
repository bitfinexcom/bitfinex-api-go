package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

// GetNonce - getting unique nonce
func GetNonce() (string, error) {
	bi, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}

	return bi.String(), nil
}
