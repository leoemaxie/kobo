# Kobo Core

The backend services powering Kobo.

## Getting Started

1. Set up your `.env` file using `.env.sample`.
2. Run migrations: `make migrate-up`
3. Start the API: `make run`
4. Start the background worker: `make worker`

## Architecture

Please refer to `docs/ARCHITECTURE.md` for package boundaries and architectural constraints.
