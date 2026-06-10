# ADR 001: Dragon Market Architecture

**Status:** Accepted  
**Date:** 2026-06-10

## Context

Dragon Market is a secure marketplace where Guilds list and trade magical items. The PRD requires:

- Fixed-price sales (Common/Rare) and time-bound auctions (Legendary)
- Wallet reservations, daily spend caps, and an append-only ledger
- Resilience to flaky external oracle pricing
- Transactional integrity under concurrent bids

## Decision

We use a **modular monolith** with **DDD-inspired bounded contexts** and **Clean Architecture** layering:

| Layer | Location | Responsibility |
|-------|----------|----------------|
| Delivery | `cmd/` | Fiber HTTP, middleware, background workers |
| Application | `service/{bc}/` | Use cases, DTOs, orchestration |
| Domain | `domain/{bc}/` | Entities, enums, repository interfaces + GORM impl |
| Infrastructure | `pkg/` | Postgres client, oracle client, HTTP helpers |

Bounded contexts:

- **Wallet** — balance, reserve/release/debit/spend, ledger, daily cap
- **Marketplace** — item catalog, listing, limit-order purchase, oracle cache reads
- **Trading** — auctions, bids, anti-snipe, settlement

## Database: PostgreSQL

**Why PostgreSQL?**

- ACID transactions with row-level locking (`SELECT … FOR UPDATE`) for auctions and wallets
- Strong consistency for financial ledger and bid ordering
- Mature tooling (GORM, docker-compose) for interview delivery speed

Schema is managed via GORM `AutoMigrate` at startup for simplicity; production would use versioned SQL migrations.

## Key trade-offs

| Choice | Benefit | Cost |
|--------|---------|------|
| GORM in domain repos | Fast delivery, less boilerplate | Domain coupled to ORM tags |
| AutoMigrate vs Flyway | Zero extra infra | Less migration audit trail |
| `X-Guild-ID` header auth | Matches PRD scope | Not production-grade identity |
| Mock oracle + sync worker | Isolates display prices from trades | Not a real HTTP oracle integration |
| Fiber over stdlib | Ergonomics, middleware | Extra dependency |

## Financial model

```
available = balance − reserved
```

- **Bids:** `Reserve` → outbid `Release` → settle `Debit` + seller `Credit`
- **Limit orders:** immediate `Spend` (debit) + seller `Credit`
- **Ledger:** append-only; idempotency keys on writes

Oracle prices are **display-only**; trades use fixed list prices or auction bids.

## Background workers

1. **Auction closer** — settles or cancels expired auctions
2. **Oracle sync** — refreshes cached reference prices every 30s; invalid readings are discarded

## What we would add with more time

- Versioned SQL migrations and integration tests with testcontainers
- Real oracle HTTP client with circuit breaker and timeouts
- `POST /items/{id}/buy` load tests for concurrent bids/purchases
- OpenTelemetry metrics and structured logging
- OAuth2 / mTLS for guild identity instead of header trust
