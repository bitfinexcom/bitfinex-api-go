package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
)

type PulseService struct {
	requestFactory
	Synchronous
}

// PublicPulseProfile get pulse profile by nickname
func (ps *PulseService) PublicPulseProfile(nickname string) (*pulseprofile.PulseProfile, error) {
	if (len(nickname)) == 0 {
		return nil, fmt.Errorf("nickname is required argument")
	}

	req := NewRequestWithMethod(path.Join("pulse", "profile", nickname), "GET")
	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pp, err := pulseprofile.NewFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

// PublicPulseHistory returns public pulse messages
func (ps *PulseService) PublicPulseHistory(limit, end string) ([]*pulse.Pulse, error) {
	req := NewRequestWithMethod(path.Join("pulse", "hist"), "GET")
	req.Params = make(url.Values)
	req.Params.Add("limit", limit)
	req.Params.Add("end", end)

	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pph, err := pulse.NewFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pph, nil
}

// AddPulse submits pulse messages
func (ps *PulseService) AddPulse(p *pulse.Pulse) (*pulse.Pulse, error) {
	tl := len(p.Title)
	if tl < 16 || tl > 120 {
		return nil, fmt.Errorf("Title length min 16 and max 120 characters. Got:%d", tl)
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := ps.requestFactory.NewAuthenticatedRequestWithBytes(bitfinex.PermissionWrite, path.Join("pulse", "add"), payload)
	if err != nil {
		return nil, err
	}

	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pm, err := pulse.NewSingleFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pm, nil
}
