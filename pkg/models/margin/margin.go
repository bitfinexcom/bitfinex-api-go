package margin

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type InfoBase struct {
	UserProfitLoss float64
	UserSwaps      float64
	MarginBalance  float64
	MarginNet      float64
	MarginRequired float64
}

type InfoUpdate struct {
	Symbol          string
	TradableBalance float64
	GrossBalance    float64
	Buy             float64
	Sell            float64
}

// FromRaw returns either a InfoBase or InfoUpdate, since
// the Margin Info is split up into a base and per symbol parts.
func FromRaw(raw []interface{}) (o interface{}, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	typ, ok := raw[0].(string)
	if !ok {
		return o, fmt.Errorf("expected margin info type in first position for margin info but got %#v", raw)
	}

	if len(raw) > 1 && typ == "base" { // This should be ["base", [...]]
		data, ok := raw[1].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in second position for margin info but got %#v", raw)
		}

		return baseFromRaw(data)
	}

	if len(raw) > 2 && typ == "sym" { // This should be ["sym", SYMBOL, [...]]
		symbol, ok := raw[1].(string)
		if !ok {
			return o, fmt.Errorf("expected margin info symbol in second position for margin info update but got %#v", raw)
		}

		data, ok := raw[2].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in third position for margin info update but got %#v", raw)
		}

		return updateFromRaw(symbol, data)
	}

	return nil, fmt.Errorf("invalid margin info type in %#v", raw)
}

func updateFromRaw(symbol string, raw []interface{}) (o *InfoUpdate, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for margin info update: %#v", raw)
	}

	o = &InfoUpdate{
		Symbol:          symbol,
		TradableBalance: convert.F64ValOrZero(raw[0]),
		GrossBalance:    convert.F64ValOrZero(raw[1]),
		Buy:             convert.F64ValOrZero(raw[2]),
		Sell:            convert.F64ValOrZero(raw[3]),
	}

	return
}

func baseFromRaw(raw []interface{}) (ib *InfoBase, err error) {
	if len(raw) < 5 {
		return ib, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	ib = &InfoBase{
		UserProfitLoss: convert.F64ValOrZero(raw[0]),
		UserSwaps:      convert.F64ValOrZero(raw[1]),
		MarginBalance:  convert.F64ValOrZero(raw[2]),
		MarginNet:      convert.F64ValOrZero(raw[3]),
		MarginRequired: convert.F64ValOrZero(raw[4]),
	}

	return
}
