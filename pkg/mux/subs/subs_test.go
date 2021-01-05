package subs_test

import (
	"testing"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/mux/subs"
	"github.com/stretchr/testify/assert"
)

func TestLimitReached(t *testing.T) {
	cases := map[string]struct {
		limit    int
		expected bool
		subs     []event.Subscribe
	}{
		"limit unreached": {
			limit:    20,
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
		},
		"limit reached": {
			limit:    2,
			expected: true,
			subs: []event.Subscribe{
				{Event: "foo"},
				{Event: "bar"},
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			s := subs.New()
			s.SubsLimit = v.limit
			for _, e := range v.subs {
				s.Add(e)
			}

			got := s.LimitReached()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestAdded(t *testing.T) {
	cases := map[string]struct {
		expected bool
		subs     []event.Subscribe
		pld      event.Subscribe
	}{
		"not added": {
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "bar"},
		},
		"added": {
			expected: true,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "foo"},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			s := subs.New()
			for _, e := range v.subs {
				s.Add(e)
			}

			got := s.Added(v.pld)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := map[string]struct {
		expected bool
		subs     []event.Subscribe
		pld      event.Subscribe
	}{
		"removed": {
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "foo"},
		},
		"not removed": {
			expected: true,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "bar"},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			s := subs.New()
			for _, e := range v.subs {
				s.Add(e)
			}

			s.Remove(v.pld)
			got := s.Added(v.subs[0])
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestGetAll(t *testing.T) {
	cases := map[string]struct {
		expected []event.Subscribe
		pld      []event.Subscribe
	}{
		"removed": {
			expected: []event.Subscribe{{Event: "foo"}, {Event: "bar"}},
			pld:      []event.Subscribe{{Event: "foo"}, {Event: "bar"}},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			s := subs.New()
			for _, e := range v.pld {
				s.Add(e)
			}

			got := s.GetAll()
			assert.Equal(t, v.expected, got)
		})
	}
}
