# Dragon Market

Secure Go marketplace for Guilds trading Common/Rare items (limit orders) and Legendary items (auctions).

## Stack

- **Go 1.26** + **Fiber** HTTP
- **PostgreSQL 16** + **GORM**
- **Docker Compose** for local runtime

## Quick start

```bash
cp .env.example .env
docker compose up --build
```

API: `http://localhost:8080`  
Postgres: `localhost:5432`

Health check:

```bash
curl http://localhost:8080/health
```

## Configuration

See `.env.example`. Key variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Server listen address |
| `POSTGRES_*` | — | Database connection |
| `AUCTION_DURATION` | `24h` | Legendary auction window |
| `AUCTION_CLOSER_INTERVAL` | `1m` | Expired auction worker |
| `ORACLE_SYNC_INTERVAL` | `30s` | Oracle price cache refresh |

## API

Auth for writes: header `X-Guild-ID`  
Idempotency: optional `Idempotency-Key` on `POST /items`, `POST /items/{id}/purchase`, `POST /items/{id}/bid`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness |
| POST | `/items` | Register item for sale |
| GET | `/items` | List items (+ oracle reference prices) |
| GET | `/items/{id}` | Item details |
| POST | `/items/{id}/purchase` | Buy Common/Rare at list price |
| POST | `/items/{id}/bid` | Place auction bid |
| DELETE | `/items/{id}/bid/{bid_id}` | Withdraw non-winning bid |
| GET | `/auctions` | Active auctions |
| GET | `/auctions/{id}` | Auction + bids |
| GET | `/guilds/{id}/wallet` | Wallet balance |

Import `postman/Dragon-Market.postman_collection.json` for examples.

## Dev seed guilds

| Guild | ID | Balance |
|-------|-----|---------|
| Iron Vanguard | `11111111-1111-1111-1111-111111111111` | 1,000,000 |
| Shadow Syndicate | `22222222-2222-2222-2222-222222222222` | 1,000,000 |
| Crystal Forge | `33333333-3333-3333-3333-333333333333` | 1,000,000 |

Daily purchase cap: **500,000** per guild.

## Architecture

```
cmd/           HTTP handlers, router, workers, main
service/{bc}/  use cases (marketplace, trading, wallet)
domain/{bc}/   entities + repository ports/implementations
pkg/           postgres, oracle, fiberx
```

Dependency rule: `cmd → service → domain ← pkg`

See [docs/adr/001-architecture.md](docs/adr/001-architecture.md).

## Tests

```bash
go test ./...
```

## Local run (without Docker)

```bash
# Postgres must be running and .env configured
go run ./cmd
```
