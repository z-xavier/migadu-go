package migadu

import (
	"context"
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

// ListRewrites lists all rewrites for a domain.
// It returns the rewrites and any error encountered.
func (c *Client) ListRewrites(ctx context.Context, domain string) ([]*Rewrite, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddPath(rewritesPath).
		Build()
	if err != nil {
		return nil, err
	}

	resp, err := doRequest[struct {
		Rewrites []*Rewrite `json:"rewrites,omitempty"`
	}](c, ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Rewrites, nil
}

// GetRewrite retrieves a single rewrite given its name.
// It returns a pointer to a Rewrite struct and any error encountered.
func (c *Client) GetRewrite(ctx context.Context, domain, name string) (*Rewrite, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodGet).
		AddRestfulPath(rewritesPath, name).
		Build()
	if err != nil {
		return nil, err
	}

	return doRequest[Rewrite](c, ctx, req)
}

// CreateRewrite creates a rewrite using all fields supported by the API.
func (c *Client) CreateRewrite(ctx context.Context, domain string, rewrite CreateRewriteRequest) (*Rewrite, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPost).
		AddPath(rewritesPath).
		SetHeaderContentTypeJson().
		SetBodyJson(rewrite).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Rewrite](c, ctx, req)
}

// UpdateRewrite updates only fields explicitly set on update.
func (c *Client) UpdateRewrite(ctx context.Context, domain, name string, update UpdateRewriteRequest) (*Rewrite, error) {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return nil, err
	}
	req, err := builder.
		SetMethod(http.MethodPut).
		AddRestfulPath(rewritesPath, name).
		SetHeaderContentTypeJson().
		SetBodyJson(update).
		Build()
	if err != nil {
		return nil, err
	}
	return doRequest[Rewrite](c, ctx, req)
}

// DeleteRewrite deletes a rewrite by name.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, domain, name string) error {
	builder, err := c.getDomainReqBuilder(domain)
	if err != nil {
		return err
	}
	req, err := builder.
		SetMethod(http.MethodDelete).
		AddRestfulPath(rewritesPath, name).
		Build()
	if err != nil {
		return err
	}
	if _, err = doRequest[struct{}](c, ctx, req); err != nil {
		return err
	}
	return nil
}
