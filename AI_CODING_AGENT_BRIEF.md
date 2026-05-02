# OpenGraphy Go Rebuild Brief

## Context

Build a Go implementation inspired by the public behavior and layout of `https://opengraph.dev`, but do not copy its logo, brand name, proprietary icon, exact visual identity, or copyrighted copy. The product should become a free multilingual Open Graph diagnostic and preview tool that demonstrates the company's technical capability.

Reference files in this workspace:

- `index.html`, `index.png`: captured homepage HTML and screenshot.
- `panel.htm`, `panel.png`: captured preview panel HTML and screenshot.
- `_reference/opengraph-dev-mirror/`: downloaded public resources from `https://opengraph.dev` for study only.

Important observation from the downloaded site:

- The homepage is Nuxt SSR/static HTML with Tailwind-style classes.
- The form submits via `GET /panel?url=...`.
- The panel page loads first, then the browser calls `POST /server/check/` on the same origin with JSON `{"url":"https://example.com"}`.
- The target website is fetched by the web server, not directly by the user's browser.
- The browser renders previews and uses `new Image()` to read preview image dimensions where possible.

## Product Goal

Create a Go-powered SaaS-style tool named OpenGraphy or another non-infringing name. It should provide secure website metadata fetching, accurate parsing, trustworthy platform previews, diagnostics, repair suggestions, multilingual UI, cookie consent, and Redis-backed temporary caching.

The frontend may be more complete than OpenGraph.dev. Pixel-level similarity can be used for structure and spacing inspiration, but avoid copying brand assets and exact text. Use a distinct logo, color palette, typography treatment, copy, and illustrations.

## Technical Stack

- Backend: Go.
- HTTP router: use standard `net/http`, Chi, Echo, or Gin. Prefer simple maintainable routing.
- Templates: Go `html/template` or a lightweight server-side rendering approach.
- Frontend: static HTML/CSS/JS served by Go. Keep homepage mostly static.
- CSS: Tailwind build, hand-written CSS, or a minimal utility system. Do not require unnecessary SPA complexity.
- Cache/statistics: Redis.
- Optional CDN/object cache for generated preview images; all generated/fetched preview artifacts expire after 4 hours.
- Deployment must run behind the existing reverse proxy. Do not expose container ports directly to host. Use Docker network `proxy` if compose is introduced later.

## Required Pages

1. Homepage `/`
   - Static-first page.
   - Hero with URL input, submit button, and short value proposition.
   - Show online users and total visit counter from Redis.
   - Explain Open Graph benefits, supported platforms, and diagnostic capabilities.
   - Include language selector and cookie consent.
   - Must not depend on JavaScript for primary content.

2. Preview panel `/panel?url=...`
   - If URL is missing or invalid, redirect or show an error state.
   - Show input bar with submitted URL.
   - Show platform previews:
     - Facebook-style card
     - X/Twitter card
     - LinkedIn card
     - WhatsApp/iMessage-like card
     - Discord/Slack-like card
     - Pinterest/Reddit-style thumbnail where useful
   - Show editable metadata fields:
     - title
     - description
     - canonical URL
     - site name
     - type
     - image URL
     - image width/height
     - locale
     - Twitter card fields
   - Show generated HTML meta tag code with copy button.
   - Show diagnostic report and repair suggestions.

3. API endpoint `POST /api/check`
   - Request: `{"url":"https://example.com"}`
   - Response: structured JSON with metadata, fetch details, image details, diagnostics, preview model, cache status.
   - The backend must fetch the target website. Do not fetch target HTML directly in the user's browser.

4. Optional preview image endpoint
   - `GET /preview/{cacheKey}.png` or CDN URL.
   - Cache artifacts for 4 hours only.

## Metadata Fetching Requirements

Implement secure server-side fetching:

- Allow only `http` and `https`.
- Normalize and canonicalize URLs.
- Reject private, loopback, link-local, multicast, and reserved IP ranges after DNS resolution.
- Re-check resolved IPs on redirects to prevent SSRF.
- Limit redirects, for example max 5.
- Set total timeout and per-request timeout.
- Limit response size, for example 2 MB for HTML.
- Use a realistic user-agent, but identify the service.
- Prefer `GET`; optionally use `HEAD` for image metadata.
- Reject unsupported content types for HTML parsing.
- Do not send user cookies or credentials to target sites.
- Do not persist user-submitted URLs beyond the temporary cache window.

## Metadata Parsing Requirements

Parse initial HTML response only. Do not rely on target-site JavaScript execution for MVP.

Extract:

- `<title>`
- `meta[name=description]`
- Open Graph tags: `og:title`, `og:description`, `og:image`, `og:image:secure_url`, `og:image:width`, `og:image:height`, `og:url`, `og:type`, `og:site_name`, `og:locale`
- Twitter tags: `twitter:card`, `twitter:title`, `twitter:description`, `twitter:image`, `twitter:url`, `twitter:site`
- canonical link
- favicon candidates
- theme color

Resolution rules:

- Prefer Open Graph values for social previews.
- Fall back to Twitter values.
- Fall back to standard title/description.
- Resolve relative image and canonical URLs against the final fetched URL.
- Preserve duplicate `og:image` entries and pick the best candidate for previews.

