package fundingoffer_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundingoffer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		fosr := fundingoffer.SubmitRequest{
			Type:   "LIMIT",
			Symbol: "fETH",
			Amount: 0.29797676,
			Rate:   0.0002,
			Period: 2,
		}
		got, err := fosr.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"fon\", null, {\"type\":\"LIMIT\",\"symbol\":\"fETH\",\"amount\":\"0.29797676\",\"rate\":\"0.0002\",\"period\":2}]"
		assert.Equal(t, expected, string(got))
	})
}
