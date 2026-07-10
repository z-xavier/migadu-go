package migadu

import (
	"context"
	"net/http"
)

// Forwarding represents an external forwarding address for a mailbox.
type Forwarding struct {
	Address            string  `json:"address,omitempty"`
	BlockedAt          *string `json:"blocked_at,omitempty"`
	ConfirmationSentAt *string `json:"confirmation_sent_at,omitempty"`
	ConfirmedAt        *string `json:"confirmed_at,omitempty"`
	ExpiresOn          *string `json:"expires_on,omitempty"`
	IsActive           bool    `json:"is_active,omitempty"`
	RemoveUponExpiry   *bool   `json:"remove_upon_expiry,omitempty"`
}

// CreateForwardingRequest contains fields accepted by the forwarding create endpoint.
type CreateForwardingRequest struct {
	Address          string  `json:"address"`
	ExpiresOn        *string `json:"expires_on,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
	RemoveUponExpiry *bool   `json:"remove_upon_expiry,omitempty"`
}

// UpdateForwardingRequest uses pointers so zero values can be sent explicitly.
type UpdateForwardingRequest struct {
	ExpiresOn        *string `json:"expires_on,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
	RemoveUponExpiry *bool   `json:"remove_upon_expiry,omitempty"`
}

// ListForwardings lists all external forwarding addresses on a mailbox.
func (c *Client) ListForwardings(ctx context.Context, domain, mailbox string) ([]*Forwarding, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddPath(forwardingsPath).
		Build()
	if err != nil {
		return nil, err
	}
	resp, err := doRequest[struct {
		Forwardings []*Forwarding `json:"forwardings"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Forwardings, nil
}

// GetForwarding retrieves an external forwarding address on a mailbox.
func (c *Client) GetForwarding(ctx context.Context, domain, mailbox, address string) (*Forwarding, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(forwardingsPath, address).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Forwarding](c, ctx, req)
}

// CreateForwarding adds an external forwarding address to a mailbox.
func (c *Client) CreateForwarding(ctx context.Context, domain, mailbox string, forwarding CreateForwardingRequest) (*Forwarding, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddRestfulPath(mailboxesPath, mailbox).
		AddPath(forwardingsPath).
		SetHeaderContentTypeJson().
		SetBodyJson(forwarding).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Forwarding](c, ctx, req)
}

// UpdateForwarding updates an external forwarding address on a mailbox.
func (c *Client) UpdateForwarding(ctx context.Context, domain, mailbox, address string, update UpdateForwardingRequest) (*Forwarding, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(forwardingsPath, address).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Forwarding](c, ctx, req)
}

// DeleteForwarding removes an external forwarding address from a mailbox.
func (c *Client) DeleteForwarding(ctx context.Context, domain, mailbox, address string) error {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(forwardingsPath, address).
		Build()
	if err != nil {
		return err
	}
	_, err = doRequest[struct{}](c, ctx, req)
	return err
}
