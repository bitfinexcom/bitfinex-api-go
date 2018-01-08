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
/*
	// raw is an interface type, but we only care about len & index 0
	s := make([]int, len(raw))
	for i, v := range raw {
		s[i] = v.(int)
	}
*/
	return len(raw) > 0 && raw[0].(int) == 1, nil
}
