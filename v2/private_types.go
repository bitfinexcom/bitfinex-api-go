package bitfinex

import (
	"encoding/json"
	"fmt"
)

// OrderNewRequest represents an order to be posted to the bitfinex websocket
// service.
type OrderNewRequest struct {
	GID           int64   `json:"gid"`
	CID           int64   `json:"cid"`
	Type          string  `json:"type"`
	Symbol        string  `json:"symbol"`
	Amount        float64 `json:"amount,string"`
	Price         float64 `json:"price,string"`
	PriceTrailing float64 `json:"price_trailing,string,omitempty"`
	PriceAuxLimit float64 `json:"price_aux_limit,string,omitempty"`
	Hidden        bool    `json:"hidden,omitempty"`
	PostOnly      bool    `json:"postonly,omitempty"`
}

// MarshalJSON converts the order object into the format required by the bitfinex
// websocket service.
func (o *OrderNewRequest) MarshalJSON() ([]byte, error) {
	aux := struct {
		GID           int64   `json:"gid"`
		CID           int64   `json:"cid"`
		Type          string  `json:"type"`
		Symbol        string  `json:"symbol"`
		Amount        float64 `json:"amount,string"`
		Price         float64 `json:"price,string"`
		PriceTrailing float64 `json:"price_trailing,string,omitempty"`
		PriceAuxLimit float64 `json:"price_aux_limit,string,omitempty"`
		Hidden        int     `json:"hidden,omitempty"`
		PostOnly      int     `json:"postonly,omitempty"`
	}{
		GID:           o.GID,
		CID:           o.CID,
		Type:          o.Type,
		Symbol:        o.Symbol,
		Amount:        o.Amount,
		Price:         o.Price,
		PriceTrailing: o.PriceTrailing,
		PriceAuxLimit: o.PriceAuxLimit,
	}

	if o.Hidden {
		aux.Hidden = 1
	}

	if o.PostOnly {
		aux.PostOnly = 1
	}

	body := []interface{}{0, "on", nil, aux}
	return json.Marshal(&body)
}

// OrderCancelRequest represents an order cancel request.
// An order can be cancelled using the internal ID or a
// combination of Client ID (CID) and the daten for the given
// CID.
type OrderCancelRequest struct {
	ID      *int64  `json:"id,omitempty"`
	CID     *int64  `json:"cid,omitempty"`
	CIDDate *string `json:"cid_date,omitempty"`
}

// MarshalJSON converts the order cancel object into the format required by the
// bitfinex websocket service.
func (o *OrderCancelRequest) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID      *int64  `json:"id,omitempty"`
		CID     *int64  `json:"cid,omitempty"`
		CIDDate *string `json:"cid_date,omitempty"`
	}{
		ID:      o.ID,
		CID:     o.CID,
		CIDDate: o.CIDDate,
	}

	body := []interface{}{0, "oc", nil, aux}
	return json.Marshal(&body)
}

// TODO: MultiOrderCancelRequest represents an order cancel request.

type Heartbeat struct {
	//ChannelIDs []int64
}

// OrderType represents the types orders the bitfinex platform can handle.
type OrderType string

const (
	OrderTypeMarket               = "MARKET"
	OrderTypeExchangeMarket       = "EXCHANGE MARKET"
	OrderTypeLimit                = "LIMIT"
	OrderTypeExchangeLimit        = "EXCHANGE LIMIT"
	OrderTypeStop                 = "STOP"
	OrderTypeExchangeStop         = "EXCHANGE STOP"
	OrderTypeTrailingStop         = "TRAILING STOP"
	OrderTypeExchangeTrailingStop = "EXCHANGE TRAILING STOP"
	OrderTypeFOK                  = "FOK"
	OrderTypeExchangeFOK          = "EXCHANGE FOK"
	OrderTypeStopLimit            = "STOP LIMIT"
)

// OrderStatus represents the possible statuses an order can be in.
type OrderStatus string

const (
	OrderStatusActive          OrderStatus = "ACTIVE"
	OrderStatusExecuted        OrderStatus = "EXECUTED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
)

// Order as returned from the bitfinex websocket service.
type Order struct {
	ID            int64
	GID           int64
	CID           int64
	Symbol        string
	MTSCreated    int64
	MTSUpdated    int64
	Amount        float64
	AmountOrig    float64
	Type          string
	TypePrev      string
	Flags         int64
	Status        OrderStatus
	Price         float64
	PriceAvg      float64
	PriceTrailing float64
	PriceAuxLimit float64
	Notify        bool
	Hidden        bool
	PlacedID      int64
}

