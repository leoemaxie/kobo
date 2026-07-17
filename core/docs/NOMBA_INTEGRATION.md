> **⚠️ WARNING / DISCLAIMER:**
> For backend stability and backward compatibility, the technical implementation details—such as Go package names (`internal/nomba`), database schemas (e.g., column names like `nomba_reference`, `nomba_account_ref`), and internal/external webhook endpoints (like `/webhooks/nomba`)—remain named `nomba`. However, conceptually and operationally under the hood, these components communicate with and process transactions for **Monnify** (not Nomba). All documentation and references within this file have been updated to describe **Monnify** integration flows. Do not modify the underlying code/db structure names to avoid breaking system stability.

# Monnify Integration Reference — for `internal/monnify`

This is the ground-truth contract for everything Kobo's `internal/monnify`
package talks to. It supersedes any assumptions made in `openapi.yaml` or
`ARCHITECTURE.md` about Monnify's own API shape — those documents describe
Kobo's API to its integrators, this document describes Monnify's actual API to
Kobo. Sourced from https://developer.monnify.com (fetched directly, not from
training data) — re-verify against the live docs before final submission if
anything here looks ambiguous.

## Environments

| Environment | Base URL |
|---|---|
| Sandbox (all hackathon work) | `https://sandbox.api.monnify.com/v1` |
| Production (post-certification, after KYC) | `https://api.monnify.com/v1` |

Store both in config; never hardcode. `platform/config/config.go` should read
`MONNIFY_BASE_URL`, `MONNIFY_CLIENT_ID`, `MONNIFY_CLIENT_SECRET`, `MONNIFY_ACCOUNT_ID`,
`MONNIFY_WEBHOOK_SECRET` from env.

### Test instruments (sandbox only)

| Instrument | Value |
|---|---|
| Test card (success) | `5060 6666 6666 6666 666`, any future expiry, any CVV |
| Test card (insufficient funds) | `5060 6666 6666 6666 674` |
| Test bank for inbound transfer | Wema Bank, account `0000000000` accepts any inbound transfer |

Use the test bank/account in `scripts/seed.go` and in reconciliation
integration tests to simulate inbound transfers without a real bank rail.

## Authentication

Monnify uses **OAuth2 client_credentials**, not a static API key header. This
changes the original assumption in earlier drafts — there is no single
bearer token to store in `.env`; tokens are short-lived and must be
refreshed.

### Obtain token

```
POST {base_url}/auth/token/issue
Headers:
  Content-Type: application/json
  accountId: <MONNIFY_ACCOUNT_ID>
Body:
  {
    "grant_type": "client_credentials",
    "client_id": "<MONNIFY_CLIENT_ID>",
    "client_secret": "<MONNIFY_CLIENT_SECRET>"
  }
```

Response:
```json
{
  "code": "00",
  "description": "Success",
  "data": {
    "businessId": "...",
    "access_token": "eyJ...",
    "refresh_token": "01h4...",
    "expiresAt": "2026-06-30T14:33:00Z"
  }
}
```

- `access_token` expires after **30 minutes**.
- Refresh proactively, 5 minutes before expiry, via `POST /auth/token/refresh` (same `accountId` header, body `{"grant_type": "refresh_token", "refresh_token": "..."}`, auth header `Authorization: Bearer <current_token>`).
- Every Monnify API call after this requires both:
  - `Authorization: Bearer <access_token>`
  - `accountId: <MONNIFY_ACCOUNT_ID>`

### Implementation note for `internal/monnify/client.go`

Build a `tokenManager` that wraps the access token in a mutex-guarded struct,
checks expiry before every outbound call, and refreshes transparently. Do not
let individual request methods (`CreateVirtualAccount`, `FetchTransactions`,
etc.) handle token refresh themselves — that logic belongs in one place,
wrapped around the underlying HTTP client (e.g. as `http.RoundTripper` or a
`doAuthenticated` helper). Every response envelope follows the same shape:
`{"code": "00", "description": "...", "data": {...}}`. Code `"00"` means
success; anything else is an error. Map this to Kobo's internal error type in
one shared response-decoding helper, not per-endpoint.

