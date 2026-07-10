package migadu

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type stringList []string

func (s *stringList) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*s = nil
		return nil
	}
	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		*s = values
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if strings.TrimSpace(value) == "" {
		*s = []string{}
		return nil
	}
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	*s = values
	return nil
}

type stringValue string

func (s *stringValue) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		*s = stringValue(value)
		return nil
	}
	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	*s = stringValue(number.String())
	return nil
}

// Domain represents a domain in the Migadu API.
type Domain struct {
	Name                             string   `json:"name,omitempty"`
	ActivatedAt                      *string  `json:"activated_at,omitempty"`
	DeactivatedAt                    *string  `json:"deactivated_at,omitempty"`
	Tags                             []string `json:"tags,omitempty"`
	State                            string   `json:"state,omitempty"`
	Description                      string   `json:"description,omitempty"`
	CanSend                          bool     `json:"can_send,omitempty"`
	CanReceive                       bool     `json:"can_receive,omitempty"`
	CanAccess                        bool     `json:"can_access,omitempty"`
	MXProxyEnabled                   bool     `json:"mx_proxy_enabled,omitempty"`
	SpamAggressiveness               string   `json:"spam_aggressiveness,omitempty"`
	SubjectRewritingEnabled          bool     `json:"subject_rewriting_enabled,omitempty"`
	JunkSubjectKeywordSpam           bool     `json:"junk_subject_keyword_spam,omitempty"`
	SenderDenylist                   []string `json:"sender_denylist,omitempty"`
	SenderAllowlist                  []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist                []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations             []string `json:"catchall_destinations,omitempty"`
	HostedDNS                        bool     `json:"hosted_dns,omitempty"`
	MailboxDefaultIncomingLimit      int      `json:"mailbox_default_incoming_limit,omitempty"`
	MailboxDefaultOutgoingLimit      int      `json:"mailbox_default_outgoing_limit,omitempty"`
	MailboxDefaultStorageLimit       int      `json:"mailbox_default_storage_limit,omitempty"`
	MailboxDefaultSendingEnabled     bool     `json:"mailbox_default_sending_enabled,omitempty"`
	MailboxDefaultReceivingEnabled   bool     `json:"mailbox_default_receiving_enabled,omitempty"`
	MailboxDefaultImapEnabled        bool     `json:"mailbox_default_imap_enabled,omitempty"`
	MailboxDefaultPop3Enabled        bool     `json:"mailbox_default_pop3_enabled,omitempty"`
	MailboxDefaultManagesieveEnabled bool     `json:"mailbox_default_managesieve_enabled,omitempty"`
}

// UnmarshalJSON accepts the alternate string and numeric representations shown in Migadu's domain documentation.
func (d *Domain) UnmarshalJSON(data []byte) error {
	type domainAlias Domain
	wire := struct {
		*domainAlias
		Tags                 stringList  `json:"tags"`
		SpamAggressiveness   stringValue `json:"spam_aggressiveness"`
		SenderDenylist       stringList  `json:"sender_denylist"`
		SenderAllowlist      stringList  `json:"sender_allowlist"`
		RecipientDenylist    stringList  `json:"recipient_denylist"`
		CatchallDestinations stringList  `json:"catchall_destinations"`
	}{domainAlias: (*domainAlias)(d)}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	d.Tags = []string(wire.Tags)
	d.SpamAggressiveness = string(wire.SpamAggressiveness)
	d.SenderDenylist = []string(wire.SenderDenylist)
	d.SenderAllowlist = []string(wire.SenderAllowlist)
	d.RecipientDenylist = []string(wire.RecipientDenylist)
	d.CatchallDestinations = []string(wire.CatchallDestinations)
	return nil
}

// CreateDomainRequest contains fields accepted by the domain create endpoint.
type CreateDomainRequest struct {
	Name                             string   `json:"name"`
	CreateDefaultAddresses           *bool    `json:"create_default_addresses,omitempty"`
	Tags                             []string `json:"tags,omitempty"`
	Description                      string   `json:"description,omitempty"`
	CanAccess                        *bool    `json:"can_access,omitempty"`
	MXProxyEnabled                   *bool    `json:"mx_proxy_enabled,omitempty"`
	SpamAggressiveness               string   `json:"spam_aggressiveness,omitempty"`
	SubjectRewritingEnabled          *bool    `json:"subject_rewriting_enabled,omitempty"`
	JunkSubjectKeywordSpam           *bool    `json:"junk_subject_keyword_spam,omitempty"`
	SenderDenylist                   []string `json:"sender_denylist,omitempty"`
	SenderAllowlist                  []string `json:"sender_allowlist,omitempty"`
	RecipientDenylist                []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations             []string `json:"catchall_destinations,omitempty"`
	HostedDNS                        *bool    `json:"hosted_dns,omitempty"`
	MailboxDefaultIncomingLimit      *int     `json:"mailbox_default_incoming_limit,omitempty"`
	MailboxDefaultOutgoingLimit      *int     `json:"mailbox_default_outgoing_limit,omitempty"`
	MailboxDefaultStorageLimit       *int     `json:"mailbox_default_storage_limit,omitempty"`
	MailboxDefaultSendingEnabled     *bool    `json:"mailbox_default_sending_enabled,omitempty"`
	MailboxDefaultReceivingEnabled   *bool    `json:"mailbox_default_receiving_enabled,omitempty"`
	MailboxDefaultImapEnabled        *bool    `json:"mailbox_default_imap_enabled,omitempty"`
	MailboxDefaultPop3Enabled        *bool    `json:"mailbox_default_pop3_enabled,omitempty"`
	MailboxDefaultManagesieveEnabled *bool    `json:"mailbox_default_managesieve_enabled,omitempty"`
}

