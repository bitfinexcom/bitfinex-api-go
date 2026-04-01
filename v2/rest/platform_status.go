package rest

type PlatformService struct {
	Synchronous
}

// Retrieves the current status of the platform
// see https://docs.bitfinex.com/reference#rest-public-platform-status for more info
func (p *PlatformService) Status() (bool, error) {
	raw, err := p.Request(NewRequestWithMethod("platform/status", "GET"))
	if err != nil {
		return false, err
	}
	if len(raw) == 0 {
		return false, nil
	}
	if val, ok := raw[0].(float64); ok && int(val) == 1 {
		return true, nil
	}
	return false, nil
}
