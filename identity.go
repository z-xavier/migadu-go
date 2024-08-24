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

// ListIdentities lists all the identities for the given mailbox local part name.
// Ir returns a pointer to an array of Identity structs and any error encountered.
func (c *Client) ListIdentities(ctx context.Context, mailbox string) ([]*Identity, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddPath(identitiesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest[struct {
		Indentities []*Identity `json:"identities,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Indentities, nil
}

// GetIdentity  retrieves a single identity given its mailbox name and local part name.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, mailbox, localPart string) (*Identity, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return DoRequest[Identity](c, ctx, req)
}

// NewIdentity creates a new identity given the mailbox, local part name and a display name.
// It returns a pointer to am Identity struct and any error encountered.
func (c *Client) NewIdentity(ctx context.Context, mailbox, localPart, displayName string) (*Identity, error) {
	identity := Identity{LocalPart: localPart, Name: displayName}
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPost).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
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
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPut).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(identity).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Identity](c, ctx, req)
}

// DeleteIdentity deletes an identity given a pointer to an Identity struct.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, mailbox, localPart string) error {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodDelete).
		AddRestfulPath(mailboxesPath, mailbox).
		AddRestfulPath(identitiesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
