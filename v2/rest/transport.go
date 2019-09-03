package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type HttpTransport struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	httpDo     func(c *http.Client, req *http.Request) (*http.Response, error)
}

func (h HttpTransport) Request(req Request) ([]interface{}, error) {
	var raw []interface{}

	rel, err := url.Parse(req.RefURL)
	if err != nil {
		return nil, err
	}
	if req.Params != nil {
		rel.RawQuery = req.Params.Encode()
	}
	if req.Data == nil {
		req.Data = []byte("{}")
	}
	body := bytes.NewReader(req.Data)

	u := h.BaseURL.ResolveReference(rel)
	httpReq, err := http.NewRequest(req.Method, u.String(), body)
	for k, v := range req.Headers {
		httpReq.Header.Add(k, v)
	}
	if err != nil {
		return nil, err
	}
	err = h.do(httpReq, &raw)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// Do executes API request created by NewRequest method or custom *http.Request.
func (h HttpTransport) do(req *http.Request, v interface{}) (error) {
	resp, err := h.httpDo(h.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	err = checkResponse(response)
	if err != nil {
		return err
	}

	if v != nil {
		err = json.Unmarshal(response.Body, v)
		if err != nil {
			return err
		}
	}

	return nil
}
