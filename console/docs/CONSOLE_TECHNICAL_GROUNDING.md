# Kobo Console — Technical Grounding: DB, Schema, and Sessions

This supersedes the "Console ↔ Core communication" section of
`CONSOLE_ARCHITECTURE.md` with a concrete decision: **one Postgres database,
shared by both Kobo Core (Go) and Kobo Console (SvelteKit)**. The only table
either app writes across the boundary is `api_integrators` — Console owns
writes to it, Core only reads from it. Every other table stays strictly
owned by one side. This document is the ground truth for schema and session
implementation; hand it directly to a coding agent alongside
`CONSOLE_ARCHITECTURE.md`.

---

## 1. One database, two schemas, one shared table

```
Postgres instance: kobo
│
├── schema: core          (owned by Go backend, migrated via golang-migrate)
│     identities
│     virtual_accounts
│     ledger_entries
│     idempotency_keys
│     exceptions
│     identity_events
│
├── schema: console        (owned by SvelteKit, migrated via Drizzle Kit)
│     users
│     sessions
│     email_verification_tokens
│     billing_records
│     admin_audit_log
│
└── api_integrators         (shared table, lives in `console` schema,
                              Core connects with a role that has SELECT-only
                              grant on this one table)
```

**Why one DB instead of two:** avoids cross-database queries/replication for
the one table that genuinely needs to be shared, and keeps hackathon
infrastructure simple — one Postgres instance to provision, one connection
string family, one backup target. The schema separation (`core` vs
`console`) is what preserves the "don't touch each other's tables" rule
without needing two physical databases.

**Why `api_integrators` lives in the `console` schema, not `core`:** Console
is the system of record for integrator identity, billing plan, and
credentials — that's its job. Core needs to read `api_integrators` only to
validate an inbound request's API key and resolve which integrator
namespace to scope the query to. Core never creates, updates, or deletes
rows here.

### Postgres roles and grants (set up once, referenced by both apps' connection configs)

```sql
-- Run once during infra setup, not part of either app's migrations.

CREATE ROLE kobo_core_app LOGIN PASSWORD '...';
CREATE ROLE kobo_console_app LOGIN PASSWORD '...';

CREATE SCHEMA core AUTHORIZATION kobo_core_app;
CREATE SCHEMA console AUTHORIZATION kobo_console_app;

-- Core's app role gets full rights on its own schema, read-only on the one
-- shared table.
GRANT ALL ON SCHEMA core TO kobo_core_app;
GRANT ALL ON ALL TABLES IN SCHEMA core TO kobo_core_app;
GRANT USAGE ON SCHEMA console TO kobo_core_app;
GRANT SELECT ON console.api_integrators TO kobo_core_app;

-- Console's app role gets full rights on its own schema, including the
-- shared table (it owns writes to api_integrators).
GRANT ALL ON SCHEMA console TO kobo_console_app;
GRANT ALL ON ALL TABLES IN SCHEMA console TO kobo_console_app;
```

This is enforced at the database level, not just by convention in
application code — if Core's Go code somehow tried to `UPDATE
api_integrators`, Postgres itself rejects it. That's a stronger guarantee
than "pase don't do this" in a doc, and worth the five minutes of setup.

Both apps' `DATABASE_URL` connection strings point at the same Postgres
instance and database, differing only in the role/credentials and (for
`search_path`) which schema they default to.

---

## 2. Drizzle schema

All Console-owned tables plus the shared `api_integrators` table, since
Console is the one that migrates it.

```typescript
// src/lib/server/db/schema.ts
import {
  pgSchema,
  uuid,
  text,
  timestamp,
  boolean,
  bigint,
  jsonb,
  pgEnum,
  uniqueIndex,
  index,
} from 'drizzle-orm/pg-core';

// All Console tables live in the `console` Postgres schema, matching the
// grants set up in section 1. Drizzle's pgSchema() maps directly to this.
export const consoleSchema = pgSchema('console');

// ---------------------------------------------------------------------------
// Shared table: api_integrators
// Console has full read/write. Core (Go/sqlc) has SELECT-only via Postgres
// grants, and defines its own read-only view of this same table in its own
// sqlc queries — see core/internal/platform/db/queries/integrators.sql
// (to be added to the Go side, mirroring this shape).
// ---------------------------------------------------------------------------

export const environmentEnum = consoleSchema.enum('environment', ['sandbox', 'production']);
export const integratorStatusEnum = consoleSchema.enum('integrator_status', [
  'active',
  'suspended',
]);

