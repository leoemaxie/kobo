# Kobo - B2B Ledger & Virtual Account Infrastructure

Kobo is a robust virtual account management system and B2B ledger designed as a foundation for modern fintech applications. It provides identity-anchored virtual accounts, real-time bank notifications (Push), and a comprehensive ecosystem for seamless integration.

This repository is a monorepo containing the entire Kobo infrastructure, including the backend API, the developer portal, our documentation site, official SDKs, and reference applications.

---

## Table of Contents
- [🏗️ Ecosystem Overview](#️-ecosystem-overview)
- [🔑 Demo Credentials](#-demo-credentials)
- [🚀 Getting Started](#-getting-started)
  - [1. Start the Core API](#1-start-the-core-api)
  - [2. Start the Developer Console](#2-start-the-developer-console)
  - [3. Explore the Documentation](#3-explore-the-documentation)
- [💻 Integrating with Kobo](#-integrating-with-kobo)
  - [Available SDKs](#available-sdks)
  - [Reference Application](#reference-application)
- [🛠️ Contributing](#️-contributing)

---

## 🏗️ Ecosystem Overview

The Kobo platform is composed of several isolated, specialized components, each located in its own directory:

| Component | Path | Description | Stack |
| --- | --- | --- | --- |
| **Kobo Core** | [`/core`](./core) | The central API server and background worker powering the system. It handles the core ledger, webhooks, and database migrations. | Go, PostgreSQL |
| **Kobo Console** | [`/console`](./console) | The developer portal. Integrators use this to manage API keys, billing, and monitor their virtual accounts. | SvelteKit, Tailwind v4 |
| **Documentation Portal** | [`/www`](./www) | The official documentation site featuring guides, API references, and OpenAPI models. | Docusaurus, React |
| **Official SDKs** | [`/sdks`](./sdks) | Zero-dependency, type-safe client libraries for integrating with the Kobo API. | Go, Java, TypeScript |
| **Reference Apps** | [`/apps`](./apps) | Example applications, such as a "School Fees" app, demonstrating how to integrate Kobo into real-world products. | SvelteKit, Tailwind v4 |

---

## 🔑 Demo Credentials

To make exploring the platform as easy as possible, we have provided a set of pre-configured accounts and access keys. You can find the complete list of [Demo Credentials here](https://docs.google.com/document/d/1TGdU9InmpFJVxtCkCXqXaHRgiJ1z3VNTqSfS8jvnf6Y/edit?usp=sharing).

---

## 🚀 Getting Started

To get the full Kobo platform running locally, you'll need to set up the individual services. We recommend starting with **Kobo Core**.

### 1. Start the Core API
The Core API is the foundation of Kobo.
1. Navigate to [`/core`](./core).
2. Set up your `.env` (copying from `.env.sample`).
3. Run migrations: `make migrate-up`.
4. Start the server: `make run` and `make worker`.

*See the [Core README](./core/README.md) for detailed prerequisites (like PostgreSQL).*

### 2. Start the Developer Console
Once the core is running, you can spin up the Kobo Console to manage your account.
1. Navigate to [`/console`](./console).
2. Install dependencies: `pnpm install`.
3. Set up your `.env` file.
4. Start the dev server: `pnpm run dev`.

*See the [Console README](./console/README.md) for further details.*

### 3. Explore the Documentation
To read the full API spec and architectural guides:
1. Navigate to [`/www`](./www).
2. Install dependencies: `pnpm install`.
3. Generate OpenAPI models: `pnpm run gen:docs`.
4. Start the docs server: `pnpm start`.

*See the [Docs README](./www/README.md) for more info.*

---

## 💻 Integrating with Kobo

If you are an external developer looking to build on top of Kobo, you don't need to run this entire monorepo. Instead, you can integrate with Kobo using our official SDKs!

### Available SDKs
- **TypeScript**: `npm install @kobo/typescript-sdk` *(Check the SDK docs for the exact package name)*
- **Go**: `go get github.com/leoemaxie/kobo/sdks/go`
- **Java**: Maven/Gradle integration available.

Check the [`/sdks` directory](./sdks/README.md) for documentation on how to authenticate, create identities, generate virtual accounts, and handle errors robustly.

### Reference Application
Want to see Kobo in action? Check out our [School Fees reference app](./apps/school-fees) in `/apps/school-fees`. It demonstrates a complete end-to-end integration for a school management system provisioning virtual accounts for students!

---

## 🛠️ Contributing

We welcome contributions to any part of the Kobo ecosystem!
- **Core backend changes**: Please refer to `core/docs/ARCHITECTURE.md`.
- **Database schema changes**: Migrations live in `core/migrations/`. Do not use Drizzle in the `console` directory to modify schemas.
- **Documentation**: All guides and OpenAPI specs are generated from `core/openapi.yaml` and styled in `www`.
