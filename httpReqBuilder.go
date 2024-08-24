package migadu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var httpReqPool = sync.Pool{
	New: func() interface{} {
		return &HttpReqBuilder{}
	},
}

type HttpReqBuilder struct {
	method    string
	host      string
	path      string
	values    *url.Values
	header    *http.Header
	body      io.Reader
	basicAuth *basicAuth
}

type basicAuth struct {
	username, password string
}

func NewReqBuilder() *HttpReqBuilder {
	return httpReqPool.Get().(*HttpReqBuilder)
}

func (b *HttpReqBuilder) free() {
	b.method = ""
	b.host = ""
	b.path = ""
	b.values = nil
	b.header = nil
	b.body = nil
	b.basicAuth = nil
	httpReqPool.Put(b)
}

func (b *HttpReqBuilder) SetMethod(method string) *HttpReqBuilder {
	b.method = method
	return b
}

func (b *HttpReqBuilder) SetHost(host string) *HttpReqBuilder {
	b.host += host
	return b
}

func (b *HttpReqBuilder) AddPath(path string) *HttpReqBuilder {
	b.path += fmt.Sprintf("/%s", url.PathEscape(strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")))
	return b
}

func (b *HttpReqBuilder) AddRestfulPath(key, value string) *HttpReqBuilder {
	b.path += fmt.Sprintf("/%s/%s", url.PathEscape(key), url.PathEscape(value))
	return b
}

func (b *HttpReqBuilder) AddValues(key, value string) *HttpReqBuilder {
	if b.values == nil {
		b.values = &url.Values{}
	}
	b.values.Add(key, value)
	return b
}

func (b *HttpReqBuilder) SetValues(key, value string) *HttpReqBuilder {
	if b.values == nil {
		b.values = &url.Values{}
	}
	b.values.Set(key, value)
	return b
}

func (b *HttpReqBuilder) AddHeader(key, value string) *HttpReqBuilder {
	if b.header == nil {
		b.header = &http.Header{}
	}
	b.header.Add(key, value)
	return b
}

func (b *HttpReqBuilder) SetHeaderContentTypeJson() *HttpReqBuilder {
	return b.SetHeader("Content-Type", "application/json")
}

func (b *HttpReqBuilder) SetHeader(key, value string) *HttpReqBuilder {
	if b.header == nil {
		b.header = &http.Header{}
	}
	b.header.Set(key, value)
	return b
}

func (b *HttpReqBuilder) SetBasicAuth(username, password string) *HttpReqBuilder {
	b.basicAuth = &basicAuth{username, password}
	return b
}

func (b *HttpReqBuilder) SetBodyJson(body interface{}) *HttpReqBuilder {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		jsonStr = []byte{}
	}
	b.body = bytes.NewBuffer(jsonStr)
	return b
}

func (b *HttpReqBuilder) SetBody(body io.Reader) *HttpReqBuilder {
	b.body = body
	return b
}

func (b *HttpReqBuilder) Build() (*http.Request, error) {
	defer b.free()

	if b.method == "" {
		return nil, fmt.Errorf("method is required")
	}
	if b.host == "nil" {
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
		return nil, errors.Join(err, fmt.Errorf("error creating request"))
	}

	if b.header != nil {
		req.Header = *b.header
	}
	if b.basicAuth != nil {
		req.SetBasicAuth(b.basicAuth.username, b.basicAuth.password)
	}
	return req, nil
}
