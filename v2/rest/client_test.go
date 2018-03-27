package rest

import (
	"testing"
)

func TestReadParams(t *testing.T) {
	params := ReadParams()

	if params == nil {
		t.Error("Params should exist!")
	}

	m := make(map[string]interface{})
	m["limit"] = 2
	params = ReadParams(m)

	if params["limit"] != 2 {
		t.Error("Limit should be 1")
	}
}
