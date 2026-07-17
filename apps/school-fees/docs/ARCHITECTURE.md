# Kobo Reference App — School Fees — Architecture

High-level architecture for `apps/school-fees`, the reference integrator app
that demonstrates Kobo core end-to-end.

---

## What this app is for, and what it deliberately is not

**It is:** a small, real, working demonstration that an outside team can
integrate against the Kobo API using nothing but a sandbox API key and the
public `openapi.yaml` contract, and get a working fee-collection product
without touching Kobo's internals, database, or Monnify credentials directly.

---

## Tech stack (locked)

- **Framework:** SvelteKit, fullstack mode — same pattern as the Console
  (SSR pages + server routes), for tooling and mental-model consistency
  across the monorepo, not because this app needs SvelteKit's full power.
- **Auth:** Basic session-based auth, same shape as Console's (opaque
  server-side session tokens, httpOnly cookies) — but a **separate,
  smaller** implementation, not a shared library, since this app's auth
  model is simpler (parents only, no roles, no superadmin) and it must stay
  architecturally isolated from Console per the "separate integrator, zero
  special access" story.
- **Database:** Its own PostgreSQL database (`school_fees`), entirely
  separate from both `core` and `console` schemas. This app has no special
  access to Kobo's database, no shared tables, nothing. It only ever talks
  to Kobo through the public API, exactly like any other integrator would
  have to.
- **Kobo API access:** Server-side only. All calls to Kobo core
  (`sandbox.api.kobo.triumphsystems.tech` initially) happen from SvelteKit server routes
  (`+page.server.ts` / `+server.ts`), using a sandbox API key stored as a
  server-only environment variable. The browser never sees the Kobo API key
  or talks to Kobo directly — every request goes: browser → this app's
  SvelteKit server → Kobo API → back. This is both the more realistic
  integrator pattern and the safer one, since it means Kobo credentials are
  never exposed in client-side JS or devtools network tabs.
- **Styling:** Tailwind CSS, minimal — this app's UI should look clean but
  is explicitly not where design effort goes. It can borrow the Kobo brand
  accent color once that's finalized, but doesn't need its own design
  system.
- **Deployment:** Same target class as Console (Vercel/Fly.io), separate
  deployment, separate subdomain — e.g. `fees.kobo.triumphsystems.tech`
  or similar, clearly labeled as a demo/reference app if judges land on it
  directly.

---

## High-level architecture diagram

```
                    Parent's browser
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │   School Fees App (SvelteKit)          │
        │   apps/school-fees                      │
        │                                          │
        │  Routes                                  │
        │   /login  /signup                        │
        │   /dashboard  (parent's student list)     │
        │   /students/[id]  (statement + history)    │
        │   /admin/students  (register/close)         │
        │                                          │
        │  Server-side only                         │
        │   - session auth (own, separate from Console)│
        │   - Kobo API client (holds KOBO_API_KEY,    │
        │     KOBO_API_SECRET as server env vars)      │
        └───────────┬─────────────────┬───────────────┘
                    │                 │
        ┌───────────▼───────┐   ┌─────▼─────────────────────┐
        │  school_fees DB     │   │   Kobo API (sandbox)         │
        │  (Postgres, own      │   │   sandbox.api.kobo.dev/v1     │
        │  schema, Drizzle)     │   │                                │
        │                       │   │   POST /v1/identities           │
        │  parents               │   │   GET /v1/identities/{id}        │
        │  parent_students        │   │   GET /v1/accounts/{id}/statement │
        │  (mapping: which        │   │   GET /v1/accounts/{id}/transactions│
        │   parent can see which  │   │   POST /v1/identities/{id}/close     │
        │   Kobo identity_id)      │   └────────────────────────────────────┘
        └─────────────────────┘
```

---

## Core data model (high level)

This app's database is deliberately thin. It does not duplicate anything
Kobo already tracks (balances, transactions, lifecycle state) — it only
stores what Kobo has no reason to know: which parent is allowed to see
which student, and basic school-side identifiers.

- **`parents`** — email, password hash, session-backed login. No roles
  beyond "parent" and a minimal "school admin" flag for the one admin
  screen (registering/closing students) — not a full RBAC system, this is
  a demo.
- **`students`** — school-side record: name, class, and critically the
  `kobo_identity_id` returned when this app called `POST /v1/identities`.
  This app never stores balance or transaction data locally — every time a
  parent views a student's account, this app calls Kobo's statement/
  transactions endpoints live and renders the response. This is deliberate:
  it's the clearest possible demonstration that Kobo, not this app, is the
  system of record for financial data.
- **`parent_students`** — join table, which parent can view which student
  (a parent may have more than one child at the school).
- **`sessions`** — same shape as Console's sessions table, but this app's
  own copy, not shared.

No `payments`, `balances`, or `transactions` tables in this app's own
database — that would defeat the point of the demo.

---

## The three real screens

### 1. Register a student (admin-facing, minimal)

