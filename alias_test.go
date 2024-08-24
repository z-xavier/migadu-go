package migadu

import (
	"context"
	"os"
	"reflect"
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
		want    *[]Alias
		wantErr bool
	}{
		{
			name: "Test List Aliases",
			fields: fields{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &[]Alias{},
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
			if !reflect.DeepEqual(aliases, tt.want) {
				t.Errorf("ListAliases() got = %v, want %v", aliases[0], tt.want)
			}
		})
	}
}
