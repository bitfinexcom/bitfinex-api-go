package websocket

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
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

func (s *SubscriptionRequest) String() string {
	if s.Key == "" && s.Channel != "" && s.Symbol != "" {
		return fmt.Sprintf("%s %s", s.Channel, s.Symbol)
	}
	if s.Channel != "" && s.Symbol != "" && s.Precision != "" && s.Frequency != "" {
		return fmt.Sprintf("%s %s %s %s", s.Channel, s.Symbol, s.Precision, s.Frequency)
	}
	if s.Channel != "" && s.Symbol != "" {
		return fmt.Sprintf("%s %s", s.Channel, s.Symbol)
	}
	return ""
}

type UnsubscribeRequest struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}

type messageFactory func(chanID int64, raw []interface{}) (interface{}, error)

type subscription struct {
	ChanID  int64
	pending bool
	Public  bool

	Request *SubscriptionRequest

	// heartbeat timer
	hbInterval time.Duration
	hb         *time.Timer
	die        chan bool

	parentDisconnect chan error
}

func (s *subscription) activate() {
	s.heartbeat()
}

func (s *subscription) timeout() {
	s.parentDisconnect <- fmt.Errorf("heartbeat timed out on channel %d", s.ChanID)
}

func (s *subscription) heartbeat() {
	if s.hb != nil {
		s.hb.Stop()
	}
	close(s.die) // terminate previous hb timeout
	s.die = make(chan bool)
	s.hb = time.AfterFunc(s.hbInterval, s.timeout)
}

func (s *subscription) stopHeartbeatTimeout() {
	if s.hb != nil {
		s.hb.Stop()
	}
}

func isPublic(request *SubscriptionRequest) bool {
	switch request.Channel {
	case ChanBook:
		return true
	case ChanCandles:
		return true
	case ChanTicker:
		return true
	case ChanTrades:
		return true
	}
	return false
}

func newSubscription(request *SubscriptionRequest, interval time.Duration, parentDisconnect chan error) *subscription {
	return &subscription{
		Request:          request,
		pending:          true,
		Public:           isPublic(request),
		die:              make(chan bool), // kill pending heartbeats
		hbInterval:       interval,
		parentDisconnect: parentDisconnect, // disconnect parent in hb timeout
	}
}

func (s subscription) SubID() string {
	return s.Request.SubID
}

func (s subscription) Pending() bool {
	return s.pending
}

func newSubscriptions(heartbeatTimeout time.Duration) *subscriptions {
	return &subscriptions{
		subsBySubID:         make(map[string]*subscription),
		subsByChanID:        make(map[int64]*subscription),
		hbTimeout:           heartbeatTimeout,
		hbParentDisconnect:  make(chan error),
		hbChannelDisconnect: make(chan error),
	}
}

type subscriptions struct {
	lock sync.Mutex

	subsBySubID  map[string]*subscription // subscription map indexed by subscription ID
	subsByChanID map[int64]*subscription  // subscription map indexed by channel ID

	hbTimeout           time.Duration
	hbChannelDisconnect chan error // message sent when a subscription fails a heartbeat check
	hbParentDisconnect  chan error // parent listens to this channel to receive disconnect events
}

func (s *subscriptions) Empty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.subsBySubID) == 0
}

// SubscriptionSet is a typed version of an array of subscription pointers, intended to meet the sortable interface.
// We need to sort Reset()'s return values for tests with more than 1 subscription (range map order is undefined)
type SubscriptionSet []*subscription

func (s SubscriptionSet) Len() int {
	return len(s)
}
func (s SubscriptionSet) Less(i, j int) bool {
	return strings.Compare(s[i].SubID(), s[j].SubID()) < 0
}
func (s SubscriptionSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Reset clears all subscriptions from the currently managed list, and returns
// a slice of the existing subscriptions prior to reset.  Returns nil if no subscriptions exist.
func (s *subscriptions) Reset() []*subscription {
	s.lock.Lock()
	defer s.lock.Unlock()
	var subs []*subscription
	if len(s.subsBySubID) > 0 {
		subs = make([]*subscription, 0, len(s.subsBySubID))
		for _, sub := range s.subsBySubID {
			sub.stopHeartbeatTimeout()
			subs = append(subs, sub)
		}
		sort.Sort(SubscriptionSet(subs))
	}
	close(s.hbChannelDisconnect)
	close(s.hbParentDisconnect)
	s.subsBySubID = make(map[string]*subscription)
	s.subsByChanID = make(map[int64]*subscription)
	s.hbChannelDisconnect = make(chan error)
	s.hbParentDisconnect = make(chan error)
	go s.listenHeartbeats()
	return subs
}

func (s *subscriptions) listenHeartbeats() {
	if err := <-s.hbChannelDisconnect; err != nil {
		s.hbParentDisconnect <- err
	}
}

func (s *subscriptions) ListenDisconnect() <-chan error {
	return s.hbParentDisconnect
}

func (s *subscriptions) add(sub *SubscriptionRequest) *subscription {
	s.lock.Lock()
	defer s.lock.Unlock()
	subscription := newSubscription(sub, s.hbTimeout, s.hbChannelDisconnect)
	s.subsBySubID[sub.SubID] = subscription
	return subscription
}

func (s *subscriptions) removeByChannelID(chanID int64) error {
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

func (s *subscriptions) removeBySubscriptionID(subID string) error {
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
		sub.activate()
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

func (s *subscriptions) heartbeat(chanID int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsByChanID[chanID]; ok {
		sub.heartbeat()
	}
	return fmt.Errorf("could not find channel ID to update heartbeat %d", chanID)
}
