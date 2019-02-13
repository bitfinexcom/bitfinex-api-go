package websocket

import (
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"sync"
)

type messageFactory interface {
	Build(chanID int64, objType string, raw []interface{}) (interface{}, error)
	BuildSnapshot(chanID int64, raw [][]float64) (interface{}, error)
}

type TickerFactory struct {
	*subscriptions
}

func newTickerFactory(subs *subscriptions) *TickerFactory {
	return &TickerFactory{
		subscriptions: subs,
	}
}

func (f *TickerFactory) Build(chanID int64, objType string, raw []interface{}) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err == nil {
		tick, err := bitfinex.NewTickerFromRaw(sub.Request.Symbol, raw)
		return tick, err
	}
	return nil, err
}

func (f *TickerFactory) BuildSnapshot(chanID int64, raw [][]float64) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err == nil {
		return bitfinex.NewTickerSnapshotFromRaw(sub.Request.Symbol, raw)
	}
	return nil, err
}

type TradeFactory struct {
	*subscriptions
}

func newTradeFactory(subs *subscriptions) *TradeFactory {
	return &TradeFactory{
		subscriptions: subs,
	}
}

func (f *TradeFactory) Build(chanID int64, objType string, raw []interface{}) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if "tu" == objType {
		return nil, nil // do not process TradeUpdate messages on public feed, only need to process TradeExecution (first copy seen)
	}
	if err == nil {
		trade, err := bitfinex.NewTradeFromRaw(sub.Request.Symbol, raw)
		return trade, err
	}
	return nil, err
}

func (f *TradeFactory) BuildSnapshot(chanID int64, raw [][]float64) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err == nil {
		return bitfinex.NewTradeSnapshotFromRaw(sub.Request.Symbol, raw)
	}
	return nil, err
}

type BookFactory struct {
	*subscriptions
	orderbooks     map[string]*Orderbook
	manageBooks    bool
	lock           sync.Mutex
}

func newBookFactory(subs *subscriptions, obs map[string]*Orderbook, manageBooks bool) *BookFactory {
	return &BookFactory{
		subscriptions: subs,
		orderbooks: obs,
		manageBooks: manageBooks,
	}
}

func (f *BookFactory) Build(chanID int64, objType string, raw []interface{}) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err == nil {
		update, err := bitfinex.NewBookUpdateFromRaw(sub.Request.Symbol, sub.Request.Precision, raw)
		if f.manageBooks {
			if orderbook, ok := f.orderbooks[sub.Request.Symbol]; ok {
				orderbook.UpdateWith(update)
			}
		}
		return update, err
	}
	return nil, err
}


func (f *BookFactory) BuildSnapshot(chanID int64, raw [][]float64) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	update, err2 := bitfinex.NewBookUpdateSnapshotFromRaw(sub.Request.Symbol, sub.Request.Precision, raw)
	if err2 != nil {
		return nil, err2
	}
	if err == nil {
		if f.manageBooks {
			f.lock.Lock()
			defer f.lock.Unlock()
			// create new orderbook
			f.orderbooks[sub.Request.Symbol] = &Orderbook{
				symbol: sub.Request.Symbol,
				bids:   make([]*bitfinex.BookUpdate, 0),
				asks:   make([]*bitfinex.BookUpdate, 0),
			}
			f.orderbooks[sub.Request.Symbol].SetWithSnapshot(update)
		}
		return update, err
	}
	return nil, err
}

type CandlesFactory struct {
	*subscriptions
}

func newCandlesFactory(subs *subscriptions) *CandlesFactory {
	return &CandlesFactory{
		subscriptions: subs,
	}
}

func (f *CandlesFactory) Build(chanID int64, objType string, raw []interface{}) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err != nil {
		return nil, err
	}
	sym, res, err := extractSymbolResolutionFromKey(sub.Request.Key)
	if err != nil {
		return nil, err
	}
	candle, err := bitfinex.NewCandleFromRaw(sym, res, raw)
	return candle, err
}

func (f *CandlesFactory) BuildSnapshot(chanID int64, raw [][]float64) (interface{}, error) {
	sub, err := f.subscriptions.lookupByChannelID(chanID)
	if err != nil {
		return nil, err
	}
	sym, res, err := extractSymbolResolutionFromKey(sub.Request.Key)
	if err != nil {
		return nil, err
	}
	snap, err := bitfinex.NewCandleSnapshotFromRaw(sym, res, raw)
	return snap, err
}
