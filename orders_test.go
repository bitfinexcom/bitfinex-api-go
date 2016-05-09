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

    expectedId := 448411365
    if orders[0].Id != expectedId {
        t.Error("Expected", expectedId)
        t.Error("Actual ", orders[0].Id)
    }

}