// OrderFromRaw takes the raw list of values as returned from the websocket
// service and tries to convert it into an Order.
func orderFromRaw(raw []interface{}) (o Order, err error) {
	if len(raw) < 26 {
		return o, fmt.Errorf("data slice too short for order: %#v", raw)
	}

	// TODO: API docs say ID, GID, CID, MTS_CREATE, MTS_UPDATE are int but API returns float
	o = Order{
		ID:            int64(f64ValOrZero(raw[0])),
		GID:           int64(f64ValOrZero(raw[1])),
		CID:           int64(f64ValOrZero(raw[2])),
		Symbol:        sValOrEmpty(raw[3]),
		MTSCreated:    int64(f64ValOrZero(raw[4])),
		MTSUpdated:    int64(f64ValOrZero(raw[5])),
		Amount:        f64ValOrZero(raw[6]),
		AmountOrig:    f64ValOrZero(raw[7]),
		Type:          sValOrEmpty(raw[8]),
		TypePrev:      sValOrEmpty(raw[9]),
		Flags:         i64ValOrZero(raw[12]),
		Status:        OrderStatus(sValOrEmpty(raw[13])),
		Price:         f64ValOrZero(raw[16]),
		PriceAvg:      f64ValOrZero(raw[17]),
		PriceTrailing: f64ValOrZero(raw[18]),
		PriceAuxLimit: f64ValOrZero(raw[19]),
		Notify:        bValOrFalse(raw[23]),
		Hidden:        bValOrFalse(raw[24]),
		PlacedID:      i64ValOrZero(raw[25]),
	}

	return
}

// OrderSnapshotFromRaw takes a raw list of values as returned from the websocket
// service and tries to convert it into an OrderSnapshot.
func orderSnapshotFromRaw(raw []interface{}) (os OrderSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := orderFromRaw(l)
				if err != nil {
					return os, err
				}
				os = append(os, o)
			}
		}
	default:
		return os, fmt.Errorf("not an order snapshot")
	}

	return
}

// OrderSnapshot is a collection of Orders that would usually be sent on
// inital connection.
type OrderSnapshot []Order

// OrderUpdate is an Order that gets sent out after every change to an
// order.
type OrderUpdate Order

// OrderNew gets sent out after an Order was created successfully.
type OrderNew Order

// OrderCancel gets sent out after an Order was cancelled successfully.
type OrderCancel Order

type PositionStatus string

const (
	PositionStatusActive PositionStatus = "ACTIVE"
	PositionStatusClosed PositionStatus = "CLOSED"
)

type Position struct {
	Symbol               string
	Status               PositionStatus
	Amount               float64
	BasePrice            float64
	MarginFunding        float64
	MarginFundingType    int64
	ProfitLoss           float64
	ProfitLossPercentage float64
	LiquidationPrice     float64
	Leverage             float64
}

func positionFromRaw(raw []interface{}) (o Position, err error) {
	if len(raw) < 10 {
		return o, fmt.Errorf("data slice too short for position: %#v", raw)
	}

	o = Position{
		Symbol:               sValOrEmpty(raw[0]),
		Status:               PositionStatus(sValOrEmpty(raw[1])),
		Amount:               f64ValOrZero(raw[2]),
		BasePrice:            f64ValOrZero(raw[3]),
		MarginFunding:        f64ValOrZero(raw[4]),
		MarginFundingType:    i64ValOrZero(raw[5]),
		ProfitLoss:           f64ValOrZero(raw[6]),
		ProfitLossPercentage: f64ValOrZero(raw[7]),
		LiquidationPrice:     f64ValOrZero(raw[8]),
		Leverage:             f64ValOrZero(raw[9]),
	}

	return
}

type PositionSnapshot []Position
type PositionNew Position
type PositionUpdate Position
type PositionCancel Position

func positionSnapshotFromRaw(raw []interface{}) (ps PositionSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				p, err := positionFromRaw(l)
				if err != nil {
					return ps, err
				}
				ps = append(ps, p)
			}
		}
	default:
		return ps, fmt.Errorf("not a position snapshot")
	}

	return
}

type Trade struct {
	ID          int64
	Pair        string
	MTSCreate   int64
	OrderID     int64
	ExecAmout   float64
	ExecPrice   float64
	OrderType   string
	OrderPrice  float64
	Maker       bool
	Fee         float64
	FeeCurrency string
}

