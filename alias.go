package migadu

import (
	"context"
	"net/http"
	"strings"
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

// aliasJSON is used when a new/updated alias object to the API.
type aliasJSON struct {
	Alias
	DestinationsJSON string `json:"destinations,omitempty"`
}

// convertDestinationsField takes a slice of strings and joins them into a comma separated line.
func (a *aliasJSON) convertDestinationsField() {
	a.DestinationsJSON = strings.Join(a.Destinations, ",")
	a.Destinations = nil
}

// ListAliases lists all the aliases for the domain configured on the client.
// Ir returns a pointer to an array of Alias structs and any error encountered.
func (c *Client) ListAliases(ctx context.Context) ([]*Alias, error) {

	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddPath(aliasesPath).
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
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddRestfulPath(aliasesPath, localPart).
		Build()
	if err != nil {
		return nil, err
	}

	return DoRequest[Alias](c, ctx, req)
}

// NewAlias creates a new alias given the local part name and its destinations.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) NewAlias(ctx context.Context, localPart string, destinations []string) (*Alias, error) {
	aliasJSON := aliasJSON{Alias: Alias{LocalPart: localPart, Destinations: destinations}}
	aliasJSON.convertDestinationsField()
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPost).
		AddPath(aliasesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(aliasJSON).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Alias](c, ctx, req)
}

// UpdateAlias updates an alias in place given a pointer to an Alias struct.
// It returns a pointer to a new Alias struct and any error encountered.
func (c *Client) UpdateAlias(ctx context.Context, localPart string, a *Alias) (*Alias, error) {
	aliasJSON := aliasJSON{Alias: *a}
	aliasJSON.convertDestinationsField()
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPut).
		AddRestfulPath(aliasesPath, localPart).
		SetHeaderContentTypeJson().
		SetBodyJson(aliasJSON).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Alias](c, ctx, req)
}

// DeleteAlias deletes an alias given a pointer to an Alias struct.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, localPart string) error {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodDelete).
		AddRestfulPath(aliasesPath, localPart).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
