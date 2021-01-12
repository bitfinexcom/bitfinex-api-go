package wallet

import (
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/convert"
)

type Wallet struct {
	Type              string
	Currency          string
	Balance           float64
	UnsettledInterest float64
	BalanceAvailable  float64
}

type Update Wallet

type Snapshot struct {
	Snapshot []*Wallet
}

type transformerFn func(raw []interface{}) (w *Wallet, err error)

func FromRaw(raw []interface{}) (w *Wallet, err error) {
	if len(raw) < 4 {
		err = fmt.Errorf("data slice too short for wallet: %#v", raw)
		return
	}

	w = &Wallet{
		Type:              convert.SValOrEmpty(raw[0]),
		Currency:          convert.SValOrEmpty(raw[1]),
		Balance:           convert.F64ValOrZero(raw[2]),
		UnsettledInterest: convert.F64ValOrZero(raw[3]),
	}

	return
}

func FromWsRaw(raw []interface{}) (w *Wallet, err error) {
	if len(raw) < 5 {
		err = fmt.Errorf("data slice too short for wallet: %#v", raw)
		return
	}

	w = &Wallet{
		Type:              convert.SValOrEmpty(raw[0]),
		Currency:          convert.SValOrEmpty(raw[1]),
		Balance:           convert.F64ValOrZero(raw[2]),
		UnsettledInterest: convert.F64ValOrZero(raw[3]),
		BalanceAvailable:  convert.F64ValOrZero(raw[4]),
	}

	return
}

// UpdateFromRaw reds "wu" type message from authenticated data
// sream and maps it to wallet.Update data structure
func UpdateFromRaw(raw []interface{}) (Update, error) {
	w, err := FromWsRaw(raw)
	if err != nil {
		return Update{}, err
	}

	return Update(*w), nil
}

func SnapshotFromRaw(raw []interface{}, transformer transformerFn) (s *Snapshot, err error) {
	if len(raw) == 0 {
		return s, fmt.Errorf("data slice too short for wallet: %#v", raw)
	}

	ws := make([]*Wallet, 0)
	switch raw[0].(type) {
	case []interface{}:
		for _, v := range raw {
			if l, ok := v.([]interface{}); ok {
				w, err := transformer(l)
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
