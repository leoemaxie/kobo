# Kobo — Backend Architecture & Folder Structure

This document defines the project layout, package responsibilities, and build
conventions for the Kobo Go backend. 

## Tech stack

- Language: Go 1.22+
- HTTP routing: `chi` (github.com/go-chi/chi/v5)
- Database: PostgreSQL 15+
- DB access: `sqlc` (generates typed Go from raw SQL) + `pgx/v5` as the driver
- Migrations: `golang-migrate`
- Config: environment variables loaded via `envconfig` or a thin custom loader, no `.env` parsing magic in production code paths
- Validation: `go-playground/validator` for request payload validation
- Logging: `log/slog` (standard library structured logging, Go 1.21+)
- Background jobs (reconciliation sweep, webhook retry, KYC tier checks): a simple internal worker pool driven by `time.Ticker`, no external job queue for v1. Revisit only if volume demands it.
- Testing: standard library `testing` + `testify/assert` for assertions, `httptest` for handler tests
- API contract: OpenAPI 3.1 spec (see `openapi.yaml`) is the source of truth. Handlers must match it exactly. Do not let the spec drift from the implementation.

## Top-level layout

```
core/
  cmd/
    api/
      main.go                 # entrypoint: wires config, db, router, starts HTTP server
    worker/
      main.go                 # entrypoint: reconciliation sweep + background jobs, separate binary
  internal/
    identity/
      model.go                # Identity struct, DisplayProfile struct
      service.go               # business logic: register, update, close, reopen
      repository.go             # Postgres queries (via sqlc-generated code)
      repository_test.go
      service_test.go
    account/
      model.go                 # VirtualAccount struct, AccountState enum
      lifecycle.go              # the state machine: ValidTransition(), Transition()
      lifecycle_test.go         # exhaustive table-driven tests of every transition
      service.go                # provisioning orchestration, calls monnify client
      repository.go
    ledger/
      model.go                  # LedgerEntry struct
      service.go                 # balance calculation, statement generation
      repository.go
    reconciliation/
      engine.go                  # core matching logic: webhook event -> ledger entry
      idempotency.go               # dedup by Monnify transaction reference
      sweep.go                     # fallback polling job against Transactions API
      engine_test.go
    exceptions/
      model.go                   # MisdirectedPayment / Exception struct
      service.go                  # flagging + resolution logic
      repository.go
    monnify/
      client.go                   # Monnify API client: provisioning, transactions, signature verify
      webhook.go                   # webhook payload parsing + signature verification
      types.go                     # Monnify API request/response types (external contract, kept separate from internal models)
      client_test.go               # uses a mock HTTP transport, never hits real Monnify in tests
    api/
      router.go                   # chi router setup, mounts all route groups
      middleware/
        auth.go                    # API key + HMAC request signature verification
        logging.go
        recover.go
      handlers/
        identities.go               # POST /v1/identities, GET, PATCH, /close, /reopen
        accounts.go                  # GET /v1/accounts/{id}/transactions, /statement
        exceptions.go                 # GET /v1/exceptions, POST /v1/exceptions/{id}/resolve
        webhooks.go                   # internal endpoint Monnify calls into
        health.go                     # GET /healthz
      dto/
        identity_dto.go               # request/response JSON shapes, separate from internal models
        account_dto.go
        exception_dto.go
      errors.go                       # shared API error type + HTTP status mapping
    platform/
      db/
        db.go                          # pgx pool setup
        sqlc/                           # sqlc-generated code lives here, do not hand-edit
        queries/                        # raw .sql files sqlc reads, hand-written, version controlled
          identities.sql
          accounts.sql
          ledger.sql
          exceptions.sql
      config/
        config.go                       # typed config struct loaded from env
      telemetry/
        logger.go                       # slog setup
  migrations/
    0001_init.up.sql
    0001_init.down.sql
    0002_lifecycle_states.up.sql
    0002_lifecycle_states.down.sql
    ...                                  # one migration per schema change, never edit an already-applied migration
  openapi.yaml                           # API contract, source of truth, see below
  docs/
    ARCHITECTURE.md                      # this file
    RECONCILIATION.md                    # detailed write-up of edge cases from the concept note, kept in sync with reconciliation/engine.go
    LIFECYCLE.md                         # state machine diagram + transition table, kept in sync with account/lifecycle.go
    AUTHENTICATION.md                  # detailed write-up of authentication flows, kept in sync with api/middleware/auth.go
    MONNIFY_INTEGRATION.md               # detailed write-up of monnify integration flows, kept in sync with monnify/client.go
  scripts/
    seed.go                              # seeds a sandbox integrator + test identities for local dev
  Makefile                               # make run, make test, make migrate-up, make sqlc-generate
  go.mod
  go.sum
  .env.sample
  README.md
```

