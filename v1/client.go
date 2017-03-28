// Package bitfinex is the official client to access to bitfinex.com API
package bitfinex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bitfinexcom/bitfinex-api-go/utils"
)

const (
	// BaseURL is the v1 REST endpoint.
	BaseURL = "https://api.bitfinex.com/v1/"
	// WebSocketURL the v1 Websocket endpoint.
	WebSocketURL = "wss://api.bitfinex.com/ws/"
)

// Client manages all the communication with the Bitfinex API.
type Client struct {
	// Base URL for API requests.
	BaseURL                *url.URL
	WebSocketURL           string
	WebSocketTLSSkipVerify bool

	// Auth data
	APIKey    string
	APISecret string

	// Services
	Pairs         *PairsService
	Stats         *StatsService
	Ticker        *TickerService
	Account       *AccountService
	Balances      *BalancesService
	Offers        *OffersService
	Credits       *CreditsService
	Deposit       *DepositService
	Lendbook      *LendbookService
	MarginInfo    *MarginInfoService
	MarginFunding *MarginFundingService
	OrderBook     *OrderBookService
	Orders        *OrderService
	Trades        *TradesService
	Positions     *PositionsService
	History       *HistoryService
	WebSocket     *WebSocketService
	Wallet        *WalletService
}

// NewClient creates new Bitfinex.com API client.
func NewClient() *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{BaseURL: baseURL, WebSocketURL: WebSocketURL}
	c.Pairs = &PairsService{client: c}
	c.Stats = &StatsService{client: c}
	c.Account = &AccountService{client: c}
	c.Ticker = &TickerService{client: c}
	c.Balances = &BalancesService{client: c}
	c.Offers = &OffersService{client: c}
	c.Credits = &CreditsService{client: c}
	c.Deposit = &DepositService{client: c}
	c.Lendbook = &LendbookService{client: c}
	c.MarginInfo = &MarginInfoService{client: c}
	c.MarginFunding = &MarginFundingService{client: c}
	c.OrderBook = &OrderBookService{client: c}
	c.Orders = &OrderService{client: c}
	c.History = &HistoryService{client: c}
	c.Trades = &TradesService{client: c}
	c.Positions = &PositionsService{client: c}
	c.Wallet = &WalletService{client: c}
	c.WebSocket = NewWebSocketService(c)
	c.WebSocketTLSSkipVerify = false

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
	var req *http.Request
	u := c.BaseURL.ResolveReference(rel)
	req, err = http.NewRequest(method, u.String(), nil)

	if err != nil {
		return nil, err
	}

	return req, nil
}

// newAuthenticatedRequest creates new http request for authenticated routes.
func (c *Client) newAuthenticatedRequest(m string, refURL string, data map[string]interface{}) (*http.Request, error) {
	req, err := c.newRequest(m, refURL, nil)
	if err != nil {
		return nil, err
	}

	nonce := utils.GetNonce()
	payload := map[string]interface{}{
		"request": "/v1/" + refURL,
		"nonce":   nonce,
	}

	for k, v := range data {
		payload[k] = v
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(p)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-BFX-APIKEY", c.APIKey)
	req.Header.Add("X-BFX-PAYLOAD", encoded)
	req.Header.Add("X-BFX-SIGNATURE", c.signPayload(encoded))

	return req, nil
}

func (c *Client) signPayload(payload string) string {
	sig := hmac.New(sha512.New384, []byte(c.APISecret))
	sig.Write([]byte(payload))
	return hex.EncodeToString(sig.Sum(nil))
}

// Auth sets api key and secret for usage is requests that requires authentication.
func (c *Client) Auth(key string, secret string) *Client {
	c.APIKey = key
	c.APISecret = secret

	return c
}

var httpDo = func(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// Do executes API request created by NewRequest method or custom *http.Request.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := httpDo(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = checkResponse(response)
	if err != nil {
		// Return response in case caller need to debug it.
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

// Response is wrapper for standard http.Response and provides
// more methods.
type Response struct {
	Response *http.Response
	Body     []byte
}

// newResponse creates new wrapper.
func newResponse(r *http.Response) *Response {
	body, err := ioutil.ReadAll(r.Body)
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

// ErrorResponse is the custom error type that is returned if the API returns an
// error.
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
