package msg_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/msg"
	"github.com/stretchr/testify/assert"
)

func TestIsEvent(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected bool
	}{
		"event type": {
			pld:      []byte(`{}`),
			expected: true,
		},
		"not event type": {
			pld:      []byte(`[]`),
			expected: false,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{
				Data: v.pld,
			}

			got := m.IsEvent()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestIsRaw(t *testing.T) {
	cases := map[string]struct {
		pld      []byte
		expected bool
	}{
		"raw type": {
			pld:      []byte(`[]`),
			expected: true,
		},
		"not raw type": {
			pld:      []byte(`{}`),
			expected: false,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			m := msg.Msg{
				Data: v.pld,
			}

			got := m.IsRaw()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	m := msg.Msg{
		Data: []byte(`{"event":"info","version":2,"serverId":"dbea77ee-4740-4a82-84f3-c6bc1b5abb9a","platform":{"status":1}}`),
	}

	expected := event.Info{
		Subscribe: event.Subscribe{
			Event: "info",
		},
		Version:  2,
		ServerID: "dbea77ee-4740-4a82-84f3-c6bc1b5abb9a",
		Platform: event.Platform{Status: 1},
	}

	got, err := m.ProcessEvent()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
