package migadu

import (
	"context"
	"net/http"
)

// Alias represents an alias in the Migadu API.
type Alias struct {
	Address          string   `json:"address,omitempty"`
	Destinations     []string `json:"destinations,omitempty"`
	DomainName       string   `json:"domain_name,omitempty"`
	Expireable       bool     `json:"expireable,omitempty"`
	ExpiresOn        string   `json:"expires_on,omitempty"`
	IsInternal       bool     `json:"is_internal,omitempty"`
	LocalPart        string   `json:"local_part,omitempty"`
	RemoveUponExpiry bool     `json:"remove_upon_expiry,omitempty"`
}

// CreateAliasRequest contains fields accepted by the alias create endpoint.
type CreateAliasRequest struct {
	LocalPart    string   `json:"local_part"`
	Destinations []string `json:"destinations"`
	IsInternal   *bool    `json:"is_internal,omitempty"`
}

// UpdateAliasRequest uses pointers so zero values can be sent explicitly.
type UpdateAliasRequest struct {
	Destinations     *[]string `json:"destinations,omitempty"`
	IsInternal       *bool     `json:"is_internal,omitempty"`
	ExpiresOn        *string   `json:"expires_on,omitempty"`
	RemoveUponExpiry *bool     `json:"remove_upon_expiry,omitempty"`
}

// ListAliases lists all aliases for a domain.
// It returns the aliases and any error encountered.
func (c *Client) ListAliases(ctx context.Context, domain string) ([]*Alias, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddPath(aliasesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := doRequest[struct {
		Aliases []*Alias `json:"address_aliases,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Aliases, nil
}

// GetAlias retrieves a single alias given its local part name.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) GetAlias(ctx context.Context, domain, localPart string) (*Alias, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(aliasesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return doRequest[Alias](c, ctx, req)
}

// CreateAlias creates an alias using all fields supported by the API.
func (c *Client) CreateAlias(ctx context.Context, domain string, alias CreateAliasRequest) (*Alias, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddPath(aliasesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(alias).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Alias](c, ctx, req)
}

// UpdateAlias updates only fields explicitly set on update.
func (c *Client) UpdateAlias(ctx context.Context, domain, localPart string, update UpdateAliasRequest) (*Alias, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(aliasesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Alias](c, ctx, req)
}

// DeleteAlias deletes an alias by local part.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, domain, localPart string) error {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(aliasesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = doRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
