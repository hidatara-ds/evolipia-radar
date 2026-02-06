# Implementation Quick Start Guide

This guide provides a quick reference for implementing the enhancement plan.

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional)
- [migrate](https://github.com/golang-migrate/migrate) CLI

## Step-by-Step Implementation

### Phase 1: Foundation (PRs 1-3)
**Goal:** User authentication working end-to-end

1. **Run migration:**
   ```bash
   migrate -path migrations -database $DATABASE_URL up
   ```

2. **Add dependencies:**
   ```bash
   go get github.com/golang-jwt/jwt/v5
   go get golang.org/x/crypto/bcrypt
   ```

3. **Test auth flow:**
   ```bash
   # Register
   curl -X POST http://localhost:8080/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"Test123!"}'
   
   # Login
   curl -X POST http://localhost:8080/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"Test123!"}'
   
   # Use token
   curl http://localhost:8080/v1/users/me \
     -H "Authorization: Bearer <access_token>"
   ```

### Phase 2: User Features (PRs 4-7)
**Goal:** Bookmarks, history, feedback working

1. **Test bookmarks:**
   ```bash
   curl -X POST http://localhost:8080/v1/bookmarks \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"item_id":"<item-uuid>"}'
   ```

2. **Test read history:**
   ```bash
   curl -X POST http://localhost:8080/v1/history \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"item_id":"<item-uuid>"}'
   ```

### Phase 3: Personalization (PRs 8-9)
**Goal:** Personalized feeds with explanations

1. **Set preferences:**
   ```bash
   curl -X PATCH http://localhost:8080/v1/users/me/preferences \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"preferred_topics":["llm","mlops"]}'
   ```

2. **Get personalized feed:**
   ```bash
   curl http://localhost:8080/v1/feed?topic=llm \
     -H "Authorization: Bearer <token>"
   ```

### Phase 4: Digests & Alerts (PRs 10-14)
**Goal:** Email digests and keyword alerts

1. **Configure email (env vars):**
   ```bash
   export SMTP_HOST=smtp.sendgrid.net
   export SMTP_PORT=587
   export SMTP_USER=apikey
   export SMTP_PASS=<sendgrid-api-key>
   ```

2. **Subscribe to digest:**
   ```bash
   curl -X POST http://localhost:8080/v1/digests/subscriptions \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{
       "frequency":"daily",
       "send_at":"09:00:00",
       "timezone":"America/New_York"
     }'
   ```

3. **Add keyword alert:**
   ```bash
   curl -X POST http://localhost:8080/v1/alerts/keywords \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"keyword":"GPT-5"}'
   ```

### Phase 5: Observability (PRs 15-17)
**Goal:** Logging, metrics, health checks

1. **Check metrics:**
   ```bash
   curl http://localhost:8080/metrics
   ```

2. **Check health:**
   ```bash
   curl http://localhost:8080/healthz
   curl http://localhost:8080/healthz/ready
   ```

### Phase 6: CI/CD (PRs 19-20)
**Goal:** Automated testing and deployment

1. **Test CI locally:**
   ```bash
   # Install act (GitHub Actions locally)
   brew install act  # macOS
   
   # Run CI workflow
   act pull_request
   ```

2. **Test Docker build:**
   ```bash
   docker build -f Dockerfile.api -t radar-api .
   docker build -f Dockerfile.worker -t radar-worker .
   ```

## Environment Variables Reference

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/radar?sslmode=disable

# API Server
PORT=8080

# JWT
JWT_SECRET=<32-byte-secret>
JWT_ACCESS_EXPIRY=3600  # 1 hour
JWT_REFRESH_EXPIRY=2592000  # 30 days

# Email (for digests/alerts)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=<api-key>
SMTP_FROM=noreply@evolipia-radar.com

# Observability
LOG_LEVEL=info  # debug, info, warn, error
METRICS_ENABLED=true
TRACING_ENABLED=false
TRACING_ENDPOINT=http://localhost:14268/api/traces

# Worker
WORKER_CRON="*/10 * * * *"  # Every 10 minutes
DIGEST_CRON="0 * * * *"  # Every hour
ALERT_CRON="*/15 * * * *"  # Every 15 minutes

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_HOUR=100
RATE_LIMIT_AUTH_PER_MINUTE=5
```

## Testing Checklist

After each PR, verify:

- [ ] Migration runs successfully (up and down)
- [ ] All tests pass: `go test ./...`
- [ ] Linter passes: `golangci-lint run`
- [ ] API endpoints work (test with curl/Postman)
- [ ] No breaking changes to existing endpoints
- [ ] Environment variables documented
- [ ] README updated if needed

## Common Issues

### Migration fails
- Check PostgreSQL version (15+)
- Verify DATABASE_URL format
- Ensure user has CREATE TABLE permissions

### JWT errors
- Verify JWT_SECRET is set (min 32 bytes)
- Check token expiry times
- Verify token format in Authorization header

### Email not sending
- Check SMTP credentials
- Verify SMTP_HOST and SMTP_PORT
- Test with a simple SMTP client first

### Personalization not working
- Verify user is authenticated (token in header)
- Check user preferences are set
- Verify personalization_weight is reasonable (0.1-0.3)

## Next Steps

1. Review `ENHANCEMENT_PLAN.md` for detailed specs
2. Follow PR implementation order (1-20)
3. Test each PR before moving to next
4. Update documentation as you go
5. Deploy to staging after Phase 3
6. Deploy to production after Phase 6

