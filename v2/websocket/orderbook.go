package websocket

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"sort"
	"sync"
	"strings"
	"hash/crc32"
)

type Orderbook struct {
	lock sync.Mutex

	symbol string
	bids   []*bitfinex.BookUpdate
	asks   []*bitfinex.BookUpdate
}

func (ob *Orderbook) SetWithSnapshot(bs *bitfinex.BookUpdateSnapshot) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	ob.bids = make([]*bitfinex.BookUpdate, 0)
	ob.asks = make([]*bitfinex.BookUpdate, 0)
	for _, order := range bs.Snapshot {
    if (order.Side == bitfinex.Bid) {
			ob.bids = append(ob.bids, order)
		} else {
			ob.asks = append(ob.asks, order)
		}
	}
}

func (ob *Orderbook) UpdateWith(bu *bitfinex.BookUpdate) {
	ob.lock.Lock()
	defer ob.lock.Unlock()

	side := &ob.asks
	if (bu.Side == bitfinex.Bid) {
		side = &ob.bids
	}

	// check if first in book
	if (len(*side) == 0) {
		*side = append(*side, bu)
		return
	}

	// match price level
	for index, sOrder := range *side {
		if (sOrder.Price == bu.Price) {
			if (index+1 > len(*(side))) {
				return
			}
			if (bu.Count <= 0) {
				// delete if count is equal to zero
				*side = append((*side)[:index], (*side)[index+1:]...)
				return
			} else {
				// remove now and we will add in the code below
				*side = append((*side)[:index], (*side)[index+1:]...)
			}
		}
	}
	*side = append(*side, bu)
	// add to the orderbook and sort lowest to highest
	sort.Slice(*side, func(i, j int) bool {
		if (i >= len(*(side)) || j >= len(*(side))) {
			return false
		}
		if bu.Side == bitfinex.Ask {
			return (*side)[i].Price < (*side)[j].Price
		} else {
			return (*side)[i].Price > (*side)[j].Price
		}
	})
}

func (ob *Orderbook) Checksum() (uint32) {
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

