package migadu

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestDocumentedEndpoints(t *testing.T) {
	falseValue := false
	zero := 0
	emptyList := []string{}
	tests := []struct {
		name         string
		method       string
		path         string
		response     string
		expectedJSON map[string]any
		call         func(context.Context, *Client) error
	}{
		{name: "domains index", method: "GET", path: "/v1/domains", response: `{"domains":[]}`, call: func(ctx context.Context, c *Client) error { _, err := c.ListDomains(ctx); return err }},
		{name: "domains show", method: "GET", path: "/v1/domains/other.example", response: `{}`, call: func(ctx context.Context, c *Client) error { _, err := c.GetDomain(ctx, "other.example"); return err }},
		{name: "domains create", method: "POST", path: "/v1/domains", response: `{}`, expectedJSON: map[string]any{"name": "other.example", "hosted_dns": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateDomain(ctx, CreateDomainRequest{Name: "other.example", HostedDNS: &falseValue})
			return err
		}},
		{name: "domains update", method: "PATCH", path: "/v1/domains/other.example", response: `{}`, expectedJSON: map[string]any{"can_access": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateDomain(ctx, "other.example", UpdateDomainRequest{CanAccess: &falseValue})
			return err
		}},
		{name: "domains records", method: "GET", path: "/v1/domains/other.example/records", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetDomainRecords(ctx, "other.example")
			return err
		}},
		{name: "domains diagnostics", method: "GET", path: "/v1/domains/other.example/diagnostics", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetDomainDiagnostics(ctx, "other.example")
			return err
		}},
		{name: "domains activate", method: "GET", path: "/v1/domains/other.example/activate", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.ActivateDomain(ctx, "other.example")
			return err
		}},
		{name: "domains usage", method: "GET", path: "/v1/domains/other.example/usage", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetDomainUsage(ctx, "other.example")
			return err
		}},

		{name: "mailboxes index", method: "GET", path: "/v1/domains/example.com/mailboxes", response: `{"mailboxes":[]}`, call: func(ctx context.Context, c *Client) error { _, err := c.ListMailboxes(ctx, "example.com"); return err }},
		{name: "mailboxes show", method: "GET", path: "/v1/domains/example.com/mailboxes/demo", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetMailbox(ctx, "example.com", "demo")
			return err
		}},
		{name: "mailboxes create", method: "POST", path: "/v1/domains/example.com/mailboxes", response: `{}`, expectedJSON: map[string]any{"local_part": "demo", "forwarding_to": "outside@example.net", "wildcard_sender": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateMailbox(ctx, "example.com", CreateMailboxRequest{LocalPart: "demo", ForwardingTo: "outside@example.net", WildcardSender: &falseValue})
			return err
		}},
		{name: "mailboxes update", method: "PUT", path: "/v1/domains/example.com/mailboxes/demo", response: `{}`, expectedJSON: map[string]any{"may_send": false, "recipient_denylist": []any{}}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateMailbox(ctx, "example.com", "demo", UpdateMailboxRequest{MaySend: &falseValue, RecipientDenylist: &emptyList})
			return err
		}},
		{name: "mailboxes delete", method: "DELETE", path: "/v1/domains/example.com/mailboxes/demo", response: ``, call: func(ctx context.Context, c *Client) error { return c.DeleteMailbox(ctx, "example.com", "demo") }},

		{name: "identities index", method: "GET", path: "/v1/domains/example.com/mailboxes/demo/identities", response: `{"identities":[]}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.ListIdentities(ctx, "example.com", "demo")
			return err
		}},
		{name: "identities show", method: "GET", path: "/v1/domains/example.com/mailboxes/demo/identities/example", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetIdentity(ctx, "example.com", "demo", "example")
			return err
		}},
		{name: "identities create", method: "POST", path: "/v1/domains/example.com/mailboxes/demo/identities", response: `{}`, expectedJSON: map[string]any{"local_part": "example", "password": "secret"}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateIdentity(ctx, "example.com", "demo", CreateIdentityRequest{LocalPart: "example", Password: "secret"})
			return err
		}},
		{name: "identities update", method: "PUT", path: "/v1/domains/example.com/mailboxes/demo/identities/example", response: `{}`, expectedJSON: map[string]any{"footer_active": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateIdentity(ctx, "example.com", "demo", "example", UpdateIdentityRequest{FooterActive: &falseValue})
			return err
		}},
		{name: "identities delete", method: "DELETE", path: "/v1/domains/example.com/mailboxes/demo/identities/example", response: ``, call: func(ctx context.Context, c *Client) error {
			return c.DeleteIdentity(ctx, "example.com", "demo", "example")
		}},

		{name: "forwardings index", method: "GET", path: "/v1/domains/example.com/mailboxes/demo/forwardings", response: `{"forwardings":[]}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.ListForwardings(ctx, "example.com", "demo")
			return err
		}},
		{name: "forwardings show", method: "GET", path: "/v1/domains/example.com/mailboxes/demo/forwardings/outside@example.net", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetForwarding(ctx, "example.com", "demo", "outside@example.net")
			return err
		}},
		{name: "forwardings create", method: "POST", path: "/v1/domains/example.com/mailboxes/demo/forwardings", response: `{}`, expectedJSON: map[string]any{"address": "outside@example.net", "is_active": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateForwarding(ctx, "example.com", "demo", CreateForwardingRequest{Address: "outside@example.net", IsActive: &falseValue})
			return err
		}},
		{name: "forwardings update", method: "PUT", path: "/v1/domains/example.com/mailboxes/demo/forwardings/outside@example.net", response: `{}`, expectedJSON: map[string]any{"remove_upon_expiry": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateForwarding(ctx, "example.com", "demo", "outside@example.net", UpdateForwardingRequest{RemoveUponExpiry: &falseValue})
			return err
		}},
		{name: "forwardings delete", method: "DELETE", path: "/v1/domains/example.com/mailboxes/demo/forwardings/outside@example.net", response: ``, call: func(ctx context.Context, c *Client) error {
			return c.DeleteForwarding(ctx, "example.com", "demo", "outside@example.net")
		}},

		{name: "aliases index", method: "GET", path: "/v1/domains/example.com/aliases", response: `{"address_aliases":[]}`, call: func(ctx context.Context, c *Client) error { _, err := c.ListAliases(ctx, "example.com"); return err }},
		{name: "aliases show", method: "GET", path: "/v1/domains/example.com/aliases/demo", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetAlias(ctx, "example.com", "demo")
			return err
		}},
		{name: "aliases create", method: "POST", path: "/v1/domains/example.com/aliases", response: `{}`, expectedJSON: map[string]any{"local_part": "demo", "destinations": []any{"target@example.com"}, "is_internal": false}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateAlias(ctx, "example.com", CreateAliasRequest{LocalPart: "demo", Destinations: []string{"target@example.com"}, IsInternal: &falseValue})
			return err
		}},
		{name: "aliases update", method: "PUT", path: "/v1/domains/example.com/aliases/demo", response: `{}`, expectedJSON: map[string]any{"destinations": []any{}}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateAlias(ctx, "example.com", "demo", UpdateAliasRequest{Destinations: &emptyList})
			return err
		}},
		{name: "aliases delete", method: "DELETE", path: "/v1/domains/example.com/aliases/demo", response: ``, call: func(ctx context.Context, c *Client) error { return c.DeleteAlias(ctx, "example.com", "demo") }},

		{name: "rewrites index", method: "GET", path: "/v1/domains/example.com/rewrites", response: `{"rewrites":[]}`, call: func(ctx context.Context, c *Client) error { _, err := c.ListRewrites(ctx, "example.com"); return err }},
		{name: "rewrites show", method: "GET", path: "/v1/domains/example.com/rewrites/demo", response: `{}`, call: func(ctx context.Context, c *Client) error {
			_, err := c.GetRewrite(ctx, "example.com", "demo")
			return err
		}},
		{name: "rewrites create", method: "POST", path: "/v1/domains/example.com/rewrites", response: `{}`, expectedJSON: map[string]any{"name": "demo", "local_part_rule": "demo-*", "destinations": []any{"target@example.com"}, "order_num": float64(0)}, call: func(ctx context.Context, c *Client) error {
			_, err := c.CreateRewrite(ctx, "example.com", CreateRewriteRequest{Name: "demo", LocalPartRule: "demo-*", Destinations: []string{"target@example.com"}, OrderNum: &zero})
			return err
		}},
		{name: "rewrites update", method: "PUT", path: "/v1/domains/example.com/rewrites/demo", response: `{}`, expectedJSON: map[string]any{"order_num": float64(0)}, call: func(ctx context.Context, c *Client) error {
			_, err := c.UpdateRewrite(ctx, "example.com", "demo", UpdateRewriteRequest{OrderNum: &zero})
			return err
		}},
		{name: "rewrites delete", method: "DELETE", path: "/v1/domains/example.com/rewrites/demo", response: ``, call: func(ctx context.Context, c *Client) error { return c.DeleteRewrite(ctx, "example.com", "demo") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := performRequest(t, tt.response, tt.call)
			if request.Method != tt.method || request.Path != tt.path {
				t.Fatalf("request = %s %s, want %s %s", request.Method, request.Path, tt.method, tt.path)
			}
			if request.Username != "admin@example.com" || request.Password != "secret" {
				t.Fatalf("Basic Auth = %q/%q", request.Username, request.Password)
			}
			if tt.expectedJSON == nil {
				return
			}
			if request.ContentType != "application/json" {
				t.Fatalf("Content-Type = %q", request.ContentType)
			}
			var body map[string]any
			if err := json.Unmarshal(request.Body, &body); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			for key, want := range tt.expectedJSON {
				if got := body[key]; !reflect.DeepEqual(got, want) {
					t.Errorf("body[%q] = %#v, want %#v; body = %s", key, got, want, request.Body)
				}
			}
		})
	}
}

