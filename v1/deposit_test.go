package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDepositNew(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
            "result":"success",
            "method":"bitcoin",
            "currency":"BTC",
            "address":"1A2wyHKJ4KWEoahDHVxwQy3kdd6g1qiSYV"
        }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	deposit, err := NewClient().Deposit.New("bitcoin", "trading", 0)

	if err != nil {
		t.Error(err)
	}
	success, err := deposit.Success()

	if err != nil || !success {
		t.Error("Expected", true)
		t.Error("Actual ", success)
		t.Error("With message", err)
	}
}
