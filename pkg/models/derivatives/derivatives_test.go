package derivatives_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/derivatives"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDerivativeStatusFromWsRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{float64(1591614631576)}

		d, err := derivatives.NewDerivativeStatusFromWsRaw("tBTCF0:USTF0", payload)
		require.NotNil(t, err)
		require.Nil(t, d)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			float64(1591614631576),
			nil,
			9271.1234567,
			9275.3,
			nil,
			1391472.27686063,
			nil,
			1594656000000,
			-0.00011968,
			3144,
			nil,
			0,
			nil,
			nil,
			9276.06,
			nil,
			nil,
			3813.72957182,
		}

		d, err := derivatives.NewDerivativeStatusFromWsRaw("tBTCF0:USTF0", payload)
		require.Nil(t, err)

		expected := &derivatives.DerivativeStatus{
			Symbol:               "tBTCF0:USTF0",
			MTS:                  1591614631576,
			Price:                9271.1234567,
			SpotPrice:            9275.3,
			InsuranceFundBalance: 1.39147227686063e+06,
			FundingAccrued:       -0.00011968,
			MarkPrice:            9276.06,
			OpenInterest:         3813.72957182,
		}
		assert.Equal(t, expected, d)
	})
}

func TestNewDerivativeStatusFromRaw(t *testing.T) {
	t.Run("insufficient arguments", func(t *testing.T) {
		payload := []interface{}{"tBTCF0:USTF0"}

		d, err := derivatives.NewDerivativeStatusFromRaw(payload)
		require.NotNil(t, err)
		require.Nil(t, d)
	})

	t.Run("valid arguments", func(t *testing.T) {
		payload := []interface{}{
			"tBTCF0:USTF0",
			float64(1591614631576),
			nil,
			9271.1234567,
			9275.3,
			nil,
			1391472.27686063,
			nil,
			1594656000000,
			-0.00011968,
			3144,
			nil,
			0,
			nil,
			nil,
			9276.06,
			nil,
			nil,
			3813.72957182,
		}

		d, err := derivatives.NewDerivativeStatusFromRaw(payload)
		require.Nil(t, err)

		expected := &derivatives.DerivativeStatus{
			Symbol:               "tBTCF0:USTF0",
			MTS:                  1591614631576,
			Price:                9271.1234567,
			SpotPrice:            9275.3,
			InsuranceFundBalance: 1.39147227686063e+06,
			FundingAccrued:       -0.00011968,
			MarkPrice:            9276.06,
			OpenInterest:         3813.72957182,
		}
		assert.Equal(t, expected, d)
	})
}
