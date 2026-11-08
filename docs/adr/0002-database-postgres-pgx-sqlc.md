# 0002. Database: Postgres + pgx + sqlc

- Status: Proposed
- Date: 2026-11-08

## Context

We store users, dietary profiles, recipes (structured ingredients and steps), shares,
annotations and sourdough data (starters, feeding logs). The data is relational with some
semi-structured parts. We want type-safe queries without an ORM's hidden behaviour.

## Decision

We will use PostgreSQL as the only primary datastore.

- Driver: `github.com/jackc/pgx/v5` with `pgxpool`, used natively (not via `database/sql`).
- Queries: hand-written SQL compiled to Go with `sqlc` (pgx/v5 output). Generated code is
  committed, and `make generate` regenerates it.
- Semi-structured fields (for example LLM generation metadata) may use `jsonb`. Core entities
  stay normalized.
- Repository interfaces live with the domain package that uses them. sqlc code stays behind
  them, so domain logic is tested with hand-written fakes.
- No ORM.

## Consequences

- Compile-time checked queries. SQL stays reviewable and explainable.
- Adds sqlc as a dev tool and a codegen step. CI must check that generated code is current.
- Integration tests need a real Postgres (testcontainers or a CI service). These are local
  dependencies, not live external calls. They run behind a build tag or a `-short` skip.
- Dynamic filtering is clunkier in sqlc. We accept a few hand-built queries where needed.
- New dependencies (pgx, sqlc) still need approval when they are added to go.mod.

## Status

Proposed.
