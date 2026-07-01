# Kobo SDKs

This directory contains the official software development kits (SDKs) for the Kobo API. These SDKs are designed to provide a native, type-safe, and zero-dependency developer experience for integrating with the Kobo platform across multiple programming languages.

## Available SDKs

Currently, we officially support and maintain the following SDKs:

- **[TypeScript (Node.js/Browser)](typescript/README.md)**: A zero-dependency SDK using the native `fetch` API, providing fully typed interfaces for all resources (ES2022).
- **[Go](go/README.md)**: A zero-dependency SDK utilizing the standard `net/http` package, featuring context-aware methods and strongly typed models.
- **[Java](java/README.md)**: A robust SDK tailored for Java applications, leveraging standard Java networking and concurrency models.

## Core Design Principles

All Kobo SDKs adhere to the following architectural guidelines to ensure a consistent and robust developer experience:

1. **Zero External Dependencies**: Where possible, SDKs rely solely on language-native standard libraries to minimize integration friction and security risks.
2. **Type Safety**: Full typing of requests, responses, and errors, ensuring that developers get immediate feedback from their IDEs or compilers.
3. **Consistent Error Handling**: All SDKs translate non-2xx HTTP responses into structured API errors with programmatic error codes, rather than just error message strings. This allows for stable, predictable error branching logic.
4. **Resource-Scoped Sub-Clients**: The API surface is organized into logical domains (e.g., `client.identities`, `client.accounts`, `client.exceptions`) that cleanly map to the underlying REST API structure.
5. **Environment Configuration**: Easy switching between `Live` and `Sandbox` environments to support testing and production workflows.

## Common Architecture

Across all implementations, you initialize a core client using your public and secret keys, and then interact with resource-specific sub-clients. 

Example generalized flow:
1. Initialize the client with API keys.
2. Call a sub-client method (e.g., `client.identities.create(...)`).
3. Handle any structured API errors based on the API `code` property.

For specific installation instructions, usage examples, and advanced configurations, please refer to the README file in the respective language's subdirectory.
