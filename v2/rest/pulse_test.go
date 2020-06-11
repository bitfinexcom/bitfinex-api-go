package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
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
			payload, _ := json.Marshal(respData)
			w.Write(payload)
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
			payload, _ := json.Marshal(respData)
			w.Write(payload)
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

func TestPublicPulseHistory(t *testing.T) {
	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respData := []interface{}{
				[]interface{}{"id"},
			}
			payload, _ := json.Marshal(respData)
			w.Write(payload)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pp, err := c.Pulse.PublicPulseHistory("", "")
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("valid response data no profile", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respData := []interface{}{
				[]interface{}{
					"id",
					float64(1591614631576),
					nil,
					"uid",
					nil,
					"title",
					"content",
					nil,
					nil,
					1,
					1,
					nil,
					[]interface{}{"tag1", "tag2"},
					[]interface{}{"attach1", "attach2"},
					nil,
					5,
					nil,
					nil,
					nil,
				},
			}
			payload, _ := json.Marshal(respData)
			w.Write(payload)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pph, err := c.Pulse.PublicPulseHistory("", "")
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:          "id",
			MTS:         1591614631576,
			UserID:      "uid",
			Title:       "title",
			Content:     "content",
			IsPin:       1,
			IsPublic:    1,
			Tags:        []string{"tag1", "tag2"},
			Attachments: []string{"attach1", "attach2"},
			Likes:       5,
		}

		assert.Equal(t, expected, pph[0])
	})
}

func TestAddPulse(t *testing.T) {
	t.Run("invalid payload", func(t *testing.T) {
		p := &pulse.Pulse{Title: "foo"}
		c := NewClient()
		pm, err := c.Pulse.AddPulse(p)
		require.NotNil(t, err)
		require.Nil(t, pm)
	})

	t.Run("valid payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respData := []interface{}{
				"id",
				float64(1591614631576),
				nil,
				"uid",
				nil,
				"title",
				"content",
				nil,
				nil,
				1,
				1,
				nil,
				[]interface{}{"tag1", "tag2"},
				[]interface{}{"attach1", "attach2"},
				nil,
				5,
				nil,
				nil,
				nil,
			}
			payload, _ := json.Marshal(respData)
			w.Write(payload)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pm, err := c.Pulse.AddPulse(&pulse.Pulse{Title: "foo bar baz qux 123"})
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:          "id",
			MTS:         1591614631576,
			UserID:      "uid",
			Title:       "title",
			Content:     "content",
			IsPin:       1,
			IsPublic:    1,
			Tags:        []string{"tag1", "tag2"},
			Attachments: []string{"attach1", "attach2"},
			Likes:       5,
		}

		assert.Equal(t, expected, pm)
	})
}
