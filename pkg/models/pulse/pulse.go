package pulse

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

// Pulse message data structure
type Pulse struct {
	ID               string                     `json:"id,omitempty"`
	Parent           string                     `json:"parent,omitempty"`
	MTS              int64                      `json:"mts,omitempty"`
	UserID           string                     `json:"userId,omitempty"`
	Title            string                     `json:"title,omitempty"`
	Content          string                     `json:"content,omitempty"`
	IsPin            int                        `json:"isPin"`
	IsPublic         int                        `json:"isPublic"`
	CommentsDisabled int                        `json:"commentsDisabled"`
	Tags             []string                   `json:"tags,omitempty"`
	Attachments      []string                   `json:"attachments,omitempty"`
	Likes            int                        `json:"likes,omitempty"`
	PulseProfile     *pulseprofile.PulseProfile `json:"pulseProfile,omitempty"`
	Comments         int
}

var pulseFields = map[string]int{
	"ID":  0,
	"Mts": 1,
	// "PLACEHOLDER": 2,
	"UserID": 3,
	// "PLACEHOLDER": 4,
	"Title":   5,
	"Content": 6,
	// "PLACEHOLDER": 7,
	// "PLACEHOLDER": 8,
	"IsPin":            9,
	"IsPublic":         10,
	"CommentsDisabled": 11,
	"Tags":             12,
	"Attachments":      13,
	"Meta":             14,
	"Likes":            15,
	// "PLACEHOLDER": 16,
	// "PLACEHOLDER": 17,
	"PulseProfile": 18,
	"Comments":     19,
	// "PLACEHOLDER": 20,
	// "PLACEHOLDER": 21,
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
	p.CommentsDisabled = convert.ToInt(raw[pulseFields["CommentsDisabled"]])
	p.Likes = convert.ToInt(raw[pulseFields["Likes"]])
	p.Comments = convert.ToInt(raw[pulseFields["Comments"]])

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