export const apiIntegrators = consoleSchema.table('api_integrators', {
  id: uuid('id').primaryKey().defaultRandom(),
  name: text('name').notNull(),
  status: integratorStatusEnum('status').notNull().default('active'),
  productionAccessGranted: boolean('production_access_granted').notNull().default(false),
  productionAccessGrantedAt: timestamp('production_access_granted_at', {
    withTimezone: true,
  }),
  productionAccessGrantedBy: uuid('production_access_granted_by'), // references users.id, nullable, set by superadmin
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
  updatedAt: timestamp('updated_at', { withTimezone: true }).notNull().defaultNow(),
});

// ---------------------------------------------------------------------------
// api_credentials: one row per generated key/secret pair, scoped to exactly
// one environment. Never deleted, only revoked (revokedAt set).
// ---------------------------------------------------------------------------

export const apiCredentials = consoleSchema.table(
  'api_credentials',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    integratorId: uuid('integrator_id')
      .notNull()
      .references(() => apiIntegrators.id),
    environment: environmentEnum('environment').notNull(),
    keyId: text('key_id').notNull().unique(), // public identifier, safe to log, e.g. "kobo_sandbox_a1b2c3"
    secretHash: text('secret_hash').notNull(), // SHA-256 of the raw secret; raw secret is never stored
    label: text('label'), // optional integrator-supplied name, e.g. "CI pipeline key"
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    createdBy: uuid('created_by').notNull(), // references users.id
    rotatedAt: timestamp('rotated_at', { withTimezone: true }),
    revokedAt: timestamp('revoked_at', { withTimezone: true }),
    revokedBy: uuid('revoked_by'), // references users.id, nullable
    revokedReason: text('revoked_reason'),
  },
  (table) => ({
    // Hot path: Core's auth middleware looks up by keyId on every request.
    keyIdIdx: uniqueIndex('idx_api_credentials_key_id').on(table.keyId),
    integratorEnvIdx: index('idx_api_credentials_integrator_env').on(
      table.integratorId,
      table.environment,
    ),
  }),
);

// ---------------------------------------------------------------------------
// users
// ---------------------------------------------------------------------------

export const userRoleEnum = consoleSchema.enum('user_role', ['owner', 'member', 'superadmin']);

export const users = consoleSchema.table(
  'users',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    // Nullable: superadmins are not tied to an integrator.
    integratorId: uuid('integrator_id').references(() => apiIntegrators.id),
    email: text('email').notNull().unique(),
    passwordHash: text('password_hash').notNull(), // argon2
    role: userRoleEnum('role').notNull().default('member'),
    emailVerifiedAt: timestamp('email_verified_at', { withTimezone: true }),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    updatedAt: timestamp('updated_at', { withTimezone: true }).notNull().defaultNow(),
  },
  (table) => ({
    emailIdx: uniqueIndex('idx_users_email').on(table.email),
  }),
);

// ---------------------------------------------------------------------------
// sessions — server-side, revocable. See section 3 for the auth flow this backs.
// ---------------------------------------------------------------------------

export const sessions = consoleSchema.table(
  'sessions',
  {
    id: text('id').primaryKey(), // random session token, NOT a JWT — see section 3
    userId: uuid('user_id')
      .notNull()
      .references(() => users.id),
    expiresAt: timestamp('expires_at', { withTimezone: true }).notNull(),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    // Set to a non-null revocation reason to kill a session instantly
    // (e.g. superadmin-forced logout) without waiting for expiry.
    revokedAt: timestamp('revoked_at', { withTimezone: true }),
  },
  (table) => ({
    userIdx: index('idx_sessions_user').on(table.userId),
  }),
);

// ---------------------------------------------------------------------------
// email_verification_tokens — short-lived, single-use.
// ---------------------------------------------------------------------------

export const emailVerificationTokens = consoleSchema.table('email_verification_tokens', {
  id: text('id').primaryKey(), // random token, sent in the Unsend verification link
  userId: uuid('user_id')
    .notNull()
    .references(() => users.id),
  expiresAt: timestamp('expires_at', { withTimezone: true }).notNull(),
  usedAt: timestamp('used_at', { withTimezone: true }),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});

// ---------------------------------------------------------------------------
// billing_records — synced periodically from Kobo Core's usage data.
// See CONSOLE_ARCHITECTURE.md "Console <-> Core communication" for the sync job.
// ---------------------------------------------------------------------------

