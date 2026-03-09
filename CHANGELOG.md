# Changelog - deploy-supabase Branch

## [1.0.0] - 2026-03-09

### Added
- One-shot worker execution (replaces daemon/scheduler)
- Scrape logs tracking table for monitoring
- Data retention function (45 days for items, 60 days for logs)
- GitHub Actions workflows:
  - Scheduled scraper (3x/day: 07:00, 12:00, 19:00 WIB)
  - PR check (build, vet, test)
  - Weekly cleanup (data retention)
- Supabase setup SQL (RLS policies, views, functions)
- Comprehensive documentation:
  - SUPABASE_SETUP.md - Complete setup guide
  - Updated README.md - Branch overview
  - Updated docs/README.md - Architecture & workflows

### Removed
- API server (`cmd/api/`) - Flutter queries Supabase directly
- HTTP handlers (`internal/http/`) - No API layer needed
- Temporal workflows (`internal/workflows/`, `internal/activities/`)
- ML pipeline scaffold (`internal/mlpipeline/`)
- Kubernetes manifests (`k8s/`)
- Terraform configs (`terraform/`)
- Web UI (`web/`) - Replaced by Flutter app
- Docker configs (docker-compose.yml, Dockerfiles)
- Observability stack (Prometheus, Grafana, Jaeger)
- Legacy archives and scripts
- Mobile deployment scripts (deploy-mobile.*)
- Test infrastructure (`tests/`)
- Unused documentation files

### Changed
- Worker from daemon to one-shot execution
- Database connection to use Supabase PostgreSQL
- Scheduling from internal cron to GitHub Actions
- Architecture from API-based to direct Supabase access
- Simplified Makefile (removed complex targets)

### Dependencies Removed
- github.com/gin-gonic/gin (no HTTP server)
- github.com/robfig/cron/v3 (no internal scheduler)
- Temporal dependencies (no workflows)
- Prometheus/observability dependencies

### Migration Notes
- This branch is optimized for Supabase deployment
- Not compatible with previous API-based architecture
- Requires Supabase account and GitHub Actions setup
- Flutter app (separate repo) needed for frontend

### Breaking Changes
- No API endpoints available
- No web UI included
- No Docker deployment
- No Kubernetes support
- Requires migration to Supabase database
