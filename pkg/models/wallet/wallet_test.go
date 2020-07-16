package wallet_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWalletFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"exchange"}

		w, err := wallet.FromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, w)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"exchange",
			"SAN",
			19.76,
			0,
		}

		w, err := wallet.FromRaw(payload)
		require.Nil(t, err)

		expected := &wallet.Wallet{
			Type:              "exchange",
			Currency:          "SAN",
			Balance:           19.76,
			UnsettledInterest: 0,
		}

		assert.Equal(t, expected, w)
	})
}

func TestNewWalletFromWsRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"exchange"}

		invc, err := wallet.FromWsRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, invc)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"exchange",
			"SAN",
			19.76,
			0,
			12.1234,
		}

		w, err := wallet.FromWsRaw(payload)
		require.Nil(t, err)

		expected := &wallet.Wallet{
			Type:              "exchange",
			Currency:          "SAN",
			Balance:           19.76,
			UnsettledInterest: 0,
			BalanceAvailable:  12.1234,
		}

		assert.Equal(t, expected, w)
	})
}

func TestWalletSnapshotFromRaw(t *testing.T) {
	t.Run("rest success", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				"exchange",
				"SAN",
				19.76,
				0,
			},
			[]interface{}{
				"exchange",
				"SAN",
				20.76,
				2,
			},
		}

		ss, err := wallet.SnapshotFromRaw(payload, wallet.FromRaw)
		require.Nil(t, err)

		assert.Equal(t, &wallet.Snapshot{
			Snapshot: []*wallet.Wallet{
				{
					Type:              "exchange",
					Currency:          "SAN",
					Balance:           19.76,
					UnsettledInterest: 0,
				},
				{
					Type:              "exchange",
					Currency:          "SAN",
					Balance:           20.76,
					UnsettledInterest: 2,
				},
			},
		}, ss)
	})

	t.Run("rest fail", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{"exchange"},
			[]interface{}{
				"exchange",
				"SAN",
				20.76,
				2,
			},
		}

		ss, err := wallet.SnapshotFromRaw(payload, wallet.FromRaw)
		require.NotNil(t, err)
		require.Nil(t, ss)
	})

	t.Run("ws success", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{
				"exchange",
				"SAN",
				19.76,
				0,
				12.1234,
			},
			[]interface{}{
				"exchange",
				"SAN",
				20.76,
				2,
				12.12345,
			},
		}

		ss, err := wallet.SnapshotFromRaw(payload, wallet.FromWsRaw)
		require.Nil(t, err)

		assert.Equal(t, &wallet.Snapshot{
			Snapshot: []*wallet.Wallet{
				{
					Type:              "exchange",
					Currency:          "SAN",
					Balance:           19.76,
					UnsettledInterest: 0,
					BalanceAvailable:  12.1234,
				},
				{
					Type:              "exchange",
					Currency:          "SAN",
					Balance:           20.76,
					UnsettledInterest: 2,
					BalanceAvailable:  12.12345,
				},
			},
		}, ss)
	})

	t.Run("ws fail", func(t *testing.T) {
		payload := []interface{}{
			[]interface{}{"exchange"},
			[]interface{}{
				"exchange",
				"SAN",
				20.76,
				2,
				12.12345,
			},
		}

		ss, err := wallet.SnapshotFromRaw(payload, wallet.FromWsRaw)
		require.NotNil(t, err)
		require.Nil(t, ss)
	})
}
