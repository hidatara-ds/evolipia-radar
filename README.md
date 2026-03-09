# Evolipia Radar

News aggregator with AI scoring, powered by Neon PostgreSQL and GitHub Actions.

## Quick Setup

### 1. Database (Neon.tech)

1. Sign up: https://neon.tech
2. Create project: `evolipia-radar`
3. Copy connection string
4. Run migrations (see SETUP.md)

### 2. GitHub Actions

1. Add secret `DATABASE_URL` with Neon connection string
2. Push to trigger workflow
3. Worker runs 3x/day: 07:00, 12:00, 19:00 WIB

### 3. Flutter App

Use Neon's HTTP API or direct PostgreSQL connection via `postgres` package.

## Architecture

```
GitHub Actions (3x/day)
    ↓
Worker (Go) → Neon PostgreSQL
    ↓
Flutter App (Direct Query)
```

## Stack

- **Database**: Neon PostgreSQL (serverless)
- **Worker**: Go (one-shot execution)
- **Scheduler**: GitHub Actions
- **Frontend**: Flutter (separate repo)

## Development

```bash
# Build
go build -o worker ./cmd/worker

# Run
export DATABASE_URL="postgresql://..."
go run ./cmd/worker
```

## License

MIT
