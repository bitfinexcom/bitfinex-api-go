package order

import (
	"encoding/json"
	"fmt"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// NewRequest represents an order to be posted to the bitfinex websocket
// service.
type NewRequest struct {
	GID           int64                  `json:"gid"`
	CID           int64                  `json:"cid"`
	Type          string                 `json:"type"`
	Symbol        string                 `json:"symbol"`
	Amount        float64                `json:"amount,string"`
	Price         float64                `json:"price,string"`
	Leverage      int64                  `json:"lev,omitempty"`
	PriceTrailing float64                `json:"price_trailing,string,omitempty"`
	PriceAuxLimit float64                `json:"price_aux_limit,string,omitempty"`
	PriceOcoStop  float64                `json:"price_oco_stop,string,omitempty"`
	Hidden        bool                   `json:"hidden,omitempty"`
	PostOnly      bool                   `json:"postonly,omitempty"`
	Close         bool                   `json:"close,omitempty"`
	OcoOrder      bool                   `json:"oco_order,omitempty"`
	TimeInForce   string                 `json:"tif,omitempty"`
	AffiliateCode string                 `json:"-"`
	Meta          map[string]interface{} `json:"meta,omitempty"`
}

// MarshalJSON converts the order object into the format required by the bitfinex
// websocket service.
func (nr *NewRequest) MarshalJSON() ([]byte, error) {
	jsonOrder, err := nr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"on\", null, %s]", string(jsonOrder))), nil
}

// EnrichedPayload returns enriched representation of order struct for submission
func (nr *NewRequest) EnrichedPayload() interface{} {
	pld := struct {
		GID           int64                  `json:"gid"`
		CID           int64                  `json:"cid"`
		Type          string                 `json:"type"`
		Symbol        string                 `json:"symbol"`
		Amount        float64                `json:"amount,string"`
		Price         float64                `json:"price,string"`
		Leverage      int64                  `json:"lev,omitempty"`
		PriceTrailing float64                `json:"price_trailing,string,omitempty"`
		PriceAuxLimit float64                `json:"price_aux_limit,string,omitempty"`
		PriceOcoStop  float64                `json:"price_oco_stop,string,omitempty"`
		TimeInForce   string                 `json:"tif,omitempty"`
		Flags         int                    `json:"flags,omitempty"`
		Meta          map[string]interface{} `json:"meta,omitempty"`
	}{
		GID:           nr.GID,
		CID:           nr.CID,
		Type:          nr.Type,
		Symbol:        nr.Symbol,
		Amount:        nr.Amount,
		Price:         nr.Price,
		Leverage:      nr.Leverage,
		PriceTrailing: nr.PriceTrailing,
		PriceAuxLimit: nr.PriceAuxLimit,
		PriceOcoStop:  nr.PriceOcoStop,
		TimeInForce:   nr.TimeInForce,
	}

	if nr.Hidden {
		pld.Flags = pld.Flags + common.OrderFlagHidden
	}

	if nr.PostOnly {
		pld.Flags = pld.Flags + common.OrderFlagPostOnly
	}

	if nr.OcoOrder {
		pld.Flags = pld.Flags + common.OrderFlagOCO
	}

	if nr.Close {
		pld.Flags = pld.Flags + common.OrderFlagClose
	}

	if nr.Meta == nil {
		pld.Meta = make(map[string]interface{})
	}

	if nr.AffiliateCode != "" {
		pld.Meta["aff_code"] = nr.AffiliateCode
	}

	return pld
}

func (nr *NewRequest) ToJSON() ([]byte, error) {
	return json.Marshal(nr.EnrichedPayload())
}

type UpdateRequest struct {
	ID            int64                  `json:"id"`
	GID           int64                  `json:"gid,omitempty"`
	Price         float64                `json:"price,string,omitempty"`
	Amount        float64                `json:"amount,string,omitempty"`
	Leverage      int64                  `json:"lev,omitempty"`
	Delta         float64                `json:"delta,string,omitempty"`
	PriceTrailing float64                `json:"price_trailing,string,omitempty"`
	PriceAuxLimit float64                `json:"price_aux_limit,string,omitempty"`
	Hidden        bool                   `json:"hidden,omitempty"`
	PostOnly      bool                   `json:"postonly,omitempty"`
	TimeInForce   string                 `json:"tif,omitempty"`
	Meta          map[string]interface{} `json:"meta,omitempty"`
}

// MarshalJSON converts the order object into the format required by the bitfinex
// websocket service.
func (ur *UpdateRequest) MarshalJSON() ([]byte, error) {
	aux, err := ur.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"ou\", null, %s]", string(aux))), nil
}

func (ur *UpdateRequest) EnrichedPayload() interface{} {
	pld := struct {
		ID            int64                  `json:"id"`
		GID           int64                  `json:"gid,omitempty"`
		Price         float64                `json:"price,string,omitempty"`
		Amount        float64                `json:"amount,string,omitempty"`
		Leverage      int64                  `json:"lev,omitempty"`
		Delta         float64                `json:"delta,string,omitempty"`
		PriceTrailing float64                `json:"price_trailing,string,omitempty"`
		PriceAuxLimit float64                `json:"price_aux_limit,string,omitempty"`
		Hidden        bool                   `json:"hidden,omitempty"`
		PostOnly      bool                   `json:"postonly,omitempty"`
		TimeInForce   string                 `json:"tif,omitempty"`
		Flags         int                    `json:"flags,omitempty"`
		Meta          map[string]interface{} `json:"meta,omitempty"`
	}{
		ID:            ur.ID,
		GID:           ur.GID,
		Amount:        ur.Amount,
		Leverage:      ur.Leverage,
		Price:         ur.Price,
		PriceTrailing: ur.PriceTrailing,
		PriceAuxLimit: ur.PriceAuxLimit,
		Delta:         ur.Delta,
		TimeInForce:   ur.TimeInForce,
	}

	if ur.Meta == nil {
		pld.Meta = make(map[string]interface{})
	}

	if ur.Hidden {
		pld.Flags = pld.Flags + common.OrderFlagHidden
	}

	if ur.PostOnly {
		pld.Flags = pld.Flags + common.OrderFlagPostOnly
	}

	return pld
}

func (ur *UpdateRequest) ToJSON() ([]byte, error) {
	return json.Marshal(ur.EnrichedPayload())
}

// CancelRequest represents an order cancel request.
// An order can be cancelled using the internal ID or a
// combination of Client ID (CID) and the daten for the given
// CID.
type CancelRequest struct {
	ID      int64  `json:"id,omitempty"`
	CID     int64  `json:"cid,omitempty"`
	CIDDate string `json:"cid_date,omitempty"`
}

func (cr *CancelRequest) ToJSON() ([]byte, error) {
	resp := struct {
		ID      int64  `json:"id,omitempty"`
		CID     int64  `json:"cid,omitempty"`
		CIDDate string `json:"cid_date,omitempty"`
	}{
		ID:      cr.ID,
		CID:     cr.CID,
		CIDDate: cr.CIDDate,
	}

	return json.Marshal(resp)
}

// MarshalJSON converts the order cancel object into the format required by the
// bitfinex websocket service.
func (cr *CancelRequest) MarshalJSON() ([]byte, error) {
	b, err := cr.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[0, \"oc\", null, %s]", string(b))), nil
}
