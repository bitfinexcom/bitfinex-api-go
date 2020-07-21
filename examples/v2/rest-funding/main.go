package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

// Set BFX_API_KEY and BFX_API_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")
	c := rest.NewClient().Credentials(key, secret)

	offers(c)
	offerHistory(c)
	loans(c)
	loansHistory(c)
	activeCredits(c)
	creditsHistory(c)
	fundingTrades(c)
	keepFunding(c)

	/********* submit a new funding offer ***********/
	fo, err := c.Funding.SubmitOffer(&fundingoffer.SubmitRequest{
		Type:   "LIMIT",
		Symbol: "fUSD",
		Amount: 1000,
		Rate:   0.012,
		Period: 7,
		Hidden: true,
	})
	if err != nil {
		panic(err)
	}
	newOffer := fo.NotifyInfo.(*fundingoffer.New)

	/********* cancel funding offer ***********/
	fc, err := c.Funding.CancelOffer(&fundingoffer.CancelRequest{
		ID: newOffer.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(fc)
}

func offers(c *rest.Client) {
	// active funding offers
	snap, err := c.Funding.Offers("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snap.Snapshot {
		fmt.Println(item)
	}
}

func offerHistory(c *rest.Client) {
	// funding offer history
	snapHist, err := c.Funding.OfferHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapHist.Snapshot {
		fmt.Println(item)
	}
}

func loans(c *rest.Client) {
	// active loans
	snapLoans, err := c.Funding.Loans("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapLoans.Snapshot {
		fmt.Println(item)
	}
}

func loansHistory(c *rest.Client) {
	napLoansHist, err := c.Funding.LoansHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napLoansHist.Snapshot {
		fmt.Println(item)
	}
}

func activeCredits(c *rest.Client) {
	// active credits
	snapCredits, err := c.Funding.Credits("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapCredits.Snapshot {
		fmt.Println(item)
	}
}

func creditsHistory(c *rest.Client) {
	napCreditsHist, err := c.Funding.CreditsHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napCreditsHist.Snapshot {
		fmt.Println(item)
	}
}

func fundingTrades(c *rest.Client) {
	napTradesHist, err := c.Funding.Trades("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napTradesHist.Snapshot {
		fmt.Println(item)
	}
}

func keepFunding(c *rest.Client) {
	// keep funding
	resp, err := c.Funding.KeepFunding(rest.KeepFundingRequest{
		Type: "credit",
		ID:   12345, // Insert correct ID
	})
	if err != nil {
		log.Fatalf("KeepFunding error: %s", err)
	}

	spew.Dump(resp)
}
