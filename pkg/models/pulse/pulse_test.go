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

		pm, err := pulse.SnapshotFromRaw(payload)
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

	t.Run("has pulse profile", func(t *testing.T) {
		payload := []interface{}{
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

		pm, err := pulse.SnapshotFromRaw(payload)
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
