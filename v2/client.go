// Package bitfinex is the official client to access to bitfinex.com API
package bitfinex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
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
	BaseURL = "https://api.bitfinex.com/v2/"

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
	Orders    *OrderService
	Platform  *PlatformService
	Positions *PositionService
	Trades    *TradeService
}

func NewClient() *Client {
	return NewClientWithHTTP(http.DefaultClient)
}

func NewClientWithHTTP(h *http.Client) *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{BaseURL: baseURL, HTTPClient: h}

	c.Websocket = newBfxWebsocket(c, WebSocketURL)
	c.Orders = &OrderService{client: c}
	c.Platform = &PlatformService{client: c}
	c.Positions = &PositionService{client: c}
	c.Trades = &TradeService{client: c}

	return c
}

// NewRequest create new API request. Relative url can be provided in refURL.
func (c *Client) newRequest(method string, refURL string, params url.Values, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(refURL)
	if err != nil {
		return nil, err
	}
	if params != nil {
		rel.RawQuery = params.Encode()
	}

	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), body)

	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAuthenticatedRequest creates new http request for authenticated routes
func (c *Client) newAuthenticatedRequest(m string, refURL string, data map[string]interface{}) (*http.Request, error) {
	refURL = "auth/r/" + refURL

	if data == nil {
		data = map[string]interface{}{}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	rd := bytes.NewReader(b)
	req, err := c.newRequest(m, refURL, nil, rd)
	if err != nil {
		return nil, err
	}

	nonce := utils.GetNonce()
	message := "/api/v2/" + refURL + nonce + string(b)
	sig := c.sign(message)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("bfx-nonce", nonce)
	req.Header.Add("bfx-signature", sig)
	req.Header.Add("bfx-apikey", c.APIKey)

	return req, nil
}

func (c *Client) sign(msg string) string {
	sig := hmac.New(sha512.New384, []byte(c.APISecret))
	sig.Write([]byte(msg))
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

	var raw []interface{}
	// Try to decode error message
	errorResponse := &ErrorResponse{Response: r}
	err := json.Unmarshal(r.Body, &raw)
	if err != nil {
		errorResponse.Message = "Error decoding response error message. " +
			"Please see response body for more information."
		return errorResponse
	}

	if len(raw) < 3 {
		errorResponse.Message = fmt.Sprintf("Expected response to have three elements but got %#v", raw)
		return errorResponse
	}

	if str, ok := raw[0].(string); !ok || str != "error" {
		errorResponse.Message = fmt.Sprintf("Expected first element to be \"error\" but got %#v", raw)
		return errorResponse
	}

	code, ok := raw[1].(float64)
	if !ok {
		errorResponse.Message = fmt.Sprintf("Expected second element to be error code but got %#v", raw)
		return errorResponse
	}
	errorResponse.Code = int(code)

	msg, ok := raw[2].(string)
	if !ok {
		errorResponse.Message = fmt.Sprintf("Expected third element to be error message but got %#v", raw)
		return errorResponse
	}
	errorResponse.Message = msg

	return errorResponse
}

// In case if API will wrong response code
// ErrorResponse will be returned to caller
type ErrorResponse struct {
	Response *Response
	Message  string `json:"message"`
	Code     int    `json:"code"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v (%d)",
		r.Response.Response.Request.Method,
		r.Response.Response.Request.URL,
		r.Response.Response.StatusCode,
		r.Message,
		r.Code,
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
