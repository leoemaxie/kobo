# Reconciliation Engine

This document is the canonical reference for how Kobo matches Nomba inbound
transfer events to identity ledger entries. It is kept in sync with
`core/internal/reconciliation/engine.go`, `sweep.go`, and
`idempotency.go`. If code and document disagree, treat this as the spec.

Reconciliation accuracy is one of the four named judging criteria for the
Nomba Hackathon Infrastructure Track. Every edge case below has a defined,
tested behavior. There is no "we'll handle it manually" path.

---

## The core matching insight

Because each identity has its own dedicated Nomba virtual account number,
matching an inbound transfer to an identity is trivial by design:

```
incoming transfer → account number → virtual_accounts table → identity
```

There is no fuzzy matching, no name matching, no amount-based correlation.
Account number is a unique key in `virtual_accounts`. If it resolves, it is
a match. If it does not resolve, it is an exception.

The hard part is not the match itself. The hard part is making that match
**reliable** under real-world delivery conditions: duplicate webhook
deliveries, delayed or missing webhooks, out-of-order arrival, accounts that
have since moved to CLOSING or CLOSED, and amounts arriving in Naira that
need converting to kobo before any ledger write.

Everything below addresses that.

---

## Data sources

Kobo's reconciliation engine has two sources of truth for inbound transfers,
used in a primary/fallback relationship. Both sources eventually describe the
same underlying bank events. Kobo never double-counts because idempotency is
enforced at the `nomba_reference` level before any ledger write regardless
of which source produced the event.

### Primary: Nomba webhooks

Nomba calls `POST /webhooks/nomba` on Kobo's server for each `payment_success`
event. This is the fastest path from a real-world bank transfer to a Kobo
ledger entry — typically seconds after the transfer clears.

The webhook handler (`core/internal/api/handlers/webhooks.go`) must:
1. Verify the Nomba signature before touching any business logic. See
   `docs/NOMBA_INTEGRATION.md` for the exact HMAC construction — it is a
   specific concatenated-field hash, not a generic body hash.
2. Immediately return `200 OK` once the event is durably written to the
   processing queue or directly processed. Do not return non-2XX for
   business-logic failures (account not found, duplicate, etc.) — return 200
   and handle those cases in code. Returning non-2XX causes Nomba to retry
   up to 5 times over ~53 minutes, which makes legitimate failures and late
   retries indistinguishable.
3. Check idempotency before any ledger write.
4. Route to the reconciliation engine.

### Fallback: Nomba Transactions API polling (sweep)

`core/internal/reconciliation/sweep.go` runs as a background job in
`cmd/worker`, polling `GET /transactions/virtual?virtual_account=<number>`
for accounts that have no matching webhook event within a configurable
window.

The sweep window should be at least 60 minutes — slightly wider than Nomba's
~53-minute maximum webhook retry interval — to avoid the sweep backfilling
transfers that a late webhook delivery is still going to cover. Setting it
too narrow causes harmless double-processing attempts (caught by idempotency);
setting it too wide means a genuinely missing webhook goes unnoticed longer
than it should.

The sweep does not replace webhooks. It is a safety net for events that
Nomba's webhook delivery dropped entirely, which the retry mechanism cannot
recover from because Nomba's retries only fire on Kobo's non-2XX responses,
not on Nomba's own internal delivery failures.

---

## Full reconciliation flow