## Virtual Account Provisioning

```
POST {base_url}/accounts/virtual
Headers:
  Authorization: Bearer <access_token>
  accountId: <MONNIFY_ACCOUNT_ID>
  Content-Type: application/json
Body:
  {
    "accountRef": "<16-64 char unique reference>",
    "accountName": "<8-64 char account holder name>",
    "bvn": "<optional, 11 digits>",
    "expiryDate": "<optional, e.g. 2026-01-30 12:15:00, be careful with this>",
    "expectedAmount": "<optional decimal, e.g. 200.00>"
  }
```

Field mapping to Kobo's internal model:

| Monnify field | Kobo concept |
|---|---|
| `accountRef` | Use Kobo's internal `identity.id` (UUID) here. This is what ties the Monnify account back to a Kobo identity without a separate lookup table. |
| `accountName` | Maps to `Identity.DisplayProfile.DisplayName`, must be 8-64 chars — pad or validate on Kobo's side before calling Monnify, since Monnify will reject short names (e.g. very short student names) at the API level. |
| `bvn` | Optional. Leave unset for tier_1 (no-KYC) identities. Only populate when the integrator has actually collected a BVN for a higher-tier identity. |
| `expectedAmount` | Do not set this for most Kobo use cases — it is an optional cap and Kobo's own ledger/lifecycle layer is what should enforce KYC-tier-based limits, not Monnify's expectedAmount field. Leave unset unless a specific product need calls for it. |

Response `data` shape:
```json
{
  "createdAt": "2026-06-30T07:09:06.900Z",
  "accountHolderId": "uuid",
  "accountRef": "the accountRef you sent",
  "bvn": "...",
  "accountName": "...",
  "bankName": "Monnify MFB",
  "bankAccountNumber": "9391076543",
  "bankAccountName": "Monnify/Ifeoluwa Adeboye",
  "currency": "NGN",
  "callbackUrl": "...",
  "expired": false
}
```

Map `bankAccountNumber` to `VirtualAccountSummary.account_number` and
`bankName` to `VirtualAccountSummary.bank_name` in Kobo's API responses (see
`openapi.yaml`). Note `bankAccountName` often comes back prefixed with a
Monnify merchant tag (e.g. `"Monnify/Ifeoluwa Adeboye"`) — do not pass this
through to integrators as-is; store Kobo's own clean `accountName` for
display and keep the raw Monnify value only in the internal record for
debugging/audit.

On non-`"00"` response codes (400/401/403/404/429/500 per Monnify's spec), map
to the PENDING -> FAILED transition in `account/lifecycle.go` and store the
raw Monnify error code/description as `Identity.failure_reason`.

## Reconciliation Sources

Kobo's reconciliation engine (per `ARCHITECTURE.md` Section on
`internal/reconciliation`) has two real Monnify data sources, not one:
webhooks (primary signal) and the Transactions API (fallback / sweep).

### Webhooks (primary)

Monnify calls a URL you configure on the Monnify dashboard (Developer ->
Webhook Setup). There is no API call to register a webhook URL
programmatically for this hackathon — it's configured manually on the
dashboard, so this is a one-time setup step, not something `internal/monnify`
needs to automate.

**Relevant event for Kobo:** `payment_success` — triggered when a payment is
credited via virtual account transfer, card, or PayByTransfer. This is the
event `handlers/webhooks.go` should branch on; ignore other event types
(`payout_success`, `payment_failed`, `payment_reversal`, `payout_failed`,
`payout_refund`) for v1 since Kobo only handles inbound attribution.

**Headers on every webhook request** (case-insensitive, normalize to
lowercase before reading):
```
monnify-signature: <base64 HMAC-SHA256>
monnify-signature-algorithm: HmacSHA256
monnify-signature-version: 1.0.0
monnify-timestamp: <RFC-3339 timestamp>
```

