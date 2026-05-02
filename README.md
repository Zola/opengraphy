# OpenGraphy

OpenGraphy is an original Go-based Open Graph / Twitter Card metadata checker, social preview renderer, diagnostics tool, and public preview gallery.

It is inspired by common social preview checker workflows, but it does not copy third-party logos, HTML, CSS, JavaScript, images, fonts, or trademarks.

## Features

- Server-side metadata fetching, closer to social crawler behavior than browser-side fetches.
- Open Graph and Twitter Card parsing.
- Facebook, X/Twitter, LinkedIn, WhatsApp, Discord, and Pinterest-style preview cards.
- Diagnostic score and repair suggestions.
- Redis 4-hour metadata cache.
- Anonymous online users, total visits, total checks, cache hit/miss stats.
- Public wall of qualified previews.
- PK voting: winner +1, loser -1, initial rating 1400.
- Leaderboard backed by Redis and synced to SQLite hourly.
- Multilingual UI with language cookie.
- Cookie consent banner and privacy page.

## Local Development

This workspace currently may not have Go installed locally. The intended development path is Docker:

```bash
docker network create proxy
docker compose up --build
```

Then open the service through the configured reverse proxy / Nginx container.

## Environment

Copy `.env.example` to `.env` and set:

- `PUBLIC_BASE_URL`
- `RATE_SALT`
- `BOT_USER_AGENT`

## Redis Keys

- `og:metadata:{sha256(url)}`: metadata cache, TTL 4 hours.
- `og:work:{id}`: qualified public preview work.
- `og:works:recent`: sorted set by latest seen timestamp.
- `og:works:rating`: sorted set by PK rating.
- `stats:visits:total`
- `stats:checks:total`
- `stats:online`
- `stats:locale:{locale}`
- `stats:cache:hits`
- `stats:cache:misses`
- `rl:s:{session}` and `rl:ip:{hash}`: short-lived rate limits.

## Privacy Design

- Submitted metadata check results are cached for 4 hours.
- Plain IP addresses are not stored.
- Rate limiting uses salted hashes with short TTL.
- Anonymous session cookies support presence and rate limiting.
- Gallery stores only qualified public website preview data: URL, final URL, title, description, preview image, score, rating, wins/losses.
- No third-party analytics or advertising cookies are included.

## SSRF Protection

The fetcher:

- Allows only `http` and `https`.
- Blocks localhost, loopback, private, link-local, multicast, and reserved networks.
- Resolves DNS before fetching.
- Re-validates redirect targets.
- Limits redirects to 5.
- Limits HTML response body to 2MB.
- Uses an 8 second timeout.
- Does not execute JavaScript.

## Adding Languages

Add translations in `internal/i18n/i18n.go`. Language choice uses:

1. `lang` cookie.
2. `Accept-Language`.
3. Fallback to English.

## Deployment Notes

The compose file intentionally does not expose the app container directly to the host. Services communicate on an internal network and the external `proxy` network for reverse-proxy integration.

For `opengraph.aston.tw`:

1. Point DNS `opengraph.aston.tw` to the server running the public reverse proxy.
2. Create or keep the shared Docker network:

```bash
docker network create proxy
```

3. Copy `.env.example` to `.env` and set a strong `RATE_SALT`.
4. Start this app:

```bash
docker compose up -d --build
```

5. In the existing public reverse proxy, route `opengraph.aston.tw` to the Docker network alias:

```nginx
server {
    listen 80;
    server_name opengraph.aston.tw;

    location / {
        proxy_pass http://opengraphy-web:80;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

If TLS is terminated at the public reverse proxy, issue the certificate there. This app's internal Nginx listens on port 80 only and stays behind the reverse proxy.

Validate compose changes with:

```bash
docker compose config
```
