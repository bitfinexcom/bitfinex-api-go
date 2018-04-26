package utils

import (
	"strconv"
	"time"
)

// v2 types

type NonceGenerator interface {
	GetNonce() string
}

type EpochNonceGenerator struct {
	nonce uint64
}

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func (u *EpochNonceGenerator) GetNonce() string {
	return strconv.FormatUint(time.Now().UnixNano(), 10)
}

func NewEpochNonceGenerator() *EpochNonceGenerator {
	return &EpochNonceGenerator{
		nonce: uint64(time.Now().UnixNano()),
	}
}

// v1 support

var nonce uint64

func init() {
	nonce = uint64(time.Now().UnixNano())
}

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func GetNonce() string {
	return strconv.FormatUint(time.Now().UnixNano(), 10)
}
