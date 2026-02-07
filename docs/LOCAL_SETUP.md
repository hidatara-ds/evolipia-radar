# Panduan Setup & Run Lokal

Panduan lengkap untuk menjalankan evolipia-radar di local environment.

## 📋 Prerequisites

Sebelum mulai, pastikan sudah terinstall:

1. **Go 1.21+** 
   ```bash
   go version
   # Output: go version go1.21.x ...
   ```

2. **PostgreSQL 15+** (atau gunakan Docker)
   ```bash
   psql --version
   # Output: psql (PostgreSQL) 15.x ...
   ```

3. **Docker & Docker Compose** (opsional, untuk PostgreSQL)
   ```bash
   docker --version
   docker-compose --version
   ```

4. **migrate CLI** (untuk database migrations)
   ```bash
   # Install migrate
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Verify
   migrate -version
   ```

## 🚀 Setup Step-by-Step

### Step 1: Clone & Navigate ke Project

```bash
# Jika belum clone, clone dulu
# git clone https://github.com/hidatara-ds/evolipia-radar.git

cd evolipia-radar
```

### Step 2: Install Go Dependencies

```bash
# Download semua dependencies
go mod download

# Atau jika ada dependency baru
go mod tidy

# Verify dependencies terinstall
go list -m all
```

### Step 3: Setup Database

**Opsi A: Menggunakan Docker (Recommended)**

```bash
# Start PostgreSQL dengan docker-compose
docker-compose up -d postgres

# Tunggu beberapa detik sampai PostgreSQL ready
# Cek status
docker ps

# Output harusnya ada container "radar-postgres" running
```

**Opsi B: PostgreSQL Lokal**

Jika sudah punya PostgreSQL lokal, buat database:

```bash
# Login ke PostgreSQL
psql -U postgres

# Di dalam psql, buat database
CREATE DATABASE radar;

# Keluar dari psql
\q
```

### Step 4: Setup Environment Variables (Opsional)

Buat file `.env` atau set environment variables:

```bash
# Windows PowerShell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"
$env:PORT="8080"

# Windows CMD
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/radar?sslmode=disable
set PORT=8080

# Linux/Mac
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"
export PORT="8080"
```

**Default values (jika tidak set):**
- `DATABASE_URL`: `postgres://postgres:postgres@localhost:5432/radar?sslmode=disable`
- `PORT`: `8080`
- `WORKER_CRON`: `*/10 * * * *` (setiap 10 menit)
- `MAX_FETCH_BYTES`: `2000000` (2MB)
- `FETCH_TIMEOUT_SECONDS`: `8`

### Step 5: Run Database Migrations

```bash
# Menggunakan Makefile (recommended)
make migrate-up

# Atau manual
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
```

**Expected output:**
```
Running migrations...
1/u init_schema (XX.XXXXms)
```

**Verify migration:**
```bash
# Cek apakah tables sudah dibuat
psql -U postgres -d radar -c "\dt"

# Harusnya ada tables:
# - sources
# - items
# - signals
# - scores
# - summaries
# - fetch_runs
```

### Step 6: Seed Default Sources (Opsional)

```bash
# Run seed script untuk menambahkan default sources
go run scripts/seed_default_sources.go

# Atau jika ada Makefile target
# make seed
```

### Step 7: Run API Server

**Terminal 1 - API Server:**

```bash
# Menggunakan Makefile
make run-api

# Atau manual
go run ./cmd/api
```

**Expected output:**
```
API server starting on port 8080
```

**Test API server:**
```bash
# Di terminal lain, test health check
curl http://localhost:8080/healthz

# Expected response:
# {"status":"ok"}
```

### Step 8: Run Worker (Terminal Terpisah)

**Terminal 2 - Worker:**

```bash
# Menggunakan Makefile
make run-worker

# Atau manual
go run ./cmd/worker
```

**Expected output:**
```
Worker started with cron schedule: */10 * * * *
Running initial ingestion...
Found 0 enabled sources
Ingestion completed successfully
```

**Note:** Worker akan:
- Run sekali saat start
- Kemudian run setiap 10 menit (sesuai cron schedule)
- Fetch data dari enabled sources
- Process dan store items ke database

## 🧪 Testing Endpoints

Setelah API server running, test beberapa endpoints:

### 1. Health Check
```bash
curl http://localhost:8080/healthz
```

### 2. List Sources
```bash
curl http://localhost:8080/v1/sources
```

### 3. Get Feed (Top Daily)
```bash
curl http://localhost:8080/v1/feed?date=today
```

### 4. Get Rising Items
```bash
curl http://localhost:8080/v1/rising?window=2h
```

### 5. Search Items
```bash
curl "http://localhost:8080/v1/search?q=llm&limit=10"
```

### 6. Test Source Connection
```bash
curl -X POST http://localhost:8080/v1/sources/test \
  -H "Content-Type: application/json" \
  -d '{
    "type": "rss_atom",
    "category": "news",
    "url": "https://openai.com/blog/rss.xml"
  }'
```