func tradeFromRaw(raw []interface{}) (o Trade, err error) {
	if len(raw) < 11 {
		return o, fmt.Errorf("data slice too short for trade: %#v", raw)
	}

	o = Trade{
		ID:          i64ValOrZero(raw[0]),
		Pair:        sValOrEmpty(raw[1]),
		MTSCreate:   i64ValOrZero(raw[2]),
		OrderID:     i64ValOrZero(raw[3]),
		ExecAmout:   f64ValOrZero(raw[4]),
		ExecPrice:   f64ValOrZero(raw[5]),
		OrderType:   sValOrEmpty(raw[6]),
		OrderPrice:  f64ValOrZero(raw[7]),
		Maker:       bValOrFalse(raw[8]),
		Fee:         f64ValOrZero(raw[9]),
		FeeCurrency: sValOrEmpty(raw[10]),
	}

	return
}

type TradeUpdate Trade
type TradeSnapshot []Trade
type HistoricalTradeSnapshot TradeSnapshot

func tradeSnapshotFromRaw(raw []interface{}) (ts TradeSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				t, err := tradeFromRaw(l)
				if err != nil {
					return ts, err
				}
				ts = append(ts, t)
			}
		}
	default:
		return ts, fmt.Errorf("not a trade snapshot")
	}

	return
}

type TradeExecution struct {
	ID         int64
	Pair       string
	MTSCreate  int64
	OrderID    int64
	ExecAmout  float64
	ExecPrice  float64
	OrderType  string
	OrderPrice float64
	Maker      bool
}

func tradeExecutionFromRaw(raw []interface{}) (o TradeExecution, err error) {
	if len(raw) < 9 {
		return o, fmt.Errorf("data slice too short for trade execution: %#v", raw)
	}

	o = TradeExecution{
		ID:         i64ValOrZero(raw[0]),
		Pair:       sValOrEmpty(raw[1]),
		MTSCreate:  i64ValOrZero(raw[2]),
		OrderID:    i64ValOrZero(raw[3]),
		ExecAmout:  f64ValOrZero(raw[4]),
		ExecPrice:  f64ValOrZero(raw[5]),
		OrderType:  sValOrEmpty(raw[6]),
		OrderPrice: f64ValOrZero(raw[7]),
		Maker:      bValOrFalse(raw[8]),
	}

	return
}

type Wallet struct {
	Type              string
	Currency          string
	Balance           float64
	UnsettledInterest float64
	BalanceAvailable  *float64
}

func walletFromRaw(raw []interface{}) (o Wallet, err error) {
	if len(raw) < 5 {
		return o, fmt.Errorf("data slice too short for wallet: %#v", raw)
	}

	o = Wallet{
		Type:              sValOrEmpty(raw[0]),
		Currency:          sValOrEmpty(raw[1]),
		Balance:           f64ValOrZero(raw[2]),
		UnsettledInterest: f64ValOrZero(raw[3]),
		BalanceAvailable:  f64pValOrNil(raw[4]),
	}

	return
}

type WalletUpdate Wallet
type WalletSnapshot []Wallet

func walletSnapshotFromRaw(raw []interface{}) (ws WalletSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := walletFromRaw(l)
				if err != nil {
					return ws, err
				}
				ws = append(ws, o)
			}
		}
	default:
		return ws, fmt.Errorf("not an wallet snapshot")
	}

	return
}

type BalanceInfo struct {
	TotalAUM   float64
	NetAUM     float64
	WalletType string
	Currency   string
}

func balanceInfoFromRaw(raw []interface{}) (o BalanceInfo, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for balance info: %#v", raw)
	}

	o = BalanceInfo{
		TotalAUM:   f64ValOrZero(raw[0]),
		NetAUM:     f64ValOrZero(raw[1]),
		WalletType: sValOrEmpty(raw[2]),
		Currency:   sValOrEmpty(raw[3]),
	}

	return
}

type BalanceUpdate BalanceInfo

