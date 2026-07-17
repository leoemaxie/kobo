# Account Lifecycle

This document is the canonical reference for Kobo's virtual account lifecycle
state machine. It is kept in sync with `core/internal/account/lifecycle.go`
and the `IdentityState` enum in `openapi.yaml`. If those three sources ever
disagree, treat this document as the spec and the code as the bug.

---

## State diagram

```
                        ┌─────────────────────────────────────────┐
                        │                                         │
              [provisioning_failed]                        [reopen_requested]
                        │                                         │
                        ▼                                         │
PENDING ──[provisioned]──► ACTIVE ◄──[kyc_tier_upgraded]── LIMITED
                             │                                    │
                     [kyc_tier_breached]              [close_requested]
                             │                                    │
                             ▼                                    │
                           LIMITED                                │
                             │                                    │
                     [close_requested]                            │
                             │                                    │
                             └──────────────┬─────────────────────┘
                                            │
                                            ▼
                                         CLOSING ──[funds_swept]──► CLOSED
                                                                       │
                                                               [reopen_requested]
                                                                       │
                                                                       ▼
                                                                    ACTIVE
                        ▲
                        │
PENDING ──[provisioning_failed]──► FAILED (terminal)
```

Plain-English version for anyone reading this without the ASCII art rendering:

- Every identity enters at **PENDING** when Kobo submits a provisioning
  request to Monnify.
- On Monnify confirmation, it becomes **ACTIVE** and the virtual account number
  is usable.
- If Monnify rejects provisioning for any reason, it moves to **FAILED** and
  stays there. No path out of FAILED exists in v1.
- An ACTIVE account becomes **LIMITED** when cumulative inflow approaches the
  KYC tier ceiling. It still accepts funds and reconciles normally while
  LIMITED, but the integrator is notified that a KYC upgrade is needed to
  avoid disruption. Verify the exact Monnify ceiling values and signaling
  before finalising the `kyc_tier_breached` trigger logic — this is flagged
  as an open item in `docs/MONNIFY_INTEGRATION.md`.
- Either ACTIVE or LIMITED can be moved to **CLOSING** by the integrator
  calling `POST /v1/identities/{id}/close`. Closure is a two-step process:
  CLOSING freezes new attribution while in-flight settlements clear; **CLOSED**
  is only reached once the balance is zero and funds are swept to the
  integrator-specified destination.
