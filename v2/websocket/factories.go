package websocket

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/book"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/candle"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/derivatives"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/trade"
)

type messageFactory interface {
	Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error)
	BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error)
}

type TickerFactory struct {
	*subscriptions
}

func newTickerFactory(subs *subscriptions) *TickerFactory {
	return &TickerFactory{
		subscriptions: subs,
	}
}

func (f *TickerFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error) {
	return ticker.FromRaw(sub.Request.Symbol, raw)
}

func (f *TickerFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error) {
	return ticker.SnapshotFromRaw(sub.Request.Symbol, raw)
}

type TradeFactory struct {
	*subscriptions
}

func newTradeFactory(subs *subscriptions) *TradeFactory {
	return &TradeFactory{
		subscriptions: subs,
	}
}

func (f *TradeFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error) {
	if "tu" == objType {
		return nil, nil // do not process TradeUpdate messages on public feed, only need to process TradeExecution (first copy seen)
	}
	return trade.FromRaw(sub.Request.Symbol, raw)
}

func (f *TradeFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error) {
	return trade.SnapshotFromRaw(sub.Request.Symbol, raw)
}

type BookFactory struct {
	*subscriptions
	orderbooks  map[string]*Orderbook
	manageBooks bool
	lock        sync.Mutex
}

func newBookFactory(subs *subscriptions, obs map[string]*Orderbook, manageBooks bool) *BookFactory {
	return &BookFactory{
		subscriptions: subs,
		orderbooks:    obs,
		manageBooks:   manageBooks,
	}
}

func ConvertBytesToJsonNumberArray(raw_bytes []byte) ([]interface{}, error) {
	var raw_json_number []interface{}
	d := json.NewDecoder(strings.NewReader(string(raw_bytes)))
	d.UseNumber()
	str_conv_err := d.Decode(&raw_json_number)
	if str_conv_err != nil {
		return nil, str_conv_err
	}
	return raw_json_number, nil
}

func (f *BookFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error) {
	update, err := book.FromRaw(sub.Request.Symbol, sub.Request.Precision, raw)
	if f.manageBooks {
		f.lock.Lock()
		defer f.lock.Unlock()
		if orderbook, ok := f.orderbooks[sub.Request.Symbol]; ok {
			orderbook.UpdateWith(update)
		}
	}
	return update, err
}

func (f *BookFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error) {
	update, err := book.SnapshotFromRaw(sub.Request.Symbol, sub.Request.Precision, raw)
	if err != nil {
		return nil, err
	}

	if f.manageBooks {
		f.lock.Lock()
		defer f.lock.Unlock()
		// create new orderbook
		f.orderbooks[sub.Request.Symbol] = &Orderbook{
			symbol: sub.Request.Symbol,
			bids:   make([]*book.Book, 0),
			asks:   make([]*book.Book, 0),
		}
		f.orderbooks[sub.Request.Symbol].SetWithSnapshot(update)
	}

	return update, nil
}

type CandlesFactory struct {
	*subscriptions
}

func newCandlesFactory(subs *subscriptions) *CandlesFactory {
	return &CandlesFactory{
		subscriptions: subs,
	}
}

func (f *CandlesFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error) {
	sym, res, err := extractSymbolResolutionFromKey(sub.Request.Key)
	if err != nil {
		return nil, err
	}
	candle, err := candle.FromRaw(sym, res, raw)
	return candle, err
}

func (f *CandlesFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error) {
	sym, res, err := extractSymbolResolutionFromKey(sub.Request.Key)
	if err != nil {
		return nil, err
	}
	snap, err := candle.SnapshotFromRaw(sym, res, raw)
	return snap, err
}

type StatsFactory struct {
	*subscriptions
}

func newStatsFactory(subs *subscriptions) *StatsFactory {
	return &StatsFactory{
		subscriptions: subs,
	}
}

func (f *StatsFactory) Build(sub *subscription, objType string, raw []interface{}, raw_bytes []byte) (interface{}, error) {
	splits := strings.Split(sub.Request.Key, ":")
	if len(splits) != 3 {
		return nil, fmt.Errorf("unable to parse key to symbol %s", sub.Request.Key)
	}
	symbol := splits[1] + ":" + splits[2]
	d, err := derivatives.FromWsRaw(symbol, raw)
	return d, err
}

func (f *StatsFactory) BuildSnapshot(sub *subscription, raw [][]interface{}, raw_bytes []byte) (interface{}, error) {
	// no snapshots
	return nil, nil
}
