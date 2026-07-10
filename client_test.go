package migadu

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type capturedRequest struct {
	Method      string
	Path        string
	ContentType string
	Username    string
	Password    string
	Body        []byte
}

func performRequest(t *testing.T, response string, call func(context.Context, *Client) error) capturedRequest {
	t.Helper()
	var captured capturedRequest
	client, err := New("admin@example.com", "secret", "example.com")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	client.BaseURL = "https://api.test"
	client.HTTPClient = doerFunc(func(r *http.Request) (*http.Response, error) {
		var body []byte
		if r.Body != nil {
			var err error
			body, err = io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
		}
		username, password, _ := r.BasicAuth()
		captured = capturedRequest{
			Method: r.Method, Path: r.URL.Path, ContentType: r.Header.Get("Content-Type"),
			Username: username, Password: password, Body: body,
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(response)),
		}, nil
	})
	if err := call(context.Background(), client); err != nil {
		t.Fatalf("API call error = %v", err)
	}
	return captured
}

func TestNewDoesNotMakeRequest(t *testing.T) {
	client, err := New("admin@example.com", "secret", "example.com")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if client.Domain != "example.com" || client.BaseURL != APIHost {
		t.Fatalf("New() client = %+v", client)
	}

	accountClient, err := NewClient("admin@example.com", "secret")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if accountClient.Domain != "" {
		t.Fatalf("NewClient() domain = %q", accountClient.Domain)
	}
}

func TestNewValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		new  func() error
		want error
	}{
		{name: "email", new: func() error { _, err := NewClient("", "key"); return err }, want: ErrEmailRequired},
		{name: "API key", new: func() error { _, err := NewClient("admin@example.com", ""); return err }, want: ErrAPIKeyRequired},
		{name: "domain", new: func() error { _, err := New("admin@example.com", "key", ""); return err }, want: ErrDomainRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.new(); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

func TestDoRequestReturnsAPIErrorAndClosesBody(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader(`{"error":"dns_check_failed","message":"DNS checks failed"}`)}
	client := &Client{HTTPClient: doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusUnprocessableEntity, ContentLength: -1, Body: body}, nil
	})}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DoRequest[Domain](client, context.Background(), req)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusUnprocessableEntity || apiErr.Code != "dns_check_failed" || apiErr.Message != "DNS checks failed" {
		t.Fatalf("APIError = %+v", apiErr)
	}
	if !body.closed {
		t.Fatal("response body was not closed")
	}
}

func TestDoRequestAcceptsEmptySuccessBody(t *testing.T) {
	client := &Client{HTTPClient: doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
	})}
	req, _ := http.NewRequest(http.MethodDelete, "https://example.com", nil)
	result, err := DoRequest[struct{}](client, context.Background(), req)
	if err != nil || result == nil {
		t.Fatalf("DoRequest() result = %v, error = %v", result, err)
	}
}

func TestDoRequestAppliesTimeout(t *testing.T) {
	client := &Client{
		Timeout: time.Millisecond,
		HTTPClient: doerFunc(func(req *http.Request) (*http.Response, error) {
			<-req.Context().Done()
			return nil, req.Context().Err()
		}),
	}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_, err := DoRequest[struct{}](client, context.Background(), req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want deadline exceeded", err)
	}
}

func TestDoRequestUsesContextWhenTimeoutDisabled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := &Client{
		HTTPClient: doerFunc(func(req *http.Request) (*http.Response, error) {
			return nil, req.Context().Err()
		}),
	}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_, err := DoRequest[struct{}](client, ctx, req)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
}

func TestRequestBuilderReturnsBuildErrors(t *testing.T) {
	if _, err := newReqBuilder().SetMethod(http.MethodGet).Build(); err == nil {
		t.Fatal("Build() accepted an empty host")
	}
	_, err := newReqBuilder().
		SetHost("https://example.com").
		SetMethod(http.MethodPost).
		SetBodyJson(make(chan int)).
		Build()
	if err == nil || !strings.Contains(err.Error(), "encode JSON request body") {
		t.Fatalf("Build() error = %v", err)
	}
}
