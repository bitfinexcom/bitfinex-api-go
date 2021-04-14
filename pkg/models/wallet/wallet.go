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
