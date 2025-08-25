package wallet

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

type Wallet struct {
	Type              string
	Currency          string
	Balance           float64
	UnsettledInterest float64
	BalanceAvailable  float64
	LastChange        string
	TradeDetails      map[string]interface{}
}

type Update Wallet

type Snapshot struct {
	Snapshot []*Wallet
}

func FromRaw(raw []interface{}) (w *Wallet, err error) {
	if len(raw) < 7 {
		err = fmt.Errorf("data slice too short for wallet: %#v", raw)
		return
	}

	w = &Wallet{
		Type:              convert.SValOrEmpty(raw[0]),
		Currency:          convert.SValOrEmpty(raw[1]),
		Balance:           convert.F64ValOrZero(raw[2]),
		UnsettledInterest: convert.F64ValOrZero(raw[3]),
		BalanceAvailable:  convert.F64ValOrZero(raw[4]),
		LastChange:        convert.SValOrEmpty(raw[5]),
	}

	if meta, ok := raw[6].(map[string]interface{}); ok {
		w.TradeDetails = meta
	}

	return
}

// UpdateFromRaw reds "wu" type message from authenticated data
// sream and maps it to wallet.Update data structure
func UpdateFromRaw(raw []interface{}) (Update, error) {
	w, err := FromRaw(raw)
	if err != nil {
		return Update{}, err
	}

	return Update(*w), nil
}

func SnapshotFromRaw(raw []interface{}) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for wallet: %#v", raw)
	}

	ws := make([]*Wallet, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				w, err := FromRaw(l)
				if err != nil {
					return s, err
				}
				ws = append(ws, w)
			}
		}
	default:
		return s, fmt.Errorf("not an wallet snapshot")
	}
	s = &Snapshot{Snapshot: ws}

	return
}

type Movement struct {
	ID int64
	common.MovementType
	Currency                string
	CurrencyName            string
	StartedAt               int64
	LastUpdatedAt           int64
	Status                  string
	Amount                  float64
	Fee                     float64
	DestinationAddress      string
	AddrTag                 string
	TransactionId           string
	WithdrawTransactionNote string
}

type MovementSnapshot struct {
	Snapshot []*Movement
}

// MovementFromRaw ...
// [
//    ID,
//    CURRENCY,
//    CURRENCY_NAME,
//    null,
//    null,
//    MTS_STARTED,
//    MTS_UPDATED,
//    null,
//    null,
//    STATUS,
//    null,
//    null,
//    AMOUNT,
//    FEES,
//    null,
//    null,
//    DESTINATION_ADDRESS,
//    null,
//    null,
//    null,
//    TRANSACTION_ID,
//    WITHDRAW_TRANSACTION_NOTE
//  ]
//
func MovementFromRaw(raw []interface{}) (w *Movement, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("parse movement err: %v", err1)
		}
	}()

	if len(raw) < 22 {
		err = fmt.Errorf("data slice too short for movement: %#v", raw)
		return
	}

	w = &Movement{
		ID:                      convert.I64ValOrZero(raw[0]),
		Currency:                convert.SValOrEmpty(raw[1]),
		CurrencyName:            convert.SValOrEmpty(raw[2]),
		StartedAt:               convert.I64ValOrZero(raw[5]),
		LastUpdatedAt:           convert.I64ValOrZero(raw[6]),
		Status:                  convert.SValOrEmpty(raw[9]),
		Fee:                     -convert.F64ValOrZero(raw[13]),
		DestinationAddress:      convert.SValOrEmpty(raw[16]),
		AddrTag:                 convert.SValOrEmpty(raw[17]),
		TransactionId:           convert.SValOrEmpty(raw[20]),
		WithdrawTransactionNote: convert.SValOrEmpty(raw[21]),
	}

	amount := convert.F64ValOrZero(raw[12])
	if amount > 0 {
		w.MovementType = common.Deposit
		w.Amount = amount
	} else {
		w.MovementType = common.Withdraw
		w.Amount = -amount
	}

	return
}

func MovementSnapshotFromRaw(raw []interface{}) (s *MovementSnapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for movements: %#v", raw)
	}

	movements := make([]*Movement, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				w, err := MovementFromRaw(l)
				if err != nil {
					return s, err
				}
				movements = append(movements, w)
			}
		}
	default:
		return s, fmt.Errorf("not an movements snapshot")
	}

	s = &MovementSnapshot{Snapshot: movements}
	return
}
