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

type messageFactory func(chanID int64, raw []interface{}) (interface{}, error)

type subscription struct {
	ChanID  int64
	pending bool
	Public  bool

	Request *SubscriptionRequest

	// heartbeats
	hbInterval time.Duration
	hbDeadline time.Time
	hbCancel   chan bool
	hbTimeout  chan time.Time
	hbReset    chan time.Time

	// to parent
	parentDisconnect chan error
}

func (s *subscription) Activate(t time.Time) {
	s.hbDeadline = t.Add(s.hbInterval)
	go s.timeoutHearbeat()
}

func (s *subscription) Heartbeat(t time.Time) {
	s.hbReset <- t
}

// timeout returns a cancel channel to cancel the heartbeat timeout created by this goroutine
func (s *subscription) timeoutHearbeat() {
	for {
		go func() {
			time.Sleep(s.hbInterval)
			s.hbTimeout <- time.Now()
		}()
		select {
		case t := <-s.hbReset: // heartbeat received, deadline moved
			s.hbDeadline = t.Add(s.hbInterval)
		case <-s.hbCancel: // underlying connection closed
			return
		case t := <-s.hbTimeout: // heartbeat timeout expired
			if t.After(s.hbDeadline) {
				s.parentDisconnect <- fmt.Errorf("channel %d timed out heartbeat %s", s.ChanID, s.hbInterval.String())
				return
			}
		}
	}
}

func (s *subscription) Close() {
	close(s.hbCancel)
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
		ChanID:           -1,
		Request:          request,
		pending:          true,
		Public:           isPublic(request),
		hbInterval:       interval,
		hbCancel:         make(chan bool),
		hbReset:          make(chan time.Time),
		hbTimeout:        make(chan time.Time),
		parentDisconnect: parentDisconnect,
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
		subsBySubID:         make(map[string]*subscription),
		subsByChanID:        make(map[int64]*subscription),
		hbTimeout:           heartbeatTimeout,
		hbParentDisconnect:  make(chan error),
		hbChannelDisconnect: make(chan error),
	}
	return subs
}

type subscriptions struct {
	lock sync.Mutex

	subsBySubID  map[string]*subscription // subscription map indexed by subscription ID
	subsByChanID map[int64]*subscription  // subscription map indexed by channel ID

	hbTimeout           time.Duration
	hbParentDisconnect  chan error // parent listens to this channel to receive disconnect events
	hbChannelDisconnect chan error // parent listens to this channel to recieve heartbeat timeouts
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

func (s *subscriptions) forwardDisconnects() {
	if e, ok := <-s.hbChannelDisconnect; ok {
		s.hbParentDisconnect <- e
	}
}

// Close will stop pending heartbeat timeouts. This is a terminal call and should only be called once.
func (s *subscriptions) Close() {
	s.close() // eat return
	close(s.hbParentDisconnect)
	close(s.hbChannelDisconnect)
}

func (s *subscriptions) close() []*subscription {
	s.lock.Lock()
	var subs []*subscription
	if len(s.subsBySubID) > 0 {
		subs = make([]*subscription, 0, len(s.subsBySubID))
		for _, sub := range s.subsBySubID {
			sub.Close()
			subs = append(subs, sub)
		}
		sort.Sort(SubscriptionSet(subs))
	}
	s.lock.Unlock()
	return subs
}

func (s *subscriptions) eatSiblingDisconnects() {
	for {
		select {
		case <-s.hbChannelDisconnect:
			// eat
		case <-time.After(s.hbTimeout):
			return // no longer hungry
		}
	}
}

// Reset clears all subscriptions from the currently managed list, and returns
// a slice of the existing subscriptions prior to reset.  Returns nil if no subscriptions exist.
func (s *subscriptions) Reset() []*subscription {
	subs := s.close()
	s.lock.Lock()
	s.subsBySubID = make(map[string]*subscription)
	s.subsByChanID = make(map[int64]*subscription)
	// drain any excess disconnect messages from last reset.
	// sibling channels may also send disconnects after the first,
	// which may disrupt the reconnect process. this could be improved
	s.eatSiblingDisconnects()

	go s.forwardDisconnects()
	s.lock.Unlock()
	return subs
}

// ListenDisconnect returns an error channel which receives a message when a heartbeat has expired a channel.
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
	sub.Close()
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
	sub.Close()
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
			log.Printf("activated subscription %s %s for channel %d", sub.Request.Channel, sub.Request.Symbol, chanID)
		}
		sub.pending = false
		sub.ChanID = chanID
		s.subsByChanID[chanID] = sub
		sub.Activate(time.Now())
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
		sub.Heartbeat(time.Now())
	}
	return fmt.Errorf("could not find channel ID to update heartbeat %d", chanID)
}
