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
	values    *url.Values
	header    *http.Header
	cookies   []*http.Cookie
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

func (b *httpReqBuilder) AddValues(key, value string) *httpReqBuilder {
	if b.values == nil {
		b.values = &url.Values{}
	}
	b.values.Add(key, value)
	return b
}

func (b *httpReqBuilder) SetValues(key, value string) *httpReqBuilder {
	if b.values == nil {
		b.values = &url.Values{}
	}
	b.values.Set(key, value)
	return b
}

func (b *httpReqBuilder) AddHeader(key, value string) *httpReqBuilder {
	if b.header == nil {
		b.header = &http.Header{}
	}
	b.header.Add(key, value)
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

func (b *httpReqBuilder) SetCookies(cookie ...*http.Cookie) *httpReqBuilder {
	if b.cookies == nil {
		b.cookies = make([]*http.Cookie, 0)
	}
	b.cookies = append(b.cookies, cookie...)
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

func (b *httpReqBuilder) SetBody(body io.Reader) *httpReqBuilder {
	b.body = body
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
	if b.values != nil {
		parse.RawQuery = b.values.Encode()
	}
	req, err := http.NewRequest(b.method, parse.String(), b.body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if b.cookies != nil {
		for _, cookie := range b.cookies {
			req.AddCookie(cookie)
		}
	}

	if b.header != nil {
		req.Header = *b.header
	}
	if b.basicAuth != nil {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}
	return req, nil
}
