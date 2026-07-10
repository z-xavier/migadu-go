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
	APIHost        = "https://api.migadu.com"
	V1Path         = "v1"
	DefaultTimeout = 30 * time.Second

	DomainsPath     = "domains"
	AliasesPath     = "aliases"
	RewritesPath    = "rewrites"
	IdentitiesPath  = "identities"
	MailboxesPath   = "mailboxes"
	ForwardingsPath = "forwardings"
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
	Email      string
	APIKey     string
	Domain     string
	BaseURL    string
	Cookies    []*http.Cookie
	Timeout    time.Duration
	HTTPClient HTTPDoer
}

func (c *Client) getV1ReqBuilder() *httpReqBuilder {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = APIHost
	}
	builder := newReqBuilder().
		SetHost(baseURL).
		AddPath(V1Path).
		SetBasicAuth(c.Email, c.APIKey)
	if len(c.Cookies) > 0 {
		builder.SetCookies(c.Cookies...)
	}
	return builder
}

func (c *Client) getDomainReqBuilder(domain string) (*httpReqBuilder, error) {
	if strings.TrimSpace(domain) == "" {
		return nil, ErrDomainRequired
	}
	return c.getV1ReqBuilder().AddRestfulPath(DomainsPath, domain), nil
}

func (c *Client) getConfiguredDomainReqBuilder() (*httpReqBuilder, error) {
	return c.getDomainReqBuilder(c.Domain)
}

func newClient(email, apiKey string) (*Client, error) {
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, ErrAPIKeyRequired
	}
	return &Client{
		Email:      email,
		APIKey:     apiKey,
		BaseURL:    APIHost,
		Timeout:    DefaultTimeout,
		HTTPClient: http.DefaultClient,
	}, nil
}

// NewClient creates an account-level Migadu API client without making a network request.
func NewClient(email, apiKey string) (*Client, error) {
	return newClient(email, apiKey)
}

// New creates a domain-scoped Migadu API client without making a network request.
func New(email string, apiKey string, domain string) (*Client, error) {
	if strings.TrimSpace(domain) == "" {
		return nil, ErrDomainRequired
	}
	client, err := newClient(email, apiKey)
	if err != nil {
		return nil, err
	}
	client.Domain = domain
	return client, nil
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

// DoRequest executes and decodes a Migadu API request.
func DoRequest[T any](c *Client, ctx context.Context, req *http.Request) (*T, error) {
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