```
Nomba event arrives (webhook) OR sweep job runs
             │
             ▼
    ┌─────────────────┐
    │ Signature verify │  (webhooks only; sweep uses OAuth token, already verified)
    └────────┬────────┘
             │ fail → 401, stop
             │ pass ↓
    ┌─────────────────────────────────────┐
    │ Parse nomba_reference + account_number│
    │ Convert amount: Naira × 100 → kobo  │
    └────────┬────────────────────────────┘
             │
             ▼
    ┌───────────────────────┐
    │  Idempotency check    │  SELECT FROM idempotency_keys WHERE nomba_reference = ?
    └────────┬──────────────┘
             │ found → return existing ledger entry, stop (deduplicated)
             │ not found ↓
    ┌─────────────────────────────────────┐
    │ Resolve account number → identity   │  SELECT FROM virtual_accounts WHERE account_number = ?
    └────────┬────────────────────────────┘
             │ not found → exception: payment_to_unknown_account, stop
             │ found ↓
    ┌──────────────────────────────┐
    │  Check AcceptsInboundFunds() │  per lifecycle.go
    └────────┬─────────────────────┘
             │ state = ACTIVE or LIMITED  → proceed to ledger write
             │ state = CLOSING            → exception: payment_during_closing, stop
             │ state = CLOSED or FAILED   → exception: payment_to_closed_account, stop
             │ state = PENDING            → exception: payment_to_unknown_account, stop
             │                              (account number shouldn't exist yet)
             ↓
    ┌─────────────────────────────────────────────────────────────┐
    │  Within a single DB transaction:                            │
    │    1. INSERT INTO ledger_entries (...)                      │
    │    2. INSERT INTO idempotency_keys (nomba_reference, ...)   │
    │  If idempotency INSERT fails on unique violation:           │
    │    → another concurrent worker beat us to it, treat as      │
    │      duplicate, roll back, return existing entry            │
    └─────────────────────────────────────────────────────────────┘
             │
             ▼
    Emit transaction.received webhook to integrator
    (async, outside the DB transaction)
```

---

## Edge cases

Every case below has a defined outcome. None of them are "log and investigate
later." The reconciliation engine either writes a ledger entry or creates an
exception row — there is no silent discard path.

---

### 1. Duplicate webhook delivery

**What happens:** Nomba retries a webhook after a network blip, sending the
same `payment_success` event twice with the same `transactionId`.

**Kobo's behaviour:** The second delivery hits the idempotency check, finds
the `nomba_reference` already in `idempotency_keys`, and returns the original
`ledger_entry_id` without writing anything. The webhook handler returns `200`
to Nomba. No duplicate ledger entry is created.

**Implementation:** The `idempotency_keys.nomba_reference` column is a
`PRIMARY KEY` (unique constraint). The insert in step 2 of the DB transaction
fails with a unique violation if the reference was already processed. The
engine must catch this specific error and treat it as a success, not as an
application error.

**Idempotency key:** Use Nomba's `transactionId` field from the webhook
payload (`data.transaction.transactionId`). See the open item in
`docs/NOMBA_INTEGRATION.md` about confirming this is the same identifier
space as the sweep endpoint's `id` field — verify in sandbox before
finalising, since the dedup key must be the same across both sources for the
idempotency table to work for both.

---

### 2. Webhook arrives before sweep; sweep arrives later for the same transfer

**What happens:** A webhook is delivered and processed (ledger entry written,
idempotency key inserted). The sweep job runs later and picks up the same
transfer from the Transactions API.

**Kobo's behaviour:** Sweep hits idempotency check, finds the reference, stops.
No duplicate. No exception. The `first_seen_via` column on `idempotency_keys`
records `webhook` for the original, so there is an audit trail of which source
was first.

---

### 3. Webhook delayed or missing; sweep fires first

**What happens:** Nomba's internal delivery fails silently (different from a
Kobo non-2XX response, which would trigger retries). The sweep job runs after
the configured window and finds the transfer in the Transactions API.

