package bitfinex

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestWalletTransfer(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
          "status":"success",
          "message":"1.0 USD transfered from Exchange to Deposit"
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	response, err := NewClient().Wallet.Transfer(10.0, "BTC", "1WalletA", "1WalletB")

	if err != nil {
		t.Error(err)
	}

	if response[0].Status != "success" {
		t.Error("Expected", "success")
		t.Error("Actual ", response[0].Status)
	}
}

func TestWithdrawCrypto(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
          "status":"success",
          "message":"Your withdrawal request has been successfully submitted.",
          "withdrawal_id":586829
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	response, err := NewClient().Wallet.WithdrawCrypto(10.0, "bitcoin", WALLET_DEPOSIT, "1WalletABC")

	if err != nil {
		t.Error(err)
	}

	if response[0].Status != "success" {
		t.Error("Expected", "success")
		t.Error("Actual ", response[0].Status)
	}
	if response[0].WithdrawalID != 586829 {
		t.Error("Expected", 586829)
		t.Error("Actual ", response[0].WithdrawalID)
	}

}

func TestWithdrawWire(t *testing.T) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		msg := `[{
          "status":"success",
          "message":"Your withdrawal request has been successfully submitted.",
          "withdrawal_id":586829
        }]`
		resp := http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(msg)),
			StatusCode: 200,
		}
		return &resp, nil
	}

	intermediaryBank := BankAccount{
		AccountName:   "Bank Account Name",
		AccountNumber: "IBAN12355678976543",
		BankName:      "HongKong Bank",
		BankAddress:   "Bank Address",
		BankCity:      "Bank City",
		BankCountry:   "Bank Country",
		SwiftCode:     "SWIFT",
	}
	beneficiaryBank := BankAccount{
		AccountName:   "Bank Account Name",
		AccountNumber: "IBAN12355678976543",
		BankName:      "HongKong Bank",
		BankAddress:   "Bank Address",
		BankCity:      "Bank City",
		BankCountry:   "Bank Country",
		SwiftCode:     "SWIFT",
	}
	response, err := NewClient().Wallet.WithdrawWire(10.0, true, WALLET_DEPOSIT, beneficiaryBank, intermediaryBank, "Wire MESSAGE")

	if err != nil {
		t.Error(err)
	}

	if response[0].Status != "success" {
		t.Error("Expected", "success")
		t.Error("Actual ", response[0].Status)
	}
	if response[0].WithdrawalID != 586829 {
		t.Error("Expected", 586829)
		t.Error("Actual ", response[0].WithdrawalID)
	}

}
