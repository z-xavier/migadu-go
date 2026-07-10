package migadu

import (
	"context"
	"fmt"
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

// ListAliases lists all the aliases for the domain configured on the client.
// It returns the aliases and any error encountered.
func (c *Client) ListAliases(ctx context.Context) ([]*Alias, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddPath(AliasesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest[struct {
		Aliases []*Alias `json:"address_aliases,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Aliases, nil
}

// GetAlias retrieves a single alias given its local part name.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) GetAlias(ctx context.Context, localPart string) (*Alias, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(AliasesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return DoRequest[Alias](c, ctx, req)
}

// NewAlias creates a new alias given the local part name and its destinations.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) NewAlias(ctx context.Context, localPart string, destinations []string) (*Alias, error) {
	return c.CreateAlias(ctx, CreateAliasRequest{LocalPart: localPart, Destinations: destinations})
}

// CreateAlias creates an alias using all fields supported by the API.
func (c *Client) CreateAlias(ctx context.Context, alias CreateAliasRequest) (*Alias, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddPath(AliasesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(alias).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Alias](c, ctx, req)
}

// UpdateAlias updates an alias in place given a pointer to an Alias struct.
// It returns a pointer to a new Alias struct and any error encountered.
func (c *Client) UpdateAlias(ctx context.Context, localPart string, a *Alias) (*Alias, error) {
	if a == nil {
		return nil, fmt.Errorf("alias is required")
	}
	update := UpdateAliasRequest{
		Destinations: &a.Destinations, IsInternal: &a.IsInternal,
		RemoveUponExpiry: &a.RemoveUponExpiry,
	}
	if a.ExpiresOn != "" {
		update.ExpiresOn = &a.ExpiresOn
	}
	return c.UpdateAliasWithRequest(ctx, localPart, update)
}

// UpdateAliasWithRequest updates only fields explicitly set on update.
func (c *Client) UpdateAliasWithRequest(ctx context.Context, localPart string, update UpdateAliasRequest) (*Alias, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(AliasesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Alias](c, ctx, req)
}

// DeleteAlias deletes an alias by local part.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, localPart string) error {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(AliasesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
