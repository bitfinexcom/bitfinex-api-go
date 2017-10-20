package bitfinex

type PlatformService struct {
	client *Client
}

// Status indicates whether the platform is currently operative or not.
func (p *PlatformService) Status() (bool, error) {
	req, err := p.client.newRequest("GET", "platform/status", nil, nil)
	if err != nil {
		return false, err
	}

	var s []int
	_, err = p.client.do(req, &s)
	if err != nil {
		return false, err
	}

	return len(s) > 0 && s[0] == 1, nil
}
