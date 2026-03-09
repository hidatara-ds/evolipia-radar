@echo off
REM Quick deployment script untuk mobile PWA (Windows)

echo.
echo 🚀 Evolipia Radar - Mobile Deployment
echo ======================================
echo.

REM Check git status
git status --short > nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Git repository tidak ditemukan
    pause
    exit /b 1
)

echo Pilih platform deployment:
echo 1) Fly.io (Recommended - Fast ^& Free)
echo 2) Render.com (Easy - Auto deploy from GitHub)
echo 3) Railway.app (Simple - Good free tier)
echo.
set /p choice="Pilihan (1-3): "

if "%choice%"=="1" goto flyio
if "%choice%"=="2" goto render
if "%choice%"=="3" goto railway
goto invalid

:flyio
echo.
echo 📦 Deploying to Fly.io...
echo.

REM Check if flyctl installed
where flyctl >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ flyctl belum terinstall
    echo Install: https://fly.io/docs/hands-on/install-flyctl/
    echo Windows: iwr https://fly.io/install.ps1 -useb ^| iex
    pause
    exit /b 1
)

REM Check if logged in
flyctl auth whoami >nul 2>nul
if %errorlevel% neq 0 (
    echo 🔐 Login ke Fly.io...
    flyctl auth login
)

REM Check if app exists
flyctl status >nul 2>nul
if %errorlevel% neq 0 (
    echo 🆕 Membuat app baru...
    flyctl launch --no-deploy
)

REM Deploy
echo 🚀 Deploying...
flyctl deploy

echo.
echo ✅ Deployment selesai!
flyctl open
goto end

:render
echo.
echo 📦 Setup Render.com deployment...
echo.
echo Langkah manual:
echo 1. Push code ke GitHub: git push
echo 2. Buka https://render.com
echo 3. New ^> Blueprint
echo 4. Connect repository ini
echo 5. Render akan auto-detect render.yaml
echo 6. Klik 'Apply'
echo.
set /p push="Push ke GitHub sekarang? (y/n): "
if /i "%push%"=="y" (
    git push
    echo ✅ Pushed! Sekarang setup di render.com
)
goto end

:railway
echo.
echo 📦 Deploying to Railway...
echo.

REM Check if railway installed
where railway >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing Railway CLI...
    npm install -g @railway/cli
)

REM Login
railway whoami >nul 2>nul
if %errorlevel% neq 0 (
    echo 🔐 Login ke Railway...
    railway login
)

REM Init if needed
if not exist "railway.json" (
    echo 🆕 Membuat project baru...
    railway init
)

REM Deploy
echo 🚀 Deploying...
railway up

echo.
echo ✅ Deployment selesai!
railway open
goto end

:invalid
echo ❌ Pilihan tidak valid
pause
exit /b 1

:end
echo.
echo 🎉 Deployment selesai!
echo.
echo 📱 Test di HP:
echo 1. Buka URL aplikasi di browser HP
echo 2. Chrome: Menu (⋮) ^> Add to Home Screen
echo 3. Safari: Share ^> Add to Home Screen
echo.
echo 🔧 Jangan lupa:
echo - Set environment variables di dashboard platform
echo - Update API_BASE di web/index.html kalau API beda domain
echo - Test PWA dengan Lighthouse di Chrome DevTools
echo.
pause
