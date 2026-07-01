# Kobo TypeScript SDK

Official TypeScript SDK for the Kobo API. This SDK is built with zero external dependencies and uses the native `fetch` API. It provides fully typed interfaces for all Kobo API resources.

## Features
- Zero external dependencies
- Native `fetch` API under the hood
- Fully typed with TypeScript (ES2022)
- Resource-scoped sub-clients (Identities, Accounts, Exceptions)

## Installation

```bash
npm install @kobo/sdk
# or
yarn add @kobo/sdk
# or
pnpm add @kobo/sdk
```

## Quick Start

```typescript
import { KoboClient } from "@kobo/sdk";

// Initialize the client
const kobo = new KoboClient("kobo_live_pk_...", "kobo_live_sk_...");

// Or use the sandbox environment
// const kobo = KoboClient.sandbox("kobo_test_pk_...", "kobo_test_sk_...");

async function run() {
  try {
    // Check API health
    const health = await kobo.health();
    console.log("Health:", health);

    // Create an identity
    const identity = await kobo.identities.create({
      external_reference: "customer_12345",
      display_name: "Jane Doe",
      kyc_tier_hint: "tier_1",
      metadata: { source: "web_signup" }
    });
    
    console.log("Created Identity:", identity.id);

    // Fetch an identity
    const fetchedIdentity = await kobo.identities.get(identity.id);
    console.log("Identity State:", fetchedIdentity.state);

    // List transactions for an account
    if (fetchedIdentity.virtual_account) {
        // Wait for account to be provisioned, then list txns
        const txns = await kobo.accounts.listTransactions(fetchedIdentity.id);
        console.log("Transactions:", txns.data);
    }
    
  } catch (error) {
    if (error.name === "KoboError") {
      console.error(`API Error: ${error.code} - ${error.message}`);
    } else {
      console.error("Unknown Error:", error);
    }
  }
}

run();
```

## Error Handling

The SDK throws a `KoboError` when the API returns a non-2xx response. You should branch your error handling logic on the `code` property, not the `message` string.

```typescript
import { KoboError } from "@kobo/sdk";

try {
  await kobo.identities.get("invalid-id");
} catch (error) {
  if (error instanceof KoboError) {
    if (error.code === "not_found") {
      console.log("Identity not found.");
    } else {
      console.log(`API Error: ${error.code}`);
    }
  }
}
```

## License

MIT
