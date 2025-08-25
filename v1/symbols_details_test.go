package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSymbolsDetailsGet(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
           "pair":"btcusd",
           "price_precision":5,
           "initial_margin":"30.0",
           "minimum_margin":"15.0",
           "maximum_order_size":"2000.0",
           "minimum_order_size":"0.01",
           "expiration":"NA"
        }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	symbolsDetails, err := NewClient().SymbolsDetails.GetSymbolsDetails()
	if err != nil {
		t.Error(err)
	}
	if len(symbolsDetails) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(symbolsDetails))
	}
}
