package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSymbolsGet(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[
  			"btcusd",
  			"ltcusd",
  			"ltcbtc",
        ]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	symbols, err := NewClient().Symbols.GetSymbols()
	if err != nil {
		t.Error(err)
	}
	if len(symbols) != 3 {
		t.Error("Expected", 3)
		t.Error("Actual ", len(symbols))
	}
}
