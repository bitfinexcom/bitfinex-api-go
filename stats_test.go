package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStatsAll(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "period":1,
            "volume":"7967.96766158"
          },{
            "period":7,
            "volume":"55938.67260266"
          },{
            "period":30,
            "volume":"275148.09653645"
          }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	stats, err := NewClient().Stats.All("btcusd", "10", "")

	if err != nil {
		t.Error(err)
	}

	if len(stats) != 3 {
		t.Error("Expected", 3)
		t.Error("Actual ", len(stats))
	}
}
