# Frontend and Tailwind Plan

OpenGraphy can move to Tailwind CSS gradually without changing the Docker networking model or exposing any host ports.

## Goals

- Keep the app server-side rendered with Go templates.
- Use Tailwind for layout, spacing, typography, responsive behavior, and component polish.
- Keep custom CSS only for reusable semantic components that are clearer than long utility lists.
- Avoid a large one-shot rewrite. Migrate one page or component group at a time.

## Preferred Direction

1. Add a small frontend build pipeline with `package.json`, `tailwind.config.js`, and an input CSS file.
2. Generate a committed output file under `web/static/css/`, so production Docker does not need Node.
3. Keep the current `/static` cache strategy and bump asset versions after style changes.
4. Start migration from the highest-value surfaces:
   - homepage hero and primary form
   - preview result page cards
   - gallery voting and leaderboard
   - cookie consent and language controls

## Design Rules

- The UI should feel like a practical SaaS tool for developers and website owners.
- Prefer clear hierarchy, calm spacing, high readability, and dense-but-scannable information.
- Keep card radius at 8px or below unless there is a product reason to change it.
- Do not use decorative gradient blobs or one-color-only palettes.
- Mobile layouts must be checked for long URLs, translated labels, and generated metadata code.

## Migration Rules

- Do not modify `docker-compose.yml` just to add Tailwind.
- Do not introduce a runtime Node dependency in the production container.
- Do not replace all CSS in one commit.
- Each migration step should include visual QA for desktop and mobile.
- Keep class names readable in Go templates; use partials for repeated component markup.

## Proposed Scripts

When Tailwind is added, use:

```sh
npm run build:css
```

The script should compile from `web/assets/css/input.css` to `web/static/css/app.css`.

