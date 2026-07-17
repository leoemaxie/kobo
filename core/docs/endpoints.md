# Core API Endpoints

This document outlines the available endpoints in the Kobo Core API.

## Public Routes

- `GET /healthz`: Health check endpoint.

## API v1 Routes

Base URL: `/v1`

### Schema

- `GET /v1/`: Returns the OpenAPI JSON schema.

### Admin Routes

- `POST /v1/admin/integrators`: Provision a new integrator.
- `POST /v1/admin/billing/checkout`: Create a billing checkout.

### Protected Routes

These routes require authentication, a valid request ID, and client IP tracking.

#### Identities
- `POST /v1/identities`: Create a new identity.
- `GET /v1/identities/{id}`: Retrieve a specific identity by ID.
- `PATCH /v1/identities/{id}`: Update an existing identity.
- `POST /v1/identities/{id}/close`: Close a specific identity.
- `POST /v1/identities/{id}/reopen`: Reopen a closed identity.

#### Ledger & Accounts
- `GET /v1/accounts/{accountId}/transactions`: Retrieve transactions for a specific account.
- `GET /v1/accounts/{accountId}/statement`: Retrieve a statement for a specific account.

#### Exceptions
- `GET /v1/exceptions`: List open exceptions.
- `POST /v1/exceptions/{exceptionId}/resolve`: Resolve a specific exception.

### Webhooks

- `POST /v1/webhooks/monnify`: Handle incoming Monnify webhooks.
