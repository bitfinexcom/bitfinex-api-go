package fundinginfo

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type FundingInfo struct {
	Symbol       string
	YieldLoan    float64
	YieldLend    float64
	DurationLoan float64
	DurationLend float64
}

func FromRaw(raw []interface{}) (fi *FundingInfo, err error) {
	if len(raw) < 3 { // "sym", symbol, data
		return fi, fmt.Errorf("data slice too short for funding info: %#v", raw)
	}

	sym, ok := raw[1].(string)
	if !ok {
		return fi, fmt.Errorf("expected symbol in second position of funding info: %v", raw)
	}

	data, ok := raw[2].([]interface{})
	if !ok {
		return fi, fmt.Errorf("expected list in third position of funding info: %v", raw)
	}

	if len(data) < 4 {
		return fi, fmt.Errorf("data too short: %#v", data)
	}

	fi = &FundingInfo{
		Symbol:       sym,
		YieldLoan:    convert.F64ValOrZero(data[0]),
		YieldLend:    convert.F64ValOrZero(data[1]),
		DurationLoan: convert.F64ValOrZero(data[2]),
		DurationLend: convert.F64ValOrZero(data[3]),
	}

	return
}
