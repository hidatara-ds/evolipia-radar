# 📱 Panduan Deploy Aplikasi ke HP

Aplikasi Evolipia Radar sudah mobile-ready dan bisa dijalankan di HP dengan beberapa cara.

## 🎯 Opsi 1: PWA (Progressive Web App) - RECOMMENDED

Cara paling mudah tanpa perlu publish ke app store. User bisa "install" langsung dari browser.

### Langkah Deploy:

#### A. Deploy Backend ke Server

Pilih salah satu platform hosting:

**1. Railway.app (Gratis + Mudah)**
```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Deploy
railway init
railway up
```

**2. Fly.io (Gratis tier bagus)**
```bash
# Install flyctl
# Windows: iwr https://fly.io/install.ps1 -useb | iex

# Login
fly auth login

# Deploy
fly launch
fly deploy
```

**3. Render.com (Gratis, auto-deploy dari GitHub)**
- Push code ke GitHub
- Buat akun di render.com
- New > Web Service
- Connect repository
- Build command: `go build -o api ./cmd/api`
- Start command: `./api`
- Add PostgreSQL database

#### B. Setup Database

Semua platform di atas punya PostgreSQL managed:
- Railway: Add PostgreSQL plugin
- Fly.io: `fly postgres create`
- Render: Add PostgreSQL database

Set environment variable `DATABASE_URL` ke connection string.

#### C. Jalankan Migrasi

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations (ganti URL dengan database kamu)
migrate -path migrations -database "postgres://user:pass@host:5432/dbname?sslmode=require" up
```

#### D. Test di HP

1. Buka browser di HP (Chrome/Safari)
2. Akses URL aplikasi (misal: `https://evolipia-radar.fly.dev`)
3. Klik menu browser (⋮) > "Add to Home Screen" / "Install App"
4. Icon muncul di home screen seperti app native!

### Fitur PWA yang Sudah Ada:

✅ Offline support (service worker)
✅ Install prompt
✅ App icon & splash screen
✅ Standalone mode (fullscreen tanpa browser bar)
✅ Mobile-first responsive design
✅ Touch gestures (pull to refresh)

---

## 🚀 Opsi 2: Native Mobile App (React Native / Flutter)

Kalau mau app native di Play Store / App Store:

### React Native Expo (Lebih Mudah)

```bash
# Install Expo CLI
npm install -g expo-cli

# Create project
npx create-expo-app evolipia-radar-mobile
cd evolipia-radar-mobile

# Install dependencies
npm install axios @react-navigation/native @react-navigation/bottom-tabs
```

Struktur minimal:
```
mobile/
├── App.js              # Entry point
├── screens/
│   ├── FeedScreen.js   # Feed panel
│   ├── RisingScreen.js # Rising panel
│   ├── SearchScreen.js # Search panel
│   └── ChatScreen.js   # AI Chat panel
├── components/
│   └── ItemCard.js     # Reusable item card
└── services/
    └── api.js          # API client
```

### Flutter (Performa Lebih Baik)

```bash
# Install Flutter: https://flutter.dev/docs/get-started/install

# Create project
flutter create evolipia_radar_mobile
cd evolipia_radar_mobile

# Add dependencies di pubspec.yaml
dependencies:
  http: ^1.1.0
  provider: ^6.1.1
```

---

## 📦 Opsi 3: Capacitor (Web → Native)

Convert web app yang sudah ada jadi native app:

```bash
# Install Capacitor
npm install @capacitor/core @capacitor/cli
npx cap init

# Add platforms
npx cap add android
npx cap add ios

# Copy web assets
npx cap copy

# Open in Android Studio / Xcode
npx cap open android
npx cap open ios
```

---

## 🎨 Opsi 4: WebView App (Paling Cepat)

Buat wrapper sederhana yang load URL aplikasi:

### Android (Kotlin)

```kotlin
// MainActivity.kt
class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        val webView = WebView(this)
        webView.settings.javaScriptEnabled = true
        webView.loadUrl("https://your-app-url.com")
        
        setContentView(webView)
    }
}
```

### iOS (Swift)

```swift
// ViewController.swift
import UIKit
import WebKit

class ViewController: UIViewController {
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let webView = WKWebView(frame: view.bounds)
        view.addSubview(webView)
        
        if let url = URL(string: "https://your-app-url.com") {
            webView.load(URLRequest(url: url))
        }
    }
}
```

