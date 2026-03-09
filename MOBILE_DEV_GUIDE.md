# 📱 Mobile Development Guide

Quick reference untuk developer yang mau develop/maintain mobile PWA.

## 🚀 Quick Commands

```bash
# Test PWA locally
make test-pwa
# atau
./test-pwa.sh      # Linux/Mac
test-pwa.bat       # Windows

# Deploy to production
make deploy-mobile
# atau
./deploy-mobile.sh # Linux/Mac
deploy-mobile.bat  # Windows

# Build production binary
make build-mobile

# Build Docker image
make docker-build-mobile

# Run Lighthouse audit
make lighthouse
```

## 📁 File Structure

```
web/
├── index.html          # Main UI (mobile-first)
├── style.css           # Responsive styles
├── sw.js               # Service worker (offline support)
└── manifest.json       # PWA manifest (install config)

assets/
├── icon.png            # 192x192 app icon
├── icon1.png           # 512x512 app icon
└── maskot1.png         # Chat mascot

Deployment:
├── fly.toml            # Fly.io config
├── render.yaml         # Render.com config
├── railway.json        # Railway config
├── Dockerfile.api      # Production Docker build
└── .dockerignore       # Docker build optimization
```

## 🔧 Development Workflow

### 1. Local Development

```bash
# Start database
docker-compose up -d postgres

# Run migrations
make migrate-up

# Start API server
make run-api

# Open browser
open http://localhost:8080
```

### 2. Test PWA Features

```bash
# Run PWA test script
make test-pwa

# Check in Chrome DevTools:
# - Application > Manifest
# - Application > Service Workers
# - Lighthouse > PWA audit
```

### 3. Test on Mobile (Same Network)

```bash
# Get your local IP
ip addr show | grep 'inet '  # Linux
ipconfig                      # Windows

# Open on phone
http://YOUR_IP:8080

# Try "Add to Home Screen"
```

### 4. Deploy to Production

```bash
# Option 1: Interactive script
make deploy-mobile

# Option 2: Direct platform deploy
flyctl deploy              # Fly.io
git push                   # Render (auto-deploy)
railway up                 # Railway
```

## 🎨 UI Development

### Component Structure

```javascript
// Main panels
- Feed Panel      (#panel-feed)
- Rising Panel    (#panel-rising)
- Search Panel    (#panel-search)
- Chat Panel      (#panel-chat)
- Settings Panel  (#panel-settings)

// Navigation
- Bottom nav bar with 5 tabs
- Active state management
- Panel switching logic
```

### Styling Guidelines

```css
/* Mobile-first approach */
- Base styles for mobile (320px+)
- Tablet breakpoint: 768px
- Desktop breakpoint: 1024px

/* Key classes */
.panel          - Main content panels
.nav-btn        - Bottom navigation buttons
.item-card      - News item cards
.chat-msg       - Chat messages
.modal-overlay  - Modal dialogs
```

### Adding New Features

1. **Add UI in index.html**
   ```html
   <section id="panel-new" class="panel">
     <!-- Your content -->
   </section>
   ```

2. **Add navigation button**
   ```html
   <button class="nav-btn" data-panel="new">
     <span class="icon">🆕</span>
     <span>New</span>
   </button>
   ```

3. **Add panel logic**
   ```javascript
   function onShowPanel(panelId) {
     if (panelId === 'new') loadNewFeature();
   }
   ```

4. **Add API endpoint** (if needed)
   ```go
   // internal/http/handlers/handlers.go
   v1.GET("/new-endpoint", h.NewHandler)
   ```

## 🔌 API Integration

### Calling Backend API

```javascript
// API base URL (auto-detect)
const API_BASE = window.location.hostname === 'localhost' 
  ? '' 
  : 'https://your-api.com';

// Make request
async function loadData() {
  const response = await fetch(apiUrl('/v1/endpoint'));
  const data = await response.json();
  return data;
}
```

### Available Endpoints

```
GET  /healthz                    - Health check
GET  /v1/feed                    - Daily feed
GET  /v1/rising                  - Rising trends
GET  /v1/items/:id               - Item details
GET  /v1/search?q=query          - Search
GET  /v1/sources                 - List sources
POST /v1/sources                 - Create source
POST /v1/sources/test            - Test source
PATCH /v1/sources/:id/enable     - Enable/disable
```

## 🎯 PWA Features

### Service Worker (sw.js)

