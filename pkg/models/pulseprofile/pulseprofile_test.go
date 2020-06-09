package pulseprofile_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"abc123"}
		pp, err := pulseprofile.NewFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, pp)
	})

	t.Run("sufficient arguments", func(t *testing.T) {
		payload := []interface{}{
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

		pp, err := pulseprofile.NewFromRaw(payload)
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