// marginInfoFromRaw returns either a MarginInfoBase or MarginInfoUpdate, since
// the Margin Info is split up into a base and per symbol parts.
func marginInfoFromRaw(raw []interface{}) (o interface{}, err error) {
	if len(raw) < 2 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	typ, ok := raw[0].(string)
	if !ok {
		return o, fmt.Errorf("expected margin info type in first position for margin info but got %#v", raw)
	}

	if len(raw) == 2 && typ == "base" { // This should be ["base", [...]]
		data, ok := raw[1].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in second position for margin info but got %#v", raw)
		}

		return marginInfoBaseFromRaw(data)
	} else if len(raw) == 3 && typ == "sym" { // This should be ["sym", SYMBOL, [...]]
		symbol, ok := raw[1].(string)
		if !ok {
			return o, fmt.Errorf("expected margin info symbol in second position for margin info update but got %#v", raw)
		}

		data, ok := raw[2].([]interface{})
		if !ok {
			return o, fmt.Errorf("expected margin info array in third position for margin info update but got %#v", raw)
		}

		return marginInfoUpdateFromRaw(symbol, data)
	}

	return nil, fmt.Errorf("invalid margin info type in %#v", raw)
}

type MarginInfoUpdate struct {
	Symbol          string
	TradableBalance float64
}

func marginInfoUpdateFromRaw(symbol string, raw []interface{}) (o MarginInfoUpdate, err error) {
	if len(raw) < 1 {
		return o, fmt.Errorf("data slice too short for margin info update: %#v", raw)
	}

	o = MarginInfoUpdate{
		Symbol:          symbol,
		TradableBalance: f64ValOrZero(raw[0]),
	}

	return
}

type MarginInfoBase struct {
	UserProfitLoss float64
	UserSwaps      float64
	MarginBalance  float64
	MarginNet      float64
}

func marginInfoBaseFromRaw(raw []interface{}) (o MarginInfoBase, err error) {
	if len(raw) < 4 {
		return o, fmt.Errorf("data slice too short for margin info base: %#v", raw)
	}

	o = MarginInfoBase{
		UserProfitLoss: f64ValOrZero(raw[0]),
		UserSwaps:      f64ValOrZero(raw[1]),
		MarginBalance:  f64ValOrZero(raw[2]),
		MarginNet:      f64ValOrZero(raw[3]),
	}

	return
}

type FundingInfo struct {
	Symbol       string
	YieldLoan    float64
	YieldLend    float64
	DurationLoan float64
	DurationLend float64
}

func fundingInfoFromRaw(raw []interface{}) (o FundingInfo, err error) {
	if len(raw) < 3 { // "sym", symbol, data
		return o, fmt.Errorf("data slice too short for funding info: %#v", raw)
	}

	sym, ok := raw[1].(string)
	if !ok {
		return o, fmt.Errorf("expected symbol in second position of funding info: %v", raw)
	}

	data, ok := raw[2].([]interface{})
	if !ok {
		return o, fmt.Errorf("expected list in third position of funding info: %v", raw)
	}

	if len(data) < 4 {
		return o, fmt.Errorf("data too short: %#v", data)
	}

	o = FundingInfo{
		Symbol:       sym,
		YieldLoan:    f64ValOrZero(data[0]),
		YieldLend:    f64ValOrZero(data[1]),
		DurationLoan: f64ValOrZero(data[2]),
		DurationLend: f64ValOrZero(data[3]),
	}

	return
}

type OfferStatus string

const (
	OfferStatusActive          OfferStatus = "ACTIVE"
	OfferStatusExecuted        OfferStatus = "EXECUTED"
	OfferStatusPartiallyFilled OfferStatus = "PARTIALLY FILLED"
	OfferStatusCanceled        OfferStatus = "CANCELED"
)

type Offer struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	MTSUpdated int64
	Amout      float64
	AmountOrig float64
	Type       string
	Flags      interface{}
	Status     OfferStatus
	Rate       float64
	Period     int64
	Notify     bool
	Hidden     bool
	Insure     bool
	Renew      bool
	RateReal   float64
}

func offerFromRaw(raw []interface{}) (o Offer, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short for offer: %#v", raw)
	}

	o = Offer{
		ID:         i64ValOrZero(raw[0]),
		Symbol:     sValOrEmpty(raw[1]),
		MTSCreated: i64ValOrZero(raw[2]),
		MTSUpdated: i64ValOrZero(raw[3]),
		Amout:      f64ValOrZero(raw[4]),
		AmountOrig: f64ValOrZero(raw[5]),
		Type:       sValOrEmpty(raw[6]),
		Flags:      raw[9],
		Status:     OfferStatus(sValOrEmpty(raw[10])),
		Rate:       f64ValOrZero(raw[14]),
		Period:     i64ValOrZero(raw[15]),
		Notify:     bValOrFalse(raw[16]),
		Hidden:     bValOrFalse(raw[17]),
		Insure:     bValOrFalse(raw[18]),
		Renew:      bValOrFalse(raw[19]),
		RateReal:   f64ValOrZero(raw[20]),
	}

	return
}

