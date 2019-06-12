package websocket

import (
	"fmt"
	"github.com/op/go-logging"
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
	DMS         int      `json:"dms,omitempty"` // dead man switch

	// unauthenticated
	Channel   string `json:"channel,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair,omitempty"`
}

const MaxChannels = 25

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

type HeartbeatDisconnect struct {
	Subscription *subscription
	Error        error
}

type UnsubscribeRequest struct {
	Event  string `json:"event"`
	ChanID int64  `json:"chanId"`
}

type subscription struct {
	ChanID     int64
	SocketId   SocketId
	pending    bool
	Public     bool

	Request    *SubscriptionRequest

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

func newSubscription(socketId SocketId, request *SubscriptionRequest) *subscription {
	return &subscription{
		ChanID:  -1,
		SocketId: socketId,
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

func newSubscriptions(heartbeatTimeout time.Duration, log *logging.Logger) *subscriptions {
	subs := &subscriptions{
		subsBySubID:  make(map[string]*subscription),
		subsByChanID: make(map[int64]*subscription),
		subsBySocketId: make(map[SocketId]SubscriptionSet),
		hbTimeout:    heartbeatTimeout,
		hbShutdown:   make(chan struct{}),
		hbDisconnect: make(chan HeartbeatDisconnect),
		hbSleep:      heartbeatTimeout / time.Duration(4),
		log:          log,
		lock:         &sync.RWMutex{},
	}
	go subs.control()
	return subs
}

// nolint
type heartbeat struct {
	ChanID int64
	*time.Time
}

type subscriptions struct {
	lock         *sync.RWMutex
	log          *logging.Logger

	subsBySocketId map[SocketId]SubscriptionSet // subscripts map indexed by socket id
	subsBySubID  map[string]*subscription // subscription map indexed by subscription ID
	subsByChanID map[int64]*subscription  // subscription map indexed by channel ID

	hbActive     bool
	hbDisconnect chan HeartbeatDisconnect // disconnect parent due to heartbeat timeout
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
func (s SubscriptionSet) RemoveByChannelId(chanId int64) SubscriptionSet {
	rIndex := -1
	for i, sub := range s {
		if sub.ChanID == chanId {
			rIndex = i
			break
		}
	}
	if rIndex >= 0 {
		return append(s[:rIndex], s[rIndex+1:]...)
	}
	return s
}

func (s SubscriptionSet) RemoveBySubscriptionId(subID string) SubscriptionSet {
	rIndex := -1
	for i, sub := range s {
		if sub.SubID() == subID {
			rIndex = i
			break
		}
	}
	if rIndex >= 0 {
		return append(s[:rIndex], s[rIndex+1:]...)
	}
	return s
}

func (s *subscriptions) heartbeat(chanID int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sub, ok := s.subsByChanID[chanID]; ok {
		sub.hbDeadline = time.Now().Add(s.hbTimeout)
	}
}

func (s *subscriptions) sweep(exp time.Time) {
	s.lock.RLock()
	if !s.hbActive {
		s.lock.RUnlock()
		return
	}
	disconnects := make([]HeartbeatDisconnect, 0)
	for _, sub := range s.subsByChanID {
		if exp.After(sub.hbDeadline) {
			s.hbActive = false
			hbErr := HeartbeatDisconnect{
				Subscription: sub,
				Error: fmt.Errorf("heartbeat disconnect on channel %d expired at %s (%s timeout)", sub.ChanID, sub.hbDeadline, s.hbTimeout),
			}
			disconnects = append(disconnects, hbErr)
		}
	}
	s.lock.RUnlock()
	for _, dis := range disconnects {
		s.hbDisconnect <- dis
	}
}

func (s *subscriptions) control() {
	for {
		select {
		case <-s.hbShutdown:
			return
		default:
		}
		s.sweep(time.Now())
		time.Sleep(s.hbSleep)
	}
}

// Close is terminal. Do not call heartbeat after close.
func (s *subscriptions) Close() {
	s.ResetAll()
	close(s.hbShutdown)
}


// Reset clears all subscriptions assigned to the given socket ID, and returns
// a slice of the existing subscriptions prior to reset
func (s *subscriptions) ResetSocketSubscriptions(socketId SocketId) []*subscription {
	var retSubs []*subscription
	s.lock.Lock()
	if set, ok := s.subsBySocketId[socketId]; ok {
		for _, sub := range set {
			retSubs = append(retSubs, sub)
			// remove from chanId array
			delete(s.subsByChanID, sub.ChanID)
			// remove from subId array
			delete(s.subsBySubID, sub.SubID())
		}
	}
	s.subsBySocketId[socketId] = make(SubscriptionSet, 0)
	s.lock.Unlock()
	return retSubs
}

// Removes all tracked subscriptions
func (s *subscriptions) ResetAll() {
	s.lock.Lock()
	s.subsBySubID = make(map[string]*subscription)
	s.subsByChanID = make(map[int64]*subscription)
	s.subsBySocketId = make(map[SocketId]SubscriptionSet)
	s.lock.Unlock()
}

// ListenDisconnect returns an error channel which receives a message when a heartbeat has expired a channel.
func (s *subscriptions) ListenDisconnect() <-chan HeartbeatDisconnect {
	return s.hbDisconnect
}

func (s *subscriptions) add(socketId SocketId, sub *SubscriptionRequest) *subscription {
	s.lock.Lock()
	defer s.lock.Unlock()
	subscription := newSubscription(socketId, sub)
	s.subsBySubID[sub.SubID] = subscription
	if _, ok := s.subsBySocketId[socketId]; !ok {
		s.subsBySocketId[socketId] = make(SubscriptionSet, 0)
	}
	s.subsBySocketId[socketId] = append(s.subsBySocketId[socketId], subscription)
	return subscription
}

func (s *subscriptions) removeByChannelID(chanID int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	// remove from socketId map
	sub, ok := s.subsByChanID[chanID]
	if !ok {
		return fmt.Errorf("could not find channel ID %d", chanID)
	}
	delete(s.subsByChanID, chanID)
	delete(s.subsBySubID, sub.SubID())
	// remove from socket map
	if _, ok := s.subsBySocketId[sub.SocketId]; ok {
		s.subsBySocketId[sub.SocketId] = s.subsBySocketId[sub.SocketId].RemoveByChannelId(chanID)
	}
	return nil
}

// nolint:megacheck
func (s *subscriptions) removeBySubscriptionID(subID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sub, ok := s.subsBySubID[subID]
	if !ok {
		return fmt.Errorf("could not find subscription ID %s to remove", subID)
	}
	// exists, remove both indices
	delete(s.subsBySubID, subID)
	delete(s.subsByChanID, sub.ChanID)
	// remove from socket map
	if _, ok := s.subsBySocketId[sub.SocketId]; ok {
		s.subsBySocketId[sub.SocketId] = s.subsBySocketId[sub.SocketId].RemoveBySubscriptionId(subID)
	}
	return nil
}

func (s *subscriptions) activate(subID string, chanID int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if sub, ok := s.subsBySubID[subID]; ok {
		if chanID != 0 {
			s.log.Info("activated subscription %s %s for channel %d", sub.Request.Channel, sub.Request.Symbol, chanID)
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
	s.lock.RLock()
	defer s.lock.RUnlock()
	if sub, ok := s.subsByChanID[chanID]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("could not find subscription for channel ID %d", chanID)
}

func (s *subscriptions) lookupBySubscriptionID(subID string) (*subscription, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if sub, ok := s.subsBySubID[subID]; ok {
		return sub, nil
	}
	return nil, fmt.Errorf("could not find subscription ID %s", subID)
}

func (s *subscriptions) lookupBySocketId(socketId SocketId) (*SubscriptionSet, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if set, ok := s.subsBySocketId[socketId]; ok {
		return &set, nil
	}
	return nil, fmt.Errorf("could not find subscription with socketId %d", socketId)
}
