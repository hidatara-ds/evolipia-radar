#!/bin/bash
# Test PWA locally before deployment

set -e

echo "🧪 Testing PWA Configuration"
echo "=============================="
echo ""

# Check if required files exist
echo "📁 Checking required files..."
files=(
    "web/index.html"
    "web/manifest.json"
    "web/sw.js"
    "web/style.css"
    "assets/icon.png"
    "assets/icon1.png"
    "assets/maskot1.png"
)

missing=0
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file - MISSING!"
        missing=$((missing + 1))
    fi
done

if [ $missing -gt 0 ]; then
    echo ""
    echo "❌ $missing file(s) missing!"
    exit 1
fi

echo ""
echo "📋 Validating manifest.json..."
if command -v jq &> /dev/null; then
    if jq empty web/manifest.json 2>/dev/null; then
        echo "  ✅ Valid JSON"
        
        # Check required fields
        name=$(jq -r '.name' web/manifest.json)
        start_url=$(jq -r '.start_url' web/manifest.json)
        display=$(jq -r '.display' web/manifest.json)
        
        echo "  📱 Name: $name"
        echo "  🔗 Start URL: $start_url"
        echo "  📺 Display: $display"
        
        if [ "$display" != "standalone" ]; then
            echo "  ⚠️  Warning: display should be 'standalone' for fullscreen"
        fi
    else
        echo "  ❌ Invalid JSON!"
        exit 1
    fi
else
    echo "  ⚠️  jq not installed, skipping JSON validation"
fi

echo ""
echo "🖼️  Checking icon sizes..."
if command -v identify &> /dev/null; then
    icon_size=$(identify -format "%wx%h" assets/icon.png 2>/dev/null || echo "unknown")
    icon1_size=$(identify -format "%wx%h" assets/icon1.png 2>/dev/null || echo "unknown")
    
    echo "  📐 icon.png: $icon_size (should be 192x192)"
    echo "  📐 icon1.png: $icon1_size (should be 512x512)"
else
    echo "  ⚠️  ImageMagick not installed, skipping size check"
fi

echo ""
echo "🔍 Checking service worker registration..."
if grep -q "serviceWorker.register" web/index.html; then
    echo "  ✅ Service worker registration found"
else
    echo "  ❌ Service worker registration not found!"
    exit 1
fi

echo ""
echo "🌐 Starting local server..."
echo ""

# Check if server is already running
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  Port 8080 already in use"
    echo ""
    echo "📱 Open in browser:"
    echo "   http://localhost:8080"
    echo ""
    echo "🧪 Test PWA:"
    echo "   1. Open Chrome DevTools (F12)"
    echo "   2. Go to Lighthouse tab"
    echo "   3. Select 'Progressive Web App'"
    echo "   4. Click 'Analyze page load'"
    echo "   5. Target score: 100/100"
    echo ""
    echo "📱 Test on mobile:"
    echo "   1. Get your local IP: ip addr show | grep 'inet '"
    echo "   2. Open http://YOUR_IP:8080 on phone"
    echo "   3. Try 'Add to Home Screen'"
    echo ""
else
    echo "Starting API server..."
    echo ""
    
    # Check if binary exists
    if [ -f "./api" ]; then
        ./api &
        SERVER_PID=$!
    elif [ -f "./api.exe" ]; then
        ./api.exe &
        SERVER_PID=$!
    else
        echo "Building API server..."
        go build -o api ./cmd/api
        ./api &
        SERVER_PID=$!
    fi
    
    # Wait for server to start
    sleep 3
    
    echo ""
    echo "✅ Server running on http://localhost:8080"
    echo ""
    echo "📱 Test PWA:"
    echo "   1. Open http://localhost:8080 in Chrome"
    echo "   2. Open DevTools (F12) > Lighthouse"
    echo "   3. Run PWA audit (target: 100/100)"
    echo ""
    echo "📱 Test on mobile (same network):"
    LOCAL_IP=$(hostname -I | awk '{print $1}' 2>/dev/null || echo "YOUR_IP")
    echo "   1. Open http://$LOCAL_IP:8080 on phone"
    echo "   2. Try 'Add to Home Screen'"
    echo ""
    echo "Press Ctrl+C to stop server"
    echo ""
    
    # Wait for Ctrl+C
    trap "kill $SERVER_PID 2>/dev/null; echo ''; echo '👋 Server stopped'; exit 0" INT
    wait $SERVER_PID
fi