```javascript
// Cache strategy
- Static assets: Cache-first
- API calls: Network-first with cache fallback
- Images: Cache-first with network fallback

// Update strategy
- Install: Cache new assets
- Activate: Clean old caches
- Fetch: Serve from cache or network
```

### Manifest (manifest.json)

```json
{
  "name": "App Name",
  "short_name": "Short",
  "start_url": "/",
  "display": "standalone",    // Fullscreen
  "theme_color": "#0f172a",
  "background_color": "#0f172a",
  "icons": [...]
}
```

### Install Prompt

```javascript
// Auto-trigger install prompt
let deferredPrompt;

window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault();
  deferredPrompt = e;
  // Show custom install button
});

// Trigger install
deferredPrompt.prompt();
const { outcome } = await deferredPrompt.userChoice;
```

## 🧪 Testing Checklist

### Pre-Deploy Tests

- [ ] All files exist (manifest, sw, icons)
- [ ] Manifest.json valid JSON
- [ ] Service worker registers
- [ ] Icons correct size (192x192, 512x512)
- [ ] HTTPS enabled (production)
- [ ] API endpoints work
- [ ] Offline mode works

### Lighthouse Audit

```bash
make lighthouse
# Target: 100/100 PWA score
```

### Manual Testing

- [ ] Install prompt appears
- [ ] App installs to home screen
- [ ] Standalone mode works (no browser bar)
- [ ] All features functional
- [ ] Offline mode works
- [ ] Performance acceptable

## 🐛 Common Issues

### Service Worker Not Updating

```javascript
// Force update
navigator.serviceWorker.getRegistrations()
  .then(regs => regs.forEach(reg => reg.unregister()));

// Hard reload
Ctrl+Shift+R (Chrome)
```

### Install Prompt Not Showing

- Check HTTPS enabled
- Check manifest.json valid
- Check service worker registered
- Open in Chrome (not in-app browser)

### Offline Mode Not Working

- Check service worker cache strategy
- Check files cached in install event
- Test in DevTools offline mode

### CORS Errors

```go
// Add CORS middleware in cmd/api/main.go
import "github.com/gin-contrib/cors"

router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://your-domain.com"},
    AllowMethods: []string{"GET", "POST", "PATCH"},
}))
```

## 📊 Performance Optimization

### Bundle Size

```bash
# Check sizes
du -sh web/*
du -sh assets/*

# Optimize images
# Use tools like imagemin, squoosh
```

### Caching Strategy

```javascript
// Cache static assets aggressively
const CACHE_ASSETS = [
  '/',
  '/web/index.html',
  '/web/style.css',
  '/assets/icon.png'
];

// Cache API responses with TTL
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes
```

### Lazy Loading

```javascript
// Load images on demand
<img loading="lazy" src="..." />

// Load panels on demand
function onShowPanel(panelId) {
  if (!panelLoaded[panelId]) {
    loadPanelData(panelId);
    panelLoaded[panelId] = true;
  }
}
```

## 🔐 Security

### Content Security Policy

```html
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; 
               script-src 'self' 'unsafe-inline' cdn.tailwindcss.com;
               style-src 'self' 'unsafe-inline';
               img-src 'self' data: https:;">
```

### API Key Storage

```javascript
// Store in localStorage (client-side only)
localStorage.setItem('api_key', key);

// For production: Use backend proxy
// Don't expose API keys in frontend
```

## 📚 Resources

- [PWA Checklist](https://web.dev/pwa-checklist/)
- [Service Worker API](https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API)
- [Web App Manifest](https://developer.mozilla.org/en-US/docs/Web/Manifest)
- [Lighthouse CI](https://github.com/GoogleChrome/lighthouse-ci)
- [Workbox](https://developers.google.com/web/tools/workbox) - Advanced SW library

## 🎓 Learning Path

1. **Basics**: HTML/CSS/JS fundamentals
2. **PWA**: Service workers, manifest, caching
3. **API**: REST API integration, error handling
4. **Testing**: Lighthouse, manual testing, debugging
5. **Deploy**: Platform-specific deployment
6. **Monitor**: Analytics, error tracking, performance

## 💡 Tips

- Test on real devices early and often
- Keep bundle size minimal
- Use browser DevTools extensively
- Monitor Lighthouse score regularly
- Follow mobile-first design principles
- Test on slow 3G network
- Handle offline gracefully
- Provide loading states
- Use touch-friendly UI (44px+ tap targets)

---

Happy coding! 🚀
