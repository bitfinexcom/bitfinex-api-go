package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrdersAll(t *testing.T) {
	httpDo := func(_ *http.Client, req *http.Request) (*http.Response, error) {
		msg := `
				[
					[4419360502,null,83283216761,"tIOTBTC",1508281683000,1508281731000,63938,63938,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.0000843,0,0,0,null,null,null,0,0,null],
					[4419354239,null,83265164211,"tIOTBTC",1508281665000,1508281674000,63976,63976,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.00008425,0,0,0,null,null,null,0,0,null],
					[4419339620,null,83217673277,"tIOTBTC",1508281618000,1508281653000,64014,64014,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.0000842,0,0,0,null,null,null,0,0,null]
				]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders, err := NewClientWithHttpDo(httpDo).Orders.All()

	if err != nil {
		t.Error(err)
	}

	if len(orders.Snapshot) != 3 {
		t.Fatalf("expected three orders but got %d", len(orders.Snapshot))
	}
}

func TestOrdersHistory(t *testing.T) {
	httpDo := func(_ *http.Client, req *http.Request) (*http.Response, error) {
		msg := `
				[
					[4419360502,null,83283216761,"tIOTBTC",1508281683000,1508281731000,63938,63938,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.0000843,0,0,0,null,null,null,0,0,null],
					[4419354239,null,83265164211,"tIOTBTC",1508281665000,1508281674000,63976,63976,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.00008425,0,0,0,null,null,null,0,0,null],
					[4419339620,null,83217673277,"tIOTBTC",1508281618000,1508281653000,64014,64014,"EXCHANGE LIMIT",null,null,null,null,"CANCELED",null,null,0.0000842,0,0,0,null,null,null,0,0,null]
				]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	orders, err := NewClientWithHttpDo(httpDo).Orders.GetHistoryBySymbol(bitfinex.TradingPrefix + bitfinex.IOTBTC)

	if err != nil {
		t.Error(err)
	}

	if len(orders.Snapshot) != 3 {
		t.Errorf("expected three orders but got %d", len(orders.Snapshot))
	}
}

func TestCancelOrderMulti(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/cancel/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			gotReqPld := CancelOrderMultiArgs{}
			err := json.NewDecoder(r.Body).Decode(&gotReqPld)
			require.Nil(t, err)

			expectedReqPld := CancelOrderMultiArgs{
				OrderIDs:      OrderIDs{123},
				GroupOrderIDs: GroupOrderIDs{234},
				All:           1,
			}
			assert.Equal(t, expectedReqPld, gotReqPld)

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err = w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pld := CancelOrderMultiArgs{
			OrderIDs:      OrderIDs{123},
			GroupOrderIDs: GroupOrderIDs{234},
			All:           1,
		}

		rsp, err := c.Orders.CancelOrderMulti(pld)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestCancelOrdersMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		rsp, err := c.Orders.CancelOrdersMultiOp(OrderIDs{1189428429, 1189428430})
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}

func TestOrderMultiOp(t *testing.T) {
	t.Run("calls correct resource with correct payload", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/auth/w/order/multi", r.RequestURI)
			assert.Equal(t, "POST", r.Method)

			respMock := []interface{}{1568711312683, nil, nil, nil, nil, nil, nil, nil}
			payload, _ := json.Marshal(respMock)
			_, err := w.Write(payload)
			require.Nil(t, err)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		c := NewClientWithURL(server.URL)
		pld := OrderMultiArgs{
			Ops: [][]interface{}{
				{
					"oc",
					map[string]int{"id": 1189502430},
				},
				{
					"oc_multi",
					map[string][]int{"id": OrderIDs{1189502431, 1189502432}},
				},
			},
		}
		rsp, err := c.Orders.OrderMultiOp(pld)
		require.Nil(t, err)
		assert.Equal(t, int64(1568711312683), rsp.MTS)
	})
}
