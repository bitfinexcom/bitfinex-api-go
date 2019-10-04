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
Return Candles for the public account.

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
Return Candles for the public account.

#### func (*CandleService) Last

```go
func (c *CandleService) Last(symbol string, resolution bitfinex.CandleResolution) (*bitfinex.Candle, error)
```
Return Candles for the public account.

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

#### func  NewClientWithHttpDo

```go
func NewClientWithHttpDo(httpDo func(c *http.Client, r *http.Request) (*http.Response, error)) *Client
```

#### func  NewClientWithSynchronousNonce

```go
func NewClientWithSynchronousNonce(sync Synchronous, nonce utils.NonceGenerator) *Client
```

#### func  NewClientWithSynchronousURLNonce

```go
func NewClientWithSynchronousURLNonce(sync Synchronous, url string, nonce utils.NonceGenerator) *Client
```
mock me in tests

#### func  NewClientWithURL

```go
func NewClientWithURL(url string) *Client
```

#### func  NewClientWithURLHttpDo

```go
func NewClientWithURLHttpDo(base string, httpDo func(c *http.Client, r *http.Request) (*http.Response, error)) *Client
```

#### func  NewClientWithURLHttpDoNonce

```go
func NewClientWithURLHttpDoNonce(base string, httpDo func(c *http.Client, r *http.Request) (*http.Response, error), nonce utils.NonceGenerator) *Client
```

#### func  NewClientWithURLNonce

```go
func NewClientWithURLNonce(url string, nonce utils.NonceGenerator) *Client
```

#### func (*Client) Credentials

```go
func (c *Client) Credentials(key string, secret string) *Client
```

#### func (*Client) NewAuthenticatedRequest

```go
func (c *Client) NewAuthenticatedRequest(permissionType bitfinex.PermissionType, refURL string) (Request, error)
```

#### func (*Client) NewAuthenticatedRequestWithBytes

```go
func (c *Client) NewAuthenticatedRequestWithBytes(permissionType bitfinex.PermissionType, refURL string, data []byte) (Request, error)
```

#### func (*Client) NewAuthenticatedRequestWithData

```go
func (c *Client) NewAuthenticatedRequestWithData(permissionType bitfinex.PermissionType, refURL string, data map[string]interface{}) (Request, error)
```

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
All returns all orders for the authenticated account.

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

#### func (*FundingService) Credits

```go
func (fs *FundingService) Credits(symbol string) (*bitfinex.FundingCreditSnapshot, error)
```

#### func (*FundingService) CreditsHistory

```go
func (fs *FundingService) CreditsHistory(symbol string) (*bitfinex.FundingCreditSnapshot, error)
```

#### func (*FundingService) Loans

```go
func (fs *FundingService) Loans(symbol string) (*bitfinex.FundingLoanSnapshot, error)
```

#### func (*FundingService) LoansHistory

```go
func (fs *FundingService) LoansHistory(symbol string) (*bitfinex.FundingLoanSnapshot, error)
```

#### func (*FundingService) OfferHistory

```go
func (fs *FundingService) OfferHistory(symbol string) (*bitfinex.FundingOfferSnapshot, error)
```

#### func (*FundingService) Offers

```go
func (fs *FundingService) Offers(symbol string) (*bitfinex.FundingOfferSnapshot, error)
```
All returns all ledgers for the authenticated account

#### func (*FundingService) SubmitOffer

```go
func (fs *FundingService) SubmitOffer(fo *bitfinex.FundingOfferRequest) (*bitfinex.Notification, error)
```

#### func (*FundingService) Trades

```go
func (fs *FundingService) Trades(symbol string) (*bitfinex.FundingTradeSnapshot, error)
```

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
All returns all ledgers for the authenticated account.

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
Get all active orders

#### func (*OrderService) AllHistory

```go
func (s *OrderService) AllHistory() (*bitfinex.OrderSnapshot, error)
```
Get all historical orders

#### func (*OrderService) GetByOrderId

```go
func (s *OrderService) GetByOrderId(orderID int64) (o *bitfinex.Order, err error)
```
Get an active order using its order id

#### func (*OrderService) GetBySymbol

```go
func (s *OrderService) GetBySymbol(symbol string) (*bitfinex.OrderSnapshot, error)
```
Get all active orders with the given symbol

#### func (*OrderService) GetHistoryByOrderId

```go
func (s *OrderService) GetHistoryByOrderId(orderID int64) (o *bitfinex.Order, err error)
```
Get a historical order using its order id

#### func (*OrderService) GetHistoryBySymbol

```go
func (s *OrderService) GetHistoryBySymbol(symbol string) (*bitfinex.OrderSnapshot, error)
```
Get all historical orders with the given symbol

