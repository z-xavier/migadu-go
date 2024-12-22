package migadu

import (
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestClient_New(t *testing.T) {
	type args struct {
		email  string
		apiKey string
		domain string
	}
	tests := []struct {
		name    string
		args    *args
		want    *Client
		wantErr bool
	}{
		{
			name: "Test New apiKey invalid or expired",
			args: &args{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: "foo",
				domain: "example.com",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test New success",
			args: &args{
				email:  os.Getenv("MIGADU_ADMIN_EMAIL"),
				apiKey: os.Getenv("MIGADU_API_KEY"),
				domain: os.Getenv("MIGADU_DOMAIN"),
			},
			want: &Client{
				Email:      os.Getenv("MIGADU_ADMIN_EMAIL"),
				APIKey:     os.Getenv("MIGADU_API_KEY"),
				Domain:     os.Getenv("MIGADU_DOMAIN"),
				Timeout:    DefaultTimeout,
				HTTPClient: http.DefaultClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.email, tt.args.apiKey, tt.args.domain)
			if err != nil != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
