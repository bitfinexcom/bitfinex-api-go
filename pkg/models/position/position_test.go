package position_test

import (
	"reflect"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/stretchr/testify/assert"
)

func TestFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *position.Position
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"retrieve positions": {
			pld: []interface{}{
				"tBTCUSD", "ACTIVE", 0.0195, 8565.0267019, 0, 0, -0.33455568705000516, -0.0003117550117425625,
				7045.876419249083, 3.0673001895895604, nil, 142355652, 1574002216000, 1574002216000, nil, 0, nil, 0, 0,
				map[string]interface{}{"reason": "TRADE", "order_id": 34271018124, "liq_stage": nil, "trade_price": "8565.0267019", "trade_amount": "0.0195", "order_id_oppo": 34277498022},
			},
			expected: &position.Position{
				Id:                   142355652,
				Symbol:               "tBTCUSD",
				Status:               "ACTIVE",
				Amount:               0.0195,
				BasePrice:            8565.0267019,
				MarginFunding:        0,
				MarginFundingType:    0,
				ProfitLoss:           -0.33455568705000516,
				ProfitLossPercentage: -0.0003117550117425625,
				LiquidationPrice:     7045.876419249083,
				Leverage:             3.0673001895895604,
				Flag:                 nil,
				MtsCreate:            1574002216000,
				MtsUpdate:            1574002216000,
				Type:                 "",
				Collateral:           0,
				CollateralMin:        0,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      34271018124,
					"order_id_oppo": 34277498022,
					"reason":        "TRADE",
					"trade_amount":  "0.0195",
					"trade_price":   "8565.0267019",
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"claim position": {
			pld: []interface{}{
				"tBTCUSD", "ACTIVE", -0.001, 10119, 0, 0, nil, nil, nil, nil, nil, 142031891,
				1568650294000, 1568650294000, nil, 0, nil, 0, nil,
				map[string]interface{}{
					"reason":        "TRADE",
					"order_id":      31089337681,
					"liq_stage":     nil,
					"trade_price":   "10119.0",
					"trade_amount":  "-0.001",
					"user_id_oppo":  183923,
					"order_id_oppo": 31089692490,
				},
			},
			expected: &position.Position{
				Id:                   142031891,
				Symbol:               "tBTCUSD",
				Status:               "ACTIVE",
				Amount:               -0.001,
				BasePrice:            10119,
				MarginFunding:        0,
				MarginFundingType:    0,
				ProfitLoss:           0,
				ProfitLossPercentage: 0,
				LiquidationPrice:     0,
				Leverage:             0,
				Flag:                 nil,
				MtsCreate:            1568650294000,
				MtsUpdate:            1568650294000,
				Type:                 "",
				Collateral:           0,
				CollateralMin:        0,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      31089337681,
					"order_id_oppo": 31089692490,
					"reason":        "TRADE",
					"trade_amount":  "-0.001",
					"trade_price":   "10119.0",
					"user_id_oppo":  183923,
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"position history item": {
			pld: []interface{}{
				"tBTCUSD", "CLOSED", 0, 8639.50216298, 0, 0, nil, nil, nil, nil, nil,
				142355652, 1574002216000, 1574003099000, nil, nil, nil, nil, nil, nil,
			},
			expected: &position.Position{
				Id:                   142355652,
				Symbol:               "tBTCUSD",
				Status:               "CLOSED",
				Amount:               0,
				BasePrice:            8639.50216298,
				MarginFunding:        0,
				MarginFundingType:    0,
				ProfitLoss:           0,
				ProfitLossPercentage: 0,
				LiquidationPrice:     0,
				Leverage:             0,
				Flag:                 nil,
				MtsCreate:            1574002216000,
				MtsUpdate:            1574003099000,
				Type:                 "",
				Collateral:           0,
				CollateralMin:        0,
				Meta:                 nil,
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"position audit": {
			pld: []interface{}{
				"tETHUSD", "ACTIVE", 0.9682, 182.76, 0, 0, nil, nil, nil, nil, nil,
				142358085, 1574088504000, 1574088504000, nil, nil, nil, 0, nil,
				map[string]interface{}{
					"reason":        "TRADE",
					"order_id":      34319724780,
					"liq_stage":     nil,
					"trade_price":   "182.76",
					"trade_amount":  "0.9682",
					"user_id_oppo":  411107,
					"order_id_oppo": 34326729830,
				},
			},
			expected: &position.Position{
				Id:                   142358085,
				Symbol:               "tETHUSD",
				Status:               "ACTIVE",
				Amount:               0.9682,
				BasePrice:            182.76,
				MarginFunding:        0,
				MarginFundingType:    0,
				ProfitLoss:           0,
				ProfitLossPercentage: 0,
				LiquidationPrice:     0,
				Leverage:             0,
				Flag:                 nil,
				MtsCreate:            1574088504000,
				MtsUpdate:            1574088504000,
				Type:                 "",
				Collateral:           0,
				CollateralMin:        0,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      34319724780,
					"order_id_oppo": 34326729830,
					"reason":        "TRADE",
					"trade_amount":  "0.9682",
					"trade_price":   "182.76",
					"user_id_oppo":  411107,
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws pn pu pc": {
			pld: []interface{}{
				"tETHUST", "ACTIVE", 0.2, 153.71, 0, 0, -0.07944800000000068, -0.05855181835925015,
				67.52755254906451, 1.409288545397275, nil, 142420429, nil, nil, nil, 0, nil, 0, 0,
				map[string]interface{}{
					"reason":        "TRADE",
					"order_id":      34934099168,
					"order_id_oppo": 34934090814,
					"liq_stage":     nil,
					"trade_price":   "153.71",
					"trade_amount":  "0.2",
				},
			},
			expected: &position.Position{
				Id:                   142420429,
				Symbol:               "tETHUST",
				Status:               "ACTIVE",
				Amount:               0.2,
				BasePrice:            153.71,
				MarginFunding:        0,
				MarginFundingType:    0,
				ProfitLoss:           -0.07944800000000068,
				ProfitLossPercentage: -0.05855181835925015,
				LiquidationPrice:     67.52755254906451,
				Leverage:             1.409288545397275,
				Flag:                 nil,
				MtsCreate:            0,
				MtsUpdate:            0,
				Type:                 "",
				Collateral:           0,
				CollateralMin:        0,
				Meta: map[string]interface{}{
					"liq_stage":     nil,
					"order_id":      34934099168,
					"order_id_oppo": 34934090814,
					"reason":        "TRADE",
					"trade_amount":  "0.2",
					"trade_price":   "153.71",
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := position.FromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSnapshotFromRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []interface{}
		expected *position.Snapshot
		err      func(*testing.T, error)
	}{
		"invalid pld": {
			pld:      []interface{}{},
			expected: nil,
			err: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		"rest positions snapshot": {
			pld: []interface{}{
				[]interface{}{
					"tETHUSD", "ACTIVE", -0.2, 167.01, 0, 0, nil, nil, nil, nil, nil,
					142661142, 1579552390000, 1579552390000, nil, nil, nil, nil, nil, nil,
				},
				[]interface{}{
					"tETHUSD", "ACTIVE", -0.2, 167.01, 0, 0, nil, nil, nil, nil, nil,
					142661143, 1579552390000, 1579552390000, nil, nil, nil, nil, nil, nil,
				},
			},
			expected: &position.Snapshot{
				Snapshot: []*position.Position{
					{
						Id:        142661142,
						Symbol:    "tETHUSD",
						Status:    "ACTIVE",
						Amount:    -0.2,
						BasePrice: 167.01,
						MtsCreate: 1579552390000,
						MtsUpdate: 1579552390000,
						Type:      "ps",
					},
					{
						Id:        142661143,
						Symbol:    "tETHUSD",
						Status:    "ACTIVE",
						Amount:    -0.2,
						BasePrice: 167.01,
						MtsCreate: 1579552390000,
						MtsUpdate: 1579552390000,
						Type:      "ps",
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		"ws positions snapshot": {
			pld: []interface{}{
				[]interface{}{
					"tETHUST", "ACTIVE", 0.2, 153.71, 0, 0, nil, nil, nil,
					nil, nil, 142420429, nil, nil, nil, 0, nil, 0, nil,
					map[string]interface{}{
						"reason":        "TRADE",
						"order_id":      34934099168,
						"order_id_oppo": 34934090814,
						"liq_stage":     nil,
						"trade_price":   "153.71",
						"trade_amount":  "0.2",
					},
				},
			},
			expected: &position.Snapshot{
				Snapshot: []*position.Position{
					{
						Id:        142420429,
						Symbol:    "tETHUST",
						Status:    "ACTIVE",
						Amount:    0.2,
						BasePrice: 153.71,
						Type:      "ps",
						Meta: map[string]interface{}{
							"reason":        "TRADE",
							"order_id":      34934099168,
							"order_id_oppo": 34934090814,
							"liq_stage":     nil,
							"trade_price":   "153.71",
							"trade_amount":  "0.2",
						},
					},
				},
			},
			err: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			got, err := position.SnapshotFromRaw(v.pld)
			v.err(t, err)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestNewFromRaw(t *testing.T) {
	pld := []interface{}{
		"tETHUST", "ACTIVE", 0.2, 153.71, 0, 0, -0.07944800000000068, -0.05855181835925015,
		67.52755254906451, 1.409288545397275, nil, 142420429, nil, nil, nil, 0, nil, 0, 0,
		map[string]interface{}{
			"reason":        "TRADE",
			"order_id":      34934099168,
			"order_id_oppo": 34934090814,
			"liq_stage":     nil,
			"trade_price":   "153.71",
			"trade_amount":  "0.2",
		},
	}

	expected := "position.New"
	p, err := position.NewFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(p).String()
	assert.Equal(t, expected, got)
	assert.Equal(t, "pn", p.Type)
}

func TestUpdateFromRaw(t *testing.T) {
	pld := []interface{}{
		"tETHUST", "ACTIVE", 0.2, 153.71, 0, 0, -0.07944800000000068, -0.05855181835925015,
		67.52755254906451, 1.409288545397275, nil, 142420429, nil, nil, nil, 0, nil, 0, 0,
		map[string]interface{}{
			"reason":        "TRADE",
			"order_id":      34934099168,
			"order_id_oppo": 34934090814,
			"liq_stage":     nil,
			"trade_price":   "153.71",
			"trade_amount":  "0.2",
		},
	}

	expected := "position.Update"
	p, err := position.UpdateFromRaw(pld)
	assert.Nil(t, err)

	got := reflect.TypeOf(p).String()
	assert.Equal(t, expected, got)
	assert.Equal(t, "pu", p.Type)
}
