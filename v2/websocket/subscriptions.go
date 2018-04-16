package websocket

import (
	"fmt"
	"log"
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

type subscription struct {
	ChanID  int64
	pending bool
	Public  bool

	Request *SubscriptionRequest

	hbDeadline time.Time
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

func newSubscription(request *SubscriptionRequest) *subscription {
	return &subscription{
		ChanID:  -1,
		Request: request,
		pending: true,
		Public:  isPublic(request),
	}
}

func (s subscription) SubID() string {
	return s.Request.SubID
}

func (s subscription) Pending() bool {
	return s.pending
}

func newSubscriptions(heartbeatTimeout time.Duration) *subscriptions {
	subs := &subscriptions{
		subsBySubID:  make(map[string]*subscription),
		subsByChanID: make(map[int64]*subscription),
		hbTimeout:    heartbeatTimeout,
		hbShutdown:   make(chan struct{}),
		hbDisconnect: make(chan error),
		hbSleep:      heartbeatTimeout / time.Duration(4),
	}
	go subs.control()
	return subs
}

type heartbeat struct {
	ChanID int64
	*time.Time
}

type subscriptions struct {
	lock sync.Mutex

	subsBySubID  map[string]*subscription // subscription map indexed by subscription ID
	subsByChanID map[int64]*subscription  // subscription map indexed by channel ID

	hbActive     bool
	hbDisconnect chan error // disconnect parent due to heartbeat timeout
	hbTimeout    time.Duration
	hbSleep      time.Duration
	hbShutdown   chan struct{}
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

func (s *subscriptions) heartbeat(chanID int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsByChanID[chanID]; ok {
		sub.hbDeadline = time.Now().Add(s.hbTimeout)
	}
}

func (s *subscriptions) sweep(exp time.Time) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.hbActive {
		return nil
	}
	for _, sub := range s.subsByChanID {
		if exp.After(sub.hbDeadline) {
			s.hbActive = false
			return fmt.Errorf("heartbeat disconnect on channel %d expired at %s (%s timeout)", sub.ChanID, sub.hbDeadline, s.hbTimeout)
		}
	}
	return nil
}

func (s *subscriptions) control() {
	for {
		select {
		case <-s.hbShutdown:
			return
		default:
		}
		if err := s.sweep(time.Now()); err != nil {
			s.hbDisconnect <- err
		}
		time.Sleep(s.hbSleep)
	}
}

// Close is terminal. Do not call heartbeat after close.
func (s *subscriptions) Close() {
	s.reset()
	close(s.hbShutdown)
}

func (s *subscriptions) reset() []*subscription {
	s.lock.Lock()
	var subs []*subscription
	if len(s.subsBySubID) > 0 {
		subs = make([]*subscription, 0, len(s.subsBySubID))
		for _, sub := range s.subsBySubID {
			subs = append(subs, sub)
		}
		sort.Sort(SubscriptionSet(subs))
	}
	s.lock.Unlock()
	return subs
}

// Reset clears all subscriptions from the currently managed list, and returns
// a slice of the existing subscriptions prior to reset.  Returns nil if no subscriptions exist.
func (s *subscriptions) Reset() []*subscription {
	subs := s.reset()
	s.lock.Lock()
	s.hbActive = false
	s.subsBySubID = make(map[string]*subscription)
	s.subsByChanID = make(map[int64]*subscription)
	s.lock.Unlock()
	return subs
}

// ListenDisconnect returns an error channel which receives a message when a heartbeat has expired a channel.
func (s *subscriptions) ListenDisconnect() <-chan error {
	return s.hbDisconnect
}

func (s *subscriptions) add(sub *SubscriptionRequest) *subscription {
	s.lock.Lock()
	defer s.lock.Unlock()
	subscription := newSubscription(sub)
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
		if chanID != 0 {
			//log.Printf("%#v", sub.Request)
			log.Printf("activated subscription %s %s for channel %d", sub.Request.Channel, sub.Request.Symbol, chanID)
		}
		sub.pending = false
		sub.ChanID = chanID
		sub.hbDeadline = time.Now().Add(s.hbTimeout)
		s.subsByChanID[chanID] = sub
		s.hbActive = true
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
