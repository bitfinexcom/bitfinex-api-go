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

func TestHistoryMovements(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "id":581183,
            "currency":"BTC",
            "method":"BITCOIN",
            "type":"WITHDRAWAL",
            "amount":".01",
            "description":"3QXYWgRGX2BPYBpUDBssGbeWEa5zq6snBZ, offchain transfer ",
            "status":"COMPLETED",
            "timestamp":"1443833327.0"
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	movements, err := NewClient().History.Movements("BTC", "", time.Time{}, time.Time{}, 0)

	if err != nil {
		t.Error(err)
	}

	if len(movements) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(movements))
	}

}

func TestHistoryTrades(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "price":"246.94",
            "amount":"1.0",
            "timestamp":"1444141857.0",
            "exchange":"",
            "type":"Buy",
            "fee_currency":"USD",
            "fee_amount":"-0.49388",
            "tid":11970839,
            "order_id":446913929
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	trades, err := NewClient().History.Trades("BTC", time.Time{}, time.Time{}, 0, false)

	if err != nil {
		t.Error(err)
	}

	if len(trades) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(trades))
	}
	if trades[0].TID != 11970839 {
		t.Error("Expected", 11970839)
		t.Error("Actual ", trades[0].TID)
	}

}
