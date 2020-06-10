package convert_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItfToStrSlice(t *testing.T) {
	t.Run("invalid slice arguments", func(t *testing.T) {
		payload := []interface{}{123, 234, 345}
		got, err := convert.ItfToStrSlice(payload)
		require.NotNil(t, err)
		require.Nil(t, got)
	})

	t.Run("non slice arguments", func(t *testing.T) {
		payload := "123"
		got, err := convert.ItfToStrSlice(payload)
		require.Nil(t, err)
		assert.Equal(t, []string{}, got)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{"foo", "bar", "baz"}
		got, err := convert.ItfToStrSlice(payload)
		expected := []string{"foo", "bar", "baz"}
		require.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}

func TestToInt(t *testing.T) {
	t.Run("valid argument", func(t *testing.T) {
		payload := 1234
		expected := 1234
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})

	t.Run("invalid argument", func(t *testing.T) {
		payload := "1"
		expected := 0
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})
}
