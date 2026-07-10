package migadu

import (
	"context"
	"fmt"
	"net/http"
)

// Rewrite represents a rewrite rule in the Migadu API.
type Rewrite struct {
	Destinations  []string `json:"destinations,omitempty"`
	DomainName    string   `json:"domain_name,omitempty"`
	LocalPartRule string   `json:"local_part_rule,omitempty"`
	Name          string   `json:"name,omitempty"`
	OrderNum      int      `json:"order_num,omitempty"`
}

// CreateRewriteRequest contains fields accepted by the rewrite create endpoint.
type CreateRewriteRequest struct {
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	Destinations  []string `json:"destinations"`
	OrderNum      *int     `json:"order_num,omitempty"`
}

// UpdateRewriteRequest uses pointers so zero values can be sent explicitly.
type UpdateRewriteRequest struct {
	Name          *string   `json:"name,omitempty"`
	LocalPartRule *string   `json:"local_part_rule,omitempty"`
	Destinations  *[]string `json:"destinations,omitempty"`
	OrderNum      *int      `json:"order_num,omitempty"`
}

// ListRewrites lists all the rewrites for the domain configured on the client.
// It returns the rewrites and any error encountered.
func (c *Client) ListRewrites(ctx context.Context) ([]*Rewrite, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
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
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
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
	return c.CreateRewrite(ctx, CreateRewriteRequest{Name: name, LocalPartRule: localPartRule, Destinations: destinations})
}

// CreateRewrite creates a rewrite using all fields supported by the API.
func (c *Client) CreateRewrite(ctx context.Context, rewrite CreateRewriteRequest) (*Rewrite, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddPath(RewritesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(rewrite).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Rewrite](c, ctx, req)
}

// UpdateRewrite updates a rewrite in place given a pointer to a Rewrite struct.
// It returns a pointer to a new Rewrite struct and any error encountered.
func (c *Client) UpdateRewrite(ctx context.Context, name string, r *Rewrite) (*Rewrite, error) {
	if r == nil {
		return nil, fmt.Errorf("rewrite is required")
	}
	return c.UpdateRewriteWithRequest(ctx, name, UpdateRewriteRequest{
		Name: &r.Name, LocalPartRule: &r.LocalPartRule, Destinations: &r.Destinations, OrderNum: &r.OrderNum,
	})
}

// UpdateRewriteWithRequest updates only fields explicitly set on update.
func (c *Client) UpdateRewriteWithRequest(ctx context.Context, name string, update UpdateRewriteRequest) (*Rewrite, error) {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(RewritesPath, name).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return DoRequest[Rewrite](c, ctx, req)
}

// DeleteRewrite deletes a rewrite by name.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, name string) error {
	builder, err := c.getConfiguredDomainReqBuilder()
	if err != nil {
		return err
	}
	req, err := builder.
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
