package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestTradesServiceGet(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
           "timestamp":1444266681,
           "tid":11988919,
           "price":"244.8",
           "amount":"0.03297384",
           "exchange":"bitfinex",
           "type":"sell"
       }]`

		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	trades, err := NewClient().Trades.All("ethusd", time.Time{}, 0)

	if err != nil {
		t.Error(err)
	}

	if len(trades) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(trades))
	}
}
