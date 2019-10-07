# rest
--
    import "github.com/bitfinexcom/bitfinex-api-go/v2/rest"


## Usage

```go
const (
	DERIV_TYPE = "deriv"
)
```

#### type BookService

```go
type BookService struct {
	Synchronous
}
```


#### func (*BookService) All

```go
func (b *BookService) All(symbol string, precision bitfinex.BookPrecision, priceLevels int) (*bitfinex.BookUpdateSnapshot, error)
```
Retrieve all books for the given symbol with the given precision at the given
price level see https://docs.bitfinex.com/reference#rest-public-books for more
info

#### type CandleService

```go
type CandleService struct {
	Synchronous
}
```

CandleService manages the Candles endpoint.

#### func (*CandleService) History

```go
func (c *CandleService) History(symbol string, resolution bitfinex.CandleResolution) (*bitfinex.CandleSnapshot, error)
```
Retrieves all candles (Max=1000) with the given symbol and the given candle
resolution See https://docs.bitfinex.com/reference#rest-public-candles for more
info

#### func (*CandleService) HistoryWithQuery

```go
func (c *CandleService) HistoryWithQuery(
	symbol string,
	resolution bitfinex.CandleResolution,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
) (*bitfinex.CandleSnapshot, error)
```
Retrieves all candles (Max=1000) that fit the given query criteria See
https://docs.bitfinex.com/reference#rest-public-candles for more info

#### func (*CandleService) Last

```go
func (c *CandleService) Last(symbol string, resolution bitfinex.CandleResolution) (*bitfinex.Candle, error)
```
Retrieve the last candle for the given symbol with the given resolution See
https://docs.bitfinex.com/reference#rest-public-candles for more info

#### type Client

```go
type Client struct {

	// service providers
	Candles     CandleService
	Orders      OrderService
	Positions   PositionService
	Trades      TradeService
	Tickers     TickerService
	Currencies  CurrenciesService
	Platform    PlatformService
	Book        BookService
	Wallet      WalletService
	Ledgers     LedgerService
	Stats       StatsService
	Status      StatusService
	Derivatives DerivativesService
	Funding     FundingService

	Synchronous
}
```


#### func  NewClient

```go
func NewClient() *Client
```
Create a new Rest client

#### func  NewClientWithHttpDo

```go
func NewClientWithHttpDo(httpDo func(c *http.Client, r *http.Request) (*http.Response, error)) *Client
```
Create a new Rest client with a custom http handler

#### func  NewClientWithSynchronousNonce

```go
func NewClientWithSynchronousNonce(sync Synchronous, nonce utils.NonceGenerator) *Client
```
Create a new Rest client with a synchronous HTTP handler and a custom nonce
generaotr

#### func  NewClientWithSynchronousURLNonce

```go
func NewClientWithSynchronousURLNonce(sync Synchronous, url string, nonce utils.NonceGenerator) *Client
```
Create a new Rest client with a synchronous HTTP handler and a custom base url
and nonce generator

#### func  NewClientWithURL

```go
func NewClientWithURL(url string) *Client
```
Create a new Rest client with a custom base url

#### func  NewClientWithURLHttpDo

```go
func NewClientWithURLHttpDo(base string, httpDo func(c *http.Client, r *http.Request) (*http.Response, error)) *Client
```
Create a new Rest client with a custom base url and HTTP handler

#### func  NewClientWithURLHttpDoNonce

```go
func NewClientWithURLHttpDoNonce(base string, httpDo func(c *http.Client, r *http.Request) (*http.Response, error), nonce utils.NonceGenerator) *Client
```
Create a new Rest client with a custom base url, HTTP handler and none generator

#### func  NewClientWithURLNonce

```go
func NewClientWithURLNonce(url string, nonce utils.NonceGenerator) *Client
```
Create a new Rest client with a custom nonce generator

#### func (*Client) Credentials

```go
func (c *Client) Credentials(key string, secret string) *Client
```
Set the clients credentials in order to make authenticated requests

#### func (*Client) NewAuthenticatedRequest

```go
func (c *Client) NewAuthenticatedRequest(permissionType bitfinex.PermissionType, refURL string) (Request, error)
```
Create a new authenticated GET request with the given permission type and
endpoint url For example permissionType = "r" and refUrl = "/orders" then the
target endpoint will be https://api.bitfinex.com/v2/auth/r/orders/:Symbol