type FundingOfferNew Offer
type FundingOfferUpdate Offer
type FundingOfferCancel Offer
type FundingOfferSnapshot []Offer

func fundingOfferSnapshotFromRaw(raw []interface{}) (fos FundingOfferSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := offerFromRaw(l)
				if err != nil {
					return fos, err
				}
				fos = append(fos, o)
			}
		}
	default:
		return fos, fmt.Errorf("not a funding offer snapshot")
	}

	return
}

type HistoricalOffer Offer

type CreditStatus string

const (
	CreditStatusActive          CreditStatus = "ACTIVE"
	CreditStatusExecuted        CreditStatus = "EXECUTED"
	CreditStatusPartiallyFilled CreditStatus = "PARTIALLY FILLED"
	CreditStatusCanceled        CreditStatus = "CANCELED"
)

type Credit struct {
	ID            int64
	Symbol        string
	Side          string
	MTSCreated    int64
	MTSUpdated    int64
	Amout         float64
	Flags         interface{}
	Status        CreditStatus
	Rate          float64
	Period        int64
	MTSOpened     int64
	MTSLastPayout int64
	Notify        bool
	Hidden        bool
	Insure        bool
	Renew         bool
	RateReal      float64
	NoClose       bool
	PositionPair  string
}

func creditFromRaw(raw []interface{}) (o Credit, err error) {
	if len(raw) < 22 {
		return o, fmt.Errorf("data slice too short for offer: %#v", raw)
	}

	o = Credit{
		ID:            i64ValOrZero(raw[0]),
		Symbol:        sValOrEmpty(raw[1]),
		Side:          sValOrEmpty(raw[2]),
		MTSCreated:    i64ValOrZero(raw[3]),
		MTSUpdated:    i64ValOrZero(raw[4]),
		Amout:         f64ValOrZero(raw[5]),
		Flags:         raw[6],
		Status:        CreditStatus(sValOrEmpty(raw[7])),
		Rate:          f64ValOrZero(raw[11]),
		Period:        i64ValOrZero(raw[12]),
		MTSOpened:     i64ValOrZero(raw[13]),
		MTSLastPayout: i64ValOrZero(raw[14]),
		Notify:        bValOrFalse(raw[15]),
		Hidden:        bValOrFalse(raw[16]),
		Insure:        bValOrFalse(raw[17]),
		Renew:         bValOrFalse(raw[18]),
		RateReal:      f64ValOrZero(raw[19]),
		NoClose:       bValOrFalse(raw[20]),
		PositionPair:  sValOrEmpty(raw[21]),
	}

	return
}

type HistoricalCredit Credit
type FundingCreditNew Credit
type FundingCreditUpdate Credit
type FundingCreditCancel Credit

type FundingCreditSnapshot []Credit

func fundingCreditSnapshotFromRaw(raw []interface{}) (fcs FundingCreditSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := creditFromRaw(l)
				if err != nil {
					return fcs, err
				}
				fcs = append(fcs, o)
			}
		}
	default:
		return fcs, fmt.Errorf("not a funding credit snapshot")
	}

	return
}

type LoanStatus string

const (
	LoanStatusActive          LoanStatus = "ACTIVE"
	LoanStatusExecuted        LoanStatus = "EXECUTED"
	LoanStatusPartiallyFilled LoanStatus = "PARTIALLY FILLED"
	LoanStatusCanceled        LoanStatus = "CANCELED"
)

type Loan struct {
	ID            int64
	Symbol        string
	Side          string
	MTSCreated    int64
	MTSUpdated    int64
	Amout         float64
	Flags         interface{}
	Status        LoanStatus
	Rate          float64
	Period        int64
	MTSOpened     int64
	MTSLastPayout int64
	Notify        bool
	Hidden        bool
	Insure        bool
	Renew         bool
	RateReal      float64
	NoClose       bool
}

