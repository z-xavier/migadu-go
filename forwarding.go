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
func (c *Client) ListForwardings(ctx context.Context, mailbox string) ([]*Forwarding, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(MailboxesPath, mailbox).
		AddPath(ForwardingsPath).
		Build()
	if err != nil {
		return nil, err
	}
	resp, err := DoRequest[struct {
		Forwardings []*Forwarding `json:"forwardings"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Forwardings, nil
}

// GetForwarding retrieves an external forwarding address on a mailbox.
func (c *Client) GetForwarding(ctx context.Context, mailbox, address string) (*Forwarding, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(ForwardingsPath, address).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Forwarding](c, ctx, req)
}

// CreateForwarding adds an external forwarding address to a mailbox.
func (c *Client) CreateForwarding(ctx context.Context, mailbox string, forwarding CreateForwardingRequest) (*Forwarding, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddRestfulPath(MailboxesPath, mailbox).
		AddPath(ForwardingsPath).
		SetHeaderContentTypeJson().
		SetBodyJson(forwarding).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Forwarding](c, ctx, req)
}

// UpdateForwarding updates an external forwarding address on a mailbox.
func (c *Client) UpdateForwarding(ctx context.Context, mailbox, address string, update UpdateForwardingRequest) (*Forwarding, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(ForwardingsPath, address).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Forwarding](c, ctx, req)
}

// DeleteForwarding removes an external forwarding address from a mailbox.
func (c *Client) DeleteForwarding(ctx context.Context, mailbox, address string) error {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(ForwardingsPath, address).
		Build()
	if err != nil {
		return err
	}
	_, err = DoRequest[struct{}](c, ctx, req)
	return err
}
