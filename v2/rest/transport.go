package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		req.Data = map[string]interface{}{}
	}

	b, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)

	u := h.BaseURL.ResolveReference(rel)
	httpReq, err := http.NewRequest(req.Method, u.String(), body)

	if err != nil {
		return nil, err
	}

	resp, err := h.do(httpReq, &raw)
	if err != nil {
		return nil, fmt.Errorf("could not parse response: %s", resp.Response.Status)
	}

	return raw, nil
}

// Do executes API request created by NewRequest method or custom *http.Request.
func (h HttpTransport) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := h.httpDo(h.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	err = checkResponse(response)
	if err != nil {
		return response, err
	}

	if v != nil {
		err = json.Unmarshal(response.Body, v)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}
