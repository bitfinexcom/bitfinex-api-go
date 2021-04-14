package fundingcredit_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingcredit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		flcr := fundingcredit.CancelRequest{ID: 123}
		got, err := flcr.MarshalJSON()

		require.Nil(t, err)

		expected := "[0, \"fcc\", null, {\"id\":123}]"
		assert.Equal(t, expected, string(got))
	})
}
