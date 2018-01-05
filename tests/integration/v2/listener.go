package tests

import (
	"errors"
	"log"
	"time"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

type listener struct {
	infoEvents           chan *websocket.InfoEvent
	authEvents           chan *websocket.AuthEvent
	ticks                chan *bitfinex.Ticker
	subscriptionEvents   chan *websocket.SubscribeEvent
	unsubscriptionEvents chan *websocket.UnsubscribeEvent
	walletUpdates        chan *bitfinex.WalletUpdate
	balanceUpdates       chan *bitfinex.BalanceUpdate
	walletSnapshot       chan *bitfinex.WalletSnapshot
	positionSnapshot     chan *bitfinex.PositionSnapshot
	notifications        chan *bitfinex.Notification
	positions            chan *bitfinex.PositionUpdate
	tradeUpdates         chan *bitfinex.TradeUpdate
	tradeExecutions      chan *bitfinex.TradeExecution
	cancels              chan *bitfinex.OrderCancel
	marginBase           chan *bitfinex.MarginInfoBase
	marginUpdate         chan *bitfinex.MarginInfoUpdate
	funding              chan *bitfinex.FundingInfo
	orderNew             chan *bitfinex.OrderNew
	errors               chan error
}

func newListener() *listener {
	return &listener{
		infoEvents:           make(chan *websocket.InfoEvent, 10),
		authEvents:           make(chan *websocket.AuthEvent, 10),
		ticks:                make(chan *bitfinex.Ticker, 10),
		subscriptionEvents:   make(chan *websocket.SubscribeEvent, 10),
		unsubscriptionEvents: make(chan *websocket.UnsubscribeEvent, 10),
		walletUpdates:        make(chan *bitfinex.WalletUpdate, 10),
		balanceUpdates:       make(chan *bitfinex.BalanceUpdate, 10),
		walletSnapshot:       make(chan *bitfinex.WalletSnapshot, 10),
		positionSnapshot:     make(chan *bitfinex.PositionSnapshot, 10),
		errors:               make(chan error, 10),
		notifications:        make(chan *bitfinex.Notification, 10),
		positions:            make(chan *bitfinex.PositionUpdate, 10),
		tradeUpdates:         make(chan *bitfinex.TradeUpdate, 10),
		tradeExecutions:      make(chan *bitfinex.TradeExecution, 10),
		cancels:              make(chan *bitfinex.OrderCancel, 10),
		marginBase:           make(chan *bitfinex.MarginInfoBase, 10),
		marginUpdate:         make(chan *bitfinex.MarginInfoUpdate, 10),
		orderNew:             make(chan *bitfinex.OrderNew, 10),
		funding:              make(chan *bitfinex.FundingInfo, 10),
	}
}

func (l *listener) nextInfoEvent() (*websocket.InfoEvent, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.infoEvents:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for InfoEvent")
	}
}

func (l *listener) nextAuthEvent() (*websocket.AuthEvent, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.authEvents:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for AuthEvent")
	}
}

func (l *listener) nextWalletUpdate() (*bitfinex.WalletUpdate, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.walletUpdates:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for WalletUpdate")
	}
}

func (l *listener) nextBalanceUpdate() (*bitfinex.BalanceUpdate, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.balanceUpdates:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for BalanceUpdate")
	}
}

func (l *listener) nextWalletSnapshot() (*bitfinex.WalletSnapshot, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.walletSnapshot:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for WalletSnapshot")
	}
}

func (l *listener) nextPositionSnapshot() (*bitfinex.PositionSnapshot, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.positionSnapshot:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for PositionSnapshot")
	}
}

func (l *listener) nextSubscriptionEvent() (*websocket.SubscribeEvent, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.subscriptionEvents:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for SubscribeEvent")
	}
}

func (l *listener) nextUnsubscriptionEvent() (*websocket.UnsubscribeEvent, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.unsubscriptionEvents:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for UnsubscribeEvent")
	}
}

