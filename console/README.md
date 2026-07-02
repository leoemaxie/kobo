# Kobo Console

The **Kobo Console** is a fullstack SvelteKit application that allows integrators to manage their Kobo accounts, API credentials, billing, and sandbox environments. It serves as the developer portal for interacting with the Kobo infrastructure.

## 🏗 Architecture & Stack

This application is structurally independent of the Kobo Core Go API and has its own isolated domain, authentication flow, and database schema, enforcing a strict zero-trust boundary.

- **Framework**: SvelteKit (Fullstack mode)
- **Styling**: Tailwind CSS v4 (using the dark "ClickHouse" aesthetic with `Void Black` backgrounds and `Electric Lime` accents)
- **Database**: PostgreSQL
- **ORM**: Drizzle ORM
- **Authentication**: Custom session-based auth (httpOnly cookies) using Argon2 for password hashing
- **Deployment**: Vercel / Fly.io

*For a detailed architectural breakdown, see [CONSOLE_ARCHITECTURE](docs/CONSOLE_ARCHITECTURE.md).*

## 🚀 Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/en) (v18+)
- [pnpm](https://pnpm.io/)
- PostgreSQL instance running locally or remotely

### Installation

1. Install dependencies:
   ```bash
   pnpm install
   ```

2. Set up your environment variables by copying `.env.example` to `.env` (if applicable) and configuring your PostgreSQL connection string:
   ```bash
   DATABASE_URL="postgres://user:password@localhost:5432/kobo"
   ```

3. Generate the SvelteKit types:
   ```bash
   pnpm svelte-kit sync
   ```

4. Start the development server:
   ```bash
   pnpm run dev
   ```

Navigate to `http://localhost:5173` to view the application. The application will automatically redirect unauthenticated users to the `/login` route.

## 🛡️ Superadmin Oversight

The console comes with built-in capabilities to monitor and manage integrators, grant production access, suspend fraudulent accounts, and issue manual billing adjustments. All privileged actions are permanently recorded in the immutable `admin_audit_log` table.

Superadmins can access these features at the `/admin` route.

## 📁 Key Directories

- `src/routes/`: SvelteKit file-based routing. Includes `/(app)`, `/login`, `/admin`, etc.
- `src/lib/components/`: Reusable Tailwind v4 UI components (Cards, Buttons, Inputs, Navbars).
- `src/lib/server/`: Server-only code, including Drizzle ORM database schemas (`db/schema.ts`) and session utilities (`auth/session.ts`).
- `static/`: Static assets like the Kobo logo, favicon, and the web app manifest.
