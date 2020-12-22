package status

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// FromWSRaw - based on condition will return snapshot or single record of
// derivative or liquidation data structure
func FromWSRaw(key string, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice")
	}

	_, isSnapshot := data[0].([]interface{})
	ss := strings.SplitN(key, ":", 2)
	if len(ss) < 2 {
		return nil, fmt.Errorf("unexpected key: %s", key)
	}

	if isSnapshot && ss[0] == "deriv" {
		return DerivSnapshotFromRaw(ss[1], convert.ToInterfaceArray(data))
	}

	if !isSnapshot && ss[0] == "deriv" {
		return DerivFromRaw(ss[1], data)
	}

	if isSnapshot && ss[0] == "liq" {
		return LiqSnapshotFromRaw(convert.ToInterfaceArray(data))
	}

	if !isSnapshot && ss[0] == "liq" {
		return LiqFromRaw(data)
	}

	return nil, fmt.Errorf("%s: unrecognized data slice:%#v", key, data)
}
