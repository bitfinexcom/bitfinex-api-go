package convert

import (
	"encoding/json"
	"fmt"
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

func ItfToStrSlice(in interface{}) ([]string, error) {
	ret := []string{}
	raw, ok := in.([]interface{})
	if !ok {
		return ret, nil
	}

	for _, e := range raw {
		item, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("expected slice of strings, got: %v", in)
		}
		ret = append(ret, item)
	}

	return ret, nil
}

// ToInt converts various types to integer. If fails, returns 0
func ToInt(in interface{}) int {
	var out int

	switch v := in.(type) {
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			out = i
		}
	case float64:
		out = int(v)
	default:
		if val, ok := in.(int); ok {
			out = val
		}
	}

	return out
}

func ToInterfaceArray(i []interface{}) [][]interface{} {
	newArr := make([][]interface{}, len(i))
	for index, item := range i {
		newArr[index] = item.([]interface{})
	}
	return newArr
}

func ToInterface(flt []float64) []interface{} {
	data := make([]interface{}, len(flt))
	for j, f := range flt {
		data[j] = f
	}
	return data
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

func FloatToJsonNumber(i interface{}) json.Number {
	if r, ok := i.(json.Number); ok {
		return r
	}
	return json.Number(strconv.FormatFloat(i.(float64), 'f', -1, 64))
}

func I64ValOrZero(in interface{}) (out int64) {
	switch v := in.(type) {
	case int:
		out = int64(v)
	default:
		if v, ok := in.(float64); ok {
			out = int64(v)
		}
	}

	return out
}

func IValOrZero(i interface{}) int {
	if r, ok := i.(float64); ok {
		return int(r)
	}
	return 0
}

func F64ValOrZero(in interface{}) (out float64) {
	switch v := in.(type) {
	case int:
		out = float64(v)
	default:
		if v, ok := in.(float64); ok {
			out = v
		}
	}

	return out
}

func SiMapOrEmpty(i interface{}) map[string]interface{} {
	if m, ok := i.(map[string]interface{}); ok {
		return m
	}
	return make(map[string]interface{})
}

func BValOrFalse(in interface{}) (out bool) {
	switch v := in.(type) {
	case int:
		if v == 1 {
			out = true
		}
	case string:
		if v == "1" {
			out = true
		}
	default:
		if v, ok := in.(bool); ok {
			out = v
		}
	}
	return
}

func SValOrEmpty(i interface{}) string {
	if r, ok := i.(string); ok {
		return r
	}
	return ""
}