func loanFromRaw(raw []interface{}) (o Loan, err error) {
	if len(raw) < 21 {
		return o, fmt.Errorf("data slice too short for loan: %#v", raw)
	}

	o = Loan{
		ID:            i64ValOrZero(raw[0]),
		Symbol:        sValOrEmpty(raw[1]),
		Side:          sValOrEmpty(raw[2]),
		MTSCreated:    i64ValOrZero(raw[3]),
		MTSUpdated:    i64ValOrZero(raw[4]),
		Amout:         f64ValOrZero(raw[5]),
		Flags:         raw[6],
		Status:        LoanStatus(sValOrEmpty(raw[7])),
		Rate:          f64ValOrZero(raw[11]),
		Period:        i64ValOrZero(raw[12]),
		MTSOpened:     i64ValOrZero(raw[13]),
		MTSLastPayout: i64ValOrZero(raw[14]),
		Notify:        bValOrFalse(raw[15]),
		Hidden:        bValOrFalse(raw[16]),
		Insure:        bValOrFalse(raw[17]),
		Renew:         bValOrFalse(raw[18]),
		RateReal:      f64ValOrZero(raw[19]),
		NoClose:       bValOrFalse(raw[20]),
	}

	return o, nil
}

type HistoricalLoan Loan
type FundingLoanNew Loan
type FundingLoanUpdate Loan
type FundingLoanCancel Loan

type FundingLoanSnapshot []Loan

func fundingLoanSnapshotFromRaw(raw []interface{}) (fls FundingLoanSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := loanFromRaw(l)
				if err != nil {
					return fls, err
				}
				fls = append(fls, o)
			}
		}
	default:
		return fls, fmt.Errorf("not a funding loan snapshot")
	}

	return
}

type FundingTrade struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	OfferID    int64
	Amount     float64
	Rate       float64
	Period     int64
	Maker      int64
}

func fundingTradeFromRaw(raw []interface{}) (o FundingTrade, err error) {
	if len(raw) < 8 {
		return o, fmt.Errorf("data slice too short for funding trade: %#v", raw)
	}

	o = FundingTrade{
		ID:         i64ValOrZero(raw[0]),
		Symbol:     sValOrEmpty(raw[1]),
		MTSCreated: i64ValOrZero(raw[2]),
		OfferID:    i64ValOrZero(raw[3]),
		Amount:     f64ValOrZero(raw[4]),
		Rate:       f64ValOrZero(raw[5]),
		Period:     i64ValOrZero(raw[6]),
		Maker:      i64ValOrZero(raw[7]),
	}

	return
}

type FundingTradeExecution FundingTrade
type FundingTradeUpdate FundingTrade
type FundingTradeSnapshot []FundingTrade
type HistoricalFundingTradeSnapshot FundingTradeSnapshot

func fundingTradeSnapshotFromRaw(raw []interface{}) (fts FundingTradeSnapshot, err error) {
	if len(raw) == 0 {
		return
	}

	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				o, err := fundingTradeFromRaw(l)
				if err != nil {
					return fts, err
				}
				fts = append(fts, o)
			}
		}
	default:
		return fts, fmt.Errorf("not a funding trade snapshot")
	}

	return
}

type Notification struct {
	MTS        int64
	Type       string
	MessageID  int64
	NotifyInfo interface{}
	Code       *int64
	Status     string
	Text       string
}

func notificationFromRaw(raw []interface{}) (o Notification, err error) {
	if len(raw) < 8 {
		return o, fmt.Errorf("data slice too short for notification: %#v", raw)
	}

	o = Notification{
		MTS:       i64ValOrZero(raw[0]),
		Type:      sValOrEmpty(raw[1]),
		MessageID: i64ValOrZero(raw[2]),
		//NotifyInfo: raw[4],
		Code:   i64pValOrNil(raw[5]),
		Status: sValOrEmpty(raw[6]),
		Text:   sValOrEmpty(raw[7]),
	}

	var nraw []interface{}
	nraw = raw[4].([]interface{})
	switch o.Type {
	case "on-req":
		on, err := orderFromRaw(nraw)
		if err != nil {
			return o, err
		}
		o.NotifyInfo = OrderNew(on)
	case "oc-req":
		oc, err := orderFromRaw(nraw)
		if err != nil {
			return o, err
		}
		o.NotifyInfo = OrderCancel(oc)
	case "fon-req":
		fon, err := offerFromRaw(nraw)
		if err != nil {
			return o, err
		}
		o.NotifyInfo = FundingOfferNew(fon)
	case "foc-req":
		foc, err := offerFromRaw(nraw)
		if err != nil {
			return o, err
		}
		o.NotifyInfo = FundingOfferCancel(foc)
	case "uca":
		o.NotifyInfo = raw[4]
	}

	return
}
