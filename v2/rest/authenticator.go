package rest

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type Authenticator interface {
	SetCredentials(key string, secret string)
	NewAuthenticatedPostRequest(string, map[string]interface{}) (Request, error)
}

type AuthenticatorAttributes struct {
	Key    string
	Secret string
}

type authenticator struct {
	attributes AuthenticatorAttributes
}

// type authHeaders struct {
// 	BfxNonce     string
// 	BfxSignature string
// 	BfxApikey    string
// }

func NewAuthenticator() Authenticator {
	return &authenticator{}
}

func (a *authenticator) SetCredentials(key string, secret string) {
	a.attributes.Key = key
	a.attributes.Secret = secret
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

func (a *authenticator) authHeaders(path string, data map[string]interface{}) (ah map[string]string, err error) {
	nonce := strconv.FormatInt(time.Now().Unix()*10000, 10)
	_, err = json.Marshal(data)
	if err != nil {
		log.Println("Cannot convert request body to JSON!")
		return
	}

	payload := "/api/v2/" + path + nonce + "{}" //string(jsonBody)
	secret := []byte(a.attributes.Secret)
	h := hmac.New(sha512.New384, secret)
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))

	ah = make(map[string]string)
	ah["bfx-nonce"] = nonce
	ah["bfx-signature"] = signature
	ah["bfx-apikey"] = a.attributes.Key
	ah["Content-Type"] = "application/json"
	return
}
