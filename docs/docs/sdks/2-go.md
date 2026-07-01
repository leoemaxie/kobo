---
sidebar_position: 2
title: With Go
---

# Go SDK

The official Go SDK for Kobo is designed to be highly concurrent, robust, and easy to integrate into your microservices.

## Installation

```bash
go get github.com/leoemaxie/kobo/sdks/go/kobo
```

## Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/leoemaxie/kobo/sdks/go/kobo"
)

func main() {
    // Production
    client := kobo.New("kobo_live_pk_...", "kobo_live_sk_...")

    // Sandbox
    sandboxClient := kobo.NewSandbox("kobo_test_pk_...", "kobo_test_sk_...")

    // Health Check
    if _, err := client.Health(context.Background()); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
}
```

---

## Identities

### Create Identity
```go
identity, err := client.Identities.Create(ctx, kobo.CreateIdentityRequest{
    ExternalReference: "user-12345",
    Profile: kobo.Profile{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john.doe@example.com",
    },
})
```

### Get Identity
```go
identity, err := client.Identities.Get(ctx, "idx_123abc")
```

### Update Identity
```go
identity, err := client.Identities.Update(ctx, "idx_123abc", kobo.UpdateIdentityRequest{
    Profile: kobo.UpdateProfile{
        DisplayName: kobo.String("John D."),
    },
})
```

### Close Identity
```go
identity, err := client.Identities.Close(ctx, "idx_123abc", kobo.CloseIdentityRequest{
    Reason: "User requested",
})
```

### Reopen Identity
```go
identity, err := client.Identities.Reopen(ctx, "idx_123abc")
```

---

## Accounts

### List Transactions
```go
limit := 50
page, err := client.Accounts.ListTransactions(ctx, "acc_123abc", kobo.ListTransactionsOptions{
    PaginationOptions: kobo.PaginationOptions{Limit: &limit},
})
```

### Get Statement
```go
period := "2024-03"
statement, err := client.Accounts.GetStatement(ctx, "acc_123abc", kobo.GetStatementOptions{
    Period: &period,
})
```

---

## Exceptions

### List Exceptions
```go
status := kobo.ExceptionStatusOpen
page, err := client.Exceptions.List(ctx, kobo.ListExceptionsOptions{
    Status: &status,
})
```

### Resolve Exception
```go
exception, err := client.Exceptions.Resolve(ctx, "exc_123abc", kobo.ResolveExceptionRequest{
    Action: "refund",
    Reason: "Refund requested by sender",
})
```
