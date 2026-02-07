# 🚀 Quick Start - Run Lokal

Panduan cepat untuk menjalankan aplikasi di local.

## ⚡ Quick Commands

```bash
# 1. Install dependencies
go mod download

# 2. Start PostgreSQL (Docker)
docker-compose up -d postgres

# 3. Run migrations
make migrate-up

# 4. Terminal 1 - Run API Server
make run-api

# 5. Terminal 2 - Run Worker
make run-worker
```

## ✅ Verify Setup

```bash
# Test health check
curl http://localhost:8080/healthz
# Expected: {"status":"ok"}

# Test feed
curl http://localhost:8080/v1/feed?date=today
```

## 📋 Prerequisites Checklist

- [ ] Go 1.21+ installed (`go version`)
- [ ] Docker installed (`docker --version`)
- [ ] migrate CLI installed (`migrate -version`)
- [ ] Port 8080 available
- [ ] Port 5432 available (PostgreSQL)

## 🐛 Common Issues

**Database connection failed?**
```bash
# Check PostgreSQL running
docker ps

# Restart PostgreSQL
docker-compose restart postgres
```

**Port 8080 in use?**
```bash
# Change PORT
export PORT=8081  # Linux/Mac
$env:PORT="8081"  # Windows PowerShell
```

**No enabled sources?**
```bash
# List sources
curl http://localhost:8080/v1/sources

# Enable a source (replace {id})
curl -X PATCH http://localhost:8080/v1/sources/{id}/enable \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

## 📖 Full Documentation

Lihat [LOCAL_SETUP.md](LOCAL_SETUP.md) untuk panduan lengkap.
