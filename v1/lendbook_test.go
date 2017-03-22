package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLendbookGet(t *testing.T) {
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

	book, err := NewClient().Lendbook.Get("usd", 0, 0)

	if err != nil {
		t.Error(err)
	}

	if len(book.Bids) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(book.Bids))
	}

	if book.Bids[0].Period != 30 {
		t.Error("Expected", 30)
		t.Error("Actual ", book.Bids[0].Period)
	}
}

func TestLendbookLends(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "rate":"9.8998",
            "amount_lent":"22528933.77950878",
            "amount_used":"0.0",
            "timestamp":1444264307
        }]`

		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	lends, err := NewClient().Lendbook.Lends("usd")

	if err != nil {
		t.Error(err)
	}

	if len(lends) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(lends))
	}
}
