package bitfinex

import (
	"encoding/json"
	"testing"
)

func TestMarginInfo_Unmarshal(t *testing.T) {
	// data from an actual Bitfinex's response
	data := []byte(`[{"margin_balance":"3168.43615333","tradable_balance":"-4105.85460135","unrealized_pl":"-174.0072","unrealized_swap":"-3.77879387","net_value":"2990.65015946","required_margin":"1344.0","leverage":"2.5","margin_requirement":"13.0","margin_limits":[{"on_pair":"BTCUSD","initial_margin":"30.0","margin_requirement":"15.0","tradable_balance":"-2544.520981132333333333"},{"on_pair":"LTCUSD","initial_margin":"30.0","margin_requirement":"15.0","tradable_balance":"-2544.520981132333333333"},{"on_pair":"LTCBTC","initial_margin":"30.0","margin_requirement":"15.0","tradable_balance":"-1634.516821132333333333"}],"message":"Margin requirement, leverage and tradable balance are now per pair. Values displayed in the root of the JSON message are incorrect (deprecated). You will find the correct ones under margin_limits, for each pair. Please update your code as soon as possible."}]`)

	var v []MarginInfo
	err := json.Unmarshal(data, &v)

	if err != nil {
		t.Fatalf("Failed unmarshaling MarginInfo: %s", err.Error())
	}
}
