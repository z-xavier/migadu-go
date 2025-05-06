package migadu

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestClient_ListAliases(t *testing.T) {
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
			name: "Test List Aliases",
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
			aliases, err := c.ListAliases(tt.args.ctx)
			if err != nil != tt.wantErr {
				t.Errorf("ListAliases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if aliases == nil {
				t.Errorf("ListAliases() error, aliases is nil")
				return
			}
			for i, alias := range aliases {
				fmt.Printf("ListAliases() got %d: \n %+v \n", i, alias)
			}
		})
	}
}

func TestClient_GetAlias(t *testing.T) {
	type fields struct {
		email  string
		apiKey string
		domain string
	}
	type args struct {
		ctx       context.Context
		localPart string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Get Alias",
			fields: fields{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			args: args{
				ctx:       context.Background(),
				localPart: "test",
			},
			wantErr: false,
		},
		{
			name: "Test Get Alias saint2584",
			fields: fields{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			args: args{
				ctx:       context.Background(),
				localPart: "saint2584",
			},
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
			alias, err := c.GetAlias(tt.args.ctx, tt.args.localPart)
			if err != nil != tt.wantErr {
				t.Errorf("GetAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetAlias() got = %v", alias)
		})
	}
}

func TestClient_NewAlias(t *testing.T) {
	type fields struct {
		email  string
		apiKey string
		domain string
	}
	type args struct {
		ctx          context.Context
		localPart    string
		destinations []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test New Alias",
			fields: fields{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			args: args{
				ctx:       context.Background(),
				localPart: "test",
				destinations: []string{
					"TestByMigaduGo",
				},
			},
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
			alias, err := c.NewAlias(tt.args.ctx, tt.args.localPart, tt.args.destinations)
			if err != nil != tt.wantErr {
				t.Errorf("GetAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("GetAlias() got = %v", alias)
		})
	}
}
