package fundingcredit

import (
	"encoding/json"
	"fmt"
)

type CancelRequest struct {
	ID int64
}

func (cr *CancelRequest) ToJSON() ([]byte, error) {
	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: cr.ID,
	}
	return json.Marshal(resp)
}

// MarshalJSON converts the funding credit cancel request into the format required by the
// bitfinex websocket service.
func (cr *CancelRequest) MarshalJSON() ([]byte, error) {
	b, err := cr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"fcc\", null, %s]", string(b))), nil
}