- A **CLOSED** account can be **ACTIVE** again via
  `POST /v1/identities/{id}/reopen`. This reuses the same identity ID and
  does not require a new provisioning call, unless Monnify requires a fresh
  `accountRef` on reopen (see open item #2 in `docs/MONNIFY_INTEGRATION.md`).

---

## Full transition table

Every valid (current state, event) → next state pair. Any combination not
listed here is rejected by `lifecycle.Transition()` with `ErrInvalidTransition`,
which the API layer surfaces as 409 Conflict.

| Current state | Event                  | Next state | Who triggers it                                                                 |
|---------------|------------------------|------------|---------------------------------------------------------------------------------|
| PENDING       | provisioned            | ACTIVE     | `account/service.go` on successful Monnify provisioning response                 |
| PENDING       | provisioning_failed    | FAILED     | `account/service.go` on any non-`"00"` Monnify response or timeout               |
| ACTIVE        | kyc_tier_breached      | LIMITED    | `cmd/worker` background job, when cumulative inflow approaches tier ceiling     |
| ACTIVE        | close_requested        | CLOSING    | `handlers/identities.go` on `POST /v1/identities/{id}/close`                   |
| LIMITED       | kyc_tier_upgraded      | ACTIVE     | `cmd/worker` background job, after integrator completes KYC upgrade             |
| LIMITED       | close_requested        | CLOSING    | `handlers/identities.go` on `POST /v1/identities/{id}/close`                   |
| CLOSING       | funds_swept            | CLOSED     | `cmd/worker` background job, once balance = 0 and sweep transfer confirmed      |
| CLOSED        | reopen_requested       | ACTIVE     | `handlers/identities.go` on `POST /v1/identities/{id}/reopen`                  |

### Non-transitions (explicitly called out to prevent confusion)

These look like they might be state changes but are NOT modeled in the state
machine because they do not change state:

| Action | What actually happens |
|---|---|
| Rename (`PATCH /v1/identities/{id}` with new `display_name`) | `identity_events` row of type `renamed` written; Monnify `accountName` re-submitted. Identity stays in its current state. |
| Metadata update (`PATCH /v1/identities/{id}` with new `metadata`) | `identity_events` row of type `metadata_updated` written. No state change, no Monnify call. |
| Inbound transfer received (ACTIVE or LIMITED) | `ledger_entries` row written. No state change — receiving money does not change lifecycle state. |
| KYC tier checked, no breach | No event raised. Worker logs its check and moves on. |

---

## State-by-state reference

### PENDING

**Entered from:** identity creation (`POST /v1/identities`)

**Meaning:** Kobo has accepted the registration and submitted a provisioning
request to Monnify. The virtual account does not yet exist. The identity's
`virtual_account` field in API responses is `null` while in this state.

**What can happen:**
- Monnify confirms provisioning → ACTIVE
- Monnify rejects provisioning, or the Monnify call times out after retries → FAILED
- Inbound transfer arriving while PENDING: should not happen in practice since
  no account number has been issued yet, but if a Monnify-side race produces
  one, the reconciliation engine has no account number to match against, so
  it is routed to the exceptions queue as `payment_to_unknown_account`.

**Idempotency on provisioning:** Kobo sends a per-identity `X-Idempotent-key`
header on the Monnify provisioning call. If the call fails at the network layer
and Kobo retries, the same idempotency key is reused so Monnify does not create
a second account. See `docs/MONNIFY_INTEGRATION.md`.

---

### ACTIVE

**Entered from:** PENDING (provisioned), LIMITED (kyc_tier_upgraded),
CLOSED (reopen_requested)

**Meaning:** The virtual account is live and accepting inbound transfers. Normal
operating state for a healthy identity.

**What can happen:**
- Inbound transfer → reconciled to ledger normally
- Cumulative inflow approaches KYC tier ceiling → LIMITED
- Integrator requests closure → CLOSING

---

### LIMITED

**Entered from:** ACTIVE (kyc_tier_breached)

**Meaning:** The virtual account is still live and still reconciles inbound
transfers normally. The distinction from ACTIVE is purely administrative:
the integrator and Kobo's background worker know that the account is near
or at its KYC tier ceiling and that action is required to avoid disruption.
End users (parents paying school fees, for example) experience no difference.

**Key design decision:** Kobo moves an account to LIMITED *before* Monnify
would hard-reject an inbound transfer, not after. This converts a hard
payment failure at the bank-rail level into a soft warning at the application
level. The integrator has time to trigger a KYC upgrade before any transfer
is actually rejected.

**What can happen:**
- KYC upgrade completed → ACTIVE
- Integrator requests closure → CLOSING

---

### CLOSING

**Entered from:** ACTIVE or LIMITED (close_requested)

**Meaning:** The account is being wound down. New inbound transfers that
arrive while an account is CLOSING are classified as `payment_during_closing`
in the exceptions queue, held rather than attributed to the ledger, and
resolved once the closure completes.

**Balance tracking during CLOSING:** The background worker polls the ledger
balance. When it reaches zero (all credits fully settled, all in-flight
payments resolved), it executes the sweep transfer to the
`sweep_destination` specified in the original `POST /close` request and,
on confirmation, fires the `funds_swept` event to complete the transition
to CLOSED.

**Sweep destinations (from `CloseIdentityRequest` in `openapi.yaml`):**
- `refund_to_source` — return remaining funds to the sender of the most
  recent inbound transfer (or the last known sender if balance is from
  multiple sources; integrator is notified to verify).
- `integrator_account` — sweep to a bank account reference the integrator
  has on file with Kobo.
- `successor_identity` — redirect balance to another live identity's virtual
  account (e.g. transferring a student's fee balance to a sibling who joins
  the same school).

---

### CLOSED

**Entered from:** CLOSING (funds_swept)

**Meaning:** The account is inactive, balance is zero, and the virtual account
number may no longer accept inbound transfers. Any transfer that arrives
after CLOSED is flagged as `payment_to_closed_account` in the exceptions
queue.

**Reopening:** `POST /v1/identities/{id}/reopen` transitions CLOSED → ACTIVE.
This reuses the existing identity ID and, where Monnify allows it, the same
`accountRef`. See open item #2 in `docs/MONNIFY_INTEGRATION.md` on whether
Monnify allows the same `accountRef` to be reused or requires a new one.

---

### FAILED

**Entered from:** PENDING (provisioning_failed)

**Meaning:** Monnify rejected the provisioning request, or Kobo's retries were
exhausted. No virtual account exists for this identity.

**This state is terminal in v1.** There is no retry-in-place path out of
FAILED. The integrator should create a new identity (new `POST /v1/identities`
with corrected data) if they want to retry. FAILED identities are retained
in the database for audit purposes with `failure_reason` populated.

---

## Identity events log

Every state transition, and some non-transition actions, produce an
`identity_events` row (see `identities.sql: InsertIdentityEvent` and
the `identity_events` table in `0001_init.up.sql`). This is the audit trail
that answers "when did this account close and why" or "when was this name
changed" without requiring anyone to reconstruct it from raw database diffs.

| event_type           | Triggered by |
|---|---|
| created              | `POST /v1/identities` |
| provisioned          | Successful Monnify provisioning response |
| provisioning_failed  | Failed Monnify provisioning response |
| activated            | PENDING → ACTIVE transition |
| limited              | ACTIVE → LIMITED transition |
| unlimited            | LIMITED → ACTIVE transition |
| closing_started      | ACTIVE or LIMITED → CLOSING transition |
| closed               | CLOSING → CLOSED transition |
| reopened             | CLOSED → ACTIVE transition |
| renamed              | `PATCH /v1/identities/{id}` with new `display_name` |
| metadata_updated     | `PATCH /v1/identities/{id}` with new `metadata` |

---

## Webhook events emitted to integrators

These are the lifecycle-related events Kobo pushes to integrators (not to be
confused with the Monnify webhook events Kobo receives internally).

| event name           | Fired when |
|---|---|
| identity.activated   | Account transitions to ACTIVE (including from reopen) |
| identity.limited     | Account transitions to LIMITED |
| identity.closing     | Account transitions to CLOSING |
| identity.closed      | Account transitions to CLOSED |
| identity.reopened    | Account transitions from CLOSED back to ACTIVE |

---

## Implementer notes for `core/internal/account/lifecycle.go`

- `Transition(current State, event Event) (State, error)` is the single
  function through which every state change must pass. It is pure and
  synchronous: no DB calls, no Monnify calls, no logging inside it.
- The callers of `Transition()` — `account/service.go` for provisioning
  results, `handlers/identities.go` for close/reopen requests, and
  `cmd/worker` for KYC-tier and sweep checks — are responsible for
  persisting the result to `identities.state` and writing the corresponding
  `identity_events` row in the same database transaction.
- `AcceptsInboundFunds(state State) bool` is the function the reconciliation
  engine calls to decide whether to credit the ledger or raise an exception.
  See `RECONCILIATION.md` for how this fits into the full reconciliation
  flow.
- `IsTerminal(state State) bool` returns `true` only for FAILED. CLOSED is
  explicitly not terminal because reopen is valid.