#### func (*Client) NewAuthenticatedRequestWithBytes

```go
func (c *Client) NewAuthenticatedRequestWithBytes(permissionType bitfinex.PermissionType, refURL string, data []byte) (Request, error)
```
Create a new authenticated POST request with the given permission type,endpoint
url and data (bytes) as the body For example permissionType = "r" and refUrl =
"/orders" then the target endpoint will be
https://api.bitfinex.com/v2/auth/r/orders/:Symbol

#### func (*Client) NewAuthenticatedRequestWithData

```go
func (c *Client) NewAuthenticatedRequestWithData(permissionType bitfinex.PermissionType, refURL string, data map[string]interface{}) (Request, error)
```
Create a new authenticated POST request with the given permission type,endpoint
url and data (map[string]interface{}) as the body For example permissionType =
"r" and refUrl = "/orders" then the target endpoint will be
https://api.bitfinex.com/v2/auth/r/orders/:Symbol

#### type CurrenciesService

```go
type CurrenciesService struct {
	Synchronous
}
```

TradeService manages the Trade endpoint.

#### func (*CurrenciesService) Conf

```go
func (cs *CurrenciesService) Conf(label, symbol, unit, explorer, pairs bool) ([]bitfinex.CurrencyConf, error)
```
Retreive currency and symbol service configuration data see
https://docs.bitfinex.com/reference#rest-public-conf for more info

#### type DerivativesService

```go
type DerivativesService struct {
	Synchronous
}
```

OrderService manages data flow for the Order API endpoint

#### type ErrorResponse

```go
type ErrorResponse struct {
	Response *Response
	Message  string `json:"message"`
	Code     int    `json:"code"`
}
```

In case if API will wrong response code ErrorResponse will be returned to caller

#### func (*ErrorResponse) Error

```go
func (r *ErrorResponse) Error() string
```

#### type FundingService

```go
type FundingService struct {
	Synchronous
}
```

LedgerService manages the Ledgers endpoint.

#### func (*FundingService) CancelOffer

```go
func (fs *FundingService) CancelOffer(fc *bitfinex.FundingOfferCancelRequest) (*bitfinex.Notification, error)
```
Submits a request to cancel the given offer see
https://docs.bitfinex.com/reference#cancel-funding-offer for more info

#### func (*FundingService) Credits

```go
func (fs *FundingService) Credits(symbol string) (*bitfinex.FundingCreditSnapshot, error)
```
Retreive all of the active credits used in positions see
https://docs.bitfinex.com/reference#rest-auth-funding-credits for more info

#### func (*FundingService) CreditsHistory

```go
func (fs *FundingService) CreditsHistory(symbol string) (*bitfinex.FundingCreditSnapshot, error)
```
Retreive all of the past in-active credits used in positions see
https://docs.bitfinex.com/reference#rest-auth-funding-credits-hist for more info

#### func (*FundingService) Loans

```go
func (fs *FundingService) Loans(symbol string) (*bitfinex.FundingLoanSnapshot, error)
```
Retreive all of the active funding loans see
https://docs.bitfinex.com/reference#rest-auth-funding-loans for more info

#### func (*FundingService) LoansHistory

```go
func (fs *FundingService) LoansHistory(symbol string) (*bitfinex.FundingLoanSnapshot, error)
```
Retreive all of the past in-active funding loans see
https://docs.bitfinex.com/reference#rest-auth-funding-loans-hist for more info

#### func (*FundingService) OfferHistory

```go
func (fs *FundingService) OfferHistory(symbol string) (*bitfinex.FundingOfferSnapshot, error)
```
Retreive all of the past in-active funding offers see
https://docs.bitfinex.com/reference#rest-auth-funding-offers-hist for more info

#### func (*FundingService) Offers

```go
func (fs *FundingService) Offers(symbol string) (*bitfinex.FundingOfferSnapshot, error)
```
Retreive all of the active fundign offers see
https://docs.bitfinex.com/reference#rest-auth-funding-offers for more info

#### func (*FundingService) SubmitOffer

```go
func (fs *FundingService) SubmitOffer(fo *bitfinex.FundingOfferRequest) (*bitfinex.Notification, error)
```
Submits a request to create a new funding offer see
https://docs.bitfinex.com/reference#submit-funding-offer for more info

