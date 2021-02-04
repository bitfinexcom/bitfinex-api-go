package trades

import (
	"errors"
	"strings"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// FromWSRaw acts as a relay for public trades channel to abstract complexity from msg.
// Data arrives under "trades" channel and then splits into sub types:
// ["tu", "te", "ftu", "fte"] and can also be a snapshot.
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
			return TEUFromRaw(pair, data)
		case "te":
			return TEFromRaw(pair, data)
		case "fte":
			return FTEFromRaw(pair, data)
		case "ftu":
			return FTEUFromRaw(pair, data)
		}
	}

	return TFromRaw(pair, data)
}
