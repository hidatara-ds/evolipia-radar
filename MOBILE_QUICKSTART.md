# 📱 Mobile Quick Start

Aplikasi ini sudah siap jalan di HP sebagai PWA (Progressive Web App).

## 🚀 Deploy dalam 5 Menit

### Windows:
```bash
deploy-mobile.bat
```

### Linux/Mac:
```bash
chmod +x deploy-mobile.sh
./deploy-mobile.sh
```

Pilih platform (Fly.io recommended), script akan handle semuanya.

## 📱 Install di HP

Setelah deploy:

1. **Buka URL aplikasi di browser HP**
   - Chrome/Edge: Menu (⋮) → "Add to Home Screen"
   - Safari: Share → "Add to Home Screen"

2. **Icon muncul di home screen** seperti app native!

3. **Buka dari icon** → fullscreen, no browser bar

## ✨ Fitur PWA

- ✅ Install seperti app native
- ✅ Offline support
- ✅ Push notifications (ready)
- ✅ Fast loading
- ✅ Auto-update

## 🔧 Platform Hosting (Gratis)

| Platform | Setup | Free Tier | Region |
|----------|-------|-----------|--------|
| **Fly.io** | CLI | 3 VMs, 3GB RAM | Global |
| **Render** | GitHub | 750 jam/bulan | Singapore |
| **Railway** | CLI/GitHub | $5 credit/bulan | Global |

## 📊 Monitoring

Setelah deploy, cek:
- Health: `https://your-app.com/healthz`
- Lighthouse PWA score di Chrome DevTools
- Install prompt muncul di HP

## 🐛 Troubleshooting

**"Add to Home Screen" tidak muncul:**
- Pastikan HTTPS aktif (auto di hosting)
- Buka di Chrome/Safari (bukan in-app browser)
- Cek manifest.json & service worker di DevTools

**API error di HP:**
- Update `API_BASE` di `web/index.html` dengan URL production
- Cek CORS headers kalau API beda domain

## 📚 Dokumentasi Lengkap

Lihat [docs/MOBILE_DEPLOYMENT.md](docs/MOBILE_DEPLOYMENT.md) untuk:
- Setup detail per platform
- Native app (React Native/Flutter)
- Capacitor wrapper
- Analytics & monitoring

## 🎯 Next Steps

1. Deploy backend → `deploy-mobile.bat`
2. Test di HP → Install PWA
3. (Optional) Publish ke Play Store/App Store → Pakai Capacitor

---

**Need help?** Check [docs/MOBILE_DEPLOYMENT.md](docs/MOBILE_DEPLOYMENT.md)
