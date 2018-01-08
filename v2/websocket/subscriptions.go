package websocket

import (
	"fmt"
	"sync"
)

type SubscriptionRequest struct {
	SubID string `json:"subId"`
	Event string `json:"event"`

	// authenticated
	APIKey      string   `json:"apiKey,omitempty"`
	AuthSig     string   `json:"authSig,omitempty"`
	AuthPayload string   `json:"authPayload,omitempty"`
	AuthNonce   string   `json:"authNonce,omitempty"`
	Filter      []string `json:"filter,omitempty"`

	// unauthenticated
	Channel   string `json:"channel,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair,omitempty"`
}

type UnsubscribeRequest struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}

type messageFactory func(chanID int64, raw []interface{}) (interface{}, error)

type subscription struct {
	ChanID  int64
	pending bool

	Request *SubscriptionRequest
}

func newSubscription(request *SubscriptionRequest) *subscription {
	return &subscription{
		Request: request,
		pending: true,
	}
}

func (s subscription) SubID() string {
	return s.Request.SubID
}

func (s subscription) Pending() bool {
	return s.pending
}

func newSubscriptions() *subscriptions {
	return &subscriptions{
		subsBySubID:  make(map[string]*subscription),
		subsByChanID: make(map[int64]*subscription),
	}
}

type subscriptions struct {
	lock sync.Mutex

	subsBySubID  map[string]*subscription // subscription map indexed by subscription ID
	subsByChanID map[int64]*subscription  // subscription map indexed by channel ID
}

func (s *subscriptions) add(sub *SubscriptionRequest) *subscription {
	s.lock.Lock()
	defer s.lock.Unlock()
	subscription := newSubscription(sub)
	s.subsBySubID[sub.SubID] = subscription
	return subscription
}

func (s *subscriptions) removeByChanID(chanID int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sub, ok := s.subsByChanID[chanID]
	if !ok {
		return fmt.Errorf("could not find channel ID %d", chanID)
	}
	delete(s.subsByChanID, chanID)
	if _, ok = s.subsBySubID[sub.SubID()]; ok {
		delete(s.subsBySubID, sub.SubID())
	}
	return nil
}

func (s *subscriptions) removeBySubID(subID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sub, ok := s.subsBySubID[subID]
	if !ok {
		return fmt.Errorf("could not find subscription ID %s to remove", subID)
	}
	// exists, remove both indices
	delete(s.subsBySubID, subID)
	if _, ok = s.subsByChanID[sub.ChanID]; ok {
		delete(s.subsByChanID, sub.ChanID)
	}
	return nil
}

func (s *subscriptions) activate(subID string, chanID int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsBySubID[subID]; ok {
		sub.pending = false
		sub.ChanID = chanID
		s.subsByChanID[chanID] = sub
		return nil
	}
	return fmt.Errorf("could not find subscription ID %s to activate", subID)
}

func (s *subscriptions) lookupByChannelID(chanID int64) (*subscription, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsByChanID[chanID]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("could not find subscription for channel ID %d", chanID)
}

func (s *subscriptions) lookupBySubscriptionID(subID string) (*subscription, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsBySubID[subID]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("could not find subscription ID %s", subID)
}
