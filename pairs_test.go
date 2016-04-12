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
	numPairs := len(*pairs)

	if err != nil {
		t.Error(err)
	}

	if numPairs != 5 {
		t.Error("Expected", 5)
		t.Error("Actual ", numPairs)
	}

	if (*pairs)[0] != "btcusd" {
		t.Error("Expected", "btcusd")
		t.Error("Actual ", (*pairs)[0])
	}
}
