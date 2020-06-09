package pulseprofile

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// PulseProfile data structure
type PulseProfile struct {
	ID            string
	MTS           int64
	Nickname      string
	Picture       string
	Text          string
	TwitterHandle string
}

var pulseProfileFields = map[string]int{
	"ID":            0,
	"Mts":           1,
	"Nickname":      3,
	"Picture":       5,
	"Text":          6,
	"TwitterHandle": 9,
}

// NewFromRaw takes in slice of interfaces and converts them to
// pointer to Pulse Profile
func NewFromRaw(raw []interface{}) (*PulseProfile, error) {
	if len(raw) < 14 {
		return nil, fmt.Errorf("data slice too short for PulseProfile: %#v", raw)
	}

	pp := &PulseProfile{}

	pp.ID = convert.SValOrEmpty(raw[pulseProfileFields["ID"]])
	pp.MTS = convert.I64ValOrZero(raw[pulseProfileFields["Mts"]])
	pp.Nickname = convert.SValOrEmpty(raw[pulseProfileFields["Nickname"]])
	pp.Picture = convert.SValOrEmpty(raw[pulseProfileFields["Picture"]])
	pp.Text = convert.SValOrEmpty(raw[pulseProfileFields["Text"]])
	pp.TwitterHandle = convert.SValOrEmpty(raw[pulseProfileFields["TwitterHandle"]])

	return pp, nil
}
