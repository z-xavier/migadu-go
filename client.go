package migadu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	APIHost        = "https://api.migadu.com"
	V1Path         = "v1"
	DefaultTimeout = 30 * time.Second

	DomainsPath    = "domains"
	AliasesPath    = "aliases"
	RewritesPath   = "rewrites"
	IdentitiesPath = "identities"
	MailboxesPath  = "mailboxes"
)

// httpClient implements the most basic function of http.Client.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a client for working with Migadu API.
type Client struct {
	Email      string
	APIKey     string
	Domain     string
	Cookies    []*http.Cookie
	Timeout    time.Duration
	HTTPClient httpClient
}

func (c *Client) GetV1ReqBuilder() *HttpReqBuilder {
	return NewReqBuilder().
		SetHost(APIHost).
		AddPath(V1Path).
		AddRestfulPath(DomainsPath, c.Domain).
		SetBasicAuth(c.Email, c.APIKey)
}

// testAuth tests that credentials are valid before creating a client.
// It returns any error encountered.
// https://www.migadu.com/api/#api-requests
func (c *Client) testAuth() error {
	if _, err := c.ListMailboxes(context.Background()); err != nil {
		return err
	}
	return nil
}

// New creates a new client one domain on Migadu given the admin email and API key.
// It returns a pointer to the new client and any error encountered.
func New(email string, apiKey string, domain string) (*Client, error) {
	client := &Client{
		Email:      email,
		APIKey:     apiKey,
		Domain:     domain,
		Timeout:    DefaultTimeout,
		HTTPClient: http.DefaultClient,
	}

	if err := client.testAuth(); err != nil {
		return nil, err
	}

	return client, nil
}

func DoRequest[T any](c *Client, ctx context.Context, req *http.Request) (*T, error) {
	if c.Timeout != 0 {
		ctx, cancel := context.WithTimeout(ctx, c.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		codeErr := fmt.Errorf("status code %d", resp.StatusCode)
		if resp.ContentLength > 0 {
			defer resp.Body.Close()
			if body, err := io.ReadAll(resp.Body); err == nil {
				return nil, errors.Join(codeErr, fmt.Errorf("resp body: %s", string(body)))
			}
		}
		return nil, errors.Join(codeErr, fmt.Errorf("contentLength == 0"))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result T
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
