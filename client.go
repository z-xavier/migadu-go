package migadu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultAPIHost = "https://api.migadu.com"
	v1Path         = "v1"
	DefaultTimeout = 30 * time.Second

	domainsPath     = "domains"
	aliasesPath     = "aliases"
	rewritesPath    = "rewrites"
	identitiesPath  = "identities"
	mailboxesPath   = "mailboxes"
	forwardingsPath = "forwardings"
)

var (
	ErrEmailRequired  = errors.New("email is required")
	ErrAPIKeyRequired = errors.New("API key is required")
	ErrDomainRequired = errors.New("domain is required")
)

// HTTPDoer is implemented by http.Client and test transports that can execute requests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a client for working with Migadu API.
type Client struct {
	BaseURL    string
	Timeout    time.Duration
	HTTPClient HTTPDoer
	email      string
	apiKey     string
}

func (c *Client) getV1ReqBuilder() *httpReqBuilder {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = defaultAPIHost
	}
	builder := newReqBuilder().
		SetHost(baseURL).
		AddPath(v1Path).
		SetBasicAuth(c.email, c.apiKey)
	return builder
}

func (c *Client) getDomainReqBuilder(domain string) (*httpReqBuilder, error) {
	if strings.TrimSpace(domain) == "" {
		return nil, ErrDomainRequired
	}
	return c.getV1ReqBuilder().AddRestfulPath(domainsPath, domain), nil
}

// New creates a Migadu API client without making a network request.
func New(email, apiKey string) (*Client, error) {
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrAPIKeyRequired
	}
	return &Client{
		BaseURL:    defaultAPIHost,
		Timeout:    DefaultTimeout,
		HTTPClient: http.DefaultClient,
		email:      email,
		apiKey:     apiKey,
	}, nil
}

// APIError describes a non-success response returned by the Migadu API.
type APIError struct {
	StatusCode int
	Code       string `json:"error"`
	Message    string `json:"message"`
	Body       string `json:"-"`
}

func (e *APIError) Error() string {
	detail := e.Message
	if detail == "" {
		detail = e.Code
	}
	if detail == "" {
		detail = e.Body
	}
	if detail == "" {
		return fmt.Sprintf("Migadu API returned status %d", e.StatusCode)
	}
	return fmt.Sprintf("Migadu API returned status %d: %s", e.StatusCode, detail)
}

func doRequest[T any](c *Client, ctx context.Context, req *http.Request) (*T, error) {
	if c.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}
	req = req.WithContext(ctx)
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		apiErr := &APIError{StatusCode: resp.StatusCode, Body: string(body)}
		_ = json.Unmarshal(body, apiErr)
		return nil, apiErr
	}
	var result T
	if len(strings.TrimSpace(string(body))) == 0 {
		return &result, nil
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
