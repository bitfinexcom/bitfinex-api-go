package book_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBookFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{1591614631576}
		rawNums := []interface{}{1.12345}

		b, err := book.FromRaw("tBTCUSD", "P0", payload, rawNums)
		require.NotNil(t, err)
		require.Nil(t, b)
	})

	t.Run("valid trading arguments", func(t *testing.T) {
		payload := []interface{}{
			98169.99541156, 2, 0.000202,
		}

		rawNums := []interface{}{98169.99541156, 2, 0.000202}

		b, err := book.FromRaw("tBTCUSD", "P0", payload, rawNums)
		require.Nil(t, err)

		expected := &book.Book{
			Symbol:      "tBTCUSD",
			Price:       98169.99541156,
			PriceJsNum:  "98169.99541156",
			Count:       2,
			Amount:      0.000202,
			AmountJsNum: "0.000202",
			Side:        common.Bid,
			Action:      book.BookEntry,
		}
		assert.Equal(t, expected, b)
	})

	t.Run("valid raw trading arguments", func(t *testing.T) {
		payload := []interface{}{
			34006738527, 8744.9, 0.25603413,
		}

		rawNums := []interface{}{34006738527, 8744.9, 0.25603413}

		b, err := book.FromRaw("tBTCUSD", "R0", payload, rawNums)
		require.Nil(t, err)

		expected := &book.Book{
			ID:          34006738527,
			Symbol:      "tBTCUSD",
			Price:       8744.9,
			PriceJsNum:  "8744.9",
			Amount:      0.25603413,
			AmountJsNum: "0.25603413",
			Side:        common.Bid,
			Action:      book.BookEntry,
		}
		assert.Equal(t, expected, b)
	})

	t.Run("valid funding arguments", func(t *testing.T) {
		t.Skip("skipping as implementation is missing atm")
		payload := []interface{}{
			0.0003301, 30, 1, -3862.874,
		}

		rawNums := []interface{}{0.0003301, 30, 1, -3862.874}

		b, err := book.FromRaw("fUSD", "P0", payload, rawNums)
		require.Nil(t, err)

		expected := &book.Book{}
		assert.Equal(t, expected, b)
	})

	t.Run("valid raw funding arguments", func(t *testing.T) {
		t.Skip("skipping as implementation is missing atm")
		payload := []interface{}{
			645902785, 30, 0.0003301, -3862.874,
		}

		rawNums := []interface{}{645902785, 30, 0.0003301, -3862.874}

		b, err := book.FromRaw("fUSD", "R0", payload, rawNums)
		require.Nil(t, err)

		expected := &book.Book{}
		assert.Equal(t, expected, b)
	})
}
