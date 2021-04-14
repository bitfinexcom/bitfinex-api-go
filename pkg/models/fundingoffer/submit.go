package fundingoffer

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

type SubmitRequest struct {
	Type   string
	Symbol string
	Amount float64
	Rate   float64
	Period int64
	Hidden bool
}

func (sr *SubmitRequest) ToJSON() ([]byte, error) {
	aux := struct {
		Type   string  `json:"type"`
		Symbol string  `json:"symbol"`
		Amount float64 `json:"amount,string"`
		Rate   float64 `json:"rate,string"`
		Period int64   `json:"period"`
		Flags  int     `json:"flags,omitempty"`
	}{
		Type:   sr.Type,
		Symbol: sr.Symbol,
		Amount: sr.Amount,
		Rate:   sr.Rate,
		Period: sr.Period,
	}
	if sr.Hidden {
		aux.Flags = aux.Flags + common.OrderFlagHidden
	}
	return json.Marshal(aux)
}

// MarshalJSON converts the offer submit object into the format required by the
// bitfinex websocket service.
func (sr *SubmitRequest) MarshalJSON() ([]byte, error) {
	aux, err := sr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"fon\", null, %s]", string(aux))), nil
}
