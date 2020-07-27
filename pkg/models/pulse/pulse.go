package pulse

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

// Pulse message data structure
type Pulse struct {
	ID           string                     `json:"id,omitempty"`
	MTS          int64                      `json:"mts,omitempty"`
	UserID       string                     `json:"userId,omitempty"`
	Title        string                     `json:"title,omitempty"`
	Content      string                     `json:"content,omitempty"`
	IsPin        int                        `json:"isPin"`
	IsPublic     int                        `json:"isPublic"`
	Tags         []string                   `json:"tags,omitempty"`
	Attachments  []string                   `json:"attachments,omitempty"`
	Likes        int                        `json:"likes,omitempty"`
	PulseProfile *pulseprofile.PulseProfile `json:"pulseProfile,omitempty"`
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

// FromRaw returns pointer to Pulse message
func FromRaw(raw []interface{}) (*Pulse, error) {
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
		return p, nil
	}

	profilePayload := convert.ToInterfaceArray(rawProfileItf)
	if len(profilePayload) < 1 {
		return p, nil
	}

	p.PulseProfile, err = pulseprofile.NewFromRaw(profilePayload[0])
	if err != nil {
		return nil, err
	}

	return p, nil
}

// SnapshotFromRaw returns slice of Pulse message pointers
func SnapshotFromRaw(raws []interface{}) ([]*Pulse, error) {
	if len(raws) < 1 {
		return nil, fmt.Errorf("data slice is too short for Pulse History: %#v", raws)
	}

	res := []*Pulse{}

	for _, raw := range raws {
		raw := raw.([]interface{})
		p, err := FromRaw(raw)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	return res, nil
}
