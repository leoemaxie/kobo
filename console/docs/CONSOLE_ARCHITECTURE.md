# Kobo Console — Architecture

High-level architecture for the Kobo Console: the fullstack SvelteKit
application integrators use to manage their Kobo account, API keys, billing,
and sandbox environment. This document grounds any AI coding agent working
on the console; it is deliberately high-level (not a folder-by-folder spec
like `core/docs/ARCHITECTURE.md`) so the agent has room to make idiomatic
SvelteKit decisions within clear boundaries.

This is a **separate application** from `core/` (the Go API backend). The
Console does not reimplement Kobo's business logic — it manages the things
around the API: who can call it, with what credentials, in which
environment, and what they're billed for.

---

## What the Console is for

Four responsibilities, and nothing beyond them for v1:

1. **Identity & access** — integrator signup, email verification, login, superadmin oversight.
2. **API credential management** — generate, view (once), rotate, and revoke API key/secret pairs, separately for sandbox and production.
3. **Billing visibility** — view usage, transaction volume, and charges tied to the pricing model in the concept note (provisioning fee, per-transaction fee, platform fee).
4. **Sandbox environment** — a fully isolated test mode with its own keys, its own data, and no path for sandbox activity to touch production data or real Monnify transactions.

The Console is **not** where reconciliation, ledger, or lifecycle logic
lives. It reads from and writes to its own tables (integrators, API
credentials, billing records, admin actions) and, where it needs Kobo
platform data (usage stats, transaction counts), it calls the Go API as a
client — it does not query the Go backend's Postgres database directly, even
though both may run against the same physical Postgres server. This keeps
the two systems loosely coupled: the Console can be redeployed, rebuilt, or
even replaced without touching Kobo core, and vice versa.

---

## Tech stack (locked)

