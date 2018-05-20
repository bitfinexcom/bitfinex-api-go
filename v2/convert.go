package bitfinex

import (
	"fmt"
)

func F64Slice(in []interface{}) ([]float64, error) {
	var ret []float64
	for _, e := range in {
		if item, ok := e.(float64); ok {
			ret = append(ret, item)
		} else {
			return nil, fmt.Errorf("expected slice of float64 but got: %v", in)
		}
	}

	return ret, nil
}

func i64ValOrZero(i interface{}) int64 {
	if r, ok := i.(float64); ok {
		return int64(r)
	}
	return 0
}

func iValOrZero(i interface{}) int {
	if r, ok := i.(float64); ok {
		return int(r)
	}
	return 0
}

func i64pValOrNil(i interface{}) *int64 {
	if i == nil {
		return nil
	}

	if r, ok := i.(int64); ok {
		return &r
	}
	return nil
}

func ui64ValOrZero(i interface{}) uint64 {
	if r, ok := i.(float64); ok {
		return uint64(r)
	}
	return 0
}

func f64ValOrZero(i interface{}) float64 {
	if r, ok := i.(float64); ok {
		return r
	}
	return 0.0
}

func f64pValOrNil(i interface{}) *float64 {
	if i == nil {
		return nil
	}

	if r, ok := i.(float64); ok {
		return &r
	}
	return nil
}

func bValOrFalse(i interface{}) bool {
	if r, ok := i.(bool); ok {
		return r
	}
	return false
}

func sValOrEmpty(i interface{}) string {
	if r, ok := i.(string); ok {
		return r
	}
	return ""
}
