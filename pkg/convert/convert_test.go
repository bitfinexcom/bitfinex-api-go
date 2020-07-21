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
	t.Run("valid int argument", func(t *testing.T) {
		payload := 1234
		expected := 1234
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})

	t.Run("valid string int", func(t *testing.T) {
		payload := "1"
		expected := 1
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})

	t.Run("float64", func(t *testing.T) {
		var payload float64 = 1234
		expected := 1234
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})

	t.Run("invalid string int", func(t *testing.T) {
		payload := "foo"
		expected := 0
		got := convert.ToInt(payload)
		assert.Equal(t, expected, got)
	})
}

func TestToInterface(t *testing.T) {
	payload := []float64{1.1234, 2.1234}
	expected := []interface{}{1.1234, 2.1234}
	got := convert.ToInterface(payload)
	assert.Equal(t, expected, got)
}

func TestF64ValOrZero(t *testing.T) {
	t.Run("converts int to float64", func(t *testing.T) {
		var expected float64 = 910
		got := convert.F64ValOrZero(910)
		assert.Equal(t, expected, got)
	})

	t.Run("converts float64 to float64", func(t *testing.T) {
		var expected float64 = 910.1234
		got := convert.F64ValOrZero(float64(910.1234))
		assert.Equal(t, expected, got)
	})
}

func TestI64ValOrZero(t *testing.T) {
	t.Run("converts int to int64", func(t *testing.T) {
		var expected int64 = 910
		got := convert.I64ValOrZero(910)
		assert.Equal(t, expected, got)
	})

	t.Run("converts float64 to int64", func(t *testing.T) {
		var expected int64 = 1594891800000
		got := convert.I64ValOrZero(1.5948918e+12)
		assert.Equal(t, expected, got)
	})
}

func TestBValOrFalse(t *testing.T) {
	t.Run("converts to truthy bool", func(t *testing.T) {
		expected := true
		got := convert.BValOrFalse(true)
		assert.Equal(t, expected, got)
	})

	t.Run("converts to falsy bool", func(t *testing.T) {
		expected := false
		got := convert.BValOrFalse(false)
		assert.Equal(t, expected, got)
	})

	t.Run("converts 1 to truthy bool", func(t *testing.T) {
		expected := true
		got := convert.BValOrFalse(1)
		assert.Equal(t, expected, got)
	})

	t.Run("converts \"1\" to truthy bool", func(t *testing.T) {
		expected := true
		got := convert.BValOrFalse("1")
		assert.Equal(t, expected, got)
	})

	t.Run("converts 0 to falsy bool", func(t *testing.T) {
		expected := false
		got := convert.BValOrFalse(0)
		assert.Equal(t, expected, got)
	})

	t.Run("converts \"0\" to falsy bool", func(t *testing.T) {
		expected := false
		got := convert.BValOrFalse("0")
		assert.Equal(t, expected, got)
	})
}