A simple form: student name, class. On submit, the server route calls
`POST /v1/identities` with `external_reference` set to this app's own
generated student ID and `display_name` set to the student's name, per the
`CreateIdentityRequest` shape in `core/openapi.yaml`. Store the returned
`kobo_identity_id` on the local `students` row. That's the entire
provisioning flow from this app's side — everything else (virtual account
creation, Monnify call, lifecycle state) is Kobo's problem, not this app's,
which is exactly the point being demonstrated.

### 2. Parent dashboard → student account view

Parent logs in, sees their linked student(s) (via `parent_students`),
clicks into one. The server route calls
`GET /v1/accounts/{accountId}/statement` and
`GET /v1/accounts/{accountId}/transactions` for that student's Kobo
identity, and renders the balance, payment history, and status. No local
caching of financial data — always live from Kobo.

### 3. Admin: close a student's account (e.g. graduation)

A button that calls `POST /v1/identities/{id}/close` with a
`sweep_destination`. For a hackathon demo, the simplest defensible choice is
`sweep_destination: { type: "refund_to_source" }` — return any remaining
balance to whoever last paid in, since a real "integrator bank account"
destination requires banking details this demo doesn't need to model.

---

## Why server-side-only API access matters here specifically

This decision (confirmed: exposed via SvelteKit server, not the browser)
is the single most important architectural choice in this app, and it's
worth stating why explicitly rather than just as a rule:

1. **It's the realistic pattern.** No real integrator would ever ship a
   secret API key to a browser bundle. Demonstrating this correctly is
   itself part of what makes Kobo credible as infrastructure — judges
   evaluating "developer API quality" will notice if the reference
   implementation gets this basic thing wrong.
2. **It matches how Kobo's own auth is designed.** Per `openapi.yaml`,
   Kobo requests are authenticated with an API key plus an HMAC signature
   over the request body and a timestamp (see `core/docs/MONNIFY_INTEGRATION.md`
   for the pattern this mirrors on the Monnify side, and `openapi.yaml`'s
   `HmacSignature` security scheme for Kobo's own equivalent). HMAC signing
   requires the API secret to compute the signature — that secret cannot
   exist in browser-delivered JavaScript under any circumstance.
3. **It keeps the trust boundary honest.** This app's server is the
   integrator's trusted backend. The parent's browser is untrusted, exactly
   as it would be for a real school's parents. Kobo only ever talks to a
   server that holds real credentials.

Concretely: build one small `kobo-client.ts` module under
`apps/school-fees/src/lib/server/` (note: `server/` — SvelteKit's convention
for code that's guaranteed never to ship to the client bundle) that wraps
fetch calls to Kobo's API, attaches the API key header and HMAC signature,
and is imported only from `+page.server.ts` / `+server.ts` files, never from
`.svelte` component code or anything under a non-`server/` lib path.

---

## Relationship to Console

This app and the Console are both SvelteKit, both use sessions, but they
are otherwise unrelated:

- This app's `parents` table has nothing to do with Console's `users`
  table. A parent logging into the school-fees demo is not a Kobo
  integrator account holder in Console's sense.
- This app's Kobo API key is one that would have been generated *by*
  someone using Console's `/api-keys` screen (in a real flow: "Team
  Triumph" as the integrator, generating a sandbox key for their
  school-fees demo project) and then pasted into this app's environment
  variables. For the hackathon, this can just be seeded/hardcoded in a
  `.env` file rather than requiring a live Console signup flow, since
  Console and this app are being built in parallel and don't need to block
  each other.
- Do not let this app reach into Console's or Core's database under any
  circumstances, even read-only. The entire value of this app as a
  demonstration depends on it having zero special access — it must be
  provably indistinguishable, from Kobo's point of view, from any other
  integrator calling the public API.

---

## What this document deliberately leaves open

- Exact Drizzle schema field names (small enough that an agent can design
  this directly from the "Core data model" section above).
- Specific route file layout under `src/routes/`.
- Visual design — this app should look clean but is not a design
  priority; borrow Kobo's accent color once finalized and move on.
- Whether admin (student registration/closure) gets its own login role or
  just a hardcoded `is_admin` flag on a parent row — either is fine for a
  demo, pick whichever is faster to build.

---

## Summary of hard rules for any agent building this

1. Separate app, separate database, separate deployment from both Core and
   Console. No shared tables, no shared code beyond the monorepo's general
   conventions.
2. Kobo API calls happen exclusively from SvelteKit server-side code
   (`+page.server.ts`, `+server.ts`, or a `lib/server/` module). The API
   key and secret must never reach the browser bundle.
3. This app stores no financial data locally — balances and transactions
   are always fetched live from Kobo's API, never cached or duplicated in
   `school_fees`'s own database. The only local data is the mapping from
   parent → student → `kobo_identity_id`.
4. Scope stays to three screens: register, view, close. Every feature
   beyond that is explicitly out of scope per `core/docs/COVERAGE.md`.