### 7. Create Source
```bash
curl -X POST http://localhost:8080/v1/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "OpenAI Blog",
    "type": "rss_atom",
    "category": "news",
    "url": "https://openai.com/blog/rss.xml"
  }'
```

### 8. Enable Source
```bash
# Ganti {source-id} dengan ID dari response create source
curl -X PATCH http://localhost:8080/v1/sources/{source-id}/enable \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

## 📊 Monitoring

### Check Database

```bash
# Login ke PostgreSQL
psql -U postgres -d radar

# Cek jumlah items
SELECT COUNT(*) FROM items;

# Cek sources
SELECT id, name, type, enabled, status FROM sources;

# Cek latest items
SELECT title, url, published_at, domain 
FROM items 
ORDER BY published_at DESC 
LIMIT 10;

# Cek scores
SELECT i.title, s.final, s.hot, s.relevance 
FROM items i 
JOIN scores s ON s.item_id = i.id 
ORDER BY s.final DESC 
LIMIT 10;
```

### Check Logs

API Server dan Worker akan output logs ke console. Monitor untuk:
- API requests
- Ingestion progress
- Errors

## 🛠️ Troubleshooting

### Problem: Database connection failed

**Error:**
```
Failed to connect to database: connection refused
```

**Solution:**
1. Pastikan PostgreSQL running:
   ```bash
   # Docker
   docker ps
   
   # Local PostgreSQL
   sudo systemctl status postgresql  # Linux
   # atau cek di Services (Windows)
   ```

2. Cek DATABASE_URL:
   ```bash
   echo $DATABASE_URL  # Linux/Mac
   echo %DATABASE_URL%  # Windows CMD
   $env:DATABASE_URL   # Windows PowerShell
   ```

3. Test connection:
   ```bash
   psql -U postgres -d radar -c "SELECT 1;"
   ```

### Problem: Migration failed

**Error:**
```
error: migration failed in line X
```

**Solution:**
1. Cek apakah database sudah ada:
   ```bash
   psql -U postgres -l
   ```

2. Rollback migration:
   ```bash
   make migrate-down
   ```

3. Run migration lagi:
   ```bash
   make migrate-up
   ```

### Problem: Port already in use

**Error:**
```
Failed to start server: listen tcp :8080: bind: address already in use
```

**Solution:**
1. Cek process yang menggunakan port 8080:
   ```bash
   # Windows
   netstat -ano | findstr :8080
   
   # Linux/Mac
   lsof -i :8080
   ```

2. Kill process atau ganti PORT:
   ```bash
   # Windows
   taskkill /PID <pid> /F
   
   # Linux/Mac
   kill -9 <pid>
   
   # Atau ganti PORT
   export PORT=8081  # Linux/Mac
   set PORT=8081     # Windows CMD
   $env:PORT="8081"  # Windows PowerShell
   ```

### Problem: No enabled sources

**Worker output:**
```
Found 0 enabled sources
```

**Solution:**
1. Cek sources:
   ```bash
   curl http://localhost:8080/v1/sources
   ```

2. Enable source:
   ```bash
   # Get source ID dari list sources
   curl -X PATCH http://localhost:8080/v1/sources/{source-id}/enable \
     -H "Content-Type: application/json" \
     -d '{"enabled": true}'
   ```

3. Worker akan fetch dari enabled sources pada run berikutnya

### Problem: Go dependencies error

**Error:**
```
go: cannot find module providing package ...
```

**Solution:**
```bash
# Clean module cache
go clean -modcache

# Download dependencies lagi
go mod download

# Verify
go mod verify
```

## 🎯 Quick Start Commands

**Full setup dalam satu go (copy-paste):**

```bash
# 1. Install dependencies
go mod download

# 2. Start PostgreSQL (Docker)
docker-compose up -d postgres

# 3. Wait for PostgreSQL (5 seconds)
sleep 5  # Linux/Mac
timeout /t 5  # Windows CMD

# 4. Run migrations
make migrate-up

# 5. Run API (Terminal 1)
make run-api

# 6. Run Worker (Terminal 2 - buka terminal baru)
make run-worker
```

## 📝 Next Steps

Setelah setup berhasil:

1. **Test endpoints** - Gunakan curl atau Postman
2. **Add sources** - Test dengan RSS feeds
3. **Monitor worker** - Lihat ingestion logs
4. **Check database** - Verify data masuk

## 🔗 Useful Commands

```bash
# Stop PostgreSQL (Docker)
docker-compose down

# View PostgreSQL logs
docker logs radar-postgres

# Rollback migration
make migrate-down

# Check Go version
go version

# Format code
go fmt ./...

# Run tests
go test ./...

# Build binaries
go build -o bin/api ./cmd/api
go build -o bin/worker ./cmd/worker
```

## 📚 Additional Resources

- [README.md](../README.md) - Project documentation (root)
- [ENHANCEMENT_PLAN.md](ENHANCEMENT_PLAN.md) - Future enhancements
- [DEPENDENCIES.md](DEPENDENCIES.md) - Dependencies list

---

**Happy Coding! 🚀**
