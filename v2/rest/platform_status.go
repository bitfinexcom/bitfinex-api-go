package rest

type PlatformService struct {
	Synchronous
}

// Status indicates whether the platform is currently operative or not.
func (p *PlatformService) Status() (bool, error) {
	raw, err := p.Request(NewRequestWithMethod("platform/status", "GET"))
	if err != nil {
		return false, err
	}
	return len(raw) > 0 && raw[0].(float64) == 1, nil
}
