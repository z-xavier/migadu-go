package migadu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpReqBuilder struct {
	method    string
	host      string
	path      string
	header    *http.Header
	body      io.Reader
	basicAuth *basicAuth
	err       error
}

type basicAuth struct {
	username, password string
}

func newReqBuilder() *httpReqBuilder {
	return &httpReqBuilder{}
}

func (b *httpReqBuilder) SetMethod(method string) *httpReqBuilder {
	b.method = method
	return b
}

func (b *httpReqBuilder) SetHost(host string) *httpReqBuilder {
	b.host = strings.TrimSuffix(host, "/")
	return b
}

func (b *httpReqBuilder) AddPath(path string) *httpReqBuilder {
	b.path += fmt.Sprintf("/%s", url.PathEscape(strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")))
	return b
}

func (b *httpReqBuilder) AddRestfulPath(key, value string) *httpReqBuilder {
	b.path += fmt.Sprintf("/%s/%s", url.PathEscape(key), url.PathEscape(value))
	return b
}

func (b *httpReqBuilder) SetHeaderContentTypeJson() *httpReqBuilder {
	return b.SetHeader("Content-Type", "application/json")
}

func (b *httpReqBuilder) SetHeader(key, value string) *httpReqBuilder {
	if b.header == nil {
		b.header = &http.Header{}
	}
	b.header.Set(key, value)
	return b
}

func (b *httpReqBuilder) SetBasicAuth(username, password string) *httpReqBuilder {
	b.basicAuth = &basicAuth{username, password}
	return b
}

func (b *httpReqBuilder) SetBodyJson(body interface{}) *httpReqBuilder {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		b.err = fmt.Errorf("encode JSON request body: %w", err)
		return b
	}
	b.body = bytes.NewBuffer(jsonStr)
	return b
}

func (b *httpReqBuilder) Build() (*http.Request, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.method == "" {
		return nil, fmt.Errorf("method is required")
	}
	if b.host == "" {
		return nil, fmt.Errorf("host is required")
	}
	parse, err := url.Parse(b.host + b.path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(b.method, parse.String(), b.body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if b.header != nil {
		req.Header = *b.header
	}
	if b.basicAuth != nil {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}
	return req, nil
}
