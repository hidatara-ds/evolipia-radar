@echo off
REM Test PWA locally before deployment (Windows)

echo.
echo 🧪 Testing PWA Configuration
echo ==============================
echo.

REM Check if required files exist
echo 📁 Checking required files...
set missing=0

if exist "web\index.html" (echo   ✅ web\index.html) else (echo   ❌ web\index.html - MISSING! & set /a missing+=1)
if exist "web\manifest.json" (echo   ✅ web\manifest.json) else (echo   ❌ web\manifest.json - MISSING! & set /a missing+=1)
if exist "web\sw.js" (echo   ✅ web\sw.js) else (echo   ❌ web\sw.js - MISSING! & set /a missing+=1)
if exist "web\style.css" (echo   ✅ web\style.css) else (echo   ❌ web\style.css - MISSING! & set /a missing+=1)
if exist "assets\icon.png" (echo   ✅ assets\icon.png) else (echo   ❌ assets\icon.png - MISSING! & set /a missing+=1)
if exist "assets\icon1.png" (echo   ✅ assets\icon1.png) else (echo   ❌ assets\icon1.png - MISSING! & set /a missing+=1)
if exist "assets\maskot1.png" (echo   ✅ assets\maskot1.png) else (echo   ❌ assets\maskot1.png - MISSING! & set /a missing+=1)

if %missing% gtr 0 (
    echo.
    echo ❌ %missing% file(s) missing!
    pause
    exit /b 1
)

echo.
echo 📋 Validating manifest.json...
type web\manifest.json >nul 2>&1
if %errorlevel% equ 0 (
    echo   ✅ File readable
) else (
    echo   ❌ Cannot read manifest.json!
    pause
    exit /b 1
)

echo.
echo 🔍 Checking service worker registration...
findstr /C:"serviceWorker.register" web\index.html >nul
if %errorlevel% equ 0 (
    echo   ✅ Service worker registration found
) else (
    echo   ❌ Service worker registration not found!
    pause
    exit /b 1
)

echo.
echo 🌐 Starting local server...
echo.

REM Check if server is already running
netstat -ano | findstr ":8080" | findstr "LISTENING" >nul
if %errorlevel% equ 0 (
    echo ⚠️  Port 8080 already in use
    echo.
    echo 📱 Open in browser:
    echo    http://localhost:8080
    echo.
    echo 🧪 Test PWA:
    echo    1. Open Chrome DevTools (F12^)
    echo    2. Go to Lighthouse tab
    echo    3. Select 'Progressive Web App'
    echo    4. Click 'Analyze page load'
    echo    5. Target score: 100/100
    echo.
    echo 📱 Test on mobile:
    echo    1. Get your local IP: ipconfig
    echo    2. Open http://YOUR_IP:8080 on phone
    echo    3. Try 'Add to Home Screen'
    echo.
    pause
    exit /b 0
)

REM Start server
if exist "api.exe" (
    echo Starting API server...
    start /B api.exe
) else if exist "api" (
    echo Starting API server...
    start /B api
) else (
    echo Building API server...
    go build -o api.exe ./cmd/api
    if %errorlevel% neq 0 (
        echo ❌ Build failed!
        pause
        exit /b 1
    )
    start /B api.exe
)

REM Wait for server to start
timeout /t 3 /nobreak >nul

echo.
echo ✅ Server running on http://localhost:8080
echo.
echo 📱 Test PWA:
echo    1. Open http://localhost:8080 in Chrome
echo    2. Open DevTools (F12^) ^> Lighthouse
echo    3. Run PWA audit (target: 100/100^)
echo.
echo 📱 Test on mobile (same network^):
for /f "tokens=2 delims=:" %%a in ('ipconfig ^| findstr /C:"IPv4"') do set LOCAL_IP=%%a
set LOCAL_IP=%LOCAL_IP:~1%
echo    1. Open http://%LOCAL_IP%:8080 on phone
echo    2. Try 'Add to Home Screen'
echo.
echo Press any key to stop server...
pause >nul

REM Stop server
taskkill /F /IM api.exe >nul 2>&1
taskkill /F /IM api >nul 2>&1

echo.
echo 👋 Server stopped
