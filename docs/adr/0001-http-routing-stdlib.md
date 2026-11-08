# 0001. HTTP routing: stdlib net/http, no framework

- Status: Proposed
- Date: 2026-11-08

## Context

Palai is an API-only JSON service. Since Go 1.22, `http.ServeMux` supports method matching and
path wildcards (`GET /recipes/{id}`), which covers what we previously needed a router for.
Frameworks (gin, echo, chi, fiber) add dependencies, their own context types or handler
signatures, and upgrade churn. Some still help with middleware chaining and binding.

## Decision

We will use `net/http` and `http.ServeMux` with Go 1.22+ patterns for all routing.

- Handlers are plain `http.Handler` / `http.HandlerFunc`.
- Middleware is `func(http.Handler) http.Handler`, composed by hand in `internal/httpapi`.
- Path params come from `r.PathValue`. JSON encode/decode uses `encoding/json` through small
  local helpers.
- No web framework or third-party router. Adding one needs a new ADR.

## Consequences

- Zero routing dependencies. Handlers are trivially testable with `httptest`.
- 405 and HEAD handling come from the stdlib mux.
- We write our own small helpers (JSON responses, error envelope, request validation,
  middleware chain). They must stay small and well tested.
- No built-in route groups. We group by registering prefixed patterns per feature.

## Status

Proposed.
