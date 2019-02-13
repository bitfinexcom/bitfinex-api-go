package utils

import (
	"strconv"
	"sync"
	"sync/atomic"
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
	return strconv.FormatUint(atomic.AddUint64(&u.nonce, 1), 10)
}

func NewEpochNonceGenerator() *EpochNonceGenerator {
	return &EpochNonceGenerator{
		nonce: uint64(time.Now().Unix()) * 1000,
	}
}

// v1 support

var nonce uint64

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
	return strconv.FormatUint(atomic.AddUint64(&u.nonce, 1), 10)
}

func NewEpochNonceGenerator() *EpochNonceGenerator {
	return &EpochNonceGenerator{
		nonce: uint64(time.Now().Unix()) * 1000,
	}
}

// v1 support

// var nonce uint64

// func init() {
// 	nonce = uint64(time.Now().Unix()) * 1000
// }

var offset = uint64(0)
var lock sync.Mutex

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func GetNonce() string {
	lock.Lock()
	defer lock.Unlock()
	lNonce := uint64(time.Now().Unix()) * 1000
	if offset >= 999 {
		offset = 0
	} else {
		offset++
	}
	return strconv.FormatUint(lNonce+offset, 10)
}
