package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMarginFundingNew(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
            "id":13800585,
            "currency":"USD",
            "rate":"20.0",
            "period":2,
            "direction":"lend",
            "timestamp":"1444279698.21175971",
            "is_live":true,
            "is_cancelled":false,
            "original_amount":"50.0",
            "remaining_amount":"50.0",
            "executed_amount":"0.0",
            "offer_id":13800585
        }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	offer, err := NewClient().MarginFunding.new("BTC", "loan", 10.0, 0.01, 10)

	if err != nil {
		t.Error(err)
	}

	if offer.ID != 13800585 {
		t.Error("Expected", 13800585)
		t.Error("Actual ", offer.ID)
	}
	if !offer.IsLive {
		t.Error("Expected", true)
		t.Error("Actual ", offer.IsLive)
	}
}
