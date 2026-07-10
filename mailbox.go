package migadu

import (
	"context"
	"net/http"
)

// Mailbox represents a mailbox in the Migadu API.
type Mailbox struct {
	Address               string     `json:"address,omitempty"`
	AutorespondActive     bool       `json:"autorespond_active,omitempty"`
	AutorespondBody       string     `json:"autorespond_body,omitempty"`
	AutorespondExpiresOn  string     `json:"autorespond_expires_on,omitempty"`
	AutorespondSubject    string     `json:"autorespond_subject,omitempty"`
	ChangedAt             string     `json:"changed_at,omitempty"`
	Delegations           []string   `json:"delegations,omitempty"`
	DomainName            string     `json:"domain_name,omitempty"`
	Expireable            bool       `json:"expireable,omitempty"`
	ExpiresOn             string     `json:"expires_on,omitempty"`
	FooterActive          bool       `json:"footer_active,omitempty"`
	FooterHTMLBody        string     `json:"footer_html_body,omitempty"`
	FooterPlainBody       string     `json:"footer_plain_body,omitempty"`
	Identities            []Identity `json:"identities,omitempty"`
	IsInternal            bool       `json:"is_internal,omitempty"`
	LastLoginAt           string     `json:"last_login_at,omitempty"`
	LocalPart             string     `json:"local_part,omitempty"`
	MayAccessImap         bool       `json:"may_access_imap,omitempty"`
	MayAccessManagesieve  bool       `json:"may_access_managesieve,omitempty"`
	MayAccessPop3         bool       `json:"may_access_pop3,omitempty"`
	MayReceive            bool       `json:"may_receive,omitempty"`
	MaySend               bool       `json:"may_send,omitempty"`
	Name                  string     `json:"name,omitempty"`
	Password              string     `json:"password,omitempty"`
	PasswordMethod        string     `json:"password_method,omitempty"`
	PasswordRecoveryEmail string     `json:"password_recovery_email,omitempty"`
	RecipientDenylist     []string   `json:"recipient_denylist,omitempty"`
	RemoveUponExpiry      bool       `json:"remove_upon_expiry,omitempty"`
	SenderAllowlist       []string   `json:"sender_allowlist,omitempty"`
	SenderDenylist        []string   `json:"sender_denylist,omitempty"`
	SpamAction            string     `json:"spam_action,omitempty"`
	SpamAggressiveness    string     `json:"spam_aggressiveness,omitempty"`
	StorageUsage          float64    `json:"storage_usage,omitempty"`
	WildcardSender        bool       `json:"wildcard_sender,omitempty"`
}

// CreateMailboxRequest contains fields accepted by the mailbox create endpoint.
type CreateMailboxRequest struct {
	LocalPart             string   `json:"local_part"`
	Name                  string   `json:"name,omitempty"`
	PasswordMethod        string   `json:"password_method,omitempty"`
	Password              string   `json:"password,omitempty"`
	PasswordRecoveryEmail string   `json:"password_recovery_email,omitempty"`
	ForwardingTo          string   `json:"forwarding_to,omitempty"`
	IsInternal            *bool    `json:"is_internal,omitempty"`
	WildcardSender        *bool    `json:"wildcard_sender,omitempty"`
	MaySend               *bool    `json:"may_send,omitempty"`
	MayReceive            *bool    `json:"may_receive,omitempty"`
	MayAccessImap         *bool    `json:"may_access_imap,omitempty"`
	MayAccessPop3         *bool    `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve  *bool    `json:"may_access_managesieve,omitempty"`
	SpamAction            string   `json:"spam_action,omitempty"`
	SpamAggressiveness    string   `json:"spam_aggressiveness,omitempty"`
	SenderDenylist        []string `json:"sender_denylist,omitempty"`
	SenderAllowlist       []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist     []string `json:"recipient_denylist,omitempty"`
	FooterActive          *bool    `json:"footer_active,omitempty"`
	FooterPlainBody       *string  `json:"footer_plain_body,omitempty"`
	FooterHTMLBody        *string  `json:"footer_html_body,omitempty"`
}

// UpdateMailboxRequest uses pointers so zero values can be sent explicitly.
type UpdateMailboxRequest struct {
	Name                  *string   `json:"name,omitempty"`
	IsInternal            *bool     `json:"is_internal,omitempty"`
	WildcardSender        *bool     `json:"wildcard_sender,omitempty"`
	MaySend               *bool     `json:"may_send,omitempty"`
	MayReceive            *bool     `json:"may_receive,omitempty"`
	MayAccessImap         *bool     `json:"may_access_imap,omitempty"`
	MayAccessPop3         *bool     `json:"may_access_pop3,omitempty"`
	MayAccessManagesieve  *bool     `json:"may_access_managesieve,omitempty"`
	Password              *string   `json:"password,omitempty"`
	PasswordRecoveryEmail *string   `json:"password_recovery_email,omitempty"`
	AutorespondActive     *bool     `json:"autorespond_active,omitempty"`
	AutorespondBody       *string   `json:"autorespond_body,omitempty"`
	AutorespondExpiresOn  *string   `json:"autorespond_expires_on,omitempty"`
	AutorespondSubject    *string   `json:"autorespond_subject,omitempty"`
	Delegations           *[]string `json:"delegations,omitempty"`
	ExpiresOn             *string   `json:"expires_on,omitempty"`
	RemoveUponExpiry      *bool     `json:"remove_upon_expiry,omitempty"`
	RecipientDenylist     *[]string `json:"recipient_denylist,omitempty"`
	SenderAllowlist       *[]string `json:"sender_allowlist,omitempty"`
	SenderDenylist        *[]string `json:"sender_denylist,omitempty"`
	SpamAction            *string   `json:"spam_action,omitempty"`
	SpamAggressiveness    *string   `json:"spam_aggressiveness,omitempty"`
	FooterActive          *bool     `json:"footer_active,omitempty"`
	FooterPlainBody       *string   `json:"footer_plain_body,omitempty"`
	FooterHTMLBody        *string   `json:"footer_html_body,omitempty"`
}

// ListMailboxes lists all mailboxes for a domain.
// It returns the mailboxes and any error encountered.
func (c *Client) ListMailboxes(ctx context.Context, domain string) ([]*Mailbox, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddPath(mailboxesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := doRequest[struct {
		Mailboxes []*Mailbox `json:"mailboxes,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Mailboxes, nil
}

// GetMailbox retrieves a single mailbox given its local part name.
// It returns a pointer to a Mailbox struct and any error encountered.
func (c *Client) GetMailbox(ctx context.Context, domain, localPart string) (*Mailbox, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(mailboxesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Mailbox](c, ctx, req)
}

// CreateMailbox creates a mailbox using all fields supported by the API.
func (c *Client) CreateMailbox(ctx context.Context, domain string, mailbox CreateMailboxRequest) (*Mailbox, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddPath(mailboxesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(mailbox).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Mailbox](c, ctx, req)
}

// UpdateMailbox updates only fields explicitly set on update.
func (c *Client) UpdateMailbox(ctx context.Context, domain, localPart string, update UpdateMailboxRequest) (*Mailbox, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(mailboxesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Mailbox](c, ctx, req)
}

// DeleteMailbox deletes a mailbox by local part.
// It returns any error encountered.
func (c *Client) DeleteMailbox(ctx context.Context, domain, localPart string) error {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(mailboxesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = doRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
