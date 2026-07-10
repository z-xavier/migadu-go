package migadu

import (
	"context"
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
func (c *Client) ListIdentities(ctx context.Context, domain, mailbox string) ([]*Identity, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddPath(identitiesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := doRequest[struct {
		Identities []*Identity `json:"identities,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Identities, nil
}

// GetIdentity  retrieves a single identity given its mailbox name and local part name.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, domain, mailbox, localPart string) (*Identity, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return doRequest[Identity](c, ctx, req)
}

// CreateIdentity creates an identity using all fields supported by the API.
func (c *Client) CreateIdentity(ctx context.Context, domain, mailbox string, identity CreateIdentityRequest) (*Identity, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddRestfulPath(mailboxesPath, mailbox).
		AddPath(identitiesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(identity).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Identity](c, ctx, req)
}

// UpdateIdentity updates only fields explicitly set on update.
func (c *Client) UpdateIdentity(ctx context.Context, domain, mailbox, localPart string, update UpdateIdentityRequest) (*Identity, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Identity](c, ctx, req)
}

// DeleteIdentity deletes an identity by mailbox and local part.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, domain, mailbox, localPart string) error {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = doRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