#### func (*FundingService) Trades

```go
func (fs *FundingService) Trades(symbol string) (*bitfinex.FundingTradeSnapshot, error)
```
Retreive all of the matched funding trades see
https://docs.bitfinex.com/reference#rest-auth-funding-trades-hist for more info

#### type HttpTransport

```go
type HttpTransport struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}
```


#### func (HttpTransport) Request

```go
func (h HttpTransport) Request(req Request) ([]interface{}, error)
```

#### type LedgerService

```go
type LedgerService struct {
	Synchronous
}
```

LedgerService manages the Ledgers endpoint.

#### func (*LedgerService) Ledgers

```go
func (s *LedgerService) Ledgers(currency string, start int64, end int64, max int32) (*bitfinex.LedgerSnapshot, error)
```
Retrieves all of the past ledger entreies see
https://docs.bitfinex.com/reference#ledgers for more info

#### type OrderService

```go
type OrderService struct {
	Synchronous
}
```

OrderService manages data flow for the Order API endpoint

#### func (*OrderService) All

```go
func (s *OrderService) All() (*bitfinex.OrderSnapshot, error)
```
Retrieves all of the active orders See
https://docs.bitfinex.com/reference#rest-auth-orders for more info

#### func (*OrderService) AllHistory

```go
func (s *OrderService) AllHistory() (*bitfinex.OrderSnapshot, error)
```
Retrieves all past orders See https://docs.bitfinex.com/reference#orders-history
for more info

#### func (*OrderService) GetByOrderId

```go
func (s *OrderService) GetByOrderId(orderID int64) (o *bitfinex.Order, err error)
```
Retrieve an active order by the given ID See
https://docs.bitfinex.com/reference#rest-auth-orders for more info

#### func (*OrderService) GetBySymbol

```go
func (s *OrderService) GetBySymbol(symbol string) (*bitfinex.OrderSnapshot, error)
```
Retrieves all of the active orders with for the given symbol See
https://docs.bitfinex.com/reference#rest-auth-orders for more info

#### func (*OrderService) GetHistoryByOrderId

```go
func (s *OrderService) GetHistoryByOrderId(orderID int64) (o *bitfinex.Order, err error)
```
Retrieve a single order in history with the given id See
https://docs.bitfinex.com/reference#orders-history for more info

#### func (*OrderService) GetHistoryBySymbol

```go
func (s *OrderService) GetHistoryBySymbol(symbol string) (*bitfinex.OrderSnapshot, error)
```
Retrieves all past orders with the given symbol See
https://docs.bitfinex.com/reference#orders-history for more info

#### func (*OrderService) OrderTrades

```go
func (s *OrderService) OrderTrades(symbol string, orderID int64) (*bitfinex.TradeExecutionUpdateSnapshot, error)
```
Retrieves the trades generated by an order See
https://docs.bitfinex.com/reference#orders-history for more info

#### func (*OrderService) SubmitCancelOrder

```go
func (s *OrderService) SubmitCancelOrder(oc *bitfinex.OrderCancelRequest) error
```
Submit a request to cancel an order with the given Id see
https://docs.bitfinex.com/reference#cancel-order for more info

#### func (*OrderService) SubmitOrder

```go
func (s *OrderService) SubmitOrder(order *bitfinex.OrderNewRequest) (*bitfinex.Notification, error)
```
Submit a request to create a new order see
https://docs.bitfinex.com/reference#submit-order for more info

#### func (*OrderService) SubmitUpdateOrder

```go
func (s *OrderService) SubmitUpdateOrder(order *bitfinex.OrderUpdateRequest) (*bitfinex.Notification, error)
```
Submit a request to update an order with the given id with the given changes see
https://docs.bitfinex.com/reference#order-update for more info

#### type PlatformService

```go
type PlatformService struct {
	Synchronous
}
```


#### func (*PlatformService) Status

```go
func (p *PlatformService) Status() (bool, error)
```
Retrieves the current status of the platform see
https://docs.bitfinex.com/reference#rest-public-platform-status for more info

#### type PositionService

```go
type PositionService struct {
	Synchronous
}
```

PositionService manages the Position endpoint.

#### func (*PositionService) All

```go
func (s *PositionService) All() (*bitfinex.PositionSnapshot, error)
```
Retrieves all of the active positions see
https://docs.bitfinex.com/reference#rest-auth-positions for more info

