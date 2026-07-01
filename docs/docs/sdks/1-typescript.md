---
sidebar_position: 1
title: With TypeScript
---

# TypeScript SDK

The official `@kobo/sdk` is a lightweight, zero-dependency TypeScript client for the Kobo API. It provides strong typing, automatic retries, and comprehensive error handling using the native `fetch` API.

## Installation

```bash
npm install @kobo/sdk
# or
pnpm add @kobo/sdk
```

## Initialization

Initialize the client with your Kobo API key and secret. Use `KoboClient.sandbox` for testing.

```typescript
import { KoboClient } from "@kobo/sdk";

// Production
const kobo = new KoboClient(
  process.env.KOBO_API_KEY,
  process.env.KOBO_API_SECRET
);

// Sandbox Testing
const koboSandbox = KoboClient.sandbox(
  "kobo_test_pk_...",
  "kobo_test_sk_..."
);
```

### Health Check

Useful as a connectivity check before sending real traffic:
```typescript
const health = await kobo.health();
```

---

## Identities

Identities represent end-users or businesses.

### Create Identity
Register a new identity and provision its virtual account.
```typescript
const identity = await kobo.identities.create({
  externalReference: "user-12345",
  profile: {
    firstName: "John",
    lastName: "Doe",
    email: "john.doe@example.com",
    bvn: "22222222222"
  }
});
```

### Get Identity
Fetch an identity record and its current state.
```typescript
const identity = await kobo.identities.get("idx_123abc");
```

### Update Identity
Update display profile fields.
```typescript
const updated = await kobo.identities.update("idx_123abc", {
  profile: {
    displayName: "John D."
  }
});
```

### Close Identity
Initiate closure of an identity's virtual account (asynchronous).
```typescript
await kobo.identities.close("idx_123abc", { reason: "User requested account deletion." });
```

### Reopen Identity
Reopen a CLOSED identity's account.
```typescript
await kobo.identities.reopen("idx_123abc");
```

---

## Accounts

Accounts handle ledger operations and statements.

### List Transactions
Retrieve a paginated, reconciled transaction history.
```typescript
const transactions = await kobo.accounts.listTransactions("acc_123abc", {
  limit: 50,
  cursor: "next_cursor_token"
});
```

### Get Statement
Retrieve a structured statement for a given period.
```typescript
// Omit period for current month
const statement = await kobo.accounts.getStatement("acc_123abc", { period: "2024-03" });
console.log(`Balance: ${statement.balance}`);
```

---

## Exceptions

Exceptions occur during unmatched or misdirected payments.

### List Exceptions
List misdirected-payment or unmatched-transfer cases.
```typescript
const exceptions = await kobo.exceptions.list({ status: "open", limit: 20 });
```

### Resolve Exception
Apply a resolution to a flagged exception.
```typescript
await kobo.exceptions.resolve("exc_123abc", {
  action: "refund",
  reason: "Customer requested reversal."
});
```