## Diagnostics

Return a scored report:

- Missing `og:title`, `og:description`, or `og:image`.
- Title length target: around 40-60 characters.
- Description length target: around 120-160 characters.
- Image dimensions target: 1200x630 or same 1.91:1 ratio.
- Warn if image is too small, too large, unreachable, non-HTTPS, redirected too many times, wrong content type, or lacks dimensions.
- Warn if tags appear only client-rendered and not in initial HTML.
- Warn for mixed-case Open Graph property names.
- Warn when canonical URL and `og:url` disagree.
- Warn if Twitter card is missing or inconsistent.

Each diagnostic item should include:

- severity: `pass`, `info`, `warning`, `error`
- field
- human-readable message
- suggested fix
- generated code snippet when useful

## Caching and Privacy

Use Redis:

- Cache `check:{sha256(normalizedURL)}` for 4 hours.
- Cache generated preview artifacts for 4 hours.
- Store only metadata needed for display and diagnostics.
- Do not store IP addresses, full user-agent strings, or long-term user history.
- Homepage counters:
  - `stats:visits:total`
  - `stats:online:{anonymousSessionID}` with short TTL such as 60-120 seconds
- Use anonymous session ID cookie only after cookie consent where legally appropriate.

Cookie consent:

- Show consent banner on first visit.
- Necessary cookies: language preference, consent status, anonymous session.
- No analytics/tracking cookies unless explicitly added and consented.
- Provide "Accept", "Reject optional", and "Language" controls.

## Multilingual Requirements

Support at least:

- English
- Traditional Chinese
- Simplified Chinese
- Japanese
- Korean
- Spanish
- German
- French
- Portuguese

Language selection:

- First priority: user-selected language cookie.
- Second priority: URL prefix or query if implemented.
- Third priority: `Accept-Language` from browser request headers.
- Fourth priority: user-agent locale hints only as weak fallback.
- Default: English or Traditional Chinese, depending on deployment preference.

Manual language selection must save to cookie.

## UI Direction

Do not copy OpenGraph.dev's logo or exact branding. Use a fresh identity.

Use the screenshots only as layout reference:

- Homepage:
  - top nav
  - large URL input in hero
  - educational sections
  - platform cards
  - framework snippets
  - FAQ
  - footer with language selector and external validator links
- Panel:
  - top URL banner
  - centered URL input/check control
  - grid of platform previews
  - metadata editor
  - generated code block
  - diagnostics section
  - footer

Make the panel more SaaS-like:

- Clear status badge: cached/fresh/error.
- Fetch timing and final URL.
- Diagnostic score.
- Repair checklist.
- Copy buttons for individual snippets.
- Loading, error, timeout, and empty states.

## API Shape

Example `POST /api/check` response:

```json
{
  "url": "https://example.com",
  "final_url": "https://www.example.com/",
  "cache": {"hit": false, "ttl_seconds": 14400},
  "fetch": {
    "status": 200,
    "content_type": "text/html; charset=utf-8",
    "duration_ms": 421,
    "redirects": 1
  },
  "meta": {
    "title": "Example",
    "description": "Example description",
    "og": {
      "title": "Example",
      "description": "Example description",
      "image": "https://example.com/preview.png",
      "url": "https://www.example.com/",
      "type": "website",
      "site_name": "Example"
    },
    "twitter": {
      "card": "summary_large_image",
      "title": "Example",
      "description": "Example description",
      "image": "https://example.com/preview.png"
    }
  },
  "image": {
    "url": "https://example.com/preview.png",
    "width": 1200,
    "height": 630,
    "content_type": "image/png",
    "bytes": 245000,
    "reachable": true
  },
  "diagnostics": [],
  "score": 92
}
```

## Implementation Plan

1. Create Go app skeleton with routes `/`, `/panel`, `/api/check`, `/healthz`, and static assets.
2. Implement template rendering and a distinct responsive UI.
3. Implement i18n dictionaries and language middleware.
4. Implement cookie consent and language cookie.
5. Implement Redis connection, 4-hour metadata cache, total visits, and online users.
6. Implement secure URL validator and SSRF-safe HTTP client.
7. Implement HTML metadata parser.
8. Implement image metadata fetcher with size and dimension detection.
9. Implement diagnostics and repair suggestions.
10. Implement platform preview rendering from a normalized preview model.
11. Add tests for URL validation, private IP blocking, redirects, parser fallbacks, cache TTL, and diagnostics.
12. Validate with local screenshots against `index.png` and `panel.png`, but keep original branding and visual identity distinct.

## Acceptance Criteria

- `GET /` renders without JavaScript and shows counters.
- `GET /panel?url=https://www.aston.tw` renders panel shell.
- `POST /api/check` returns metadata JSON fetched by the Go server.
- Private/internal URLs are blocked.
- Redis cache expires check results and preview artifacts after 4 hours.
- Language selection persists in cookie.
- Cookie consent works and no personal data is stored.
- Generated previews are accurate enough for Facebook, X/Twitter, LinkedIn, WhatsApp, Discord/Slack, and Pinterest/Reddit styles.
- Diagnostic report provides actionable fixes.
- No copied OpenGraph.dev logo, icon, product name lockup, or exact marketing copy remains.
