package fundingoffer_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		focr := fundingoffer.CancelRequest{ID: 123}

		got, err := focr.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"foc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}
