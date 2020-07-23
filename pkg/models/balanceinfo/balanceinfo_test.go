package balanceinfo_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceInfoFromRaw(t *testing.T) {
	t.Run("invalid arguments", func(t *testing.T) {
		payload := []interface{}{4131.85}

		b, err := balanceinfo.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, b)
	})

	t.Run("valid trading arguments", func(t *testing.T) {
		payload := []interface{}{4131.85, 4131.85}

		b, err := balanceinfo.FromRaw(payload)
		require.Nil(t, err)

		expected := &balanceinfo.BalanceInfo{
			TotalAUM: 4131.85,
			NetAUM:   4131.85,
		}
		assert.Equal(t, expected, b)
	})
}
