package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeepFunding(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		c := rest.NewClient()
		kf, err := c.Funding.KeepFunding(rest.KeepFundingRequest{Type: "foo"})
		require.NotNil(t, err)
		require.Nil(t, kf)
	})

	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/funding/keep", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := rest.KeepFundingRequest{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := rest.KeepFundingRequest{Type: "loan", ID: 123}
			assert.Equal(t, expectedReqPld, gotReqPld)

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := rest.NewClientWithURL(server.URL)
		pld := rest.KeepFundingRequest{Type: "loan", ID: 123}
		rsp, err := c.Funding.KeepFunding(pld)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}
