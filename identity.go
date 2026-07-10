package migadu

import (
	"context"
	"fmt"
	"net/http"
)

// Identity represents an identity in the Migadu API.
type Identity struct {
	Address              string `json:"address,omitempty"`
	DomainName           string `json:"domain_name,omitempty"`
	FooterActive         bool   `json:"footer_active,omitempty"`
	FooterHTMLBody       string `json:"footer_html_body,omitempty"`
	FooterPlainBody      string `json:"footer_plain_body,omitempty"`
	LocalPart            string `json:"local_part,omitempty"`
	MayAccessImap        bool   `json:"may_access_imap,omitempty"`
	MayAccessManagesieve bool   `json:"may_access_managesieve,omitempty"`
	MayAccessPop3        bool   `json:"may_access_pop3,omitempty"`
	MayReceive           bool   `json:"may_receive,omitempty"`
	MaySend              bool   `json:"may_send,omitempty"`
	Name                 string `json:"name,omitempty"`
	Password             string `json:"password,omitempty"`
}

// CreateIdentityRequest contains fields accepted by the identity create endpoint.
type CreateIdentityRequest struct {
	LocalPart            string  `json:"local_part"`
	Name                 string  `json:"name,omitempty"`
	Password             string  `json:"password,omitempty"`
	MaySend              *bool   `json:"may_send,omitempty"`
	MayReceive           *bool   `json:"may_receive,omitempty"`
	MayAccessImap        *bool   `json:"may_access_imap,omitempty"`
	MayAccessPop3        *bool   `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve *bool   `json:"may_access_managesieve,omitempty"`
	FooterActive         *bool   `json:"footer_active,omitempty"`
	FooterPlainBody      *string `json:"footer_plain_body,omitempty"`
	FooterHTMLBody       *string `json:"footer_html_body,omitempty"`
}

// UpdateIdentityRequest uses pointers so zero values can be sent explicitly.
type UpdateIdentityRequest struct {
	Name                 *string `json:"name,omitempty"`
	Password             *string `json:"password,omitempty"`
	MaySend              *bool   `json:"may_send,omitempty"`
	MayReceive           *bool   `json:"may_receive,omitempty"`
	MayAccessImap        *bool   `json:"may_access_imap,omitempty"`
	MayAccessPop3        *bool   `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve *bool   `json:"may_access_managesieve,omitempty"`
	FooterActive         *bool   `json:"footer_active,omitempty"`
	FooterPlainBody      *string `json:"footer_plain_body,omitempty"`
	FooterHTMLBody       *string `json:"footer_html_body,omitempty"`
}

// ListIdentities lists all the identities for the given mailbox local part name.
// It returns the identities and any error encountered.
func (c *Client) ListIdentities(ctx context.Context, mailbox string) ([]*Identity, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(MailboxesPath, mailbox).
		AddPath(IdentitiesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest[struct {
		Identities []*Identity `json:"identities,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Identities, nil
}

// GetIdentity  retrieves a single identity given its mailbox name and local part name.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, mailbox, localPart string) (*Identity, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(IdentitiesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return DoRequest[Identity](c, ctx, req)
}

// NewIdentity creates a new identity given the mailbox, local part name and a display name.
// It returns the created identity and any error encountered.
func (c *Client) NewIdentity(ctx context.Context, mailbox, localPart, displayName string) (*Identity, error) {
	return c.CreateIdentity(ctx, mailbox, CreateIdentityRequest{LocalPart: localPart, Name: displayName})
}

// CreateIdentity creates an identity using all fields supported by the API.
func (c *Client) CreateIdentity(ctx context.Context, mailbox string, identity CreateIdentityRequest) (*Identity, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddRestfulPath(MailboxesPath, mailbox).
		AddPath(IdentitiesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(identity).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Identity](c, ctx, req)
}

// UpdateIdentity updates an identity in place given a pointer to an Identity struct.
// It returns a pointer to a new Identity struct and any error encountered.
func (c *Client) UpdateIdentity(ctx context.Context, mailbox, localPart string, identity *Identity) (*Identity, error) {
	if identity == nil {
		return nil, fmt.Errorf("identity is required")
	}
	update := UpdateIdentityRequest{
		Name: &identity.Name, MaySend: &identity.MaySend, MayReceive: &identity.MayReceive,
		MayAccessImap: &identity.MayAccessImap, MayAccessPop3: &identity.MayAccessPop3,
		MayAccessManagesieve: &identity.MayAccessManagesieve, FooterActive: &identity.FooterActive,
		FooterPlainBody: &identity.FooterPlainBody, FooterHTMLBody: &identity.FooterHTMLBody,
	}
	if identity.Password != "" {
		update.Password = &identity.Password
	}
	return c.UpdateIdentityWithRequest(ctx, mailbox, localPart, update)
}

// UpdateIdentityWithRequest updates only fields explicitly set on update.
func (c *Client) UpdateIdentityWithRequest(ctx context.Context, mailbox, localPart string, update UpdateIdentityRequest) (*Identity, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(IdentitiesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Identity](c, ctx, req)
}

// DeleteIdentity deletes an identity by mailbox and local part.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, mailbox, localPart string) error {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(MailboxesPath, mailbox).
		AddRestfulPath(IdentitiesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