## Package boundary rules (important for agent-driven development)

These rules exist so an AI coding agent making incremental changes does not
quietly violate the architecture:

1. **`internal/monnify` is the only package allowed to import an HTTP client pointed at Monnify's API.** No other package talks to Monnify directly. This keeps Monnify's API shape from leaking into business logic, and means the rest of the system can be tested without a sandbox connection.

2. **`internal/api/handlers` contains no business logic.** Handlers parse the request, call exactly one service method, and serialize the response. If a handler has an `if` statement that isn't about HTTP concerns (auth, validation errors, status codes), that logic belongs in a `service.go` instead.

3. **`internal/api/dto` types never leak into service or repository layers.** Services operate on internal domain models (`identity.Identity`, `account.VirtualAccount`, etc). Handlers convert DTO <-> domain model at the boundary. This is what lets the OpenAPI spec change without rewriting business logic, and vice versa.

4. **State transitions only happen through `account/lifecycle.go`.** No package sets `account.State = "CLOSED"` directly anywhere else in the codebase. Every transition goes through `lifecycle.Transition(current, event)`, which returns an error for any invalid transition. This is what makes the lifecycle state machine actually enforced rather than just documented.

5. **`reconciliation/idempotency.go` is consulted before any ledger write from a webhook or sweep.** No ledger entry is created without first checking the Monnify transaction reference against the idempotency table.

6. **Every new package gets a `_test.go` file in the same PR/commit that creates it.** Particularly `lifecycle_test.go` and `engine_test.go`, since the state machine and reconciliation logic are the two pieces directly named in the judging criteria — these should have the most thorough table-driven tests in the codebase.

## Suggested build order for an agent working through this incrementally

1. `platform/db` + migrations for `identities`, `virtual_accounts`, `ledger_entries`, `exceptions`, `idempotency_keys`, `api_integrators` tables.
2. `internal/identity` (model, repository, service) with unit tests against a real local Postgres (use `testcontainers-go` or a docker-compose Postgres for integration tests).
3. `internal/account/lifecycle.go` as a pure, dependency-free state machine — build and test this in isolation before wiring it to anything else.
4. `internal/monnify` client against Monnify's sandbox, with the provisioning call only.
5. `internal/account/service.go` wiring lifecycle + monnify client + repository together.
6. `internal/api` skeleton: router, middleware/auth.go, health.go — get a deployable, authenticated, empty API running early.
7. `handlers/identities.go` wired to `internal/identity` and `internal/account` — this is the first end-to-end vertical slice (register identity -> provision account -> see it in DB).
8. `internal/reconciliation` (engine, idempotency, sweep) + `handlers/webhooks.go`.
9. `internal/ledger` (statement generation) + `handlers/accounts.go`.
10. `internal/exceptions` + `handlers/exceptions.go`.
11. Background worker (`cmd/worker`) running the reconciliation sweep and KYC-tier-check jobs on a ticker.

This order is deliberate: it gets a real, narrow, end-to-end path (identity -> provisioned account) working before any reconciliation logic exists, which gives you something demoable early and isolates the hardest part (reconciliation edge cases) into its own well-tested package once the foundation is solid.
