package bitfinex

// Package bitfinex-api-go provides structs and functions for accessing
// bitfinex.com api version 1.0
//
// Usage:
//   import "github.com/bitfinexcom/bitfinex-api-go"
//
// Create new client:
//   api := bitfinex.NewClient()
//
// For access methods that requires authentication use the next code:
//   api := bitfinex.NewClient().Auth(key, secret)
//
// Get all pairs
//   api.Pairs.All()
//
// Get account info
//   api.Account.Info()
//
// See examples dir for more info.
