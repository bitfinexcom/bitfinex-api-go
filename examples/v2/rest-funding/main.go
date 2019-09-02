package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"os"
)


// Set BFX_APIKEY and BFX_SECRET as :
//
// export BFX_API_KEY=YOUR_API_KEY
// export BFX_API_SECRET=YOUR_API_SECRET
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_KEY")
	secret := os.Getenv("BFX_SECRET")
	c := rest.NewClientWithURL("https://test.bitfinex.com/v2/").Credentials(key, secret)
	// active funding offers
	snap, err := c.Funding.Offers("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snap.Snapshot {
		fmt.Println(item)
	}
	// funding offer history
	snapHist, err := c.Funding.OfferHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapHist.Snapshot {
		fmt.Println(item)
	}
	// active loans
	snapLoans, err := c.Funding.Loans("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapLoans.Snapshot {
		fmt.Println(item)
	}
	// loans history
	napLoansHist, err := c.Funding.LoansHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napLoansHist.Snapshot {
		fmt.Println(item)
	}
	// active credits
	snapCredits, err := c.Funding.Credits("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range snapCredits.Snapshot {
		fmt.Println(item)
	}
	// credits history
	napCreditsHist, err := c.Funding.CreditsHistory("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napCreditsHist.Snapshot {
		fmt.Println(item)
	}
	// funding trades
	napTradesHist, err := c.Funding.Trades("fUSD")
	if err != nil {
		panic(err)
	}
	for _, item := range napTradesHist.Snapshot {
		fmt.Println(item)
	}

	/********* submit a new funding offer ***********/
	fo, err := c.Funding.SubmitOffer(&bitfinex.FundingOfferRequest{
		Type: "LIMIT",
		Symbol: "fUSD",
		Amount: 1000,
		Rate: 0.012,
		Period: 7,
		Hidden: true,
	})
	if err != nil {
		panic(err)
	}
	newOffer := fo.NotifyInfo.(*bitfinex.FundingOfferNew)
	/********* cancel funding offer ***********/
	fc, err := c.Funding.CancelOffer(&bitfinex.FundingOfferCancelRequest{
		Id: newOffer.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(fc)
}
