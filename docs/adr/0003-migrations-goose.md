# 0003. Migrations: goose

- Status: Proposed
- Date: 2026-11-08

## Context

The schema (ADR 0002) needs versioned, reviewable migrations that run the same way locally, in
CI and in production, and that sqlc can read as its schema source.

## Decision

We will use `github.com/pressly/goose/v3` with plain SQL migrations.

- Location: `migrations/`, files named `NNNNN_description.sql` with `-- +goose Up` / `Down`.
- Migrations are append-only. Never edit or delete a migration that has been merged; fix
  forward with a new one.
- sqlc reads `migrations/` as its schema.
- Migrations are embedded (`embed.FS`) and applied by an explicit `palai migrate` command or a
  Makefile target, never implicitly on server start.
- Go-function migrations only when SQL cannot express the change, and only with review.

## Consequences

- One source of truth for schema. Fresh databases are reproducible.
- Append-only history grows over time; we may squash with a baseline later via a new ADR.
- Down migrations are written but treated as best-effort; production recovery is fix-forward.
- Adds goose as a dependency, which needs approval when it is added to go.mod.

## Status

Proposed.
