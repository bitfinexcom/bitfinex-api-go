package utils

import (
	"strconv"
	"time"
)

// GetNonce - getting unique nonce

var nonce int64

func GetNonce() (string, error) {
	if nonce == 0 {
		nonce = time.Now().Unix() * 1000
	} else {
		nonce++
	}
	return strconv.FormatInt(nonce, 10), nil
}
