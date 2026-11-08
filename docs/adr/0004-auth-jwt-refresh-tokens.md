# 0004. Auth: JWT access tokens + refresh tokens

- Status: Proposed
- Date: 2026-11-08

## Context

Palai is API-only, consumed by first-party clients (mobile and web) and possibly third parties.
Recipes can be shared and annotated, so we need user identity and authorization. Stateless
request auth keeps the hot path cheap, but we must be able to revoke sessions.

## Decision

We will use short-lived JWT access tokens plus long-lived, opaque, rotating refresh tokens.

- Access token: signed JWT (EdDSA/Ed25519 preferred, else ES256; never `none`, and the
  algorithm is pinned on verify). TTL about 15 minutes. Claims: `sub` (user ID), `iat`, `exp`,
  `iss`, `aud`, `jti`. Sent as `Authorization: Bearer`.
- Refresh token: random 256-bit opaque value, stored only as a hash in Postgres, TTL about 30
  days. Rotated on every use. Reusing a rotated token revokes the whole token family.
- Logout and password change revoke refresh tokens. Access tokens expire naturally.
- Passwords are hashed with argon2id.
- Signing keys come from config (env/secret store), support a `kid` for rotation, and are never
  logged.
- Authorization (who can read, annotate or share a recipe) is checked in Go service code, not
  inferred from token claims beyond `sub`.

## Consequences

- No DB lookup per request for authentication. Revocation has up to one access-token TTL of
  lag.
- Refresh rotation and reuse detection need careful, well-tested code and a token table.
- Key management and rotation become an ops concern.
- The JWT library choice needs approval when it is added to go.mod.

## Status

Proposed.