---

## 🔧 Konfigurasi Tambahan untuk Mobile

### 1. HTTPS Wajib untuk PWA

PWA hanya jalan di HTTPS. Platform hosting di atas sudah auto-provide SSL.

### 2. CORS Headers

Kalau API dan frontend beda domain, tambahkan CORS di `cmd/api/main.go`:

```go
import "github.com/gin-contrib/cors"

// Setelah router := gin.Default()
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://your-frontend-domain.com"},
    AllowMethods:     []string{"GET", "POST", "PATCH"},
    AllowHeaders:     []string{"Origin", "Content-Type"},
    AllowCredentials: true,
}))
```

### 3. Environment Variables

Set di platform hosting:
```bash
DATABASE_URL=postgres://...
PORT=8080
WORKER_CRON=*/10 * * * *
```

### 4. Health Check Endpoint

Sudah ada di `/healthz` - platform hosting akan ping ini untuk monitoring.

---

## 📊 Monitoring & Analytics (Opsional)

### Google Analytics untuk PWA

Tambahkan di `web/index.html`:

```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=G-XXXXXXXXXX"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'G-XXXXXXXXXX');
</script>
```

### Sentry untuk Error Tracking

```bash
npm install @sentry/browser
```

```javascript
// Di web/index.html
import * as Sentry from "@sentry/browser";

Sentry.init({
  dsn: "https://your-sentry-dsn",
  environment: "production"
});
```

---

## 🎯 Rekomendasi Berdasarkan Kebutuhan

| Kebutuhan | Solusi | Effort | Cost |
|-----------|--------|--------|------|
| Cepat, gratis, no app store | **PWA** | ⭐ | Free |
| Performa native, offline penuh | React Native | ⭐⭐⭐ | Free (dev) |
| Best performance | Flutter | ⭐⭐⭐⭐ | Free (dev) |
| Publish ke app store | Capacitor | ⭐⭐ | $25-99/year |
| Simple wrapper | WebView | ⭐ | $25-99/year |

**Untuk MVP dan testing: Pakai PWA dulu!**

---

## 🚀 Quick Start: Deploy PWA Sekarang

```bash
# 1. Push ke GitHub
git add .
git commit -m "Ready for mobile deployment"
git push

# 2. Deploy ke Render (paling mudah)
# - Buka render.com
# - New > Web Service
# - Connect GitHub repo
# - Build: go build -o api ./cmd/api
# - Start: ./api
# - Add PostgreSQL database
# - Set DATABASE_URL env var

# 3. Test di HP
# Buka https://your-app.onrender.com
# Chrome: Menu > Add to Home Screen
# Safari: Share > Add to Home Screen

# Done! 🎉
```

---

## 📝 Checklist Sebelum Deploy

- [ ] Database migrations ready
- [ ] Environment variables configured
- [ ] HTTPS enabled (auto di hosting platform)
- [ ] Service worker registered
- [ ] Icons 192x192 dan 512x512 ready
- [ ] manifest.json configured
- [ ] Test di Chrome DevTools > Lighthouse > PWA audit
- [ ] Test install di HP (Android & iOS)

---

## 🐛 Troubleshooting

**PWA tidak muncul "Add to Home Screen":**
- Pastikan HTTPS aktif
- Cek service worker registered (DevTools > Application > Service Workers)
- Cek manifest.json valid (DevTools > Application > Manifest)
- Lighthouse audit harus pass PWA criteria

**API tidak bisa diakses dari HP:**
- Cek CORS headers
- Pastikan API URL di `web/index.html` sudah production URL
- Test API endpoint langsung di browser HP

**Database connection error:**
- Cek DATABASE_URL format: `postgres://user:pass@host:5432/dbname?sslmode=require`
- Pastikan database allow external connections
- Run migrations

---

## 📚 Resources

- [PWA Checklist](https://web.dev/pwa-checklist/)
- [Railway Docs](https://docs.railway.app/)
- [Fly.io Docs](https://fly.io/docs/)
- [Render Docs](https://render.com/docs)
- [React Native](https://reactnative.dev/)
- [Flutter](https://flutter.dev/)
- [Capacitor](https://capacitorjs.com/)
