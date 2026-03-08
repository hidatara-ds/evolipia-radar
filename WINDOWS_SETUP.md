# Windows Setup Guide

## Quick Setup (Automated)

### Option 1: PowerShell (Recommended)
```powershell
# Run in PowerShell (as Administrator if needed)
.\setup-windows.ps1
```

### Option 2: Command Prompt
```cmd
# Run in CMD
setup-windows.bat
```

## Manual Setup

### Step 1: Install golang-migrate

```bash
# In Git Bash or PowerShell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Step 2: Add migrate to PATH

**PowerShell:**
```powershell
$goPath = go env GOPATH
$env:PATH = "$goPath\bin;$env:PATH"
```

**Git Bash:**
```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

**Permanent (Windows):**
1. Open "Environment Variables" in Windows Settings
2. Add `%USERPROFILE%\go\bin` to your PATH
3. Restart terminal

### Step 3: Start PostgreSQL

```bash
docker-compose up -d postgres
```

### Step 4: Run Migrations

**PowerShell:**
```powershell
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
```

**Git Bash:**
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
```

### Step 5: Build and Run

**PowerShell:**
```powershell
# Build
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker

# Run (in separate terminals)
.\api.exe
.\worker.exe
```

**Git Bash:**
```bash
# Build
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker

# Run (in separate terminals)
./api.exe
./worker.exe
```

## Troubleshooting

### "migrate: command not found"

**Solution 1: Install migrate**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Solution 2: Add to PATH**
```bash
# Find GOPATH
go env GOPATH

# Add to PATH (PowerShell)
$env:PATH = "$(go env GOPATH)\bin;$env:PATH"

# Or add permanently via Windows Settings > Environment Variables
```

**Solution 3: Use full path**
```bash
# Find migrate location
where migrate  # CMD
Get-Command migrate  # PowerShell

# Use full path
C:\Users\YourName\go\bin\migrate.exe -path migrations -database "..." up
```

### "make: command not found"

You don't need `make` on Windows. Use the direct commands or setup scripts instead.

**Alternative: Install make**
```bash
# Using Chocolatey
choco install make

# Using Scoop
scoop install make
```

### PostgreSQL Connection Issues

**Check if PostgreSQL is running:**
```bash
docker ps
```

**Start PostgreSQL:**
```bash
docker-compose up -d postgres
```

**Check connection:**
```bash
psql -U postgres -h localhost -d radar
# Password: postgres
```

### Port Already in Use

**Find process using port 8080:**
```powershell
netstat -ano | findstr :8080
```

**Kill process:**
```powershell
taskkill /PID <process_id> /F
```

## Environment Variables

### PowerShell
```powershell
# Set for current session
$env:LLM_ENABLED = "true"
$env:LLM_API_KEY = "your_key"
$env:LLM_MODEL = "google/gemini-flash-1.5"

# Run worker
.\worker.exe
```

### Git Bash
```bash
# Set for current session
export LLM_ENABLED=true
export LLM_API_KEY=your_key
export LLM_MODEL=google/gemini-flash-1.5

# Run worker
./worker.exe
```

### Permanent (Windows)
1. Open "Environment Variables" in Windows Settings
2. Add new variables:
   - `LLM_ENABLED` = `true`
   - `LLM_API_KEY` = `your_key`
   - `LLM_MODEL` = `google/gemini-flash-1.5`
3. Restart terminal

## Running the Application

### Method 1: Separate Terminals

**Terminal 1 (API):**
```bash
./api.exe
```

**Terminal 2 (Worker):**
```bash
./worker.exe
```

### Method 2: Background Processes (PowerShell)

```powershell
# Start API in background
Start-Process -FilePath ".\api.exe" -WindowStyle Hidden

# Start Worker in background
Start-Process -FilePath ".\worker.exe" -WindowStyle Hidden

# View running processes
Get-Process | Where-Object {$_.ProcessName -like "*api*" -or $_.ProcessName -like "*worker*"}

# Stop processes
Stop-Process -Name "api"
Stop-Process -Name "worker"
```

### Method 3: Windows Service (Advanced)

Use [NSSM](https://nssm.cc/) to run as Windows service:

```bash
# Install NSSM
choco install nssm

# Install API service
nssm install evolipia-api "D:\papaengineer\evolipia-radar\api.exe"

# Install Worker service
nssm install evolipia-worker "D:\papaengineer\evolipia-radar\worker.exe"

# Start services
nssm start evolipia-api
nssm start evolipia-worker
```

## Development Workflow

### Daily Development

```powershell
# 1. Start PostgreSQL (if not running)
docker-compose up -d postgres

# 2. Build latest changes
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker

# 3. Run in separate terminals
.\api.exe
.\worker.exe

# 4. Open browser
start http://localhost:8080
```

### Testing

```powershell
# Run tests
go test ./...

# Run specific test
go test ./internal/connectors -v

# Run with coverage
go test -cover ./...
```

### Database Management

```bash
# Connect to database
psql -U postgres -h localhost -d radar

# View tables
\dt

# View sources
SELECT name, type, enabled FROM sources;

# View recent items
SELECT title, domain, created_at FROM items ORDER BY created_at DESC LIMIT 10;

# Exit
\q
```

## IDE Setup

### VS Code
1. Install "Go" extension
2. Install "PostgreSQL" extension
3. Open workspace: `code .`

### GoLand
1. Open project directory
2. Configure Go SDK (Settings > Go > GOROOT)
3. Enable Go Modules (Settings > Go > Go Modules)

## Common Commands Reference

```powershell
# Build
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker

# Run
.\api.exe
.\worker.exe

# Test
go test ./...

# Format code
go fmt ./...

# Update dependencies
go mod tidy

# Database migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" down

# Docker
docker-compose up -d postgres
docker-compose down
docker-compose logs -f postgres

# Check API
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/feed?date=today
```

## Next Steps

1. ✅ Complete setup using automated script or manual steps
2. ✅ Verify API is running: http://localhost:8080
3. ✅ Check worker logs for source fetching
4. 📖 Read `docs/ENHANCEMENTS_QUICKSTART.md` for features
5. 🚀 Enable LLM summarization (optional)
6. 🎨 Try dark mode in Settings

## Support

- **Documentation:** `docs/` directory
- **Quick Reference:** `QUICK_REFERENCE.md`
- **Troubleshooting:** This file
- **Issues:** GitHub Issues
