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

### Using Locally (Without Package Manager)

To use the SDK locally without publishing to npm:

1. Build the SDK:
   ```bash
   cd sdks/typescript
   npm run build
   ```
2. Reference the local directory in your project's `package.json`:
   ```json
   {
     "dependencies": {
       "@kobo/sdk": "file:../relative/path/to/kobo/sdks/typescript"
     }
   }
   ```
   *Note: If you are not using a package manager at all in your consumer project, you can simply copy the `dist/` directory directly into your project tree and import the compiled Javascript files.*

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

      metadata: { source: "web_signup" }
    });
    
    console.log("Created Identity:", identity.id);

    // Fetch an identity
    const fetchedIdentity = await kobo.identities.get(identity.id);
    console.log("Identity State:", fetchedIdentity.state);

    // List identities
    const identities = await kobo.identities.list({ limit: 10 });
    console.log(`Found ${identities.length} identities`);

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
