package migadu

import (
	"context"
	"net/http"
	"strings"
)

// Rewrite represents a rewrite rule in the Migadu API.
type Rewrite struct {
	Destinations  []string `json:"destinations,omitempty"`
	LocalPartRule string   `json:"local_part_rule,omitempty"`
	Name          string   `json:"name,omitempty"`
	OrderNum      int      `json:"order_num,omitempty"`
}

// rewriteJSON is used when a new/updated alias object to the API.
type rewriteJSON struct {
	Rewrite
	DestinationsJSON string `json:"destinations,omitempty"`
}

// convertDestinationsField takes a slice of strings and joins them into a comma separated line.
func (r *rewriteJSON) convertDestinationsField() {
	r.DestinationsJSON = strings.Join(r.Destinations, ",")
	r.Destinations = nil
}

// ListRewrites lists all the rewrites for the domain configured on the client.
// Ir returns a pointer to an array of Rewrite structs and any error encountered.
func (c *Client) ListRewrites(ctx context.Context) ([]*Rewrite, error) {

	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddPath(RewritesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest[struct {
		Rewrites []*Rewrite `json:"rewrites,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Rewrites, nil
}

// GetRewrite retrieves a single rewrite given its name.
// It returns a pointer to a Rewrite struct and any error encountered.
func (c *Client) GetRewrite(ctx context.Context, name string) (*Rewrite, error) {

	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddRestfulPath(RewritesPath, name).
		Build()
	if err != nil {
		return nil, err
	}

	return DoRequest[Rewrite](c, ctx, req)
}

// NewRewrite creates a new rewrite given the name, local part rule and its destinations.
// It returns a pointer to a Rewrite struct and any error encountered.
func (c *Client) NewRewrite(ctx context.Context, name string, localPartRule string, destinations []string) (*Rewrite, error) {
	rewriteJSON := rewriteJSON{Rewrite: Rewrite{Name: name, LocalPartRule: localPartRule, Destinations: destinations}}
	rewriteJSON.convertDestinationsField()
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPost).
		AddPath(RewritesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(rewriteJSON).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Rewrite](c, ctx, req)
}

// UpdateRewrite updates a rewrite in place given a pointer to a Rewrite struct.
// It returns a pointer to a new Rewrite struct and any error encountered.
func (c *Client) UpdateRewrite(ctx context.Context, name string, r *Rewrite) (*Rewrite, error) {
	rewriteJSON := rewriteJSON{Rewrite: *r}
	rewriteJSON.convertDestinationsField()
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodPut).
		AddRestfulPath(RewritesPath, name).
		SetHeaderContentTypeJson().
		SetBodyJson(rewriteJSON).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Rewrite](c, ctx, req)
}

// DeleteRewrite deletes a rewrite given a pointer to a Rewrite struct.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, name string) error {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodDelete).
		AddRestfulPath(RewritesPath, name).
		Build()
	if err != nil {
		return err
	}
	if _, err = DoRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
