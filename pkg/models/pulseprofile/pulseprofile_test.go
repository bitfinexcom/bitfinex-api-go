package pulseprofile_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProfileFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"abc123"}
		pp, err := pulseprofile.NewFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("sufficient arguments", func(t *testing.T) {
		payload := []interface{}{
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

		pp, err := pulseprofile.NewFromRaw(payload)
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
