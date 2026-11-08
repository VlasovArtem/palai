# 0005. LLM: Claude Go SDK behind RecipeGenerator interface; structured outputs

- Status: Proposed
- Date: 2026-11-08

## Context

Recipe generation uses an LLM. We need typed, validated output (ingredients with quantities and
units, steps, timings), testability without network calls, and freedom to change model or
provider later. LLM output is untrusted input.

## Decision

We will call Claude through the official Go SDK, `github.com/anthropics/anthropic-sdk-go`,
hidden behind a domain interface:

```go
type RecipeGenerator interface {
    Generate(ctx context.Context, req GenerateRequest) (GeneratedRecipe, error)
}
```

- The interface and its types live in the recipe domain package. The Claude implementation
  lives in its own package (for example `internal/llm/claude`) and is the only code that
  imports the SDK.
- We use structured outputs (a JSON schema derived from our Go types) so responses decode
  straight into `GeneratedRecipe`. We still validate the decoded value in Go (required fields,
  positive quantities, known units) and treat failures as generation errors.
- The dietary profile goes into the prompt: allergies and strict diets as hard constraints,
  preferences (for example "not sweet") as soft ones. The prompt is advisory only. Safety is enforced afterwards
  by ADR 0006.
- Model ID, max tokens and timeouts come from config. `ANTHROPIC_API_KEY` is read from the
  environment and never logged, echoed in errors, or stored.
- Calls honour `ctx` deadlines. Retries and backoff use the SDK's mechanisms, with limits.
- Tests never call the API. Domain and handler tests use hand-written fake `RecipeGenerator`s.
  The Claude adapter is tested against `httptest.Server` fixtures or recorded responses.

## Consequences

- Swapping model or provider touches one package.
- Tests are deterministic, free and offline.
- Structured outputs make parsing reliable, but schema changes must be kept in sync with Go
  types.
- Costs and latency need monitoring. Token usage should be logged as structured fields, without
  prompt content containing user PII.
- Adding the SDK to go.mod needs approval.

## Status

Proposed.