func TestDomainFlexibleResponseFields(t *testing.T) {
	request := `{
		"name":"example.com",
		"tags":"work,business",
		"sender_denylist":"one@example.com,two@example.com",
		"recipient_denylist":[],
		"spam_aggressiveness":0
	}`
	var domain Domain
	if err := json.Unmarshal([]byte(request), &domain); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if !reflect.DeepEqual(domain.SenderDenylist, []string{"one@example.com", "two@example.com"}) {
		t.Fatalf("SenderDenylist = %#v", domain.SenderDenylist)
	}
	if domain.Name != "example.com" || !reflect.DeepEqual(domain.Tags, []string{"work", "business"}) {
		t.Fatalf("Domain = %+v", domain)
	}
	if domain.SpamAggressiveness != "0" {
		t.Fatalf("SpamAggressiveness = %q", domain.SpamAggressiveness)
	}
}

func TestCreateIdentityUsesCollectionPath(t *testing.T) {
	request := performRequest(t, `{}`, func(ctx context.Context, client *Client) error {
		_, err := client.CreateIdentity(ctx, "example.com", "demo", CreateIdentityRequest{LocalPart: "example", Name: "Example"})
		return err
	})
	if request.Method != "POST" || request.Path != "/v1/domains/example.com/mailboxes/demo/identities" {
		t.Fatalf("request = %s %s", request.Method, request.Path)
	}
}
