package trades

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

// AuthFundingTrade data structure
type AuthFundingTrade struct {
	ID         int64
	Symbol     string
	MTSCreated int64
	OfferID    int64
	Amount     float64
	Rate       float64
	Period     int64
	Maker      int64
}

type AuthFundingTradeUpdate AuthFundingTrade
type AuthFundingTradeExecuted AuthFundingTrade

type AuthFundingTradeSnapshot struct {
	Snapshot []AuthFundingTrade
}

// AFTFromRaw maps raw data slice to instance of AuthFundingTrade
func AFTFromRaw(raw []interface{}) (aft AuthFundingTrade, err error) {
	if len(raw) < 8 {
		return AuthFundingTrade{}, fmt.Errorf("data slice too short for funding trade: %#v", raw)
	}

	aft = AuthFundingTrade{
		ID:         convert.I64ValOrZero(raw[0]),
		Symbol:     convert.SValOrEmpty(raw[1]),
		MTSCreated: convert.I64ValOrZero(raw[2]),
		OfferID:    convert.I64ValOrZero(raw[3]),
		Amount:     convert.F64ValOrZero(raw[4]),
		Rate:       convert.F64ValOrZero(raw[5]),
		Period:     convert.I64ValOrZero(raw[6]),
		Maker:      convert.I64ValOrZero(raw[7]),
	}

	return
}

// AFTUFromRaw maps raw data slice to instance of AuthFundingTradeUpdate
func AFTUFromRaw(raw []interface{}) (AuthFundingTradeUpdate, error) {
	aft, err := AFTFromRaw(raw)
	if err != nil {
		return AuthFundingTradeUpdate{}, err
	}

	return AuthFundingTradeUpdate(aft), nil
}

// AFTEFromRaw maps raw data slice to instance of AuthFundingTradeExecuted
func AFTEFromRaw(raw []interface{}) (AuthFundingTradeExecuted, error) {
	aft, err := AFTFromRaw(raw)
	if err != nil {
		return AuthFundingTradeExecuted{}, err
	}

	return AuthFundingTradeExecuted(aft), nil
}

// AFTSnapshotFromRaw maps raw data slice to authenticated funding trade data structures
func AFTSnapshotFromRaw(raw [][]interface{}) (AuthFundingTradeSnapshot, error) {
	if len(raw) == 0 {
		return AuthFundingTradeSnapshot{}, fmt.Errorf("data slice too short for funding trade snapshot: %#v", raw)
	}

	snap := make([]AuthFundingTrade, 0)
	for _, r := range raw {
		ft, err := AFTFromRaw(r)
		if err != nil {
			return AuthFundingTradeSnapshot{}, err
		}
		snap = append(snap, ft)
	}

	return AuthFundingTradeSnapshot{Snapshot: snap}, nil
}
