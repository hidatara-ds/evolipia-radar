@echo off
REM Evolipia Radar - Windows Setup Script (CMD)
REM Run this in Command Prompt

echo === Evolipia Radar Setup ===
echo.

REM Check if Go is installed
echo Checking Go installation...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed. Please install Go 1.21+ from https://go.dev/dl/
    exit /b 1
)
echo [OK] Go is installed

REM Check if PostgreSQL is running
echo.
echo Checking PostgreSQL...
docker ps --filter "name=postgres" --format "{{.Names}}" | findstr postgres >nul 2>&1
if %errorlevel% neq 0 (
    echo Starting PostgreSQL container...
    docker-compose up -d postgres
    timeout /t 5 /nobreak >nul
)
echo [OK] PostgreSQL is running

REM Install migrate tool
echo.
echo Installing migrate tool...
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest

REM Add GOPATH\bin to PATH
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
set PATH=%GOPATH%\bin;%PATH%

REM Run migrations
echo.
echo Running database migrations...
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
if %errorlevel% neq 0 (
    echo [ERROR] Migration failed
    exit /b 1
)
echo [OK] Migrations completed

REM Build binaries
echo.
echo Building binaries...
go build -o api.exe ./cmd/api
go build -o worker.exe ./cmd/worker
if %errorlevel% neq 0 (
    echo [ERROR] Build failed
    exit /b 1
)
echo [OK] Binaries built

echo.
echo === Setup Complete! ===
echo.
echo To start the application:
echo   1. Start API:    api.exe
echo   2. Start Worker: worker.exe
echo   3. Open browser: http://localhost:8080
echo.
echo Optional - Enable LLM Summarization:
echo   set LLM_ENABLED=true
echo   set LLM_API_KEY=your_openrouter_key
echo   set LLM_MODEL=google/gemini-flash-1.5
echo.
pause