export const billingRecords = consoleSchema.table(
  'billing_records',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    integratorId: uuid('integrator_id')
      .notNull()
      .references(() => apiIntegrators.id),
    environment: environmentEnum('environment').notNull(),
    period: text('period').notNull(), // "YYYY-MM"
    accountsProvisioned: bigint('accounts_provisioned', { mode: 'number' }).notNull().default(0),
    transactionsProcessed: bigint('transactions_processed', { mode: 'number' })
      .notNull()
      .default(0),
    amountDueKobo: bigint('amount_due_kobo', { mode: 'number' }).notNull().default(0),
    // Manual adjustments (credits/corrections) applied by a superadmin.
    // Positive = credit to integrator, negative = additional charge.
    adjustmentKobo: bigint('adjustment_kobo', { mode: 'number' }).notNull().default(0),
    adjustmentReason: text('adjustment_reason'),
    syncedAt: timestamp('synced_at', { withTimezone: true }).notNull().defaultNow(),
  },
  (table) => ({
    integratorPeriodIdx: uniqueIndex('idx_billing_integrator_period_env').on(
      table.integratorId,
      table.period,
      table.environment,
    ),
  }),
);

// ---------------------------------------------------------------------------
// admin_audit_log — append-only. Every superadmin action writes here.
// ---------------------------------------------------------------------------

export const adminActionEnum = consoleSchema.enum('admin_action', [
  'integrator_suspended',
  'integrator_reinstated',
  'production_access_granted',
  'credential_force_revoked',
  'billing_adjustment_applied',
  'user_session_revoked',
]);

export const adminAuditLog = consoleSchema.table('admin_audit_log', {
  id: uuid('id').primaryKey().defaultRandom(),
  actorUserId: uuid('actor_user_id')
    .notNull()
    .references(() => users.id), // must have role = 'superadmin' at time of action
  action: adminActionEnum('action').notNull(),
  targetIntegratorId: uuid('target_integrator_id').references(() => apiIntegrators.id),
  targetUserId: uuid('target_user_id').references(() => users.id),
  targetCredentialId: uuid('target_credential_id').references(() => apiCredentials.id),
  detail: jsonb('detail').notNull().default({}), // free-form context, e.g. { reason, previous_value, new_value }
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});
```

### Notes for the agent implementing this

- `bigint` fields use Drizzle's `{ mode: 'number' }` for kobo amounts,
  matching Core's convention of integer kobo values (see
  `core/openapi.yaml`'s `amount_kobo` fields) — do not introduce floats for
  money anywhere in the console either.
- `api_credentials.secretHash` — hash the raw secret with SHA-256 before
  storing (a fast hash is fine here, unlike password hashing, since API
  secrets are high-entropy random values, not user-chosen passwords —
  argon2 is reserved for `users.passwordHash`).
- Foreign keys from `sessions`, `email_verification_tokens`,
  `admin_audit_log` etc. into `users` are within the same `console` schema,
  so normal Drizzle relations work — no cross-schema FK complications there,
  since only `api_integrators` crosses the app boundary, and it lives fully
  inside `console`.

---

## 3. Session-based auth (no JWTs)

Confirmed decision from `CONSOLE_ARCHITECTURE.md`: sessions must be
instantly revocable (superadmin lockout is a stated requirement), so this
uses opaque server-side session tokens, not JWTs.

### Token design

- Session ID: a cryptographically random 32-byte value, base64url-encoded,
  stored as the `sessions.id` primary key (plain text, not hashed — unlike
  API credentials, session tokens are single-purpose, short-lived, and
  transmitted only via httpOnly cookies over HTTPS, so the threat model is
  different enough that hashing adds complexity without much benefit here;
  revisit if this console ever needs to log session tokens anywhere, in
  which case hash them the same way as API credentials).
- Cookie: httpOnly, `Secure`, `SameSite=Lax`, path `/`, expiry matching
  `sessions.expiresAt` (recommend 7 days, sliding — refresh on activity).

### Implementation shape

```typescript
// src/lib/server/auth/session.ts
import { randomBytes } from 'node:crypto';
import { db } from '$lib/server/db';
import { sessions, users } from '$lib/server/db/schema';
import { eq, and, isNull, gt } from 'drizzle-orm';

const SESSION_DURATION_MS = 1000 * 60 * 60 * 24 * 7; // 7 days

export function generateSessionId(): string {
  return randomBytes(32).toString('base64url');
}

export async function createSession(userId: string) {
  const id = generateSessionId();
  const expiresAt = new Date(Date.now() + SESSION_DURATION_MS);
  await db.insert(sessions).values({ id, userId, expiresAt });
  return { id, expiresAt };
}

export async function validateSession(sessionId: string) {
  const [session] = await db
    .select()
    .from(sessions)
    .where(
      and(
        eq(sessions.id, sessionId),
        isNull(sessions.revokedAt),
        gt(sessions.expiresAt, new Date()),
      ),
    )
    .innerJoin(users, eq(users.id, sessions.userId))
    .limit(1);

  return session ?? null; // null means: invalid, expired, or revoked
}

export async function revokeSession(sessionId: string) {
  await db.update(sessions).set({ revokedAt: new Date() }).where(eq(sessions.id, sessionId));
}

