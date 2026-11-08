# Palai

API-only Go service for generating, sharing and annotating cooking recipes, with dietary
profiles (strict allergies, soft preferences) and sourdough tools. Recipes are generated with
Claude. Allergy checks are deterministic Go code.

## Quick start

Requires Go 1.27+ and [golangci-lint v2](https://golangci-lint.run/).

```sh
make run                      # listens on :8080
curl localhost:8080/healthz   # {"status":"ok"}
make test lint
```

## Configuration

| Variable                 | Default | Description                       |
|--------------------------|---------|-----------------------------------|
| `PALAI_ADDR`             | `:8080` | HTTP listen address               |
| `PALAI_LOG_LEVEL`        | `info`  | `debug`, `info`, `warn`, `error`  |
| `PALAI_SHUTDOWN_TIMEOUT` | `10s`   | Graceful shutdown deadline        |

## Docs

- API: [`api/openapi.yaml`](api/openapi.yaml)
- Architecture decisions: [`docs/adr/`](docs/adr/)
- Feature specs: [`docs/specs/`](docs/specs/)
- Contributor and agent guide: [`CLAUDE.md`](CLAUDE.md)
