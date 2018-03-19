package rest

import (
	"path"
	"testing"
)

type mockNoncer struct{}

func (mn *mockNoncer) CreateNonce() string {
	return "152145189228798"
}

func TestNewAuthenticatedPostRequest(t *testing.T) {
	a := &authenticator{Noncer: &mockNoncer{}}
	a.SetCredentials("abc", "123")

	req, err := a.NewAuthenticatedPostRequest(path.Join("auth", "r", "orders", "BTCUSD", "hist"), nil)
	if err != nil {
		t.Error("NewAuthenticatedPostRequest should not throw error!")
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Error("Content type should be application/json!")
	}

	if req.Headers["bfx-nonce"] != "152145189228798" {
		t.Error("Nonce should be 152145189228798!")
	}

	if req.Headers["bfx-signature"] != "2936ac13384ea81505386777bf3089ca0aa5d3cb41e9bf2ba57077f230daef1b1dc8ce426b3ca4cb792ab58be936ed5d" {
		t.Error("Signature should be 2936ac13384ea81505386777bf3089ca0aa5d3cb41e9bf2ba57077f230daef1b1dc8ce426b3ca4cb792ab58be936ed5d!")
	}

	if req.Headers["bfx-apikey"] != "abc" {
		t.Error("Apikey should be abc!")
	}

	if req.Method != "POST" {
		t.Error("Method should be POST!")
	}
}

func TestNewAuthenticatedPostRequestWithParams(t *testing.T) {
	a := &authenticator{Noncer: &mockNoncer{}}
	a.SetCredentials("abc", "123")

	params := make(map[string]interface{})
	params["limit"] = 2
	req, err := a.NewAuthenticatedPostRequest(path.Join("auth", "r", "orders", "BTCUSD", "hist"), params)
	if err != nil {
		t.Error("NewAuthenticatedPostRequest should not throw error!")
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Error("Content type should be application/json!")
	}

	if req.Headers["bfx-nonce"] != "152145189228798" {
		t.Error("Nonce should be 152145189228798!")
	}

	if req.Headers["bfx-signature"] != "cafd073d91b816269147d46985c2ddfbeb90d79e8418c15243b85938c5bc5e758d82001195024a45ab223787593b3206" {
		t.Error("Signature should be cafd073d91b816269147d46985c2ddfbeb90d79e8418c15243b85938c5bc5e758d82001195024a45ab223787593b3206!")
	}

	if req.Headers["bfx-apikey"] != "abc" {
		t.Error("Apikey should be abc!")
	}

	if req.Method != "POST" {
		t.Error("Method should be POST!")
	}

	if req.Data["limit"] != 2 {
		t.Error("Params[limit] should be 2!")
	}
}
