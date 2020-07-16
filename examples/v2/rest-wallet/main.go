package main

import (
	"log"
	"os"

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

	getWallets(c)
	transfer(c)
	depositAddress(c)
	createDepositAddress(c)
	withdraw(c)
}

func getWallets(c *rest.Client) {
	wallets, err := c.Wallet.Wallet()
	if err != nil {
		log.Fatalf("getWallets %s", err)
	}

	spew.Dump(wallets)
}

func transfer(c *rest.Client) {
	notfication, err := c.Wallet.Transfer("exchange", "margin", "BTC", "BTC", 0.001)
	if err != nil {
		log.Fatalf("transfer %s", err)
	}

	spew.Dump(notfication)
}

func depositAddress(c *rest.Client) {
	notfication, err := c.Wallet.DepositAddress("exchange", "ethereum")
	if err != nil {
		log.Fatalf("depositAddress %s", err)
	}

	spew.Dump(notfication)
}

func createDepositAddress(c *rest.Client) {
	notfication, err := c.Wallet.DepositAddress("margin", "ethereum")
	if err != nil {
		log.Fatalf("createDepositAddress %s", err)
	}

	spew.Dump(notfication)
}

func withdraw(c *rest.Client) {
	notfication, err := c.Wallet.Withdraw("exchange", "ethereum", 0.1, "0x5B4Dbe55dE0B565db6C63405D942886140083cE8")
	if err != nil {
		log.Fatalf("withdraw %s", err)
	}

	spew.Dump(notfication)
}
