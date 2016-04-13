package bitfinex

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "testing"
    "time"
)

func TestHistoryBalance(t *testing.T) {
    httpDo = func(req *http.Request) (*http.Response, error) {
        msg := `[{
            "currency":"USD",
            "amount":"-246.94",
            "balance":"515.4476526",
            "description":"Position claimed @ 245.2 on wallet trading",
            "timestamp":"1444277602.0"
        }]`
        resp := http.Response{
            Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
            StatusCode: 200,
        }
        return &resp, nil
    }

    balance, err := NewClient().History.Balance("BTC", "", time.Time{}, time.Time{}, 0)

    if err != nil {
        t.Error(err)
    }

    if len(balance) != 1 {
        t.Error("Expected", 1)
        t.Error("Actual ", len(balance))
    }

}
