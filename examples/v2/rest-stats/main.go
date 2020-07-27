package main

import (
	"fmt"
	"log"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	c := rest.NewClient()
	positionLast(c)
	positionHistory(c)
	symbolCreditSizeLast(c)
	symbolCreditSizeHistory(c)
	fundingLast(c)
	fundingHistory(c)
}

func positionLast(c *rest.Client) {
	resp, err := c.Stats.PositionLast("tBTCUSD", common.Long)
	if err != nil {
		log.Fatalf("PositionLast: %s", err)
	}
	fmt.Println("Position Last:")
	spew.Dump(resp)
}

func positionHistory(c *rest.Client) {
	resp, err := c.Stats.PositionHistory("tBTCUSD", common.Long)
	if err != nil {
		log.Fatalf("PositionHistory: %s", err)
	}
	fmt.Println("Position history:")
	spew.Dump(resp)
}

func symbolCreditSizeLast(c *rest.Client) {
	resp, err := c.Stats.SymbolCreditSizeLast("fUSD", "tBTCUSD")
	if err != nil {
		log.Fatalf("SymbolCreditSizeLast: %s", err)
	}
	fmt.Println("Symbol Credit Size Last:")
	spew.Dump(resp)
}

func symbolCreditSizeHistory(c *rest.Client) {
	resp, err := c.Stats.SymbolCreditSizeHistory("fUSD", "tBTCUSD")
	if err != nil {
		log.Fatalf("SymbolCreditSizeHistory: %s", err)
	}
	fmt.Println("Symbol Credit Size History:")
	spew.Dump(resp)
}

func fundingLast(c *rest.Client) {
	resp, err := c.Stats.FundingLast("fUSD")
	if err != nil {
		log.Fatalf("FundingLast: %s", err)
	}
	fmt.Println("Funding Last:")
	spew.Dump(resp)
}

func fundingHistory(c *rest.Client) {
	resp, err := c.Stats.FundingHistory("fUSD")
	if err != nil {
		log.Fatalf("FundingHistory: %s", err)
	}
	fmt.Println("Funding History:")
	spew.Dump(resp)
}
