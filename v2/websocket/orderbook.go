package websocket

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strings"
	"sync"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

type Orderbook struct {
	lock sync.RWMutex

	symbol string
	bids   []*book.Book
	asks   []*book.Book
}

// return a dereferenced copy of an orderbook side. This is so consumers can access
// the book but not change the values that are used to generate the crc32 checksum
func (ob *Orderbook) copySide(side []*book.Book) []book.Book {
	var cpy []book.Book
	for i := 0; i < len(side); i++ {
		cpy = append(cpy, *side[i])
	}
	return cpy
}

func (ob *Orderbook) Symbol() string {
	return ob.symbol
}

func (ob *Orderbook) Asks() []book.Book {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	return ob.copySide(ob.asks)
}

func (ob *Orderbook) Bids() []book.Book {
	ob.lock.RLock()
	defer ob.lock.RUnlock()
	return ob.copySide(ob.bids)
}

func (ob *Orderbook) BidsAndAsks() ([]book.Book, []book.Book) {
	ob.lock.RLock()
	defer ob.lock.RUnlock()

	return ob.copySide(ob.bids), ob.copySide(ob.asks)
}

func (ob *Orderbook) SetWithSnapshot(bs *book.Snapshot) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	ob.bids = make([]*book.Book, 0)
	ob.asks = make([]*book.Book, 0)
	for _, order := range bs.Snapshot {
		if order.Side == common.Bid {
			ob.bids = append(ob.bids, order)
		} else {
			ob.asks = append(ob.asks, order)
		}
	}
}

func (ob *Orderbook) UpdateWith(b *book.Book) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	side := &ob.asks
	if b.Side == common.Bid {
		side = &ob.bids
	}

	// check if first in book
	if len(*side) == 0 {
		*side = append(*side, b)
		return
	}

	// match price level
	for index, sOrder := range *side {
		if sOrder.Price == b.Price {
			if index+1 > len(*(side)) {
				return
			}
			if b.Count <= 0 {
				// delete if count is equal to zero
				*side = append((*side)[:index], (*side)[index+1:]...)
				return
			}
			// remove now and we will add in the code below
			*side = append((*side)[:index], (*side)[index+1:]...)
		}
	}

	// price may not match at previous step
	if b.Count <= 0 {
		fmt.Printf("bitfinex matched price level %v not found at local cache, id: %v, symbol: %v.\n",
			b.Price, b.ID, b.Symbol)
		return
	}

	*side = append(*side, b)
	// add to the orderbook and sort lowest to highest
	sort.Slice(*side, func(i, j int) bool {
		if i >= len(*(side)) || j >= len(*(side)) {
			return false
		}
		if b.Side == common.Ask {
			return (*side)[i].Price < (*side)[j].Price
		}
		return (*side)[i].Price > (*side)[j].Price
	})
}

func (ob *Orderbook) Checksum() uint32 {
	ob.lock.Lock()
	defer ob.lock.Unlock()
	var checksumItems []string
	for i := 0; i < 25; i++ {
		if len(ob.bids) > i {
			// append bid
			checksumItems = append(checksumItems, (ob.bids)[i].PriceJsNum.String())
			checksumItems = append(checksumItems, (ob.bids)[i].AmountJsNum.String())
		}
		if len(ob.asks) > i {
			// append ask
			checksumItems = append(checksumItems, (ob.asks)[i].PriceJsNum.String())
			checksumItems = append(checksumItems, (ob.asks)[i].AmountJsNum.String())
		}
	}
	checksumStrings := strings.Join(checksumItems, ":")
	return crc32.ChecksumIEEE([]byte(checksumStrings))
}