// UpdateDomainRequest uses pointers so zero values can be sent explicitly.
type UpdateDomainRequest struct {
	Tags                             *[]string `json:"tags,omitempty"`
	Description                      *string   `json:"description,omitempty"`
	CanAccess                        *bool     `json:"can_access,omitempty"`
	MXProxyEnabled                   *bool     `json:"mx_proxy_enabled,omitempty"`
	SpamAggressiveness               *string   `json:"spam_aggressiveness,omitempty"`
	SubjectRewritingEnabled          *bool     `json:"subject_rewriting_enabled,omitempty"`
	JunkSubjectKeywordSpam           *bool     `json:"junk_subject_keyword_spam,omitempty"`
	SenderDenylist                   *[]string `json:"sender_denylist,omitempty"`
	SenderAllowlist                  *[]string `json:"sender_allowlist,omitempty"`
	RecipientDenylist                *[]string `json:"recipient_denylist,omitempty"`
	CatchallDestinations             *[]string `json:"catchall_destinations,omitempty"`
	HostedDNS                        *bool     `json:"hosted_dns,omitempty"`
	MailboxDefaultIncomingLimit      *int      `json:"mailbox_default_incoming_limit,omitempty"`
	MailboxDefaultOutgoingLimit      *int      `json:"mailbox_default_outgoing_limit,omitempty"`
	MailboxDefaultStorageLimit       *int      `json:"mailbox_default_storage_limit,omitempty"`
	MailboxDefaultSendingEnabled     *bool     `json:"mailbox_default_sending_enabled,omitempty"`
	MailboxDefaultReceivingEnabled   *bool     `json:"mailbox_default_receiving_enabled,omitempty"`
	MailboxDefaultImapEnabled        *bool     `json:"mailbox_default_imap_enabled,omitempty"`
	MailboxDefaultPop3Enabled        *bool     `json:"mailbox_default_pop3_enabled,omitempty"`
	MailboxDefaultManagesieveEnabled *bool     `json:"mailbox_default_managesieve_enabled,omitempty"`
}

// DNSRecord represents a DNS record returned for a domain.
type DNSRecord struct {
	Name     string `json:"name,omitempty"`
	Priority *int   `json:"priority,omitempty"`
	Type     string `json:"type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// DomainRecords contains the DNS records required by Migadu.
type DomainRecords struct {
	DomainName      string      `json:"domain_name,omitempty"`
	DKIM            []DNSRecord `json:"dkim,omitempty"`
	DMARC           *DNSRecord  `json:"dmarc,omitempty"`
	DNSVerification *DNSRecord  `json:"dns_verification,omitempty"`
	MXRecords       []DNSRecord `json:"mx_records,omitempty"`
	SPF             *DNSRecord  `json:"spf,omitempty"`
}

// DomainDiagnostics is intentionally open because Migadu does not document its response schema.
type DomainDiagnostics map[string]any

// DomainUsage represents message and storage usage for a domain.
type DomainUsage struct {
	DomainName string  `json:"domain_name,omitempty"`
	Incoming   int     `json:"incoming,omitempty"`
	Outgoing   int     `json:"outgoing,omitempty"`
	Storage    float64 `json:"storage,omitempty"`
}

// ListDomains lists all domains visible to the authenticated account.
func (c *Client) ListDomains(ctx context.Context) ([]*Domain, error) {
	req, err := c.getV1ReqBuilder().SetMethod(http.MethodGet).AddPath(domainsPath).Build()
	if err != nil {
		return nil, err
	}
	resp, err := doRequest[struct {
		Domains []*Domain `json:"domains"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Domains, nil
}

// GetDomain retrieves a domain by name.
func (c *Client) GetDomain(ctx context.Context, domain string) (*Domain, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.SetMethod(http.MethodGet).Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Domain](c, ctx, req)
}

// CreateDomain creates a domain.
func (c *Client) CreateDomain(ctx context.Context, domain CreateDomainRequest) (*Domain, error) {
	req, err := c.getV1ReqBuilder().
		SetMethod(http.MethodPost).
		AddPath(domainsPath).
		SetHeaderContentTypeJson().
		SetBodyJson(domain).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Domain](c, ctx, req)
}

// UpdateDomain updates only fields explicitly set on update.
func (c *Client) UpdateDomain(ctx context.Context, domain string, update UpdateDomainRequest) (*Domain, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPatch).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Domain](c, ctx, req)
}

// GetDomainRecords retrieves the DNS records required for a domain.
func (c *Client) GetDomainRecords(ctx context.Context, domain string) (*DomainRecords, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.SetMethod(http.MethodGet).AddPath("records").Build()
	if err != nil {
		return nil, err
	}
	return doRequest[DomainRecords](c, ctx, req)
}

// GetDomainDiagnostics runs and returns DNS diagnostics for a domain.
func (c *Client) GetDomainDiagnostics(ctx context.Context, domain string) (DomainDiagnostics, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.SetMethod(http.MethodGet).AddPath("diagnostics").Build()
	if err != nil {
		return nil, err
	}
	diagnostics, err := doRequest[DomainDiagnostics](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return *diagnostics, nil
}

// ActivateDomain asks Migadu to validate and activate a domain.
func (c *Client) ActivateDomain(ctx context.Context, domain string) (*Domain, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.SetMethod(http.MethodGet).AddPath("activate").Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Domain](c, ctx, req)
}

// GetDomainUsage retrieves current message and storage usage for a domain.
func (c *Client) GetDomainUsage(ctx context.Context, domain string) (*DomainUsage, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.SetMethod(http.MethodGet).AddPath("usage").Build()
	if err != nil {
		return nil, err
	}
	return doRequest[DomainUsage](c, ctx, req)
}
