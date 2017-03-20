package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPairsGetAll(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `["btcusd","ltcusd","ltcbtc","ethusd","ethbtc"]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	pairs, err := NewClient().Pairs.All()

	numPairs := len(pairs)

	if err != nil {
		t.Error(err)
	}

	if numPairs != 5 {
		t.Error("Expected", 5)
		t.Error("Actual ", numPairs)
	}

	if (pairs)[0] != "btcusd" {
		t.Error("Expected", "btcusd")
		t.Error("Actual ", pairs[0])
	}
}

func TestPairsAllDetailed(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "pair":"btcusd",
            "price_precision":5,
            "initial_margin":"30.0",
            "minimum_margin":"15.0",
            "maximum_order_size":"2000.0",
            "minimum_order_size":"0.01",
            "expiration":"NA"
        },{
            "pair":"ltcusd",
            "price_precision":5,
            "initial_margin":"30.0",
            "minimum_margin":"15.0",
            "maximum_order_size":"5000.0",
            "minimum_order_size":"0.1",
            "expiration":"NA"
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	pairs, err := NewClient().Pairs.AllDetailed()

	if err != nil {
		t.Error(err)
	}

	if len(pairs) != 2 {
		t.Error("Expected", 2)
		t.Error("Actual ", len(pairs))
	}

	pairMargin := pairs[0].InitialMargin
	expectedMargin := 30.0
	if (pairMargin-expectedMargin) > 0.1 || (expectedMargin-pairMargin) > 0.1 {
		t.Error("Expected", expectedMargin)
		t.Error("Actual ", pairMargin)
	}

}
