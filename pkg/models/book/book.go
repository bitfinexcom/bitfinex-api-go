package book

import (
	"encoding/json"
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
	ID          int64            // the book update ID, optional
	Symbol      string           // book symbol
	Price       float64          // updated price
	PriceJsNum  json.Number      // update price as json.Number
	Count       int64            // updated count, optional
	Amount      float64          // updated amount
	AmountJsNum json.Number      // update amount as json.Number
	Side        common.OrderSide // side
	Action      BookAction       // action (add/remove)
}

type Snapshot struct {
	Snapshot []*Book
}

func SnapshotFromRaw(symbol, precision string, raw [][]interface{}) (*Snapshot, error) {
	if len(raw) <= 0 {
		return nil, fmt.Errorf("data slice too short for book snapshot: %#v", raw)
	}
	snap := make([]*Book, len(raw))
	for i, v := range raw {
		b, err := FromRaw(symbol, precision, v)
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
func FromRaw(symbol, precision string, data []interface{}) (b *Book, err error) {
	if len(data) < 3 {
		return b, fmt.Errorf("data slice too short for book update, expected %d got %d: %#v", 3, len(data), data)
	}

	// DS: by the looks of it, it does not handle funding currency?
	var (
		price, actionCtrl float64
		id, cnt           int64
		priceNum          json.Number
		side              common.OrderSide
	)

	amount := convert.F64ValOrZero(data[2])
	amountNum := convert.FloatToJsonNumber(data[2])

	if IsRawBook(precision) {
		// [ID, price, amount]
		id = convert.I64ValOrZero(data[0])
		price = convert.F64ValOrZero(data[1])
		priceNum = convert.FloatToJsonNumber(data[1])
		actionCtrl = price
	} else {
		// [price, count, amount]
		price = convert.F64ValOrZero(data[0])
		priceNum = convert.FloatToJsonNumber(data[0])
		cnt = convert.I64ValOrZero(data[1])
		actionCtrl = float64(cnt)
	}

	if amount > 0 {
		side = common.Bid
	} else {
		side = common.Ask
	}

	var action BookAction
	if actionCtrl <= 0 {
		action = BookRemoveEntry
	} else {
		action = BookEntry
	}

	b = &Book{
		Symbol:      symbol,
		Price:       math.Abs(price),
		PriceJsNum:  priceNum,
		Count:       cnt,
		Amount:      math.Abs(amount),
		AmountJsNum: amountNum,
		Side:        side,
		Action:      action,
		ID:          id,
	}

	return
}
