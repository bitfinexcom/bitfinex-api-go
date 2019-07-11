package bitfinex

import (
	"fmt"
	"encoding/json"
	"strconv"
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

func ToInterfaceArray(i []interface{}) [][]interface{} {
	newArr := make([][]interface{}, len(i))
	for index, item := range i {
		newArr[index] = item.([]interface{})
	}
	return newArr
}

func ToFloat64Array(i [][]interface{}) ([][]float64, error) {
	newArr := make([][]float64, len(i))
	for index, item := range i {
		s, err := F64Slice(item)
		if err != nil {
			return nil, err
		}
		newArr[index] = s
	}
	return newArr, nil
}

func floatToJsonNumber(i interface{}) json.Number {
	if r, ok := i.(json.Number); ok {
		return r
	}
	return json.Number(strconv.FormatFloat(i.(float64), 'f', -1, 64))
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

func f64ValOrZero(i interface{}) float64 {
	if r, ok := i.(float64); ok {
		return r
	}
	return 0.0
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
