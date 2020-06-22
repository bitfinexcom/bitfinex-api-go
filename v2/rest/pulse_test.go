package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
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
			respMock := []interface{}{"abc123"}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
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
			assert.Equal(t, "/pulse/profile/Bitfinex", r.RequestURI)
			respMock := []interface{}{
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
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
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
			respMock := []interface{}{
				[]interface{}{"id"},
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		now := time.Now()
		millis := now.UnixNano() / 1000000
		end := bitfinex.Mts(millis)
		pp, err := c.Pulse.PublicPulseHistory(1, end)
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("valid response data no profile", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/pulse/hist?end=1591614631576&limit=1", r.RequestURI)
			limit := r.URL.Query().Get("limit")
			end := r.URL.Query().Get("end")
			assert.Equal(t, "1", limit)
			assert.Equal(t, "1591614631576", end)

			respMock := []interface{}{
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
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		end := bitfinex.Mts(1591614631576)
		pph, err := c.Pulse.PublicPulseHistory(1, end)
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

	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/pulse/add", r.RequestURI)
			respMock := []interface{}{"id"}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pm, err := c.Pulse.AddPulse(&pulse.Pulse{Title: "foo bar baz qux 123"})
		require.NotNil(t, err)
		require.Nil(t, pm)
	})

	t.Run("valid payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{
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
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
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

func TestPulseHistory(t *testing.T) {
	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/r/pulse/hist?isPublic=0", r.RequestURI)
			isPublic := r.URL.Query().Get("isPublic")
			assert.Equal(t, "0", isPublic)

			respMock := []interface{}{
				[]interface{}{"id"},
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pp, err := c.Pulse.PulseHistory(false)
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("private pulse history", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			isPublic := r.URL.Query().Get("isPublic")
			assert.Equal(t, "0", isPublic)

			respMock := []interface{}{
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
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pm, err := c.Pulse.PulseHistory(false)
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

		assert.Equal(t, expected, pm[0])
	})

	t.Run("public pulse history", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			isPublic := r.URL.Query().Get("isPublic")
			assert.Equal(t, "1", isPublic)

			respMock := []interface{}{
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
					[]interface{}{
						[]interface{}{
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
						},
					},
				},
			}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pm, err := c.Pulse.PulseHistory(true)
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
			PulseProfile: &pulseprofile.PulseProfile{
				ID:            "abc123",
				MTS:           1591614631576,
				Nickname:      "nickname",
				Picture:       "picture",
				Text:          "text",
				TwitterHandle: "twitter",
			},
		}

		assert.Equal(t, expected, pm[0])
	})
}

func TestDeletePulse(t *testing.T) {
	t.Run("response", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/pulse/del", r.RequestURI)

			respMock := []interface{}{1}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		deleted, err := c.Pulse.DeletePulse("abc123")
		require.Nil(t, err)
		assert.Equal(t, 1, deleted)
	})
}