**Payload shape for `payment_success`:**
```json
{
  "event_type": "payment_success",
  "requestId": "uuid",
  "data": {
    "merchant": {
      "walletId": "...",
      "walletBalance": 539.4,
      "userId": "..."
    },
    "terminal": {},
    "transaction": {
      "aliasAccountNumber": "9679136...",
      "fee": 0.6,
      "sessionId": "...",
      "type": "vact_transfer",
      "transactionId": "API-VACT_TRA-...",
      "aliasAccountName": "Peter/Peter Enterprise",
      "responseCode": "",
      "originatingFrom": "api",
      "transactionAmount": 120,
      "narration": "Transfer from JOHN GRASS",
      "time": "2026-02-06T10:21:56Z",
      "aliasAccountReference": "...",
      "aliasAccountType": "VIRTUAL"
    },
    "customer": {
      "bankCode": "305",
      "senderName": "JOHN GRASS",
      "bankName": "Paycom (Opay)",
      "accountNumber": "8168900XX"
    }
  }
}
```

Field mapping for `internal/reconciliation/engine.go`:

| Monnify field | Kobo use |
|---|---|
| `data.transaction.transactionId` | This is the idempotency key. Use exactly this string as the unique constraint in the `idempotency_keys` table, not `sessionId` or `aliasAccountReference`. |
| `data.transaction.aliasAccountNumber` | The virtual account number that received the funds — match against `virtual_accounts.account_number` to find the identity. |
| `data.transaction.transactionAmount` | Amount in Naira as a number (NOT kobo, despite Kobo's own internal model using kobo integers). **Multiply by 100 and round to get kobo** when writing to Kobo's ledger, since `openapi.yaml` defines `amount_kobo` as an integer in the smallest unit. |
| `data.transaction.time` | Use as `Transaction.occurred_at`. |
| `data.transaction.narration` / `data.customer.senderName` | Store for the statement/audit trail; not used for matching (matching is purely by account number, since accounts are 1:1 with identities). |

**Note on `aliasAccountType`:** confirm this is `"VIRTUAL"` before processing
as a virtual-account credit; Monnify's webhook system also reports card and
terminal payment events, which Kobo should ignore since identities are only
tied to virtual accounts in this design.

### Signature verification (exact algorithm — do not approximate)

This is **not** a generic HMAC of the raw JSON body. Monnify's signature is an
HMAC-SHA256 of a specific colon-delimited string built from selected fields
plus the timestamp header, then base64-encoded. Implement exactly this in
`internal/monnify/webhook.go`:

```go
hashingPayload := fmt.Sprintf(
    "%s:%s:%s:%s:%s:%s:%s:%s:%s",
    payload.EventType,                  // event_type
    payload.RequestID,                  // requestId
    payload.Data.Merchant.UserID,       // data.merchant.userId
    payload.Data.Merchant.WalletID,     // data.merchant.walletId
    payload.Data.Transaction.TransactionID, // data.transaction.transactionId
    payload.Data.Transaction.Type,      // data.transaction.type
    payload.Data.Transaction.Time,      // data.transaction.time
    responseCode,                       // data.transaction.responseCode, "" if literal string "null"
    monnifyTimestamp,                     // the monnify-timestamp header value, verbatim
)

mac := hmac.New(sha256.New, []byte(webhookSecret))
mac.Write([]byte(hashingPayload))
expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

// compare case-insensitively against the monnify-signature header
```

Critical implementation details an agent must not skip:
- Field order in the colon-delimited string matters exactly as shown above.
- `responseCode` needs special handling: if the JSON value is the literal
  string `"null"`, treat it as empty string `""` before hashing — this is an
  actual quirk in Monnify's own reference implementation, not a Kobo
  invention.
- Use `strings.EqualFold` (or equivalent case-insensitive compare) when
  comparing signatures, not `==`.
- Reject the request with 401 if the `monnify-timestamp` header is missing, or
  if it's older than a reasonable replay window (5 minutes is a defensible
  default, matching the convention used elsewhere in Kobo's own API auth).
- The webhook secret is configured manually on the Monnify dashboard per
  webhook URL, store it as `MONNIFY_WEBHOOK_SECRET`.

### Idempotency on outbound calls Kobo makes to Monnify

Separately from webhook idempotency (above), Monnify supports an
`X-Idempotent-key` header on outbound requests Kobo makes to Monnify (e.g.
provisioning calls), generated as a UUIDv4 per logical operation. Use this on
`CreateVirtualAccount` calls specifically, so a network failure followed by a
Kobo-side retry does not create a second Monnify account for the same
identity. Generate the idempotency key once per identity-provisioning
attempt and reuse it across retries of that same attempt, not a new key per
HTTP call.

### Webhook retry behavior (Monnify -> Kobo)

If Kobo's webhook endpoint does not return a 2XX status, Monnify retries with
exponential backoff: 2 min, ~5 min, ~11 min, 24 min, ~53 min (5 retries
total). Two implications for `handlers/webhooks.go`:

1. Always return 200 once the event is durably recorded (even if downstream
   processing like account-state updates happens async) — do not return
   non-2XX for business-logic reasons (e.g. "account not found"), only for
   genuine processing failures, otherwise legitimate failures and stale
   retries become indistinguishable.
2. Because retries can arrive up to ~53 minutes later, the reconciliation
   sweep's polling window (`internal/reconciliation/sweep.go`) should be at
   least that wide when deciding whether a missing webhook counts as
   "delayed" vs. "genuinely never sent."

### Transactions API (fallback / sweep source)

Used by `internal/reconciliation/sweep.go` to backfill any account with no
matching webhook event within the configured window.

```
GET {base_url}/transactions/virtual?virtual_account=<account_number>&dateFrom=<date>&dateTo=<date>
Headers:
  Authorization: Bearer <access_token>
  accountId: <MONNIFY_ACCOUNT_ID>
```

Response `data.results[]` shape (note: different field names than the
webhook payload — this is Monnify's transaction-history shape, not the
webhook-event shape, do not assume they're identical):

```json
{
  "id": "API-VACT_TRA-...",
  "status": "SUCCESS",
  "amount": "100.0",
  "type": "vact_transfer",
  "timeCreated": "2025-06-24T11:31:35.017Z",
  "paymentVendorReference": "...",
  "recipientAccountNumber": "8578228675",
  "recipientAccountType": "VIRTUAL",
  "senderName": "John Doe",
  "entryType": "CREDIT",
  "narration": "Transfer from John Doe"
}
```

Sweep matching logic: filter `entryType == "CREDIT"` and
`recipientAccountType == "VIRTUAL"`, match `recipientAccountNumber` against
Kobo's `virtual_accounts.account_number`, and use `id` (not
`paymentVendorReference`) as the idempotency key — confirm this is the same
identifier space as the webhook's `transactionId` field before assuming
they're interchangeable; if in doubt, store both and dedupe on whichever one
is present, since the sweep and the webhook may reference the same
underlying transaction with different reference fields depending on Monnify's
internal routing. This is the single biggest "verify against sandbox before
trusting" item in this document — write an integration test that fires a
real sandbox transfer (using the Wema Bank `0000000000` test account) and
compares the webhook payload's `transactionId` against the sweep endpoint's
`id` for the same transfer before finalizing the dedup key choice.

Pagination uses a `cursor` field returned in the response; pass it back as a
query param on the next request. See Monnify's pagination guide for the exact
param name if `cursor` alone doesn't work as expected in sandbox testing.

## Open items to verify directly against the Monnify sandbox before launch

These are flagged rather than guessed at, since getting them wrong silently
breaks reconciliation accuracy, which is a named judging criterion:

1. Confirm whether `transactionId` (webhook) and `id` (transactions API) are
   the same value for a single real transfer in sandbox — this determines
   the idempotency key strategy in `reconciliation/idempotency.go`.
2. Confirm the exact behavior when `accountRef` (Kobo's identity ID) is
   reused after an account is closed and reopened — whether Monnify allows
   re-provisioning against the same `accountRef` or requires a new one. This
   directly affects the CLOSED -> ACTIVE reopen transition in
   `account/lifecycle.go`.
3. Confirm KYC tier ceiling values and how Monnify signals an approaching or
   breached tier limit (this wasn't found in the fetched docs — it may
   require a direct question to Monnify's hackathon support channel, since
   tier-limit behavior is central to the `ACTIVE <-> LIMITED` transition
   that's explicitly named in the judging criteria).
