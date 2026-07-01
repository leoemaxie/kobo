# Kobo Documentation Portal

This repository contains the official documentation site for the **Kobo B2B Ledger and Reconciliation Engine**, built with [Docusaurus 3](https://docusaurus.io/).

It covers Kobo's Core Concepts, the generated OpenAPI Reference, and usage guides for the TypeScript, Go, and Java SDKs.

## Local Development

### Installation

```bash
pnpm install
```

### Generating the API Reference

Before starting the server, you must generate the OpenAPI documentation from the `core/openapi.yaml` spec. 

```bash
pnpm run gen:docs
```
*(If you change the OpenAPI spec, run `pnpm run clean:docs` followed by `gen:docs` again).*

### Starting the Server

```bash
pnpm start
```

This command starts a local development server and opens up a browser window at `http://localhost:3000`. Most changes are reflected live without having to restart the server.

## Build for Production

```bash
pnpm build
```

This command generates static content into the `build` directory and can be served using any static content hosting service (e.g., GitHub Pages, Vercel, Netlify).

## Architecture & Styling

- **Theme**: We use a heavily customized dark-mode Docusaurus theme tailored to the Kobo brand (Dark Onyx backgrounds, Electric Gold accents).
- **Typography**: Uses `Outfit` for headings and `Inter` for body copy via Google Fonts.
- **OpenAPI**: Powered by `docusaurus-plugin-openapi-docs` to render interactive swagger-like API models natively within MDX.
`