// Package bitfinex is the official client to access to bitfinex.com API
package bitfinex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bitfinexcom/bitfinex-api-go/utils"
)

const (
	// BaseURL is the v2 REST endpoint.
	BaseURL = "https://api.bitfinex.com/v2"

	// WebSocketURL is the v2 Websocket endpoint.
	WebSocketURL = "wss://api.bitfinex.com/ws/2"
)

var nonce int64

type Client struct {
	// Base URL for API requests.
	BaseURL *url.URL

	HTTPClient *http.Client

	// Auth data
	APIKey    string
	APISecret string

	// Services
	Websocket *bfxWebsocket
}

func NewClient() *Client {
	return NewClientWithHTTP(http.DefaultClient)
}

func NewClientWithHTTP(h *http.Client) *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{BaseURL: baseURL, HTTPClient: h}

	c.Websocket = newBfxWebsocket(c, WebSocketURL)

	return c
}

// NewRequest create new API request. Relative url can be provided in refURL.
func (c *Client) newRequest(method string, refURL string, params url.Values) (*http.Request, error) {
	rel, err := url.Parse(refURL)
	if err != nil {
		return nil, err
	}
	if params != nil {
		rel.RawQuery = params.Encode()
	}

	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), nil)

	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAuthenticatedRequest creates new http request for authenticated routes
func (c *Client) newAuthenticatedRequest(m string, refURL string, data map[string]interface{}) (*http.Request, error) {
	req, err := c.newRequest(m, refURL, nil)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"request": "/v2/" + refURL,
		"nonce":   utils.GetNonce(),
	}

	for k, v := range data {
		payload[k] = v
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	payloadEnc := base64.StdEncoding.EncodeToString(payloadJSON)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-BFX-APIKEY", c.APIKey)
	req.Header.Add("X-BFX-PAYLOAD", payloadEnc)
	req.Header.Add("X-BFX-SIGNATURE", c.signPayload(payloadEnc))

	return req, nil
}

func (c *Client) signPayload(payload string) string {
	sig := hmac.New(sha512.New384, []byte(c.APISecret))
	sig.Write([]byte(payload))
	return hex.EncodeToString(sig.Sum(nil))
}

// Credentials sets api key and secret for usage is requests that
// requires authentication
func (c *Client) Credentials(key string, secret string) *Client {
	c.APIKey = key
	c.APISecret = secret

	return c
}

var httpDo = func(c *http.Client, req *http.Request) (*http.Response, error) {
	return c.Do(req)
}

// Response is wrapper for standard http.Response and provides
// more methods.
type Response struct {
	Response *http.Response
	Body     []byte
}

// newResponse creates new wrapper.
func newResponse(r *http.Response) *Response {
	// Use a LimitReader of arbitrary size (here ~8.39MB) to prevent us from
	// reading overly large response bodies.
	lr := io.LimitReader(r.Body, 8388608)
	body, err := ioutil.ReadAll(lr)
	if err != nil {
		body = []byte(`Error reading body:` + err.Error())
	}

	return &Response{r, body}
}

// String converts response body to string.
// An empty string will be returned if error.
func (r *Response) String() string {
	return string(r.Body)
}

// checkResponse checks response status code and response
// for errors.
func checkResponse(r *Response) error {
	if c := r.Response.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Try to decode error message
	errorResponse := &ErrorResponse{Response: r}
	err := json.Unmarshal(r.Body, errorResponse)
	if err != nil {
		errorResponse.Message = "Error decoding response error message. " +
			"Please see response body for more information."
	}

	return errorResponse
}

// In case if API will wrong response code
// ErrorResponse will be returned to caller
type ErrorResponse struct {
	Response *Response
	Message  string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Response.Request.Method,
		r.Response.Response.Request.URL,
		r.Response.Response.StatusCode,
		r.Message,
	)
}

// Do executes API request created by NewRequest method or custom *http.Request.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := httpDo(c.HTTPClient, req)
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
