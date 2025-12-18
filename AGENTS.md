# Repository Guidelines

This repository tracks RunGrid, a lightweight icon organizer and launcher. Use `DESIGN.md` as the current source of truth for scope, modules, and decisions.

## Project Structure & Module Organization

- Current state: only `DESIGN.md` is present; code scaffolding is not yet committed.
- Planned layout (update when created):
  - `backend/`: Go + Wails services (scanner, icon extractor, launcher, persistence).
  - `frontend/`: UI (React/Vue), search/virtualized grids, settings.
  - `assets/`: default icons, themes, sample data.
  - `scripts/`: build/dev helpers.
  - `tests/`: shared fixtures or integration tests.

## Build, Test, and Development Commands

Tooling is not yet checked in. Once Wails scaffolding exists, expect:

- `wails dev` to run the app in development mode.
- `wails build` to produce a release build.
- `go test ./...` for backend tests.
- `npm test` or `pnpm test` for frontend tests (final choice to be documented).

Update this section as soon as commands are finalized.

## Coding Style & Naming Conventions

- Go: format with `gofmt`; prefer `goimports`. Package names are lower-case; exported symbols use `PascalCase`.
- Frontend: TypeScript preferred; 2-space indentation; components in `PascalCase`.
- CSS: `kebab-case` class names; avoid inline styles for reusable components.
- Storage: SQLite table/column names in `snake_case`.

## Testing Guidelines

- Backend uses Go’s `testing` package; file names `*_test.go`.
- Frontend tests to be added; keep tests deterministic and offline.
- Prioritize tests for scanner accuracy, icon cache behavior, search ranking, and persistence.

## Commit & Pull Request Guidelines

- No git history yet; adopt Conventional Commits (e.g., `feat: add icon cache`).
- PRs should include: summary, rationale, and UI screenshots/GIFs for visual changes.
- Link related issues or design decisions when available.

## Security & Configuration Tips

- Do not hardcode secrets or activation logic in the repo.
- Validate launch targets and keep a URL scheme allowlist.
- Keep signing keys and update credentials out of source control; document required env vars in `docs/` when added.
