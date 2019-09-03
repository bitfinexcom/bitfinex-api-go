package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"log"
	"os"
)


func main() {
	key := os.Getenv("BFX_KEY")
	secret := os.Getenv("BFX_SECRET")
	uri := "https://api.bitfinex.com/v2/"

	c := rest.NewClientWithURL(uri).Credentials(key, secret)
	// get active positions
	positions, err := c.Positions.All()
	if err != nil {
		log.Fatalf("getting wallet %s", err)
	}
	for _, p := range positions.Snapshot {
		fmt.Println(p)
	}
	// claim active position
	pClaim, err := c.Positions.Claim(&bitfinex.ClaimPositionRequest{
		Id: 36228736,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(pClaim.NotifyInfo.(*bitfinex.PositionCancel))
}