**Kobo's behaviour:** Sweep processes normally — idempotency check misses
(nothing there yet), account number resolves, ledger entry written,
idempotency key inserted with `first_seen_via = sweep`. If a late webhook
subsequently arrives (within Nomba's 5-retry window, before Nomba gives up),
it hits the idempotency check and is deduplicated. `first_seen_via = sweep`
remains on record.

---

### 4. Webhook and sweep arrive concurrently for the same transfer

**What happens:** The sweep window is configured narrow, or the job runs
while a webhook is mid-processing — two workers are racing to write the same
`nomba_reference`.

**Kobo's behaviour:** Only one succeeds. Both workers try to INSERT INTO
`idempotency_keys` within a DB transaction. One wins; the other gets a unique
violation, rolls back its transaction, and returns the existing entry. No
duplicate. The database constraint is the lock, not application-level
coordination.

**Implementation note:** This must actually be a transaction. If the ledger
insert and idempotency insert are two separate DB calls without a wrapping
transaction, a concurrent worker can see a window between the two inserts and
write a duplicate ledger entry before the idempotency key is committed.

---

### 5. Out-of-order delivery (Transactions API shows event before webhook fires)

**What happens:** The sweep job runs and picks up a transaction whose webhook
has not yet been delivered (or has been delivered but not processed). There
is no guaranteed ordering between the two sources.

**Kobo's behaviour:** Whichever source arrives first is processed. The second
is deduplicated via idempotency. Out-of-order processing is safe because
the idempotency check is on the stable `nomba_reference`, not on timestamps
or sequence numbers.

---

### 6. Partial payment (amount below any expected fee or installment)

**What happens:** A parent paying ₦150,000 school fees sends ₦30,000 as the
first installment.

**Kobo's behaviour:** The full ₦30,000 is recorded in `ledger_entries` with
`status = partial`. Kobo does not enforce or have any knowledge of what the
"expected" amount is — that is the integrator's (school-fees app's) concern.
Kobo's ledger entry has an amount and a timestamp; whether that amount
satisfies a fee structure is for the integrator to interpret by reading the
statement.

**Why Kobo doesn't enforce amounts:** The reconciliation primitive should
not have opinions about fee structures. If Kobo validated amounts against
expected values, every integrator vertical (school fees, co-op thrift,
vendor marketplace) would need a different validation rule, which defeats
the purpose of a reusable primitive. The integrator knows what the money
means; Kobo knows that the money arrived.

---

### 7. Overpayment (amount exceeds any expected balance)

**What happens:** A parent sends ₦160,000 when ₦150,000 is owed.

**Kobo's behaviour:** Full ₦160,000 is recorded with `status = overpayment`.
Kobo does not cap, reject, or return the excess. The ledger entry reflects
what actually arrived. The integrator's statement shows the surplus, and the
integrator decides what to do with it (refund, credit forward, flag for
review).

---

### 8. Payment to an account in CLOSING state

**What happens:** A transfer arrives for an identity whose closure has been
initiated but whose funds have not yet been swept (state = CLOSING).

**Kobo's behaviour:** Exception row created with `type = payment_during_closing`.
The transfer is not attributed to the ledger. The integrator is notified via
the `exception.flagged` webhook and can resolve it via
`POST /v1/exceptions/{id}/resolve` with one of:
- `return_to_sender` — funds sent back to the originating account.
- `redirect_to_successor` — if the integrator wants to re-route to another
  identity (e.g. a student transferred to a different class before their
  previous account closed).

**Why not attribute to the ledger:** Attributing to a CLOSING account would
defeat the point of the sweep-wait-for-zero balance logic. If the engine
silently credited the ledger during CLOSING, the balance would never reach
zero and the account would never actually close.

---

### 9. Payment to a CLOSED account

**What happens:** A stale bank transfer, a scheduled payment not cancelled,
or a sender who has the old account number on file, sends money to a virtual
account that is CLOSED.

**Kobo's behaviour:** Exception row created with `type = payment_to_closed_account`.
Integrator notified. Resolutions available:
- `return_to_sender`
- `redirect_to_successor` (if the identity was reopened under a new account
  number or succeeded by another identity)
- `manual_override` (integrator reviews and handles outside Kobo)

Every misdirected case is logged with a resolution outcome, so the audit
trail can always answer "where did this money end up" for any transfer.

---

### 10. Payment to an unknown account number

