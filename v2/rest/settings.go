package rest

import (
	"encoding/json"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
)

// SettingsService manages the User Settings endpoint.
type SettingsService struct {
	requestFactory
	Synchronous
}

// Read - Read user settings by keys
// see https://docs.bitfinex.com/reference/rest-auth-settings
func (s *SettingsService) Read(keys []string) ([]interface{}, error) {
	data := map[string]interface{}{"keys": keys}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionRead, "settings", data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Write - Write user settings
// see https://docs.bitfinex.com/reference/rest-auth-settings-write
func (s *SettingsService) Write(settings map[string]interface{}) ([]interface{}, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}
	req, err := s.requestFactory.NewAuthenticatedRequestWithBytes(common.PermissionWrite, "settings/set", b)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Delete - Delete user settings by keys
// see https://docs.bitfinex.com/reference/rest-auth-settings-delete
func (s *SettingsService) Delete(keys []string) ([]interface{}, error) {
	data := map[string]interface{}{"keys": keys}
	req, err := s.requestFactory.NewAuthenticatedRequestWithData(common.PermissionWrite, "settings/del", data)
	if err != nil {
		return nil, err
	}
	raw, err := s.Request(req)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
