# Kobo Go SDK

Official Go SDK for the Kobo API. This SDK is built with zero external dependencies and uses the standard `net/http` package.

## Features
- Zero external dependencies
- Fully typed models
- Context-aware methods (`context.Context`)
- Resource-scoped sub-clients (Identities, Accounts, Exceptions)

## Installation

```bash
go get github.com/leoemaxie/kobo-sdk-go
```

### Using Locally

To use the SDK locally without fetching from a remote repository, add a `replace` directive to your `go.mod` file pointing to your local clone:

```go
require github.com/leoemaxie/kobo-sdk-go v0.0.0

replace github.com/leoemaxie/kobo-sdk-go => /absolute/path/to/kobo/sdks/go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/leoemaxie/kobo-sdk-go/kobo"
)

func main() {
	// Initialize the client
	client := kobo.New("kobo_live_pk_...", "kobo_live_sk_...")

	// Or use the sandbox environment
	// client := kobo.NewSandbox("kobo_test_pk_...", "kobo_test_sk_...")

	ctx := context.Background()

	// Check API health
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Health: %s\n", health.Status)

	// Create an identity
	req := kobo.CreateIdentityRequest{
		ExternalReference: "customer_12345",
		DisplayName:       "Jane Doe",
		KYCTierHint:       kobo.Ptr(kobo.KYCTier1),
		Metadata: map[string]interface{}{
			"source": "web_signup",
		},
	}

	identity, err := client.Identities.Create(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create identity: %v", err)
	}
	fmt.Printf("Created Identity: %s\n", identity.ID)

	// Fetch an identity
	fetchedIdentity, err := client.Identities.Get(ctx, identity.ID)
	if err != nil {
		log.Fatalf("Failed to fetch identity: %v", err)
	}
	fmt.Printf("Identity State: %s\n", fetchedIdentity.State)

	// List transactions for an account
	// if fetchedIdentity.VirtualAccount != nil {
	//     opts := kobo.ListTransactionsOptions{}
	//     opts.Limit = kobo.Ptr(50)
	//     // Assuming you have the accountID
	//     // page, err := client.Accounts.ListTransactions(ctx, accountID, opts)
	// }
}
```

## Error Handling

The SDK returns a `*kobo.APIError` when the API returns a non-2xx response. You should branch your error handling logic on the `Code` property, not the `Message` string.

```go
identity, err := client.Identities.Get(ctx, "invalid-id")
if err != nil {
    if apiErr, ok := err.(*kobo.APIError); ok {
        if apiErr.Code == "not_found" {
            fmt.Println("Identity not found.")
        } else {
            fmt.Printf("API Error: %s\n", apiErr.Code)
        }
    } else {
        fmt.Printf("Unknown Error: %v\n", err)
    }
}
```

## License

MIT