#### func (*OrderService) OrderTrades

```go
func (s *OrderService) OrderTrades(symbol string, orderID int64) (*bitfinex.TradeExecutionUpdateSnapshot, error)
```
OrderTrades returns a set of executed trades related to an order.

#### func (*OrderService) SubmitCancelOrder

```go
func (s *OrderService) SubmitCancelOrder(oc *bitfinex.OrderCancelRequest) error
```

#### func (*OrderService) SubmitOrder

```go
func (s *OrderService) SubmitOrder(order *bitfinex.OrderNewRequest) (*bitfinex.Notification, error)
```

#### func (*OrderService) SubmitUpdateOrder

```go
func (s *OrderService) SubmitUpdateOrder(order *bitfinex.OrderUpdateRequest) (*bitfinex.Notification, error)
```

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
Status indicates whether the platform is currently operative or not.

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
All returns all positions for the authenticated account.

#### func (*PositionService) Claim

```go
func (s *PositionService) Claim(cp *bitfinex.ClaimPositionRequest) (*bitfinex.Notification, error)
```

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

#### func  NewRequestWithBytes

```go
func NewRequestWithBytes(refURL string, data []byte) Request
```

#### func  NewRequestWithData

```go
func NewRequestWithData(refURL string, data map[string]interface{}) (Request, error)
```

#### func  NewRequestWithDataMethod

```go
func NewRequestWithDataMethod(refURL string, data []byte, method string) Request
```

#### func  NewRequestWithMethod

```go
func NewRequestWithMethod(refURL string, method string) Request
```

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

#### func (*StatsService) CreditSizeLast

```go
func (ss *StatsService) CreditSizeLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error)
```

#### func (*StatsService) FundingHistory

```go
func (ss *StatsService) FundingHistory(symbol string) ([]bitfinex.Stat, error)
```

#### func (*StatsService) FundingLast

```go
func (ss *StatsService) FundingLast(symbol string) (*bitfinex.Stat, error)
```

#### func (*StatsService) PositionHistory

```go
func (ss *StatsService) PositionHistory(symbol string, side bitfinex.OrderSide) ([]bitfinex.Stat, error)
```

#### func (*StatsService) PositionLast

```go
func (ss *StatsService) PositionLast(symbol string, side bitfinex.OrderSide) (*bitfinex.Stat, error)
```

#### func (*StatsService) SymbolCreditSizeHistory

```go
func (ss *StatsService) SymbolCreditSizeHistory(fundingSymbol string, tradingSymbol string) ([]bitfinex.Stat, error)
```

#### func (*StatsService) SymbolCreditSizeLast

```go
func (ss *StatsService) SymbolCreditSizeLast(fundingSymbol string, tradingSymbol string) (*bitfinex.Stat, error)
```

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

#### func (*StatusService) DerivativeStatusAll

```go
func (ss *StatusService) DerivativeStatusAll() ([]*bitfinex.DerivativeStatus, error)
```

#### func (*StatusService) DerivativeStatusMulti

```go
func (ss *StatusService) DerivativeStatusMulti(symbols []string) ([]*bitfinex.DerivativeStatus, error)
```

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

#### func (*TickerService) Get

```go
func (s *TickerService) Get(symbol string) (*bitfinex.Ticker, error)
```
All returns all orders for the authenticated account.

#### func (*TickerService) GetMulti

```go
func (s *TickerService) GetMulti(symbols []string) (*[]bitfinex.Ticker, error)
```

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

#### func (*TradeService) AccountAllWithSymbol

```go
func (s *TradeService) AccountAllWithSymbol(symbol string) (*bitfinex.TradeExecutionUpdateSnapshot, error)
```

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
return account trades that fit the given conditions

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
return publicly executed trades that fit the given query conditions

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

#### func (*WalletService) DepositAddress

```go
func (ws *WalletService) DepositAddress(wallet, method string) (*bitfinex.Notification, error)
```

#### func (*WalletService) SetCollateral

```go
func (s *WalletService) SetCollateral(symbol string, amount float64) (bool, error)
```
All returns all orders for the authenticated account.

#### func (*WalletService) Transfer

```go
func (ws *WalletService) Transfer(from, to, currency, currencyTo string, amount float64) (*bitfinex.Notification, error)
```

#### func (*WalletService) Wallet

```go
func (s *WalletService) Wallet() (*bitfinex.WalletSnapshot, error)
```
All returns all orders for the authenticated account.

#### func (*WalletService) Withdraw

```go
func (ws *WalletService) Withdraw(wallet, method string, amount float64, address string) (*bitfinex.Notification, error)
```
