# Kobo Console

The **Kobo Console** is a fullstack SvelteKit application that allows integrators to manage their Kobo accounts, API credentials, billing, and sandbox environments. It serves as the developer portal for interacting with the Kobo infrastructure.

## Architecture & Stack

This application is structurally independent of the Kobo Core Go API and has its own isolated domain, authentication flow, and database schema, enforcing a strict zero-trust boundary.

- **Framework**: SvelteKit (Fullstack mode)
- **Styling**: Tailwind CSS v4 (using the dark "ClickHouse" aesthetic with `Void Black` backgrounds and `Electric Lime` accents)
- **Database**: PostgreSQL
- **ORM**: Drizzle ORM
- **Authentication**: Custom session-based auth (httpOnly cookies) using Argon2 for password hashing
- **Deployment**: Vercel / Fly.io

*For a detailed architectural breakdown, see [CONSOLE_ARCHITECTURE](docs/CONSOLE_ARCHITECTURE.md).*

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/en) (v18+)
- [pnpm](https://pnpm.io/)
- PostgreSQL instance running locally or remotely

### Installation & Running

1. **Install dependencies**:
   ```bash
   pnpm install
   ```

2. **Configure environment**:
   Copy `.env.example` to `.env` and configure your credentials:
   ```bash
   cp .env.example .env
   ```
   *(Ensure `DATABASE_URL` points to your Postgres instance and `UNSEND_API_KEY` is provided for emails).*

3. **Start the development server**:
   ```bash
   pnpm run dev
   ```
   Navigate to `http://localhost:3000`. Unauthenticated users will automatically be redirected to `/auth/login`.

### Database Setup & Drizzle ORM

This project uses [Drizzle ORM](https://orm.drizzle.team/) purely as a database client to safely read and write data. 

**Architectural Boundary**: The Console application **does not** own or manage the database schema. The Go `core` repository is the single source of truth for the database.
You should **never** run `drizzle-kit push` or `drizzle-kit generate` from this directory. Doing so will attempt to delete core tables and cause data loss.

If you need to make schema changes, write your SQL migrations in `kobo/core/migrations/`, apply them using the Go backend's tooling, and then update `src/lib/server/db/schema.ts` to reflect the new structure.

1. **Explore the database**:
   Drizzle provides a built-in GUI to view your database records locally:
   ```bash
   pnpm run db:studio
   ```

## Superadmin Oversight

The console comes with built-in capabilities to monitor and manage integrators, grant production access, suspend fraudulent accounts, and issue manual billing adjustments. All privileged actions are permanently recorded in the immutable `admin_audit_log` table.

Superadmins can access these features at the `/admin` route.

## Key Directories

- `src/routes/`: SvelteKit file-based routing. Includes `/(app)`, `/login`, `/admin`, etc.
- `src/lib/components/`: Reusable Tailwind v4 UI components (Cards, Buttons, Inputs, Navbars).
- `src/lib/server/`: Server-only code, including Drizzle ORM database schemas (`db/schema.ts`) and session utilities (`auth/session.ts`).
- `static/`: Static assets like the Kobo logo, favicon, and the web app manifest.
