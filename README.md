# Migadu API in Go

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/z-xavier/migadu-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/z-xavier/migadu-go)](https://goreportcard.com/report/github.com/z-xavier/migadu-go)

`migadu-go` is a Go library for interfacing with the [Migadu API](https://www.migadu.com/api/).

## Installing

You need a working Go environment.

```shell
go get github.com/z-xavier/migadu-go
```

## Getting Started

Example:

```go
package main

import (
    "github.com/z-xavier/migadu-go"
)

client, err := migadu.New(os.Getenv("MIGADU_ADMIN_EMAIL"), os.Getenv("MIGADU_API_KEY"), "example.com")

// Incorrect API details will return an error
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

aliases, err := client.ListAliases(context.Background())
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

for _, alias := range aliases {
    fmt.Println(alias)
}

return
}
```
