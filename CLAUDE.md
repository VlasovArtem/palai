# CLAUDE.md

## Overview

Palai is an API-only Go service for generating, sharing and annotating cooking recipes, with
dietary profiles and sourdough tools.

- Allergies are **strict**: a recipe containing a declared allergen (including derivatives, such
  as butter → dairy) must never be returned. Preferences (for example "not sweet") are **soft**.
- Recipes are generated with Claude via `github.com/anthropics/anthropic-sdk-go`, behind the
  `RecipeGenerator` interface, using structured outputs.
- Allergy checks are deterministic Go code. Never rely on the LLM for safety (ADR 0006).

## Commands

Go 1.27+, golangci-lint v2.

```sh
make build   # go build -o bin/palai ./cmd/palai
make run     # go run ./cmd/palai
make test    # go test -race ./...
make lint    # gofmt -l check + go vet ./... + golangci-lint run ./...
make fmt     # gofmt -w
make tidy    # go mod tidy (only when a go.mod change is approved)
```

Single package or test: `go test -race ./internal/config -run TestLoadFrom/defaults`.
Before saying work is done, `make lint` and `make test` must both pass.

## Layout

```
cmd/palai/          main: config, slog JSON logger, HTTP server, graceful shutdown
internal/config/    Load() from PALAI_* env vars
internal/httpapi/   stdlib ServeMux routes, handlers, middleware, JSON helpers
api/openapi.yaml    OpenAPI 3.1 contract; update it with every endpoint change
docs/adr/           Architecture Decision Records (0000-template.md)
docs/specs/         Feature specs, one per feature issue
.golangci.yml       Lint config (v2 format)
```

Planned (see ADRs): domain packages under `internal/<domain>`, `internal/llm/claude`,
`migrations/` (goose), sqlc-generated query code.

## Configuration

All config comes from env vars prefixed `PALAI_`, for example `PALAI_ADDR` (default `:8080`),
`PALAI_LOG_LEVEL`, `PALAI_SHUTDOWN_TIMEOUT`. Add new settings to `internal/config` with a
default, validation and a table-driven test case, and document them in README.md.

## Conventions

- **Errors**: wrap with context using `%w` (`fmt.Errorf("load recipe %s: %w", id, err)`).
  Compare with `errors.Is` / `errors.As`, never `==` or string matching. Don't log *and*
  return the same error; handle it once.
- **Context**: `ctx context.Context` is the first parameter of anything doing I/O or
  long-running work. Never store a context in a struct. Use `t.Context()` in tests.
- **Logging**: `log/slog` with structured fields (`slog.String("recipe_id", id)`), using the
  `*Context` variants when a ctx is available. No `fmt.Print*` or `log.Print*` in service code.
  Messages are lowercase and short.
- **HTTP**: stdlib `net/http` patterns (`GET /recipes/{id}`, `r.PathValue`). No frameworks
  (ADR 0001). Outbound requests always carry a ctx (`http.NewRequestWithContext`).
- **Tests**: table-driven with `t.Run(tt.name, ...)`, `t.Parallel()` where safe, `httptest`
  for handlers. Use **hand-written fakes** over generated mocks or mocking libraries. Fakes
  live next to the interface they implement or in the test file.
- **Interfaces**: define them where they are consumed, keep them small, and return concrete
  types.
- **Packages**: everything non-main goes under `internal/`. Exported identifiers have doc
  comments (revive enforces this).

## Design docs

- Read the relevant **`docs/adr/`** entries before changing architecture. Don't contradict an
  ADR silently; propose a new one (copy `0000-template.md`).
- Each feature has a spec in **`docs/specs/`**, written and approved before implementation.
- **`api/openapi.yaml`** is the API contract. Handlers and spec must agree; update both in the
  same change.

## Hard rules

1. **No go.mod changes without explicit approval.** This covers new dependencies, upgrades and
   Go version bumps. Ask first and say why.
2. **Never log secrets.** This includes `ANTHROPIC_API_KEY`, JWT signing keys, passwords,
   refresh tokens, `Authorization` headers and DB credentials. Don't put them in errors, test
   output or fixtures either.
3. **Migrations are append-only.** Never edit or delete a merged migration; add a new one.
4. **No live external calls in tests.** No real Claude API calls or other third-party network
   calls. Use fakes or `httptest.Server`.
5. **Allergy safety is deterministic.** Every generated or edited recipe passes the Go allergen
   check (ADR 0006), and the check fails closed.
6. Don't weaken `.golangci.yml` or add `//nolint` without a justification comment.

## Workflow

**Use the feature-implementation skill for every feature issue.**
