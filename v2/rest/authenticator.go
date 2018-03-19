package rest

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"

	"github.com/bitfinexcom/bitfinex-api-go/utils"
)

type Authenticator interface {
	SetCredentials(key string, secret string)
	SetNonceGenerator(utils.NonceGenerator)
	NewAuthenticatedPostRequest(string, map[string]interface{}) (Request, error)
}

type authenticator struct {
	Key    string
	Secret string
	utils.NonceGenerator
}

func NewAuthenticator() Authenticator {
	return &authenticator{NonceGenerator: utils.NewEpochNonceGenerator()}
}

func (a *authenticator) SetCredentials(key string, secret string) {
	a.Key = key
	a.Secret = secret
}

// For late binding Nonce Generator
func (a *authenticator) SetNonceGenerator(ng utils.NonceGenerator) {
	a.NonceGenerator = ng
}

func (a *authenticator) NewAuthenticatedPostRequest(refURL string, data map[string]interface{}) (req Request, err error) {
	authHeaders, err := a.authHeaders(refURL, data)
	if err != nil {
		return
	}

	req = Request{
		RefURL:  refURL,
		Data:    data,
		Method:  "POST",
		Headers: authHeaders,
	}

	return
}

var GetNonce = utils.GetNonce

func (a *authenticator) authHeaders(path string, data map[string]interface{}) (ah map[string]string, err error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	nonce := a.GetNonce()
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	payload := "/api/v2/" + path + nonce + string(jsonBody)
	secret := []byte(a.Secret)
	h := hmac.New(sha512.New384, secret)
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))

	ah = make(map[string]string)
	ah["bfx-nonce"] = nonce
	ah["bfx-signature"] = signature
	ah["bfx-apikey"] = a.Key
	ah["Content-Type"] = "application/json"
	return
}