#### func (*PositionService) Claim

```go
func (s *PositionService) Claim(cp *bitfinex.ClaimPositionRequest) (*bitfinex.Notification, error)
```
Submits a request to claim an active position with the given id see
https://docs.bitfinex.com/reference#claim-position for more info

#### type Request

```go
type Request struct {
	RefURL  string     // ref url
	Data    []byte     // body data
	Method  string     // http method
	Params  url.Values // query parameters
	Headers map[string]string
}
```

Request is a wrapper for standard http.Request. Default method is POST with no
data.

#### func  NewRequest

```go
func NewRequest(refURL string) Request
```
Create new POST request with an empty body as payload

#### func  NewRequestWithBytes

```go
func NewRequestWithBytes(refURL string, data []byte) Request
```
Create a new POST request with the given bytes as body

#### func  NewRequestWithData

```go
func NewRequestWithData(refURL string, data map[string]interface{}) (Request, error)
```
Create a new POST request with the given data (map[string]interface{}) as body

#### func  NewRequestWithDataMethod

```go
func NewRequestWithDataMethod(refURL string, data []byte, method string) Request
```
Create a new request with a given method (POST | GET) with bytes as body

#### func  NewRequestWithMethod

```go
func NewRequestWithMethod(refURL string, method string) Request
```
Create a new request with the given method (POST | GET)

#### type Response

```go
type Response struct {
	Response *http.Response
	Body     []byte
}
```

Response is a wrapper for standard http.Response and provides more methods.

#### func (*Response) String

```go
func (r *Response) String() string
```
String converts response body to string. An empty string will be returned if
error.

#### type StatsService

```go
type StatsService struct {
	Synchronous
}
```

TradeService manages the Trade endpoint.

#### func (*StatsService) CreditSizeHistory

```go
func (ss *StatsService) CreditSizeHistory(symbol string, side bitfinex.OrderSide) ([]bitfinex.Stat, error)
```
Retrieves platform statistics for credit size history see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) CreditSizeLast

```go
func (ss *StatsService) CreditSizeLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error)
```
Retrieves platform statistics for credit size last see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) FundingHistory

```go
func (ss *StatsService) FundingHistory(symbol string) ([]bitfinex.Stat, error)
```
Retrieves platform statistics for funding history see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) FundingLast

```go
func (ss *StatsService) FundingLast(symbol string) (*bitfinex.Stat, error)
```
Retrieves platform statistics for funding last see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) PositionHistory

```go
func (ss *StatsService) PositionHistory(symbol string, side bitfinex.OrderSide) ([]bitfinex.Stat, error)
```
Retrieves platform statistics for position history see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) PositionLast

```go
func (ss *StatsService) PositionLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error)
```
Retrieves platform statistics for position last see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) SymbolCreditSizeHistory

```go
func (ss *StatsService) SymbolCreditSizeHistory(fundingSymbol string, tradingSymbol string) ([]bitfinex.Stat, error)
```
Retrieves platform statistics for credit size history see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### func (*StatsService) SymbolCreditSizeLast

```go
func (ss *StatsService) SymbolCreditSizeLast(fundingSymbol string, tradingSymbol string) (*bitfinex.Stat, error)
```
Retrieves platform statistics for credit size last see
https://docs.bitfinex.com/reference#rest-public-stats for more info

#### type StatusService

```go
type StatusService struct {
	Synchronous
}
```

TradeService manages the Trade endpoint.

#### func (*StatusService) DerivativeStatus

```go
func (ss *StatusService) DerivativeStatus(symbol string) (*bitfinex.DerivativeStatus, error)
```
Retrieves derivative status information for the given symbol from the platform
see https://docs.bitfinex.com/reference#rest-public-status for more info

#### func (*StatusService) DerivativeStatusAll

```go
func (ss *StatusService) DerivativeStatusAll() ([]*bitfinex.DerivativeStatus, error)
```
Retrieves derivative status information for all symbols from the platform see
https://docs.bitfinex.com/reference#rest-public-status for more info

#### func (*StatusService) DerivativeStatusMulti

```go
func (ss *StatusService) DerivativeStatusMulti(symbols []string) ([]*bitfinex.DerivativeStatus, error)
```
Retrieves derivative status information for the given symbols from the platform
see https://docs.bitfinex.com/reference#rest-public-status for more info

