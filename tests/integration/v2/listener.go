package tests

import (
	"errors"
	"log"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/balanceinfo"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/fundinginfo"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/margin"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/order"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/position"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/ticker"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecution"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/tradeexecutionupdate"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/wallet"
	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v2"
	"github.com/bitfinexcom/bitfinex-api-go/v2/websocket"
)

type listener struct {
	infoEvents           chan *websocket.InfoEvent
	authEvents           chan *websocket.AuthEvent
	ticks                chan *ticker.Ticker
	subscriptionEvents   chan *websocket.SubscribeEvent
	unsubscriptionEvents chan *websocket.UnsubscribeEvent
	walletUpdates        chan *wallet.Update
	balanceUpdates       chan *balanceinfo.Update
	walletSnapshot       chan *wallet.Snapshot
	positionSnapshot     chan *position.Snapshot
	notifications        chan *bitfinex.Notification
	positions            chan *position.Update
	tradeUpdates         chan *tradeexecutionupdate.TradeExecutionUpdate
	tradeExecutions      chan *tradeexecution.TradeExecution
	cancels              chan *order.Cancel
	marginBase           chan *margin.InfoBase
	marginUpdate         chan *margin.InfoUpdate
	funding              chan *fundinginfo.FundingInfo
	orderNew             chan *order.New
	orderUpdate          chan *order.Update
	errors               chan error
}

func newListener() *listener {
	return &listener{
		infoEvents:           make(chan *websocket.InfoEvent, 10),
		authEvents:           make(chan *websocket.AuthEvent, 10),
		ticks:                make(chan *ticker.Ticker, 10),
		subscriptionEvents:   make(chan *websocket.SubscribeEvent, 10),
		unsubscriptionEvents: make(chan *websocket.UnsubscribeEvent, 10),
		walletUpdates:        make(chan *wallet.Update, 10),
		balanceUpdates:       make(chan *balanceinfo.Update, 10),
		walletSnapshot:       make(chan *wallet.Snapshot, 10),
		positionSnapshot:     make(chan *position.Snapshot, 10),
		errors:               make(chan error, 10),
		notifications:        make(chan *bitfinex.Notification, 10),
		positions:            make(chan *position.Update, 10),
		tradeUpdates:         make(chan *tradeexecutionupdate.TradeExecutionUpdate, 10),
		tradeExecutions:      make(chan *tradeexecution.TradeExecution, 10),
		cancels:              make(chan *order.Cancel, 10),
		marginBase:           make(chan *margin.InfoBase, 10),
		marginUpdate:         make(chan *margin.InfoUpdate, 10),
		orderNew:             make(chan *order.New, 10),
		orderUpdate:          make(chan *order.Update, 10),
		funding:              make(chan *fundinginfo.FundingInfo, 10),
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

func (l *listener) nextWalletUpdate() (*wallet.Update, error) {
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

func (l *listener) nextBalanceUpdate() (*balanceinfo.Update, error) {
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

func (l *listener) nextWalletSnapshot() (*wallet.Snapshot, error) {
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

func (l *listener) nextPositionSnapshot() (*position.Snapshot, error) {
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

func (l *listener) nextTick() (*ticker.Ticker, error) {
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

func (l *listener) nextTradeExecution() (*tradeexecution.TradeExecution, error) {
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

func (l *listener) nextPositionUpdate() (*position.Update, error) {
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

func (l *listener) nextTradeUpdate() (*tradeexecutionupdate.TradeExecutionUpdate, error) {
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

func (l *listener) nextOrderCancel() (*order.Cancel, error) {
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

func (l *listener) nextMarginInfoBase() (*margin.InfoBase, error) {
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

func (l *listener) nextMarginInfoUpdate() (*margin.InfoUpdate, error) {
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

func (l *listener) nextFundingInfo() (*fundinginfo.FundingInfo, error) {
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

func (l *listener) nextOrderNew() (*order.New, error) {
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

func (l *listener) nextOrderUpdate() (*order.Update, error) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		close(timeout)
	}()
	select {
	case ev := <-l.orderUpdate:
		return ev, nil
	case <-timeout:
		return nil, errors.New("timed out waiting for OrderUpdate")
	}
}

// strongly types messages and places them into a channel
func (l *listener) run(ch <-chan interface{}) {
	go func() {
		// nolint:megacheck
		for {
			select {
			case msg := <-ch:
				if msg == nil {
					return
				}
				// remove threading guarantees when mulitplexing into channels
				log.Printf("[DEBUG] WsService -> WsClient: %#v", msg)
				switch msg.(type) {
				case error:
					l.errors <- msg.(error)
				case *ticker.Ticker:
					l.ticks <- msg.(*ticker.Ticker)
				case *websocket.InfoEvent:
					l.infoEvents <- msg.(*websocket.InfoEvent)
				case *websocket.SubscribeEvent:
					l.subscriptionEvents <- msg.(*websocket.SubscribeEvent)
				case *websocket.UnsubscribeEvent:
					l.unsubscriptionEvents <- msg.(*websocket.UnsubscribeEvent)
				case *websocket.AuthEvent:
					l.authEvents <- msg.(*websocket.AuthEvent)
				case *wallet.Update:
					l.walletUpdates <- msg.(*wallet.Update)
				case *balanceinfo.Update:
					l.balanceUpdates <- msg.(*balanceinfo.Update)
				case *bitfinex.Notification:
					l.notifications <- msg.(*bitfinex.Notification)
				case *tradeexecutionupdate.TradeExecutionUpdate:
					l.tradeUpdates <- msg.(*tradeexecutionupdate.TradeExecutionUpdate)
				case *tradeexecution.TradeExecution:
					l.tradeExecutions <- msg.(*tradeexecution.TradeExecution)
				case *position.Update:
					l.positions <- msg.(*position.Update)
				case *order.Cancel:
					l.cancels <- msg.(*order.Cancel)
				case *margin.InfoBase:
					l.marginBase <- msg.(*margin.InfoBase)
				case *margin.InfoUpdate:
					l.marginUpdate <- msg.(*margin.InfoUpdate)
				case *order.New:
					l.orderNew <- msg.(*order.New)
				case *order.Update:
					l.orderUpdate <- msg.(*order.Update)
				case *fundinginfo.FundingInfo:
					l.funding <- msg.(*fundinginfo.FundingInfo)
				case *position.Snapshot:
					l.positionSnapshot <- msg.(*position.Snapshot)
				case *wallet.Snapshot:
					l.walletSnapshot <- msg.(*wallet.Snapshot)
				default:
					log.Printf("COULD NOT TYPE MSG ^")
				}
			}
		}
	}()
}
