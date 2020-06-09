package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicPulseProfile(t *testing.T) {
	t.Run("missing arguments", func(t *testing.T) {
		c := NewClient()
		pp, err := c.Pulse.PublicPulseProfile("")
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respData := []interface{}{"abc123"}
			retData, _ := json.Marshal(respData)
			w.Write(retData)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pp, err := c.Pulse.PublicPulseProfile("Bitfinex")
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("valid response data", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respData := []interface{}{
				"abc123",
				float64(1591614631576),
				nil,
				"nickname",
				nil,
				"picture",
				"text",
				nil,
				nil,
				"twitter",
				nil,
				30,
				5,
				nil,
			}
			retData, _ := json.Marshal(respData)
			w.Write(retData)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pp, err := c.Pulse.PublicPulseProfile("Bitfinex")
		require.Nil(t, err)

		expected := &pulseprofile.PulseProfile{
			ID:            "abc123",
			MTS:           1591614631576,
			Nickname:      "nickname",
			Picture:       "picture",
			Text:          "text",
			TwitterHandle: "twitter",
		}
		assert.Equal(t, expected, pp)
	})
}