#### type Synchronous

```go
type Synchronous interface {
	Request(request Request) ([]interface{}, error)
}
```


#### type TickerService

```go
type TickerService struct {
	Synchronous
}
```

TradeService manages the Trade endpoint.

#### func (*TickerService) All

```go
func (s *TickerService) All() (*[]bitfinex.Ticker, error)
```
Retrieves all tickers for all symbols see
https://docs.bitfinex.com/reference#rest-public-ticker for more info

#### func (*TickerService) Get

```go
func (s *TickerService) Get(symbol string) (*bitfinex.Ticker, error)
```
Retrieves the ticker for the given symbol see
https://docs.bitfinex.com/reference#rest-public-ticker for more info

#### func (*TickerService) GetMulti

```go
func (s *TickerService) GetMulti(symbols []string) (*[]bitfinex.Ticker, error)
```
Retrieves the tickers for the given symbols see
https://docs.bitfinex.com/reference#rest-public-ticker for more info

#### type TradeService

```go
type TradeService struct {
	Synchronous
}
```

TradeService manages the Trade endpoint.

#### func (*TradeService) AccountAll

```go
func (s *TradeService) AccountAll() (*bitfinex.TradeExecutionUpdateSnapshot, error)
```
Retrieves all matched trades for the account see
https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info

#### func (*TradeService) AccountAllWithSymbol

```go
func (s *TradeService) AccountAllWithSymbol(symbol string) (*bitfinex.TradeExecutionUpdateSnapshot, error)
```
Retrieves all matched trades with the given symbol for the account see
https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info

#### func (*TradeService) AccountHistoryWithQuery

```go
func (s *TradeService) AccountHistoryWithQuery(
	symbol string,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
) (*bitfinex.TradeExecutionUpdateSnapshot, error)
```
Queries all matched trades with group of optional parameters see
https://docs.bitfinex.com/reference#rest-auth-trades-hist for more info

#### func (*TradeService) PublicHistoryWithQuery

```go
func (s *TradeService) PublicHistoryWithQuery(
	symbol string,
	start bitfinex.Mts,
	end bitfinex.Mts,
	limit bitfinex.QueryLimit,
	sort bitfinex.SortOrder,
) (*bitfinex.TradeSnapshot, error)
```
Queries all public trades with a group of optional paramters see
https://docs.bitfinex.com/reference#rest-public-trades for more info

#### type WalletService

```go
type WalletService struct {
	Synchronous
}
```

WalletService manages data flow for the Wallet API endpoint

#### func (*WalletService) CreateDepositAddress

```go
func (ws *WalletService) CreateDepositAddress(wallet, method string) (*bitfinex.Notification, error)
```
Submits a request to create a new deposit address for the give Bitfinex wallet.
Old addresses are still valid. See
https://docs.bitfinex.com/reference#deposit-address for more info

#### func (*WalletService) DepositAddress

```go
func (ws *WalletService) DepositAddress(wallet, method string) (*bitfinex.Notification, error)
```
Retrieves the deposit address for the given Bitfinex wallet see
https://docs.bitfinex.com/reference#deposit-address for more info

#### func (*WalletService) SetCollateral

```go
func (s *WalletService) SetCollateral(symbol string, amount float64) (bool, error)
```
Update the amount of collateral for a Derivative position see
https://docs.bitfinex.com/reference#rest-auth-deriv-pos-collateral-set for more
info

#### func (*WalletService) Transfer

```go
func (ws *WalletService) Transfer(from, to, currency, currencyTo string, amount float64) (*bitfinex.Notification, error)
```
Submits a request to transfer funds from one Bitfinex wallet to another see
https://docs.bitfinex.com/reference#transfer-between-wallets for more info

#### func (*WalletService) Wallet

```go
func (s *WalletService) Wallet() (*bitfinex.WalletSnapshot, error)
```
Retrieves all of the wallets for the account see
https://docs.bitfinex.com/reference#rest-auth-wallets for more info

#### func (*WalletService) Withdraw

```go
func (ws *WalletService) Withdraw(wallet, method string, amount float64, address string) (*bitfinex.Notification, error)
```
Submits a request to withdraw funds from the given Bitfinex wallet to the given
address See https://docs.bitfinex.com/reference#withdraw for more info
