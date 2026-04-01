package bitfinex

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestAccountV1Info(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
           "maker_fees":"0.1",
           "taker_fees":"0.2",
           "fees":[{
               "pairs":"BTC",
               "maker_fees":"0.1",
               "taker_fees":"0.2"
            },{
               "pairs":"LTC",
               "maker_fees":"0.1",
               "taker_fees":"0.2"
            },{
               "pairs":"ETH",
               "maker_fees":"0.1",
               "taker_fees":"0.2"
            }]
        }]`

		resp := http.Response{
			Body:       io.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	info, err := NewClient().Account.V1Info()

	if err != nil {
		t.Error(err)
	}

	if len(info.Fees) != 3 {
		t.Error("Expected", 3)
		t.Error("Actual ", len(info.Fees))
	}
}

func TestAccountV1KeyPermission(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
            "account":{
                "read":true,
                "write":false
            },
            "history":{
                "read":true,
                "write":false
            },
            "orders":{
                "read":true,
                "write":true
            },
            "positions":{
                "read":true,
                "write":true
            },
            "funding":{
                "read":true,
                "write":true
            },
            "wallets":{
                "read":true,
                "write":true
            },
            "withdraw":{
                "read":null,
                "write":null
            }
        }`

		resp := http.Response{
			Body:       io.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}
	perm, err := NewClient().Account.V1KeyPermission()

	if err != nil {
		t.Error(err)
	}

	if !perm.Account.Read {
		t.Error("Expected", true)
		t.Error("Actual ", perm.Account.Read)
	}

	if perm.History.Write {
		t.Error("Expected", false)
		t.Error("Actual ", perm.History.Write)
	}
}