// Used by the superadmin "force logout" admin action — revokes every active
// session for a user, not just one.
export async function revokeAllSessionsForUser(userId: string) {
  await db
    .update(sessions)
    .set({ revokedAt: new Date() })
    .where(and(eq(sessions.userId, userId), isNull(sessions.revokedAt)));
}
```

### Central auth gate in `hooks.server.ts`

This is the single place both "is logged in" and "is email verified" and
"is superadmin" get enforced, per the hard rule in
`CONSOLE_ARCHITECTURE.md` ("Email verification is enforced centrally, not
per-route"). Route-level code should be able to assume `event.locals.user`
is already validated by the time it runs.

```typescript
// src/hooks.server.ts
import type { Handle } from '@sveltejs/kit';
import { validateSession } from '$lib/server/auth/session';

const PUBLIC_ROUTES = ['/login', '/signup', '/verify-email'];
const SUPERADMIN_PREFIX = '/admin';

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get('session');
  const result = sessionId ? await validateSession(sessionId) : null;

  event.locals.user = result?.users ?? null; // shape depends on your join result
  event.locals.session = result?.sessions ?? null;

  const path = event.url.pathname;
  const isPublic = PUBLIC_ROUTES.some((p) => path.startsWith(p));

  if (!isPublic && !event.locals.user) {
    return Response.redirect(new URL('/login', event.url), 302);
  }

  if (
    event.locals.user &&
    !event.locals.user.emailVerifiedAt &&
    path !== '/verify-email' &&
    !isPublic
  ) {
    return Response.redirect(new URL('/verify-email', event.url), 302);
  }

  if (path.startsWith(SUPERADMIN_PREFIX) && event.locals.user?.role !== 'superadmin') {
    return new Response('Not found', { status: 404 }); // 404, not 403 — don't reveal the admin area exists
  }

  return resolve(event);
};
```

`event.locals.user` and `event.locals.session` types should be declared in
`src/app.d.ts` so every server route gets typed access without re-checking.

### Signup → verification → session flow

```
POST /signup { email, password }
  → hash password (argon2)
  → INSERT users (emailVerifiedAt = null)
  → generate verification token, INSERT email_verification_tokens
  → send via Unsend (verification link: /verify-email?token=...)
  → create session immediately (user can browse, but hooks.server.ts
    redirects everything except /verify-email until verified)

GET /verify-email?token=...
  → look up token, check not expired / not used
  → UPDATE users SET emailVerifiedAt = now()
  → UPDATE email_verification_tokens SET usedAt = now()
  → redirect to /dashboard
```

### Superadmin seeding

Not a signup flow — a one-time script:

```typescript
// scripts/seed-superadmin.ts
// Run manually: `pnpm tsx scripts/seed-superadmin.ts`
// Never expose this as an HTTP route.
```

---

## 4. Core's read-only side of the shared table

For symmetry, Core needs a read-only sqlc query against `console.api_integrators`,
using the `kobo_core_app` role's SELECT-only grant from section 1. This
should be added to `core/internal/platform/db/queries/integrators.sql`
(new file) as a natural follow-on to this work:

```sql
-- name: GetIntegratorByAPIKeyID :one
-- Core's auth middleware calls this to resolve an inbound API key to an
-- integrator. Joins across the shared api_integrators table and the
-- Console-owned api_credentials table — both readable by kobo_core_app per
-- the grants in CONSOLE_TECHNICAL_GROUNDING.md section 1.
SELECT
    ai.id, ai.name, ai.status, ai.production_access_granted,
    ac.id AS credential_id, ac.environment, ac.secret_hash, ac.revoked_at
FROM console.api_integrators ai
JOIN console.api_credentials ac ON ac.integrator_id = ai.id
WHERE ac.key_id = $1;
```

Core's middleware then compares the inbound request's secret (hashed the
same way, SHA-256) against `secret_hash`, checks `revoked_at IS NULL`, and
checks `ai.status = 'active'` before accepting the request. If
`environment = 'production'` but `ai.production_access_granted = false`,
reject even a technically-valid key — this is the enforcement point for the
"production access is a manual gate" rule.

---

## Summary of what's now locked

1. One Postgres database, two schemas (`core`, `console`), enforced by
   Postgres role grants, not just convention.
2. `api_integrators` is the only cross-schema table; it lives in `console`,
   Core has SELECT-only access via a dedicated role.
3. Drizzle schema for all Console-owned tables is written above and ready
   to drop into `src/lib/server/db/schema.ts`.
4. Sessions are opaque, server-side, instantly revocable — no JWTs.
5. `hooks.server.ts` is the single enforcement point for auth, email
   verification, and superadmin gating.
6. Core needs one new sqlc query file to read the shared table; this is a
   small, well-scoped addition to the existing Go codebase, not a
   restructuring of it.