- **Framework:** SvelteKit, fullstack mode — SSR pages + server-side API routes (`+page.server.ts` / `+server.ts`) as the backend, no separate Node/Express server.
- **Database:** PostgreSQL, separate schema (or separate database, decide based on hosting constraints) from Kobo core's Postgres. Recommended: same Postgres instance, separate database — `kobo_console` vs `kobo_core` — so the two applications can never accidentally join across each other's tables.
- **ORM:** Drizzle ORM — typed schema, migrations, good SvelteKit fit.
- **Auth:** Custom session-based auth (SvelteKit's `hooks.server.ts` + signed httpOnly cookies) rather than a third-party auth provider, since the console needs tight control over superadmin roles and sandbox/production key scoping. Passwords hashed with `argon2`.
- **Email verification & transactional email:** Unsend, as specified.
- **API key generation:** Cryptographically random secrets (`crypto.randomBytes`), never stored in plaintext — store a hash (e.g. SHA-256) for verification, show the raw secret to the user exactly once at generation/rotation time, same pattern as Stripe/GitHub personal access tokens.
- **Styling:** Tailwind CSS, matching the design direction to be finalized separately for the marketing/docs site — the console can share the same design tokens (colors, type scale) once that's locked, but is a distinct surface from the docs site (different app, likely different subdomain: `console.kobo.dev` vs `docs.kobo.dev`).
- **Deployment target:** Vercel or Fly.io (SvelteKit adapter-agnostic decision, not architecturally significant — pick based on where Postgres is hosted for lowest latency).

---

## High-level architecture diagram

```
                        Browser (integrator or superadmin)
                                    │
                                    ▼
                    ┌───────────────────────────────┐
                    │   SvelteKit Console (fullstack) │
                    │                                 │
                    │  Routes (SSR pages)              │
                    │   /login /signup /verify-email   │
                    │   /dashboard                      │
                    │   /api-keys (sandbox + production) │
                    │   /billing                         │
                    │   /admin/*  (superadmin only)       │
                    │                                 │
                    │  Server-side logic (+server.ts,  │
                    │  +page.server.ts, hooks.server.ts)│
                    │   - session/auth middleware       │
                    │   - API key generation/rotation   │
                    │   - billing aggregation            │
                    │   - superadmin action handlers     │
                    └───────────┬─────────────┬─────────┘
                                │             │
                    ┌───────────▼───┐   ┌─────▼─────────────────┐
                    │  Console DB    │   │  Unsend (email)         │
                    │  (Postgres,    │   │  - verification emails  │
                    │  Drizzle ORM)  │   │  - key rotation alerts   │
                    │                │   │  - billing notices       │
                    │  integrators   │   └────────────────────────┘
                    │  sessions      │
                    │  api_credentials│
                    │  billing_records│
                    │  admin_audit_log│
                    └────────────────┘
                                │
                                │ (read-only API calls, not direct DB access)
                                ▼
                    ┌────────────────────────┐
                    │   Kobo Core (Go API)     │
                    │   for usage/volume stats │
                    │   the console displays    │
                    └────────────────────────┘
```

---

## Core data model (high level, not full schema)

Kept intentionally sparse here — an agent should design the exact Drizzle
schema, but these are the entities that must exist and their key
relationships:

- **`integrators`** — one row per organization/team using Kobo. Holds
  company name, billing plan, and a flag for sandbox-only vs
  production-enabled (production access likely gated until Kobo core's own
  KYC/certification step is complete, per Monnify's "Post-certification, after
  KYC" note on the production environment).
- **`users`** — individual people who can log into the console, belonging to
  one `integrator` (or, for superadmins, belonging to none — see below).
  Holds email, password hash, `email_verified_at`, and a role
  (`owner`, `member`, `superadmin`).
- **`email_verification_tokens`** — short-lived, single-use tokens sent via
  Unsend, checked at `/verify-email`.
- **`sessions`** — server-side session records backing the httpOnly cookie;
  do not use JWTs for session state here, since sessions need to be
  revocable instantly (important for the "superadmin can lock out a
  compromised account" admin task).
- **`api_credentials`** — one row per generated key/secret pair. Fields:
  `integrator_id`, `environment` (`sandbox` | `production`), `key_id`
  (public, safe to log), `secret_hash` (never the raw secret),
  `created_at`, `rotated_at`, `revoked_at`. A credential is never deleted,
  only revoked — this preserves an audit trail and matches how Kobo core's
  own API auth expects to validate keys against a stable identifier even
  after rotation.
- **`billing_records`** — usage/charge line items pulled periodically from
  Kobo core's usage API (see "Console ↔ Core communication" below), stored
  here so the console has its own billing history independent of core's
  uptime.
- **`admin_audit_log`** — append-only log of every superadmin action (key
  revocation, account suspension, manual billing adjustment, etc.), who did
  it, and when. This is non-negotiable for a console with a superadmin role
  — every privileged action must be attributable.

---

## Sandbox vs. production separation

This is the most architecturally important decision in the console, so it's
worth stating as a hard rule rather than an implementation detail:

**Every `api_credentials` row is scoped to exactly one environment
(`sandbox` or `production`), and nothing in the console UI, API, or database
query layer should ever mix the two.** Concretely:

- Sandbox and production keys are generated, displayed, and rotated through
  the same UI, but always with a visible, unmistakable environment badge
  (color-coded, per the brand styling work still to be finalized).
- Sandbox keys map to Kobo core's sandbox environment
  (`sandbox.api.kobo.dev`, per the earlier ARCHITECTURE.md/openapi.yaml
  work) and, transitively, to Monnify's sandbox
  (`sandbox.api.monnify.com/v1`). Production keys map to the production
  equivalents. The console itself does not talk to Monnify directly — it only
  manages the Kobo-side credentials that will, in turn, let Kobo core talk
  to the correct Monnify environment.
- Billing records should be queryable filtered by environment, since
  sandbox usage is free/unmetered by definition (per the Monnify hackathon's
  own sandbox-for-all-hackathon-work model) and only production usage should
  ever generate a real billing record.
- A new integrator gets sandbox access immediately upon email verification.
  Production access is a separate, gated step (see Console ↔ Core
  communication below) — do not auto-provision production credentials at
  signup.

---

## Auth & verification flow

```
Signup (email + password)
        │
        ▼
Create `users` row, email_verified_at = null
        │
        ▼
Generate verification token → send via Unsend
        │
        ▼
User clicks link → /verify-email?token=...
        │
        ▼
Mark email_verified_at, create session, redirect to /dashboard
        │
        ▼
Sandbox API credentials become generatable from /api-keys
(production credentials remain gated)
```

Unverified users can log in (to resend the verification email, check
status) but every other console route redirects to a
"verify your email" screen until `email_verified_at` is set. Implement this
check once, in `hooks.server.ts`, not scattered across individual routes.

**Superadmin accounts** are not self-service. They are seeded directly in
the database (a one-time seed script, not a signup flow) and are not tied to
any `integrator_id`. Superadmin login uses the same session mechanism but
unlocks the `/admin/*` route tree, which regular users (even integrator
`owner` role) cannot reach — enforce this in `hooks.server.ts` alongside the
email-verification check, as a single, central authorization gate rather
than per-page checks that are easy to forget on a new route.

---

## Superadmin tasks (high level; exact scope for the agent to flesh out)

- View all integrators, their environment access level, and account status.
- Suspend/reinstate an integrator (blocks all their API credentials without
  deleting anything, reversible).
- Grant production access to an integrator that has completed
  KYC/certification (this is the manual gate mentioned above — v1 does not
  need this automated).
- Force-revoke a specific `api_credentials` row (e.g. reported as leaked).
- View the `admin_audit_log`.
- Issue manual billing adjustments (credits/corrections), logged.

Every one of these actions must write an `admin_audit_log` row before (or
within the same transaction as) the action taking effect — this is not
optional, it's what makes "superadmin" a safe role to have at all.

---

## Console ↔ Core communication

The console needs two kinds of data from Kobo core that it does not own:

1. **Usage/volume data for billing** — the console should call a
   (to-be-added) internal/admin endpoint on Kobo core, authenticated with a
   console-to-core service credential (not an integrator's own API key),
   to pull transaction counts and volume per integrator on a schedule (e.g.
   nightly job), and store the result in `billing_records`. This keeps
   Kobo core's API surface (`openapi.yaml`) focused on integrator-facing
   concerns, and keeps billing aggregation as a console-side job.
2. **Credential validation** — when the console generates or rotates a key,
   Kobo core needs to know about it (it's core's `internal/api/middleware/auth.go`
   that actually validates API keys on inbound Kobo API requests, per
   `core/docs/ARCHITECTURE.md`). Two ways to do this, pick one deliberately
   rather than defaulting: either (a) the console writes credentials
   directly to a shared `api_integrators`-adjacent table that both systems
   read (tighter coupling, simpler), or (b) the console calls a core-side
   admin endpoint to register/revoke credentials (looser coupling, cleaner
   boundary, more moving parts for a hackathon timeline). Given the
   hackathon timeline, (a) is the pragmatic choice: let both apps read the
   same `api_integrators` / credential table in Postgres, but only the
   console writes to it. Document this exception clearly since it's a
   deliberate deviation from the "console never touches core's DB directly"
   rule — it applies to this one shared credentials table, not to any other
   Kobo core table.

This is a decision worth revisiting post-hackathon; for now, document it as
a named exception rather than letting it become an unstated inconsistency.

---

For an agent (or Leo) to decide during implementation, since these are
build-time details, not architectural ones:

- Exact Drizzle schema field names/types.
- Specific SvelteKit route file layout (`src/routes/...`).
- Component structure / design system implementation (pending brand and
  styling decisions still to be finalized).
- Whether billing sync from core runs as a SvelteKit scheduled function,
  a separate cron, or a lightweight worker — any of these are fine
  architecturally as long as it's a distinct, observable job, not inline in
  a page request.
- Rate limiting on login/signup/verification endpoints — should exist, exact
  mechanism (in-memory, Redis, Postgres-backed) is an implementation detail.

---

## Summary of hard rules for any agent building this

1. Console and Core are separate apps with separate primary databases;
   the one deliberate exception is the shared credentials table (see
   "Console ↔ Core communication").
2. Sandbox and production credentials are never mixed in a query, a UI
   view, or a billing calculation.
3. Raw API secrets are shown exactly once, never stored in plaintext,
   never re-displayable after generation.
4. Every superadmin action writes an audit log entry as part of the same
   operation, not as an afterthought.
5. Production API access is a manually gated step, not automatic on signup.
6. Email verification is enforced centrally (`hooks.server.ts`), not
   per-route.
