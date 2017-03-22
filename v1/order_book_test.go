package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOrderBookGet(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
           "bids":[{
           "rate":"9.1287",
           "amount":"5000.0",
           "period":30,
           "timestamp":"1444257541.0",
           "frr":"No"
       }],
           "asks":[{
           "rate":"8.3695",
           "amount":"407.5",
           "period":2,
           "timestamp":"1444260343.0",
           "frr":"No"
       }]
       }`

		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orderBook, err := NewClient().OrderBook.Get("btcusd", 0, 0, false)

	if err != nil {
		t.Error(err)
	}

	if len(orderBook.Bids) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(orderBook.Bids))
	}

}
