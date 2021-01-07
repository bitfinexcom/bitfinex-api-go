package trades

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type FundingTrade struct {
	Pair   string
	ID     int64
	MTS    int64
	Amount float64
	Rate   float64
	Period int
}

type FundingTradeUpdate FundingTrade
type FundingTradeExecute FundingTrade

type FundingTradeSnapshot struct {
	Snapshot []FundingTrade
}

func FTFromRaw(pair string, raw []interface{}) (t FundingTrade, err error) {
	if len(raw) >= 5 {
		t = FundingTrade{
			Pair:   pair,
			ID:     convert.I64ValOrZero(raw[0]),
			MTS:    convert.I64ValOrZero(raw[1]),
			Amount: convert.F64ValOrZero(raw[2]),
			Rate:   convert.F64ValOrZero(raw[3]),
			Period: convert.ToInt(raw[4]),
		}
		return
	}

	err = fmt.Errorf("data slice too short: %#v", raw)
	return
}

func FTEFromRaw(pair string, raw []interface{}) (FundingTradeExecute, error) {
	ft, err := FTFromRaw(pair, raw)
	if err != nil {
		return FundingTradeExecute{}, err
	}

	return FundingTradeExecute(ft), nil
}

func FTUFromRaw(pair string, raw []interface{}) (FundingTradeUpdate, error) {
	ft, err := FTFromRaw(pair, raw)
	if err != nil {
		return FundingTradeUpdate{}, err
	}

	return FundingTradeUpdate(ft), nil
}

func FTSnapshotFromRaw(pair string, raw [][]interface{}) (FundingTradeSnapshot, error) {
	if len(raw) == 0 {
		return FundingTradeSnapshot{}, fmt.Errorf("funding trade snapshot data slice too short:%#v", raw)
	}

	snapshot := make([]FundingTrade, 0)
	for _, v := range raw {
		t, err := FTFromRaw(pair, v)
		if err != nil {
			return FundingTradeSnapshot{}, err
		}
		snapshot = append(snapshot, t)
	}

	return FundingTradeSnapshot{Snapshot: snapshot}, nil
}
