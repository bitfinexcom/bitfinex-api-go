package main

import (
	"os"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
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

	// create order
	response, err := c.Orders.SubmitOrder(&order.NewRequest{
		Symbol: "tBTCUSD",
		CID:    time.Now().Unix() / 1000,
		Amount: 0.02,
		Type:   "EXCHANGE LIMIT",
		Price:  5000,
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 5)
	orders := response.NotifyInfo.(*order.Snapshot)
	// update orders
	for _, o := range orders.Snapshot {
		response, err = c.Orders.SubmitUpdateOrder(&order.UpdateRequest{
			ID:    o.ID,
			Price: 6000,
		})
		if err != nil {
			panic(err)
		}
		// cancel orders
		updatedOrder := response.NotifyInfo.(*order.Update)
		err := c.Orders.SubmitCancelOrder(&order.CancelRequest{
			ID: updatedOrder.ID,
		})
		if err != nil {
			panic(err)
		}
	}
}
