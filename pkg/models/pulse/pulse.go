package pulse

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

// Pulse message data structure
type Pulse struct {
	ID           string
	MTS          int64
	UserID       string
	Title        string
	Content      string
	IsPin        int
	IsPublic     int
	Tags         []string
	Attachments  []string
	Likes        int
	PulseProfile *pulseprofile.PulseProfile
}

var pulseFields = map[string]int{
	"ID":           0,
	"Mts":          1,
	"UserID":       3,
	"Title":        5,
	"Content":      6,
	"IsPin":        9,
	"IsPublic":     10,
	"Tags":         12,
	"Attachments":  13,
	"Likes":        15,
	"PulseProfile": 18,
}

func newSingleFromRaw(raw []interface{}) (*Pulse, error) {
	if len(raw) < 19 {
		return nil, fmt.Errorf("data slice too short for Pulse Message: %#v", raw)
	}

	p := &Pulse{}
	var err error

	p.ID = convert.SValOrEmpty(raw[pulseFields["ID"]])
	p.MTS = convert.I64ValOrZero(raw[pulseFields["Mts"]])
	p.UserID = convert.SValOrEmpty(raw[pulseFields["UserID"]])
	p.Title = convert.SValOrEmpty(raw[pulseFields["Title"]])
	p.Content = convert.SValOrEmpty(raw[pulseFields["Content"]])
	p.IsPin = convert.ToInt(raw[pulseFields["IsPin"]])
	p.IsPublic = convert.ToInt(raw[pulseFields["IsPublic"]])
	p.Likes = convert.ToInt(raw[pulseFields["Likes"]])

	p.Tags, err = convert.ItfToStrSlice(raw[pulseFields["Tags"]])
	if err != nil {
		return nil, err
	}

	p.Attachments, err = convert.ItfToStrSlice(raw[pulseFields["Attachments"]])
	if err != nil {
		return nil, err
	}

	rawProfile := raw[pulseFields["PulseProfile"]]
	rawProfileItf, ok := rawProfile.([]interface{})
	if !ok {
		// if, for some reazon, we can't extract profile data at given index,
		// we return whatever pulse data we have in place
		return p, nil
	}

	rawProfilePayload := convert.ToInterfaceArray(rawProfileItf)[0]
	p.PulseProfile, err = pulseprofile.NewFromRaw(rawProfilePayload)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewFromRaw returns slice of Pulse messages
func NewFromRaw(raws []interface{}) ([]*Pulse, error) {
	if len(raws) < 1 {
		return nil, fmt.Errorf("data slice is too short for Pulse History: %#v", raws)
	}

	res := []*Pulse{}

	for _, raw := range raws {
		raw := raw.([]interface{})
		p, err := newSingleFromRaw(raw)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	return res, nil
}
