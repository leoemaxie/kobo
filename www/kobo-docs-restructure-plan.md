# Agent Plan: Kobo Docs Restructure + Landing Page

Hand this to the coding agent working on the Docusaurus repo currently live
at docs.kobo.triumphsystems.tech. Goal: move to
kobo.triumphsystems.tech as root, with docs living at /docs, and a real
landing page at /.

---

## Reference sites (open these, don't just read the names)

- **resend.com** — primary reference for the landing page itself. Restrained
  single-accent color, generous whitespace, a hero that's mostly typography
  and one code snippet, feature sections that don't over-explain.
- **resend.com/docs** — reference for how the docs section should feel once
  it's a subpath rather than the whole domain.
- **plaid.com** — reference for landing page structure when the product is
  an identity/account primitive like Kobo's — look at how they explain
  "what this actually does" above the fold without jargon.
- **linear.app** — reference for typographic confidence in the hero section;
  useful if the Resend-inspired version still feels flat.
- **stripe.com/docs/api** — reference only for the /docs/api or /reference
  section, not the landing page.

---

## Step-by-step

### 1. Confirm current routing before changing anything

Have the agent check `docusaurus.config.js`/`.ts` for the current `docs`
plugin config and confirm `routeBasePath`. Since the live site already
serves `/docs/intro` as a URL, this is very likely already set to `'docs'`
(Docusaurus's default) or explicitly configured that way — the agent should
verify, not assume, before moving on. If it's set to something else (e.g.
`''`, meaning docs currently live at root), that's the actual thing to fix
in this step.

### 2. Update site URL config

In `docusaurus.config.js`/`.ts`:
```js
url: 'https://kobo.triumphsystems.tech',
baseUrl: '/',
```
This changes what Docusaurus generates for canonical URLs, sitemaps, and
any absolute links it builds — do this before touching any content so
later steps aren't testing against stale config.

### 3. Build the new landing page at `src/pages/index.tsx`

This file already exists (it's the current three-feature-grid page) — the
agent should treat this as a rebuild of that file's content and layout, not
a new file. Structure, in order:

1. **Hero** — one clear line of what Kobo is ("Identity-anchored virtual
   account infrastructure for Nigerian fintech" or similar, reuse existing
   copy if it's already good), one supporting sentence, one primary CTA
   ("Get Started" → links to `/docs/intro`) and one secondary CTA
   ("View API Reference" or similar). Follow Resend's restraint here: no
   more than this. Do not stack multiple headlines or multiple CTAs.
2. **A single code sample or diagram directly under the hero** — Resend
   does this well: show, don't just tell, immediately below the fold. For
   Kobo, this could be a minimal `POST /v1/identities` request/response
   pair, or a small version of the identity → account → ledger flow
   diagram from `core/docs/ARCHITECTURE.md`.
3. **Feature grid (existing three-feature content, restyled)** — keep the
   three features that are already there if the copy holds up, but restyle
   per Resend's card treatment: subtle borders or icon-led, not heavy
   shadowed boxes.
4. **"How it works" section** — a short 3-4 step visual walkthrough
   (Identity → Virtual Account → Reconciliation → Statement), which doubles
   as a preview of the "Core Primitives" docs page and gives people a
   reason to click into `/docs`.
5. **Footer CTA** — a final "Read the docs" / "Get an API key" pair right
   before the standard Docusaurus footer, so the page doesn't just trail
   off after the feature grid.

### 4. Apply the visual system from the earlier restyle prompt

This step assumes the color/type-scale work from the previous "Agent
Prompt: Kobo Docs Visual Refresh" has either already happened or happens
alongside this restructure — they're the same underlying `custom.css`
changes (accent color as CSS custom properties overriding Infima's
`--ifm-color-primary` family, refined type scale). If that hasn't been done
yet, do it as part of this step rather than shipping a restructured-but-
still-generic-looking site.

### 5. Update navbar and footer links

In `docusaurus.config.js`/`.ts`, the navbar config: confirm the logo links
to `/` (not `/docs/intro`), and that the "Documentation" nav item points to
`/docs/intro` or `/docs` (whichever is the actual docs landing route).
Footer links should be audited the same way — anything currently assuming
docs is the site root needs re-pointing.

### 6. Set up the redirect from the old subdomain

Since `docs.kobo.triumphsystems.tech` is currently live and may already be
linked/bookmarked, set up a redirect (at the DNS/hosting level, e.g. a
Vercel/Netlify redirect rule, or a simple CNAME + server-side redirect) from
`docs.kobo.triumphsystems.tech/*` to `kobo.triumphsystems.tech/docs/*`, so
old links don't 404. This is an infra step, not a Docusaurus code change —
flag it to Leo if the agent doesn't have hosting/DNS access.

### 7. Verify

- `npm run build` passes cleanly.
- `/` renders the new landing page, not the docs sidebar.
- `/docs/intro` (and other existing docs pages) still resolve correctly
  with no broken internal links (Docusaurus's build will actually fail on
  broken links by default, which is a useful built-in check here).
- Light and dark mode both look correct on the new landing page
  specifically, not just on docs pages.
- Old `docs.` subdomain links redirect correctly, if that infra step is
  completed.

---

## What NOT to change in this pass

- Don't touch `/docs/*` content — Basics, Core Primitives, Lifecycle,
  Reconciliation, SDKs stay as they are content-wise. Only their base path
  changes (which, per step 1, may already be correct and require no actual
  content changes at all).
- Don't introduce a different framework or eject from Docusaurus's page
  system. `src/pages/index.tsx` is the correct, supported place for a fully
  custom landing page — no need to reach for anything more complex.
- Don't restructure the sidebar/nav order within `/docs` — that's a
  separate piece of work from this URL/landing-page restructure.
