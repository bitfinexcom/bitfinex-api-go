package rest

import (
	"fmt"
	"net/url"
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

type PulseService struct {
	requestFactory
	Synchronous
}

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
