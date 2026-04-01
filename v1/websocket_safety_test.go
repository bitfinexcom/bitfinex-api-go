package bitfinex

import (
	"encoding/json"
	"testing"
)

func TestHandleDataMessage_MalformedPayload(t *testing.T) {
	ws := NewWebSocketService(&Client{})

	tests := []struct {
		name string
		msg  string
	}{
		{"string instead of array", `"not an array"`},
		{"empty object", `{}`},
		{"null", `null`},
		{"nested object array", `[{"key":"val"}]`},
		{"chanID not float", `["abc", [1,2,3]]`},
		{"empty array", `[]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("handleDataMessage panicked on %s: %v", tt.name, r)
				}
			}()
			ws.handleDataMessage([]byte(tt.msg))
		})
	}
}

func TestHandleDataMessage_ValidPayload(t *testing.T) {
	ws := NewWebSocketService(&Client{})

	// Set up a channel mapped to chanID 1
	ch := make(chan []float64, 1)
	ws.chanMap[1.0] = ch

	msg, _ := json.Marshal([]float64{1.0, 100.5, 200.3})
	ws.handleDataMessage(msg)

	select {
	case data := <-ch:
		if len(data) != 2 {
			t.Errorf("expected 2 elements, got %d", len(data))
		}
		if data[0] != 100.5 {
			t.Errorf("expected 100.5, got %f", data[0])
		}
	default:
		t.Error("expected data on channel, got nothing")
	}
}

func TestHandleDataMessage_UnmappedChannel(t *testing.T) {
	ws := NewWebSocketService(&Client{})

	// No channel mapped for chanID 999 - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("handleDataMessage panicked on unmapped channel: %v", r)
		}
	}()

	msg, _ := json.Marshal([]float64{999.0, 100.5, 200.3})
	ws.handleDataMessage(msg)
}

func TestWebSocketService_CloseNilConnection(t *testing.T) {
	ws := NewWebSocketService(&Client{})
	// ws.ws is nil since we never connected

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close() panicked on nil connection: %v", r)
		}
	}()
	ws.Close()
}
