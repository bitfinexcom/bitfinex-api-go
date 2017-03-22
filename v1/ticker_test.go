package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTickerGet(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
           "mid":"244.755",
           "bid":"244.75",
           "ask":"244.76",
           "last_price":"244.82",
           "low":"244.2",
           "high":"248.19",
           "volume":"7842.11542563",
           "timestamp":"1444253422.348340958"
        }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	tick, err := NewClient().Ticker.Get("btcusd")

	if err != nil {
		t.Error(err)
	}

	if tick.Bid != "244.75" {
		t.Error("Expected", "244.75")
		t.Error("Actual ", tick.Bid)
	}
	if tick.LastPrice != "244.82" {
		t.Error("Expected", "244.82")
		t.Error("Actual ", tick.LastPrice)
	}
}