func (l *listener) nextTick() (*bitfinex.Ticker, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.ticks:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for Ticker")
	}
}

func (l *listener) nextNotification() (*bitfinex.Notification, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.notifications:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for Notification")
	}
}

func (l *listener) nextTradeExecution() (*bitfinex.TradeExecution, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.tradeExecutions:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for TradeExecution")
	}
}

func (l *listener) nextPositionUpdate() (*bitfinex.PositionUpdate, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.positions:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for PositionUpdate")
	}
}

func (l *listener) nextTradeUpdate() (*bitfinex.TradeUpdate, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.tradeUpdates:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for TradeUpdate")
	}
}

func (l *listener) nextOrderCancel() (*bitfinex.OrderCancel, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.cancels:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for OrderCancel")
	}
}

func (l *listener) nextMarginInfoBase() (*bitfinex.MarginInfoBase, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.marginBase:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for MarginInfoBase")
	}
}

func (l *listener) nextMarginInfoUpdate() (*bitfinex.MarginInfoUpdate, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.marginUpdate:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for MarginInfoUpdate")
	}
}

func (l *listener) nextFundingInfo() (*bitfinex.FundingInfo, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.funding:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for FundingInfo")
	}
}

func (l *listener) nextOrderNew() (*bitfinex.OrderNew, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.orderNew:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for OrderNew")
	}
}

// strongly types messages and places them into a channel
func (l *listener) run(ch <-chan interface{}) {
	go func() {
		for {
			select {
			case msg := <-ch:
				if msg == nil {
					return
				}
				// remove threading guarantees when mulitplexing into channels
				log.Printf("listener raw: %#v", msg)
				switch msg.(type) {
				case error:
					l.errors <- msg.(error)
				case *bitfinex.Ticker:
					l.ticks <- msg.(*bitfinex.Ticker)
				case *websocket.InfoEvent:
					l.infoEvents <- msg.(*websocket.InfoEvent)
				case *websocket.SubscribeEvent:
					l.subscriptionEvents <- msg.(*websocket.SubscribeEvent)
				case *websocket.UnsubscribeEvent:
					l.unsubscriptionEvents <- msg.(*websocket.UnsubscribeEvent)
				case *websocket.AuthEvent:
					l.authEvents <- msg.(*websocket.AuthEvent)
				case *bitfinex.WalletUpdate:
					l.walletUpdates <- msg.(*bitfinex.WalletUpdate)
				case *bitfinex.BalanceUpdate:
					l.balanceUpdates <- msg.(*bitfinex.BalanceUpdate)
				case *bitfinex.Notification:
					l.notifications <- msg.(*bitfinex.Notification)
				case *bitfinex.TradeUpdate:
					l.tradeUpdates <- msg.(*bitfinex.TradeUpdate)
				case *bitfinex.TradeExecution:
					l.tradeExecutions <- msg.(*bitfinex.TradeExecution)
				case *bitfinex.PositionUpdate:
					l.positions <- msg.(*bitfinex.PositionUpdate)
				case *bitfinex.OrderCancel:
					l.cancels <- msg.(*bitfinex.OrderCancel)
				case *bitfinex.MarginInfoBase:
					l.marginBase <- msg.(*bitfinex.MarginInfoBase)
				case *bitfinex.MarginInfoUpdate:
					l.marginUpdate <- msg.(*bitfinex.MarginInfoUpdate)
				case *bitfinex.OrderNew:
					l.orderNew <- msg.(*bitfinex.OrderNew)
				case *bitfinex.FundingInfo:
					l.funding <- msg.(*bitfinex.FundingInfo)
				case *bitfinex.PositionSnapshot:
					l.positionSnapshot <- msg.(*bitfinex.PositionSnapshot)
				case *bitfinex.WalletSnapshot:
					l.walletSnapshot <- msg.(*bitfinex.WalletSnapshot)
				}
			}
		}
	}()
}
