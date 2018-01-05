package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var productionBaseURL = "https://api.bitfinex.com/v2/"

type Synchronous interface {
	Request(request Request) ([]interface{}, error)
}

type Client struct {
	// base members for synchronous API
	apiKey    string
	apiSecret string

	// service providers
	Orders    OrderService
	Positions PositionService
	Trades    TradeService
	Platform  PlatformService

	Synchronous
}

func NewClient() *Client {
	httpDo := func(c *http.Client, req *http.Request) (*http.Response, error) {
		return c.Do(req)
	}
	return NewClientWithHttpDo(httpDo)
}

func NewClientWithHttpDo(httpDo func(c *http.Client, r *http.Request) (*http.Response, error)) *Client {
	url, _ := url.Parse(productionBaseURL)
	sync := &HttpTransport{
		BaseURL:    url,
		httpDo:     httpDo,
		HTTPClient: http.DefaultClient,
	}
	return NewClientWithSynchronous(sync)
}

// mock me
func NewClientWithSynchronous(sync Synchronous) *Client {
	c := &Client{
		Synchronous: sync,
	}
	c.Orders = OrderService{Synchronous: c}
	return c
}

func (c Client) Credentials(key string, secret string) *Client {
	c.apiKey = key
	c.apiSecret = secret
	return &c
}

// Request is a wrapper for standard http.Request.  Default method is POST with no data.
type Request struct {
	RefURL string                 // ref url
	Data   map[string]interface{} // body data
	Method string                 // http method
	Params url.Values             // query parameters
}

// Response is a wrapper for standard http.Response and provides more methods.
type Response struct {
	Response *http.Response
	Body     []byte
}

func NewRequest(refURL string) Request {
	return NewRequestWithDataMethod(refURL, nil, "POST")
}

func NewRequestWithMethod(refURL string, method string) Request {
	return NewRequestWithDataMethod(refURL, nil, method)
}

func NewRequestWithData(refURL string, data map[string]interface{}) Request {
	return NewRequestWithDataMethod(refURL, data, "POST")
}

func NewRequestWithDataMethod(refURL string, data map[string]interface{}, method string) Request {
	return Request{
		RefURL: refURL,
		Data:   data,
		Method: method,
	}
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
