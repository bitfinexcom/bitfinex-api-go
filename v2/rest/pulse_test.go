package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
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
				"bf324e24-5a09-4317-b418-6c37281ab855",
				1591614631576,
				nil,
				"Bitfinex",
				nil,
				"image-33533a4d-a796-4afe-9b8b-690bc7075e05-1587476823358-size.png",
				"Bitfinex is the world’s leading digital asset trading platform.",
				nil,
				nil,
				"bitfinex",
				nil,
				40,
				5,
				nil,
				nil,
				nil,
				0,
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
			ID:            "bf324e24-5a09-4317-b418-6c37281ab855",
			MTS:           1591614631576,
			Nickname:      "Bitfinex",
			Picture:       "image-33533a4d-a796-4afe-9b8b-690bc7075e05-1587476823358-size.png",
			Text:          "Bitfinex is the world’s leading digital asset trading platform.",
			TwitterHandle: "bitfinex",
			Followers:     40,
			Following:     5,
			TippingStatus: 0,
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
		end := common.Mts(millis)
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
					"d139512a-6558-431a-a6aa-69646fc97d55",
					1593608548140,
					nil,
					"bf324e24-5a09-4317-b418-6c37281ab855",
					nil,
					"Take an active role in the discussion with the Comment feature!",
					"Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
					nil,
					nil,
					1,
					1,
					nil,
					[]interface{}{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
					[]interface{}{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
					[]interface{}{},
					1,
					nil,
					nil,
					nil,
					1,
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
		end := common.Mts(1591614631576)
		pph, err := c.Pulse.PublicPulseHistory(1, end)
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:               "d139512a-6558-431a-a6aa-69646fc97d55",
			MTS:              1593608548140,
			UserID:           "bf324e24-5a09-4317-b418-6c37281ab855",
			Title:            "Take an active role in the discussion with the Comment feature!",
			Content:          "Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
			IsPin:            1,
			IsPublic:         1,
			CommentsDisabled: 0,
			Tags:             []string{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
			Attachments:      []string{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
			Likes:            1,
			Comments:         1,
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
				"d139512a-6558-431a-a6aa-69646fc97d55",
				1593608548140,
				nil,
				"bf324e24-5a09-4317-b418-6c37281ab855",
				nil,
				"Take an active role in the discussion with the Comment feature!",
				"Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
				nil,
				nil,
				1,
				1,
				nil,
				[]interface{}{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
				[]interface{}{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
				[]interface{}{},
				1,
				nil,
				nil,
				nil,
				1,
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
			ID:               "d139512a-6558-431a-a6aa-69646fc97d55",
			MTS:              1593608548140,
			UserID:           "bf324e24-5a09-4317-b418-6c37281ab855",
			Title:            "Take an active role in the discussion with the Comment feature!",
			Content:          "Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
			IsPin:            1,
			IsPublic:         1,
			CommentsDisabled: 0,
			Tags:             []string{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
			Attachments:      []string{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
			Likes:            1,
			Comments:         1,
		}

		assert.Equal(t, expected, pm)
	})
}

func TestPulseHistory(t *testing.T) {
	t.Run("response data slice too short", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/r/pulse/hist", r.RequestURI)

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
		pp, err := c.Pulse.PulseHistory()
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("private pulse history", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{
				[]interface{}{
					"d139512a-6558-431a-a6aa-69646fc97d55",
					1593608548140,
					nil,
					"bf324e24-5a09-4317-b418-6c37281ab855",
					nil,
					"Take an active role in the discussion with the Comment feature!",
					"Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
					nil,
					nil,
					1,
					1,
					nil,
					[]interface{}{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
					[]interface{}{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
					[]interface{}{},
					1,
					nil,
					nil,
					nil,
					1,
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
		pm, err := c.Pulse.PulseHistory()
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:               "d139512a-6558-431a-a6aa-69646fc97d55",
			MTS:              1593608548140,
			UserID:           "bf324e24-5a09-4317-b418-6c37281ab855",
			Title:            "Take an active role in the discussion with the Comment feature!",
			Content:          "Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
			IsPin:            1,
			IsPublic:         1,
			CommentsDisabled: 0,
			Tags:             []string{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
			Attachments:      []string{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
			Likes:            1,
			Comments:         1,
		}

		assert.Equal(t, expected, pm[0])
	})

	t.Run("public pulse history", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{
				[]interface{}{
					"d139512a-6558-431a-a6aa-69646fc97d55",
					1593608548140,
					nil,
					"bf324e24-5a09-4317-b418-6c37281ab855",
					nil,
					"Take an active role in the discussion with the Comment feature!",
					"Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
					nil,
					nil,
					1,
					1,
					nil,
					[]interface{}{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
					[]interface{}{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
					[]interface{}{},
					1,
					nil,
					nil,
					[]interface{}{
						[]interface{}{
							"bf324e24-5a09-4317-b418-6c37281ab855",
							1591614631576,
							nil,
							"Bitfinex",
							nil,
							"image-33533a4d-a796-4afe-9b8b-690bc7075e05-1587476823358-size.png",
							"Bitfinex is the world’s leading digital asset trading platform.",
							nil,
							nil,
							"bitfinex",
							nil,
							40,
							5,
							nil,
							nil,
							nil,
							0,
						},
					},
					1,
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
		pm, err := c.Pulse.PulseHistory()
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:               "d139512a-6558-431a-a6aa-69646fc97d55",
			MTS:              1593608548140,
			UserID:           "bf324e24-5a09-4317-b418-6c37281ab855",
			Title:            "Take an active role in the discussion with the Comment feature!",
			Content:          "Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
			IsPin:            1,
			IsPublic:         1,
			CommentsDisabled: 0,
			Tags:             []string{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
			Attachments:      []string{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
			Likes:            1,
			Comments:         1,
			PulseProfile: &pulseprofile.PulseProfile{
				ID:            "bf324e24-5a09-4317-b418-6c37281ab855",
				MTS:           1591614631576,
				Nickname:      "Bitfinex",
				Picture:       "image-33533a4d-a796-4afe-9b8b-690bc7075e05-1587476823358-size.png",
				Text:          "Bitfinex is the world’s leading digital asset trading platform.",
				TwitterHandle: "bitfinex",
				Followers:     40,
				Following:     5,
				TippingStatus: 0,
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

func TestAddPulseComment(t *testing.T) {
	t.Run("invalid payload", func(t *testing.T) {
		p := &pulse.Pulse{Title: "foo"}
		c := NewClient()
		pm, err := c.Pulse.AddComment(p)
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
		pm, err := c.Pulse.AddComment(&pulse.Pulse{Title: "foo bar baz qux 123"})
		require.NotNil(t, err)
		require.Nil(t, pm)
	})

	t.Run("valid payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			respMock := []interface{}{
				"d139512a-6558-431a-a6aa-69646fc97d55",
				1593608548140,
				nil,
				"bf324e24-5a09-4317-b418-6c37281ab855",
				nil,
				"Take an active role in the discussion with the Comment feature!",
				"Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
				nil,
				nil,
				1,
				1,
				nil,
				[]interface{}{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
				[]interface{}{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
				[]interface{}{},
				1,
				nil,
				nil,
				nil,
				1,
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
		payload := &pulse.Pulse{Title: "foo bar baz qux 123", Parent: "foo"}
		pm, err := c.Pulse.AddComment(payload)
		require.Nil(t, err)

		expected := &pulse.Pulse{
			ID:               "d139512a-6558-431a-a6aa-69646fc97d55",
			MTS:              1593608548140,
			UserID:           "bf324e24-5a09-4317-b418-6c37281ab855",
			Title:            "Take an active role in the discussion with the Comment feature!",
			Content:          "Bitfinex Pulse, the social trading you were waiting for, keeps improving to give you a complete social trading experience.",
			IsPin:            1,
			IsPublic:         1,
			CommentsDisabled: 0,
			Tags:             []string{"BitfinexPulse", "PulseFeatures", "PulseUpdates", "L_EN", "N_Bitfinex"},
			Attachments:      []string{"image-a55eda4c-4f5a-4609-abf0-7bdfca1a8911-1593608547535-size.jpg"},
			Likes:            1,
			Comments:         1,
		}

		assert.Equal(t, expected, pm)
	})
}
