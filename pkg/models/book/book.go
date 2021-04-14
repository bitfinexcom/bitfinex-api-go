package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// BookAction represents a new/update or removal for a book entry.
type BookAction byte

const (
	//BookEntry represents a new or updated book entry.
	BookEntry BookAction = 0
	//BookRemoveEntry represents a removal of a book entry.
	BookRemoveEntry BookAction = 1
)

// Book represents an order book price update.
type Book struct {
	Symbol      string // book symbol
	ID          int64  // the book update ID, optional
	Count       int64  // updated count, optional
	Period      int64
	Price       float64 // updated price
	Amount      float64 // updated amount
	Rate        float64
	PriceJsNum  json.Number      // update price as json.Number
	AmountJsNum json.Number      // update amount as json.Number
	Side        common.OrderSide // side
	Action      BookAction       // action (add/remove)
}

type Snapshot struct {
	Snapshot []*Book
}

func SnapshotFromRaw(symbol, precision string, raw [][]interface{}, rawNumbers interface{}) (*Snapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for book snapshot: %#v", raw)
	}

	snap := make([]*Book, len(raw))
	for i, v := range raw {
		b, err := FromRaw(symbol, precision, v, rawNumbers.([]interface{})[i])
		if err != nil {
			return nil, err
		}
		snap[i] = b
	}

	return &Snapshot{Snapshot: snap}, nil
}

func IsRawBook(precision string) bool {
	return precision == "R0"
}

// FromRaw creates a new book object from raw data. Precision determines how
// to interpret the side (baked into Count versus Amount)
// raw book updates [ID, price, qty], aggregated book updates [price, amount, count]
func FromRaw(symbol, precision string, raw []interface{}, rawNumbers interface{}) (b *Book, err error) {
	if len(raw) < 3 {
		return b, fmt.Errorf("raw slice too short for book, expected %d got %d: %#v", 3, len(raw), raw)
	}

	rawBook := IsRawBook(precision)

	if len(raw) == 3 && rawBook {
		b = rawTradingPairsBook(raw, rawNumbers)
	}

	if len(raw) == 3 && !rawBook {
		b = tradingPairsBook(raw, rawNumbers)
	}

	if len(raw) >= 4 && rawBook {
		b = rawFundingPairsBook(raw, rawNumbers)
	}

	if len(raw) >= 4 && !rawBook {
		b = fundingPairsBook(raw, rawNumbers)
	}

	b.Symbol = symbol

	return
}

// FromWSRaw - based on condition will return snapshot of books or single book
func FromWSRaw(symbol, precision string, data []interface{}) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice")
	}

	_, isSnapshot := data[0].([]interface{})
	if isSnapshot {
		return SnapshotFromRaw(symbol, precision, convert.ToInterfaceArray(data), data)
	}

	return FromRaw(symbol, precision, data, data)
}

func rawTradingPairsBook(raw []interface{}, rawNumbers interface{}) *Book {
	// [ ORDER_ID, PRICE, AMOUNT ] - raw trading pairs signature
	var (
		side   common.OrderSide
		action BookAction
	)

	rawNumSlice := rawNumbers.([]interface{})
	price := convert.F64ValOrZero(raw[1])
	amount := convert.F64ValOrZero(raw[2])

	if amount > 0 {
		side = common.Bid
	} else {
		side = common.Ask
	}

	if price <= 0 {
		action = BookRemoveEntry
	} else {
		action = BookEntry
	}

	return &Book{
		Price:       math.Abs(price),
		PriceJsNum:  convert.FloatToJsonNumber(rawNumSlice[1]),
		Amount:      math.Abs(amount),
		AmountJsNum: convert.FloatToJsonNumber(rawNumSlice[2]),
		Side:        side,
		Action:      action,
		ID:          convert.I64ValOrZero(raw[0]),
	}
}

func tradingPairsBook(raw []interface{}, rawNumbers interface{}) *Book {
	// [ PRICE, COUNT, AMOUNT ] - trading pairs signature
	var (
		price    float64
		count    int64
		priceNum json.Number
		side     common.OrderSide
		action   BookAction
	)

	rawNumSlice := rawNumbers.([]interface{})
	amount := convert.F64ValOrZero(raw[2])
	amountNum := convert.FloatToJsonNumber(rawNumSlice[2])

	price = convert.F64ValOrZero(raw[0])
	priceNum = convert.FloatToJsonNumber(rawNumSlice[0])
	count = convert.I64ValOrZero(raw[1])

	if amount > 0 {
		side = common.Bid
	} else {
		side = common.Ask
	}

	if count <= 0 {
		action = BookRemoveEntry
	} else {
		action = BookEntry
	}

	return &Book{
		Price:       math.Abs(price),
		PriceJsNum:  priceNum,
		Count:       count,
		Amount:      math.Abs(amount),
		AmountJsNum: amountNum,
		Side:        side,
		Action:      action,
	}
}

func rawFundingPairsBook(raw []interface{}, rawNumbers interface{}) *Book {
	// [ ORDER_ID, PERIOD, RATE, AMOUNT ] - raw funding pairs signature
	rawNumSlice := rawNumbers.([]interface{})

	return &Book{
		ID:          convert.I64ValOrZero(raw[0]),
		Period:      convert.I64ValOrZero(raw[1]),
		Rate:        convert.F64ValOrZero(raw[2]),
		Amount:      convert.F64ValOrZero(raw[3]),
		AmountJsNum: convert.FloatToJsonNumber(rawNumSlice[3]),
	}
}

func fundingPairsBook(raw []interface{}, rawNumbers interface{}) *Book {
	// [ RATE, PERIOD, COUNT, AMOUNT ], - funding pairs signature
	rawNumSlice := rawNumbers.([]interface{})

	return &Book{
		Rate:        convert.F64ValOrZero(raw[0]),
		Period:      convert.I64ValOrZero(raw[1]),
		Count:       convert.I64ValOrZero(raw[2]),
		Amount:      convert.F64ValOrZero(raw[3]),
		AmountJsNum: convert.FloatToJsonNumber(rawNumSlice[3]),
	}
}
