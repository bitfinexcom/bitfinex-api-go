package subs

import "github.com/bitfinexcom/bitfinex-api-go/pkg/models/event"

const defaultLimit = 20

type Subs struct {
	Subs      map[event.Subscribe]bool
	SubsLimit int
}

// New returns pointer to instacne of Subscriptions
func New() *Subs {
	return &Subs{
		Subs:      make(map[event.Subscribe]bool),
		SubsLimit: defaultLimit,
	}
}

// LimitReached returns true if number of subs > subsLimit
func (s *Subs) LimitReached() bool {
	return len(s.Subs) == s.SubsLimit
}

// Added checks if given subscription is already added. Used to
// avoid duplicate subscriptions per client
func (s *Subs) Added(sub event.Subscribe) (res bool) {
	_, res = s.Subs[sub]
	return
}

// Add adds new subscription to the list
func (s *Subs) Add(sub event.Subscribe) {
	s.Subs[sub] = true
}

// Remove adds new subscription to the list
func (s *Subs) Remove(sub event.Subscribe) {
	delete(s.Subs, sub)
}

// GetAll returns all subscriptions
func (s *Subs) GetAll() (res []event.Subscribe) {
	for sub := range s.Subs {
		res = append(res, sub)
	}
	return
}
