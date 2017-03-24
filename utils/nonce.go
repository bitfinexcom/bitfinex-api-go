package utils

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// GetNonce produces a random integer >= time.Now().UnixNano().
func GetNonce() (string, error) {
	t := big.NewInt(time.Now().UnixNano())
	bi, err := rand.Int(rand.Reader, t.Add(t, big.NewInt(math.MaxInt64)))
	if err != nil {
		return "", err
	}

	return bi.String(), nil
}
