package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOfferNew(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `{
          "id":13800585,
          "currency":"USD",
          "rate":"20.0",
          "period":2,
          "direction":"lend",
          "timestamp":"1444279698.21175971",
          "is_live":true,
          "is_cancelled":true,
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

	offer, err := NewClient().Offers.New("USD", 50, 20.0, 2, LEND)

	if err != nil {
		t.Error(err)
	}

	if offer.Currency != "USD" {
		t.Error("Expected", "USD")
		t.Error("Actual ", offer.Currency)
	}

	expectedId := int64(13800585)
	if offer.OfferId != expectedId {
		t.Error("Expected", expectedId)
		t.Error("Actual ", offer.OfferId)
	}

	if offer.Id != expectedId {
		t.Error("Expected", expectedId)
		t.Error("Actual ", offer.Id)
	}
	if !offer.IsLive {
		t.Error("Expected", true)
		t.Error("Actual ", offer.IsLive)
	}

	newOffer, err := NewClient().Offers.Cancel(offer.Id)

	if err != nil {
		t.Error(err)
	}

	if newOffer.Currency != "USD" {
		t.Error("Expected", "USD")
		t.Error("Actual ", newOffer.Currency)
	}

	if !newOffer.IsCancelled {
		t.Error("Expected", true)
		t.Error("Actual ", newOffer.IsCancelled)
	}

}
