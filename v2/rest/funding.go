package rest

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"path"
)

// LedgerService manages the Ledgers endpoint.
type FundingService struct {
	requestFactory
	Synchronous
}

// All returns all ledgers for the authenticated account
func (fs *FundingService) Offers(symbol string) (*bitfinex.FundingOfferSnapshot, error) {
	req, err :=fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/offers", symbol))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	offers, err := bitfinex.NewFundingOfferSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (fs *FundingService) OfferHistory(symbol string) (*bitfinex.FundingOfferSnapshot, error) {
	req, err :=fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/offers", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	offers, err := bitfinex.NewFundingOfferSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (fs *FundingService) Loans(symbol string) (*bitfinex.FundingLoanSnapshot, error) {
	req, err := fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/loans", symbol))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	loans, err := bitfinex.NewFundingLoanSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (fs *FundingService) LoansHistory(symbol string) (*bitfinex.FundingLoanSnapshot, error) {
	req, err := fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/loans", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	loans, err := bitfinex.NewFundingLoanSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (fs *FundingService) Credits(symbol string) (*bitfinex.FundingCreditSnapshot, error) {
	req, err := fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/credits", symbol))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	loans, err := bitfinex.NewFundingCreditSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (fs *FundingService) CreditsHistory(symbol string) (*bitfinex.FundingCreditSnapshot, error) {
	req, err := fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/credits", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	loans, err := bitfinex.NewFundingCreditSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (fs *FundingService) Trades(symbol string) (*bitfinex.FundingTradeSnapshot, error) {
	req, err := fs.requestFactory.NewAuthenticatedRequest(bitfinex.PermissionRead, path.Join("funding/trades", symbol, "hist"))
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	loans, err := bitfinex.NewFundingTradeSnapshotFromRaw(raw)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (fs *FundingService) SubmitOffer(fo *bitfinex.FundingOfferRequest) (*bitfinex.Notification, error) {
	bytes, err := fo.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := fs.requestFactory.NewAuthenticatedRequestWithBytes(bitfinex.PermissionWrite, path.Join("funding/offer/submit"), bytes)
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	return bitfinex.NewNotificationFromRaw(raw)
}

func (fs *FundingService) CancelOffer(fc *bitfinex.FundingOfferCancelRequest) (*bitfinex.Notification, error) {
	bytes, err := fc.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := fs.requestFactory.NewAuthenticatedRequestWithBytes(bitfinex.PermissionWrite, "funding/offer/cancel", bytes)
	if err != nil {
		return nil, err
	}
	raw, err := fs.Request(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(raw)
	return bitfinex.NewNotificationFromRaw(raw)
}
