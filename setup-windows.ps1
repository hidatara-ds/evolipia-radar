# Evolipia Radar - Windows Setup Script
# Run this in PowerShell

Write-Host "=== Evolipia Radar Setup ===" -ForegroundColor Cyan

# Check if Go is installed
Write-Host "`nChecking Go installation..." -ForegroundColor Yellow
$goVersion = go version 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go is installed: $goVersion" -ForegroundColor Green
} else {
    Write-Host "✗ Go is not installed. Please install Go 1.21+ from https://go.dev/dl/" -ForegroundColor Red
    exit 1
}

# Check if PostgreSQL is running
Write-Host "`nChecking PostgreSQL..." -ForegroundColor Yellow
$pgCheck = docker ps --filter "name=postgres" --format "{{.Names}}" 2>$null
if ($pgCheck -match "postgres") {
    Write-Host "✓ PostgreSQL container is running" -ForegroundColor Green
} else {
    Write-Host "Starting PostgreSQL container..." -ForegroundColor Yellow
    docker-compose up -d postgres
    Start-Sleep -Seconds 5
}

# Install migrate tool if not present
Write-Host "`nChecking migrate tool..." -ForegroundColor Yellow
$migrateCheck = Get-Command migrate -ErrorAction SilentlyContinue
if ($null -eq $migrateCheck) {
    Write-Host "Installing golang-migrate..." -ForegroundColor Yellow
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    # Add GOPATH/bin to PATH if not already there
    $goPath = go env GOPATH
    $goBin = Join-Path $goPath "bin"
    if ($env:PATH -notlike "*$goBin*") {
        $env:PATH = "$goBin;$env:PATH"
        Write-Host "Added $goBin to PATH for this session" -ForegroundColor Yellow
    }
} else {
    Write-Host "✓ migrate tool is installed" -ForegroundColor Green
}

# Run migrations
Write-Host "`nRunning database migrations..." -ForegroundColor Yellow
$dbUrl = "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"
migrate -path migrations -database $dbUrl up

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Migrations completed successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Migration failed. Check the error above." -ForegroundColor Red
    exit 1
}

# Build binaries
Write-Host "`nBuilding binaries..." -ForegroundColor Yellow
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Binaries built successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== Setup Complete! ===" -ForegroundColor Cyan
Write-Host "`nTo start the application:" -ForegroundColor Yellow
Write-Host "  1. Start API:    .\api.exe" -ForegroundColor White
Write-Host "  2. Start Worker: .\worker.exe" -ForegroundColor White
Write-Host "  3. Open browser: http://localhost:8080" -ForegroundColor White
Write-Host "`nOptional - Enable LLM Summarization:" -ForegroundColor Yellow
Write-Host '  $env:LLM_ENABLED="true"' -ForegroundColor White
Write-Host '  $env:LLM_API_KEY="your_openrouter_key"' -ForegroundColor White
Write-Host '  $env:LLM_MODEL="google/gemini-flash-1.5"' -ForegroundColor White
