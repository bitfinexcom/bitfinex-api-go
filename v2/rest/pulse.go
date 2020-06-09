package rest

import (
	"fmt"
	"path"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

type PulseService struct {
	requestFactory
	Synchronous
}

func (ss *PulseService) PublicPulseProfile(nickname string) (*pulseprofile.PulseProfile, error) {
	if (len(nickname)) == 0 {
		return nil, fmt.Errorf("nickname is required argument")
	}

	req := NewRequestWithMethod(path.Join("pulse", "profile", nickname), "GET")
	raw, err := ss.Request(req)
	if err != nil {
		return nil, err
	}

	pp, err := pulseprofile.NewFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pp, nil
}
