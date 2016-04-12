package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestPositions(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[
            {
                "id":943715,
                "symbol":"btcusd",
                "status":"ACTIVE",
                "base":"246.94",
                "amount":"1.0",
                "timestamp":"1444141857.0",
                "swap":"0.0",
                "pl":"-2.22042"
            }
        ]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	client := NewClient()
	positions, err := client.Positions.All()
	if err != nil {
		t.Error(err)
	}

	if len(positions) != 1 {
		t.Error("Expected", 1)
		t.Error("Actual ", len(positions))
	}

	pos := positions[0]
	if pos.Amount != "1.0" {
		t.Error("Expected", "1.0")
		t.Error("Actual ", pos.Amount)
	}

}

func TestPositionsWhenEmpty(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
            "id":943715,
            "symbol":"btcusd",
            "status":"ACTIVE",
            "base":"246.94",
            "amount":"1.0",
            "timestamp":"1444141857.0",
            "swap":"0.0",
            "pl":"-2.2304"
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	client := NewClient()
	positions, err := client.Positions.All()

	if err != nil {
		t.Error(err)
	}

	position := positions[0]
	if position.ID != 943715 {
		t.Error("Expected", 943715)
		t.Error("Actual ", position.ID)
	}

	parsedTime, err := position.ParseTime()
	loc, _ := time.LoadLocation("Europe/Berlin")
	expectedTime := time.Date(2015, 10, 06, 16, 30, 57, 00, loc)

	if err != nil {
		t.Error(err)
	}
	if !(*parsedTime).Equal(expectedTime) {
		t.Error("Expected", &expectedTime)
		t.Error("Actual ", parsedTime)
	}

}

func TestClaimPosition(t *testing.T) {

	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
                "id":943715,
                "symbol":"btcusd",
                "status":"ACTIVE",
                "base":"246.94",
                "amount":"1.0",
                "timestamp":"1444141857.0",
                "swap":"0.0",
                "pl":"-2.22042"
            }`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	position, err := NewClient().Positions.Claim("943715", "0.5")

	if err != nil {
		t.Error(err)
	}

	if position.ID != 943715 {
		t.Error("Expected", 943715)
		t.Error("Actual ", position.ID)
	}
}
