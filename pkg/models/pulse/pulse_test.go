package pulse_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPulseFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{"abc123"},
		}

		pm, err := pulse.SnapshotFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, pm)
	})

	t.Run("missing pulse profile", func(t *testing.T) {
		payload := []interface{}{
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

		pm, err := pulse.SnapshotFromRaw(payload)
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

	t.Run("has pulse profile", func(t *testing.T) {
		payload := []interface{}{
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

		pm, err := pulse.SnapshotFromRaw(payload)
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
				Followers:     30,
				Following:     5,
			},
		}

		assert.Equal(t, expected, pm[0])
	})
}
