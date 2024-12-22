package migadu

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestClient_GetDomains(t *testing.T) {
	type fields struct {
		email  string
		apiKey string
		domain string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Get Domains",
			fields: fields{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.fields.email, tt.fields.apiKey, tt.fields.domain)
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}
			domains, err := c.GetDomains(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetDomains() got = %v", domains)
		})
	}
}
