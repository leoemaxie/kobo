# Agent Prompt: Kobo Docs Visual Refresh

Copy everything below the line into your coding agent (Claude Code, Cursor,
etc.) with access to the Kobo docs repo.

---

## Context

This is a Docusaurus documentation site for Kobo, an identity-anchored
virtual account infrastructure API built on Nomba's payment rails
(NUBAN-style dedicated accounts, reconciliation, lifecycle management).
The site is live at https://docs.kobo.triumphsystems.tech/ and already has
the correct information architecture: a landing page with a three-feature
grid, a "Documentation" section (Basics → Core Primitives → Lifecycle →
Reconciliation → SDKs), and an API Reference section.

**The content and structure are good. Do not rewrite the copy or
reorganize the nav.** The problem is purely visual: the site currently
reads as an unstyled default Docusaurus theme — generic dark landing page,
no considered accent color, no typographic hierarchy, no real brand
presence. Your job is a styling and polish pass, not a rebuild.

## Design references

Study these for direction, in this priority order:

1. **Resend's docs (resend.com/docs) — primary reference.** This is the
   closest structural match to what Kobo already has: a docs-framework
   site with a landing page that does a feature-grid-then-get-started
   flow. Study specifically: their restrained single-accent-color palette
   against near-black/near-white, generous whitespace, how sparingly they
   use monospace/code styling (only for actual code and technical terms,
   never for body copy), and their type scale (the size/weight jump
   between H1, H2, and body text should feel deliberate, not default).

2. **Plaid docs — secondary reference, for conceptual pages.** Plaid's
   subject matter (identity, accounts, webhooks) is closest to Kobo's own
   domain. Look at how they lay out conceptual/narrative docs pages (the
   equivalent of Kobo's "Core Primitives" and "Lifecycle" pages) — how
   they use diagrams inline with prose, how code blocks are introduced,
   how they handle multi-step flows visually.

3. **Stripe docs — reference for the API Reference section specifically.**
   Even though Kobo's API reference may be auto-generated from
   `openapi.yaml`, study Stripe's visual hierarchy principles: code sample
   always visible and readable, clear separation between "what this
   endpoint does" prose and the request/response example, consistent
   treatment of parameter tables.

4. **Linear's marketing site — reference for hero/landing page
   typographic confidence.** If the current hero section
   ("Virtual Account Infrastructure" or similar) feels flat, look at how
   Linear achieves visual confidence using type weight and spacing alone,
   with minimal color and no unnecessary decoration.

## What to actually do

1. **Establish a real color system.** Pick one primary accent color (do
   not default to generic fintech blue — consider something that
   differentiates Kobo, while staying credible for a financial
   infrastructure product) plus a near-black and near-white for
   dark/light mode, following Docusaurus's built-in dark mode toggle
   rather than fighting it. Define these as CSS custom properties in
   `src/css/custom.css`, overriding Docusaurus's Infima variables
   (`--ifm-color-primary` and its shade variants) rather than hardcoding
   colors throughout components.

2. **Fix the type scale.** Audit `custom.css` and any swizzled components
   for heading sizes, line-height, and font-weight. The jump from H1 to H2
   to body text should feel intentional. If the current font is
   Docusaurus's default, consider a single well-chosen sans-serif for UI
   text and a monospace font for code (system fonts are fine — do not
   over-engineer this with custom font loading unless it's already set
   up).

3. **Rebuild the landing page hero and feature grid** to match the
   restraint of the Resend reference: less visual noise, more whitespace,
   a single clear call-to-action, feature cards that use icons or subtle
   borders rather than heavy backgrounds/shadows.

4. **Add/refine the logo and favicon.** If a logo file already exists in
   `static/img/`, ensure it's used consistently across the navbar, footer,
   and favicon, at appropriate sizes for each (navbar logo, sharper
   favicon at 32x32 and 16x16, and a larger version for any social/OG
   image). If no proper logo exists yet, generate a clean wordmark
   treatment of "Kobo" using the chosen type scale and accent color as a
   placeholder, and flag that a proper vector logo should replace it.

5. **Polish the API Reference pages** for visual hierarchy per the Stripe
   reference: ensure code blocks are legible against the new color
   scheme, parameter tables are scannable, and there's clear visual
   separation between endpoint groups.

6. **Verify dark mode.** Docusaurus ships light/dark by default — confirm
   the new color system holds up in both modes, not just the one you
   design in first.

## Constraints

- Do not change any page content, copy, or navigation structure. This is
  a styling/theming pass only.
- Do not introduce a CSS framework (Tailwind, etc.) into a Docusaurus
  project unless it's already present — work within Infima's CSS variable
  system and Docusaurus's swizzling mechanism for component-level changes.
- Keep changes scoped to `src/css/custom.css`, `static/img/`, and swizzled
  components under `src/theme/` if a particular component genuinely can't
  be restyled through CSS variables alone. Do not swizzle components you
  don't need to touch.
- Test that the build still passes (`npm run build`) and that both light
  and dark mode render correctly before considering this done.

## Deliverable

A visually polished version of the existing site — same structure, same
content, same information architecture — that no longer reads as "default
Docusaurus theme" but as a considered, credible piece of financial
infrastructure branding, in the spirit of Resend's restraint and Stripe's
API-reference clarity.
