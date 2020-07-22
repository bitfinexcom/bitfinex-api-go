package order_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderNewRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		our := order.NewRequest{
			CID:    987,
			GID:    876,
			Type:   "EXCHANGE LIMIT",
			Symbol: "tBTCUSD",
			Price:  13,
			Amount: 0.001,
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"on\", null, {\"gid\":876,\"cid\":987,\"type\":\"EXCHANGE LIMIT\",\"symbol\":\"tBTCUSD\",\"amount\":\"0.001\",\"price\":\"13\"}]"
		assert.Equal(t, expected, string(got))
	})

	t.Run("MarshalJSON with extra props", func(t *testing.T) {
		our := order.NewRequest{
			CID:           987,
			GID:           876,
			Type:          "EXCHANGE LIMIT",
			Symbol:        "tBTCUSD",
			Price:         13,
			Amount:        0.001,
			Hidden:        true,
			PostOnly:      true,
			OcoOrder:      true,
			Close:         true,
			AffiliateCode: "abc",
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"on\", null, {\"gid\":876,\"cid\":987,\"type\":\"EXCHANGE LIMIT\",\"symbol\":\"tBTCUSD\",\"amount\":\"0.001\",\"price\":\"13\",\"flags\":21056,\"meta\":{\"aff_code\":\"abc\"}}]"
		assert.Equal(t, expected, string(got))
	})
}

func TestOrderUpdateRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		our := order.UpdateRequest{
			ID:     123456,
			GID:    234567,
			Price:  15.1234,
			Amount: 0.002,
			Hidden: false,
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"ou\", null, {\"id\":123456,\"gid\":234567,\"price\":\"15.1234\",\"amount\":\"0.002\"}]"
		assert.Equal(t, expected, string(got))
	})

	t.Run("MarshalJSON hidden", func(t *testing.T) {
		our := order.UpdateRequest{
			ID:     123456,
			GID:    234567,
			Price:  15.1234,
			Amount: 0.002,
			Hidden: true,
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"ou\", null, {\"id\":123456,\"gid\":234567,\"price\":\"15.1234\",\"amount\":\"0.002\",\"flags\":64}]"
		assert.Equal(t, expected, string(got))
	})

	t.Run("MarshalJSON PostOnly", func(t *testing.T) {
		our := order.UpdateRequest{
			ID:       123456,
			GID:      234567,
			Price:    15.1234,
			Amount:   0.002,
			PostOnly: true,
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"ou\", null, {\"id\":123456,\"gid\":234567,\"price\":\"15.1234\",\"amount\":\"0.002\",\"flags\":4096}]"
		assert.Equal(t, expected, string(got))
	})

	t.Run("MarshalJSON hidden and PostOnly", func(t *testing.T) {
		our := order.UpdateRequest{
			ID:       123456,
			GID:      234567,
			Price:    15.1234,
			Amount:   0.002,
			PostOnly: true,
			Hidden:   true,
		}

		got, err := our.MarshalJSON()
		require.Nil(t, err)

		expected := "[0, \"ou\", null, {\"id\":123456,\"gid\":234567,\"price\":\"15.1234\",\"amount\":\"0.002\",\"flags\":4160}]"
		assert.Equal(t, expected, string(got))
	})
}

func TestOrderCancelRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		ocr := order.CancelRequest{
			ID:      123456,
			CID:     234567,
			CIDDate: "2020-10-10",
		}
		got, err := ocr.MarshalJSON()

		require.Nil(t, err)

		expected := "[0, \"oc\", null, {\"id\":123456,\"cid\":234567,\"cid_date\":\"2020-10-10\"}]"
		assert.Equal(t, expected, string(got))
	})
}
