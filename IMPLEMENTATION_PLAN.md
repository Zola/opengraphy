# OpenGraphy Implementation Plan

## Product Shape

OpenGraphy is a Go-based Open Graph and Twitter Card checker with three public surfaces:

1. Metadata checker and preview panel.
2. Public wall of qualified website previews.
3. Facemash-style preview PK and leaderboard.

The implementation is original and does not reuse OpenGraph.dev assets, HTML, CSS, JS, logos, images, fonts, or brand identity.

## Architecture

- Go app behind reverse proxy.
- Chi routes and Go templates.
- Redis for 4-hour metadata cache, stats, presence, rate limits, gallery zsets, and fast PK updates.
- SQLite for durable gallery works and PK scores.
- Background sync writes Redis gallery state to SQLite every hour.
- Homepage remains mostly static; stats and presence load through APIs.

## Main Flows

### Check Flow

1. User submits a URL.
2. Go normalizes URL and checks rate limits.
3. Redis metadata cache is checked by `og:metadata:{sha256(url)}`.
4. On miss, Go fetches target HTML server-side with SSRF protection.
5. Parser extracts Open Graph, Twitter Card, canonical, favicon, robots, and fallback metadata.
6. Image inspector fetches only enough image data to detect dimensions.
7. Diagnostics score and repair suggestions are generated.
8. Result is cached for 4 hours.
9. If the preview is qualified, it is added to the gallery.

### Gallery Flow

- Recent section: latest 8 qualified previews by `og:works:recent`.
- Random section: 20 random previews from qualified works.
- Leaderboard: sorted by `og:works:rating`.

### PK Flow

1. `/api/pk/pair` returns two random works.
2. User chooses the better preview.
3. Winner rating +1, loser rating -1.
4. Redis updates immediately.
5. SQLite sync runs hourly for durability.

## Privacy

- No plain IP storage.
- No user URL history persisted as user history.
- Metadata cache expires after 4 hours.
- Gallery stores only qualified public website preview data.
- Necessary cookies only: language, consent, anonymous presence/session.

## Current Scope

This first implementation includes all core routes, parser, fetcher, diagnostics, Redis/SQLite gallery, i18n, consent banner, static pages, Docker and Nginx deployment files, plus tests for the riskiest pure logic.
