package bitfinex

import ()

type Ticker struct {
	Symbol          string
	Bid             float64
	BidPeriod       int64
	BidSize         float64
	Ask             float64
	AskPeriod       int64
	AskSize         float64
	DailyChange     float64
	DailyChangePerc float64
	LastPrice       float64
	Volume          float64
	High            float64
	Low             float64
}

type TickerUpdate Ticker
type TickerSnapshot []Ticker

//type Trade struct {
//ID     int64
//MTS    int64
//Amount float64
//Price  float64
//Rate   float64
//Period int64
//}
