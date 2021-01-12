package trades

import (
	"errors"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

func FromWSRaw(pair string, raw, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice for trade")
	}

	_, isSnapshot := data[0].([]interface{})
	hasType := len(raw) == 3

	if isSnapshot && strings.HasPrefix(pair, "f") {
		return FTSnapshotFromRaw(pair, convert.ToInterfaceArray(data))
	}

	if isSnapshot && strings.HasPrefix(pair, "t") {
		return TSnapshotFromRaw(pair, convert.ToInterfaceArray(data))
	}

	if hasType {
		opType, _ := raw[1].(string)

		switch opType {
		case "tu":
			return TUFromRaw(pair, data)
		case "te":
			return TEFromRaw(pair, data)
		case "fte":
			return FTEFromRaw(pair, data)
		case "ftu":
			return FTUFromRaw(pair, data)
		}
	}

	return TFromRaw(pair, data)
}
