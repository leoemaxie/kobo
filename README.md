# Kobo - Virtual Account Infrastructure

Kobo is a robust virtual account management system designed as a foundation for Nigerian fintech applications. It provides identity-anchored virtual accounts, real-time bank notifications (Push), and an SDK for seamless integration.

## Features

- **Identity-Anchored Accounts**: Every virtual account is tied to a verified identity, ensuring user accountability.
- **Push Notifications**: Built-in support for real-time Push notifications for transactions and account events.
- **SDK Library**: A comprehensive Go SDK (`pkg/kobo`) for easy client integration.
- **Resilient Architecture**: Features retry mechanisms and idempotent operations to handle network failures gracefully.

## Getting Started

### Prerequisites

- Go 1.22+

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd kobo
   ```

2. **Run the setup script**:
   This script will generate the necessary database migrations and seed the database with a test integrator.
   ```bash
   go run scripts/setup.go
   ```

### Running the Server

Start the API server:
```bash
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

## Usage

### Using the SDK

Import the SDK into your application:
```go
import "github.com/leoemaxie/kobo/pkg/kobo"
```

Initialize the client with your API credentials:
```go
client := kobo.New("your-api-key", "your-api-secret")
```

Create an identity and generate an account:
```go
identity, err := client.Identities.Create(
    context.Background(),
    "some_external_id",
    &kobo.IdentityCreateOptions{
        DisplayName: "John Doe",
    },
)

account, err := client.Accounts.Create(
    context.Background(),
    identity.ID,
    &kobo.AccountCreateOptions{
        BankName: "GTBank",
    },
)
```

### API Endpoints

Key endpoints available in the sandbox:

- `POST /v1/accounts`: Create an account for an identity.
- `GET /v1/accounts/{id}`: Get account details.
- `POST /v1/identities`: Create an identity.
- `POST /v1/webhooks/push`: Receive notifications from the Push provider.