**What happens:** Nomba sends a webhook for a transfer to an account number
that does not exist in Kobo's `virtual_accounts` table. This should not
happen under normal operation — Kobo provisions every virtual account through
Nomba — but could occur if Nomba's sandbox sends test events with stale
account numbers, or if an account number is somehow issued outside of Kobo's
provisioning flow.

**Kobo's behaviour:** Exception row created with `type = payment_to_unknown_account`,
`related_account_id = NULL`. Integrator notified. The exception queue UI
(or `GET /v1/exceptions`) shows the raw account number and amount so the
integrator can investigate and apply a `manual_override` resolution.

---

### 11. Amount unit conversion

**What happens:** Nomba's webhook payload sends `transactionAmount` as a
number in **Naira** (e.g. `120`). Kobo's internal model uses **kobo**
(e.g. `12000`).

**Kobo's behaviour:** The reconciliation engine multiplies by 100 and rounds
before any ledger write. This conversion happens in one place only:
`core/internal/reconciliation/engine.go`, not in the webhook handler or
anywhere else that touches the raw Nomba payload. The raw Naira amount is
never stored anywhere in Kobo's schema.

**Why this matters:** A conversion error here would cause every balance and
statement figure to be off by a factor of 100 — a silent, pervasive bug
with no obvious single point of failure. Centralising it means there is
exactly one line to audit.

---

### 12. Webhook arrives after account is FAILED

**What happens:** Nomba somehow delivers a credit event for an account number
that was provisioned (and thus has a real account number in Nomba) but whose
identity is in FAILED state on Kobo's side.

**Kobo's behaviour:** The account number resolves in `virtual_accounts`
(provisioning must have partially succeeded for a number to exist).
`AcceptsInboundFunds(StateFailed)` returns false, so the transfer is flagged
as `payment_to_closed_account` (the semantically closest type; "failed
provisioning" is not a distinct exception type in v1 since this case should
not occur in practice). Integrator notified.

---

## Naira-to-kobo conversion reference

| Nomba field          | Type in Nomba payload | Kobo ledger field | Conversion            |
|----------------------|-----------------------|-------------------|-----------------------|
| `transactionAmount`  | `number` (Naira)      | `amount_kobo`     | `int64(amount * 100)` |
| amount in sweep API  | `string` (e.g. "100.0") | `amount_kobo`   | parse as float64, × 100, round to int64 |

Use `math.Round` before casting to `int64` when the Transactions API returns
a string amount, to avoid float rounding errors at the kobo level.

---

## Reconciliation sweep configuration

The following values live in `platform/config/config.go` and should be tunable
via environment variables, not hardcoded:

| Config key                  | Recommended default | Reason |
|-----------------------------|---------------------|--------|
| `SWEEP_INTERVAL`            | `15m`               | How often the sweep job runs. |
| `SWEEP_LOOKBACK_WINDOW`     | `90m`               | How far back to poll per account. Wider than Nomba's 53-min retry ceiling. |
| `SWEEP_BATCH_SIZE`          | `50`                | Max accounts to check per sweep run (see `ListVirtualAccountsNeedingSweepCheck` in `virtual_accounts.sql`). |
| `WEBHOOK_REPLAY_REJECT_AGE` | `5m`                | Reject webhook events older than this to prevent replay attacks. Applied in signature verification. |

---

## What the reconciliation engine does NOT do

Explicitly out of scope, to prevent scope creep in agent-driven development:

- **Outbound transfers:** Kobo v1 reconciles inbound attribution only. Refund
  or payout flows (e.g. returning an overpayment to a parent) are out of
  scope and should be handled by the integrator directly via Nomba's payout
  APIs.
- **Amount-based matching:** Kobo does not know what a "correct" amount looks
  like for any given identity. The integrator defines fee structures.
- **Automatic refunds:** Exceptions are flagged and held; resolutions are
  applied by the integrator via the API. No money moves automatically on
  exception detection.
- **Cross-integrator reconciliation:** Each integrator's accounts are in an
  isolated namespace. The engine never looks across integrator boundaries.
