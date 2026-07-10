# Migadu API in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/z-xavier/migadu-go.svg)](https://pkg.go.dev/github.com/z-xavier/migadu-go)
[![CI](https://github.com/z-xavier/migadu-go/actions/workflows/ci.yml/badge.svg)](https://github.com/z-xavier/migadu-go/actions/workflows/ci.yml)

`migadu-go` is a Go library for interfacing with the [Migadu API](https://www.migadu.com/api/).

## Installing

You need a working Go environment.

```shell
go get github.com/z-xavier/migadu-go
```

## Getting started

Create one account-level client and pass the domain explicitly to domain-scoped operations. Constructing a client validates required arguments but does not make a network request.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    migadu "github.com/z-xavier/migadu-go"
)

func main() {
    client, err := migadu.New(
        os.Getenv("MIGADU_ADMIN_EMAIL"),
        os.Getenv("MIGADU_API_KEY"),
    )
    if err != nil {
        log.Fatal(err)
    }

    aliases, err := client.ListAliases(context.Background(), "example.com")
    if err != nil {
        log.Fatal(err)
    }
    for _, alias := range aliases {
        fmt.Println(alias.Address)
    }
}
```

The same client can operate on every domain visible to the authenticated account:

```go
domains, err := client.ListDomains(ctx)
mailboxes, err := client.ListMailboxes(ctx, "example.com")
aliases, err := client.ListAliases(ctx, "other.example")
```

Create requests expose all documented writable fields. Update requests use pointers so callers can explicitly send `false`, `0`, an empty string, or an empty list:

```go
maySend := false
mailbox, err := client.UpdateMailbox(ctx, "example.com", "demo", migadu.UpdateMailboxRequest{
    MaySend: &maySend,
})
```

Non-success responses are returned as `*migadu.APIError` with the HTTP status, Migadu error code, message, and raw response body:

```go
var apiErr *migadu.APIError
if errors.As(err, &apiErr) {
    log.Printf("Migadu error %d %s: %s", apiErr.StatusCode, apiErr.Code, apiErr.Message)
}
```

The SDK covers the documented Domains, Mailboxes, Identities, Forwardings, Aliases, and Rewrites endpoints. Every request accepts a `context.Context`; the default per-request timeout is 30 seconds and can be changed through `Client.Timeout`.
