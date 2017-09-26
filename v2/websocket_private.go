package bitfinex

import (
	"fmt"
)

func (b *bfxWebsocket) handlePrivateDataMessage(data []interface{}) (ms interface{}, err error) {
	if len(data) < 2 {
		return ms, fmt.Errorf("data message too short: %#v", data)
	}

	term, ok := data[1].(string)
	if !ok {
		return ms, fmt.Errorf("expected data term string in second position but got %#v in %#v", data[1], data)
	}

	if len(data) == 2 || term == "hb" { // Heartbeat
		// TODO: Consider adding a switch to enable/disable passing these along.
		return Heartbeat{}, nil
	}

	list, ok := data[2].([]interface{})
	if !ok {
		return ms, fmt.Errorf("expected data list in third position but got %#v in %#v", data[2], data)
	}

	ms = b.convertRaw(term, list)

	return
}

// convertRaw takes a term and the raw data attached to it to try and convert that
// untyped list into a proper type.
func (b *bfxWebsocket) convertRaw(term string, raw []interface{}) interface{} {
	// The things you do to get proper types.
	switch term {
	case "bu":
		o, err := balanceInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return BalanceUpdate(o)
	case "ps":
		o, err := positionSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "pn":
		o, err := positionFromRaw(raw)
		if err != nil {
			return err
		}
		return PositionNew(o)
	case "pu":
		o, err := positionFromRaw(raw)
		if err != nil {
			return err
		}
		return PositionUpdate(o)
	case "pc":
		o, err := positionFromRaw(raw)
		if err != nil {
			return err
		}
		return PositionCancel(o)
	case "ws":
		o, err := walletSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "wu":
		o, err := walletFromRaw(raw)
		if err != nil {
			return err
		}
		return WalletUpdate(o)
	case "os":
		o, err := orderSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "on":
		o, err := orderFromRaw(raw)
		if err != nil {
			return err
		}
		return OrderNew(o)
	case "ou":
		o, err := orderFromRaw(raw)
		if err != nil {
			return err
		}
		return OrderUpdate(o)
	case "oc":
		o, err := orderFromRaw(raw)
		if err != nil {
			return err
		}
		return OrderCancel(o)
	case "hts":
		o, err := tradeSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return HistoricalTradeSnapshot(o)
	case "te":
		o, err := tradeExecutionFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "tu":
		o, err := tradeFromRaw(raw)
		if err != nil {
			return err
		}
		return TradeUpdate(o)
	case "fte":
		o, err := fundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingTradeExecution(o)
	case "ftu":
		o, err := fundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingTradeUpdate(o)
	case "hfts":
		o, err := fundingTradeSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return HistoricalFundingTradeSnapshot(o)
	case "n":
		o, err := notificationFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fos":
		o, err := fundingOfferSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fon":
		o, err := offerFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingOfferNew(o)
	case "fou":
		o, err := offerFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingOfferUpdate(o)
	case "foc":
		o, err := offerFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingOfferCancel(o)
	case "fiu":
		o, err := fundingInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcs":
		o, err := fundingCreditSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcn":
		o, err := creditFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingCreditNew(o)
	case "fcu":
		o, err := creditFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingCreditUpdate(o)
	case "fcc":
		o, err := creditFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingCreditCancel(o)
	case "fls":
		o, err := fundingLoanSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fln":
		o, err := loanFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingLoanNew(o)
	case "flu":
		o, err := loanFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingLoanUpdate(o)
	case "flc":
		o, err := loanFromRaw(raw)
		if err != nil {
			return err
		}
		return FundingLoanCancel(o)
	//case "uac":
	case "hb":
		return Heartbeat{}
	case "ats":
		// TODO: Is not in documentation, so figure out what it is.
		return nil
	case "oc-req":
		// TODO
		return nil
	case "on-req":
		// TODO
		return nil
	case "mis": // Should not be sent anymore as of 2017-04-01
		return nil
	case "miu":
		o, err := marginInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	default:
	}

	return fmt.Errorf("term %q not recognized", term)
}
