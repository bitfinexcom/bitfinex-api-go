package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulseprofile"
)

type PulseService struct {
	requestFactory
	Synchronous
}

type Nickname string

// PublicPulseProfile returns details for a specific Pulse profile
// https://docs.bitfinex.com/reference#rest-public-pulse-profile
func (ps *PulseService) PublicPulseProfile(nickname Nickname) (*pulseprofile.PulseProfile, error) {
	if (len(nickname)) == 0 {
		return nil, fmt.Errorf("nickname is required argument")
	}

	req := NewRequestWithMethod(path.Join("pulse", "profile", string(nickname)), "GET")
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

// PublicPulseHistory returns latest pulse messages. You can specify
// an end timestamp to view older messages.
// see https://docs.bitfinex.com/reference#rest-public-pulse-hist
func (ps *PulseService) PublicPulseHistory(limit int, end common.Mts) ([]*pulse.Pulse, error) {
	req := NewRequestWithMethod(path.Join("pulse", "hist"), "GET")
	req.Params = make(url.Values)
	req.Params.Add("limit", strconv.Itoa(limit))
	req.Params.Add("end", strconv.FormatInt(int64(end), 10))

	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pph, err := pulse.SnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pph, nil
}

// AddPulse submits pulse message
// see https://docs.bitfinex.com/reference#rest-auth-pulse-add
func (ps *PulseService) AddPulse(p *pulse.Pulse) (*pulse.Pulse, error) {
	tl := len(p.Title)
	if tl < 16 || tl > 120 {
		return nil, fmt.Errorf("Title length min 16 and max 120 characters. Got:%d", tl)
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := ps.requestFactory.NewAuthenticatedRequestWithBytes(common.PermissionWrite, path.Join("pulse", "add"), payload)
	if err != nil {
		return nil, err
	}

	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pm, err := pulse.FromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

// AddComment submits pulse comment
// see https://docs.bitfinex.com/reference#rest-auth-pulse-add
func (ps *PulseService) AddComment(p *pulse.Pulse) (*pulse.Pulse, error) {
	if len(p.Parent) == 0 {
		return nil, fmt.Errorf("Pulse comment requires `Parent` parameter to be set")
	}

	return ps.AddPulse(p)
}

// PulseHistory allows you to retrieve your pulse history.
// see https://docs.bitfinex.com/reference#rest-auth-pulse-hist
func (ps *PulseService) PulseHistory() ([]*pulse.Pulse, error) {
	req, err := ps.NewAuthenticatedRequest(common.PermissionRead, path.Join("pulse", "hist"))
	if err != nil {
		return nil, err
	}

	raw, err := ps.Request(req)
	if err != nil {
		return nil, err
	}

	pph, err := pulse.SnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}

	return pph, nil
}

// DeletePulse removes your pulse message. Returns 0 if no pulse was deleted and 1 if it was
// see https://docs.bitfinex.com/reference#rest-auth-pulse-del
func (ps *PulseService) DeletePulse(pid string) (int, error) {
	payload := map[string]interface{}{"pid": pid}

	req, err := ps.NewAuthenticatedRequestWithData(common.PermissionWrite, path.Join("pulse", "del"), payload)
	if err != nil {
		return 0, err
	}

	raw, err := ps.Request(req)
	if err != nil {
		return 0, err
	}

	return convert.ToInt(raw[0]), nil
}
