package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestRequest_ErrorCheckBeforeHeaderUse(t *testing.T) {
	// Verify that an invalid method does not cause nil pointer dereference.
	// http.NewRequest should fail with an invalid method containing a space.
	h := HttpTransport{
		BaseURL:    &url.URL{Scheme: "https", Host: "api.bitfinex.com"},
		HTTPClient: http.DefaultClient,
		httpDo: func(c *http.Client, req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("should not reach here")
		},
	}

	req := Request{
		RefURL:  "/v2/test",
		Data:    []byte("{}"),
		Method:  "BAD METHOD", // contains space, triggers NewRequest error
		Headers: map[string]string{"X-Test": "value"},
	}

	_, err := h.Request(req)
	if err == nil {
		t.Fatal("expected error for invalid HTTP method, got nil")
	}
}

func TestRequest_ErrorWrapsContext(t *testing.T) {
	h := HttpTransport{
		BaseURL:    &url.URL{Scheme: "https", Host: "api.bitfinex.com"},
		HTTPClient: http.DefaultClient,
		httpDo: func(c *http.Client, req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network failure")
		},
	}

	req := Request{
		RefURL:  "/v2/test",
		Data:    []byte("{}"),
		Method:  "GET",
		Headers: make(map[string]string),
	}

	_, err := h.Request(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Verify error wrapping
	var target *url.Error
	if !errors.As(err, &target) {
		// The error should contain context from our wrapping
		if err.Error() == "network failure" {
			t.Error("error should be wrapped with context, not bare")
		}
	}
}

func TestRequest_InvalidRefURL(t *testing.T) {
	h := HttpTransport{
		BaseURL:    &url.URL{Scheme: "https", Host: "api.bitfinex.com"},
		HTTPClient: http.DefaultClient,
		httpDo: func(c *http.Client, req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("should not reach here")
		},
	}

	req := Request{
		RefURL:  "://invalid",
		Data:    []byte("{}"),
		Method:  "GET",
		Headers: make(map[string]string),
	}

	_, err := h.Request(req)
	if err == nil {
		t.Fatal("expected error for invalid ref URL, got nil")
	}
	// Should contain our wrapping context
	expected := "parse ref URL"
	if len(err.Error()) < len(expected) {
		t.Errorf("error too short, got: %s", err.Error())
	}
}
