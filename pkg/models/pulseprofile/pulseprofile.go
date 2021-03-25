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
	Followers     int64
	Following     int64
	TippingStatus int64
}

var pulseProfileFields = map[string]int{
	"ID":  0,
	"Mts": 1,
	// "PLACEHOLDER": 2,
	"Nickname": 3,
	// "PLACEHOLDER": 4,
	"Picture": 5,
	"Text":    6,
	// "PLACEHOLDER": 7,
	// "PLACEHOLDER": 8,
	"TwitterHandle": 9,
	// "PLACEHOLDER": 10,
	"Followers": 11,
	"Following": 12,
	// "PLACEHOLDER": 13,
	// "PLACEHOLDER": 14,
	// "PLACEHOLDER": 15,
	"TippingStatus": 16,
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
	pp.Followers = convert.I64ValOrZero(raw[pulseProfileFields["Followers"]])
	pp.Following = convert.I64ValOrZero(raw[pulseProfileFields["Following"]])
	pp.TippingStatus = convert.I64ValOrZero(raw[pulseProfileFields["TippingStatus"]])

	return pp, nil
}
