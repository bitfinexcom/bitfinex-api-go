package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOrdersAll(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `
        [{
           "id":448411365,
           "symbol":"btcusd",
           "exchange":"bitfinex",
           "price":"0.02",
           "avg_execution_price":"0.0",
           "side":"buy",
           "type":"exchange limit",
           "timestamp":"1444276597.0",
           "is_live":true,
           "is_cancelled":false,
           "is_hidden":false,
           "was_forced":false,
           "original_amount":"0.02",
           "remaining_amount":"0.02",
           "executed_amount":"0.0"
         }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders, err := NewClient().Orders.All()

	if err != nil {
		t.Error(err)
	}

	expectedID := int64(448411365)
	if orders[0].ID != expectedID {
		t.Error("Expected", expectedID)
		t.Error("Actual ", orders[0].ID)
	}

}

func TestCreateMulti(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
            "order_ids":[{
            "id":448383727,
            "symbol":"btcusd",
            "exchange":"bitfinex",
            "price":"0.01",
            "avg_execution_price":"0.0",
            "side":"buy",
            "type":"exchange limit",
            "timestamp":"1444274013.621701916",
            "is_live":true,
            "is_cancelled":false,
            "is_hidden":false,
            "was_forced":false,
            "original_amount":"0.01",
            "remaining_amount":"0.01",
            "executed_amount":"0.0"
         },{
            "id":448383729,
            "symbol":"btcusd",
            "exchange":"bitfinex",
            "price":"0.03",
            "avg_execution_price":"0.0",
            "side":"buy",
            "type":"exchange limit",
            "timestamp":"1444274013.661297306",
            "is_live":true,
            "is_cancelled":false,
            "is_hidden":false,
            "was_forced":false,
            "original_amount":"0.02",
            "remaining_amount":"0.02",
            "executed_amount":"0.0"
          }],
          "status":"success"
       }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	reqOrders := []SubmitOrder{{
		Symbol: "BTCUSD",
		Amount: 10.0,
		Price:  450.0,
		Type:   OrderTypeLimit,
	}, {
		Symbol: "BTCUSD",
		Amount: 10.0,
		Price:  450.0,
		Type:   OrderTypeLimit,
	}}
	response, err := NewClient().Orders.CreateMulti(reqOrders)

	if err != nil {
		t.Error(err)
	}

	if len(response.Orders) != 2 {
		t.Error("Expected", 2)
		t.Error("Actual ", len(response.Orders))
	}
}

func TestCancelMulti(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{"result":"Orders cancelled"}`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders := []int64{1000, 1001, 1002}
	response, err := NewClient().Orders.CancelMulti(orders)

	if err != nil {
		t.Error(err)
	}

	if response != "Orders cancelled" {
		t.Error("Expected", "Orders cancelled")
		t.Error("Actual ", response)
	}
}
