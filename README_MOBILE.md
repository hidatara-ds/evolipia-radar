# 📱 Evolipia Radar - Mobile Edition

> Aplikasi AI/ML tech news aggregator yang bisa jalan di HP!

## 🎯 Apa Ini?

Evolipia Radar sekarang bisa di-install di HP seperti aplikasi native, tanpa perlu download dari Play Store atau App Store. Pakai teknologi Progressive Web App (PWA).

## ✨ Fitur

- 📰 **Feed** - Berita AI/ML terpopuler hari ini
- 🔥 **Rising** - Tren yang sedang naik
- 🔍 **Search** - Cari berita/topik tertentu
- 💬 **AI Chat** - Chat dengan AI tentang berita
- ⚙️ **Settings** - Konfigurasi app
- 🌙 **Dark Mode** - Nyaman di mata
- 📴 **Offline** - Tetap bisa buka meski internet mati

## 🚀 Quick Start

### Untuk Developer:

```bash
# 1. Test PWA locally
test-pwa.bat          # Windows
./test-pwa.sh         # Linux/Mac

# 2. Deploy ke production
deploy-mobile.bat     # Windows
./deploy-mobile.sh    # Linux/Mac

# 3. Pilih platform (Fly.io/Render/Railway)
# Script akan guide kamu step-by-step
```

### Untuk User:

1. Buka URL aplikasi di browser HP
2. **Android (Chrome)**: Menu ⋮ → "Install app"
3. **iOS (Safari)**: Share ⬆️ → "Add to Home Screen"
4. Done! Icon muncul di home screen

## 📚 Dokumentasi

| Dokumen | Untuk Siapa | Isi |
|---------|-------------|-----|
| [MOBILE_QUICKSTART.md](MOBILE_QUICKSTART.md) | Developer | Deploy dalam 5 menit |
| [INSTALL_APP_HP.md](INSTALL_APP_HP.md) | User | Cara install PWA |
| [docs/MOBILE_DEPLOYMENT.md](docs/MOBILE_DEPLOYMENT.md) | Developer | Panduan lengkap deploy |
| [MOBILE_DEV_GUIDE.md](MOBILE_DEV_GUIDE.md) | Developer | Development guide |
| [docs/PWA_TESTING_CHECKLIST.md](docs/PWA_TESTING_CHECKLIST.md) | QA/Developer | Testing checklist |

## 🎯 Platform Hosting (Gratis!)

| Platform | Setup | Free Tier | Best For |
|----------|-------|-----------|----------|
| **Fly.io** | CLI | 3 VMs, 3GB RAM | Global, fast |
| **Render** | GitHub | 750 jam/bulan | Easy, auto-deploy |
| **Railway** | CLI/GitHub | $5 credit/bulan | Simple, good DX |

## 🔧 Tech Stack

- **Backend**: Go 1.21+ (Gin framework)
- **Database**: PostgreSQL 15+
- **Frontend**: Vanilla JS + Tailwind CSS
- **PWA**: Service Worker + Web Manifest
- **Deployment**: Docker + Platform hosting

## 📱 Screenshots

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   📰 Feed       │  │   🔥 Rising     │  │   🔍 Search     │
│                 │  │                 │  │                 │
│ • Top stories   │  │ • Trending now  │  │ • Find topics   │
│ • Daily digest  │  │ • Hot topics    │  │ • Filter news   │
│ • Scored items  │  │ • Momentum      │  │ • Quick access  │
└─────────────────┘  └─────────────────┘  └─────────────────┘

┌─────────────────┐  ┌─────────────────┐
│   💬 AI Chat    │  │   ⚙️ Settings   │
│                 │  │                 │
│ • Ask anything  │  │ • API config    │
│ • Get insights  │  │ • Dark mode     │
│ • Smart replies │  │ • Preferences   │
└─────────────────┘  └─────────────────┘
```

## 🎨 UI Features

- ✅ Mobile-first responsive design
- ✅ Touch-friendly (44px+ tap targets)
- ✅ Pull-to-refresh gesture
- ✅ Smooth animations
- ✅ Dark mode support
- ✅ Offline indicator
- ✅ Loading states
- ✅ Error handling

## 🔐 Security

- ✅ HTTPS enforced
- ✅ SSRF protection
- ✅ Rate limiting
- ✅ Input validation
- ✅ Secure headers
- ✅ Client-side API key storage

## 📊 Performance

Target metrics:
- First Contentful Paint: < 1.5s
- Time to Interactive: < 3.5s
- Lighthouse PWA Score: 100/100
- Bundle Size: < 500KB

## 🧪 Testing

```bash
# Test PWA configuration
make test-pwa

# Run Lighthouse audit
make lighthouse

# Build production binary
make build-mobile

# Build Docker image
make docker-build-mobile
```

## 🐛 Troubleshooting

### "Add to Home Screen" tidak muncul

**Android:**
- Buka di Chrome (bukan browser lain)
- Pastikan HTTPS aktif
- Refresh halaman

**iOS:**
- Buka di Safari (bukan Chrome iOS)
- Pastikan HTTPS aktif
- Coba tab baru

### App tidak fullscreen

- Uninstall dari home screen
- Install ulang dari browser
- Buka dari icon (bukan bookmark)

### Offline mode tidak jalan

- Clear cache app
- Reinstall app
- Cek service worker di DevTools

## 📞 Support

- 📖 [Full Documentation](docs/MOBILE_DEPLOYMENT.md)
- 🧪 [Testing Guide](docs/PWA_TESTING_CHECKLIST.md)
- 💻 [Dev Guide](MOBILE_DEV_GUIDE.md)
- 👤 [User Guide](INSTALL_APP_HP.md)

## 🎉 Benefits

✅ **Zero cost** - Free tier hosting  
✅ **Fast deployment** - 5-10 minutes  
✅ **No app store** - Direct install  
✅ **Auto-update** - Always latest version  
✅ **Cross-platform** - Android & iOS  
✅ **Small size** - ~2MB only  
✅ **Easy maintenance** - Single codebase  

## 🚧 Roadmap

### Phase 1 (Current) ✅
- [x] PWA implementation
- [x] Mobile-first UI
- [x] Offline support
- [x] Deployment scripts
- [x] Documentation

### Phase 2 (Optional)
- [ ] Push notifications
- [ ] Background sync
- [ ] Share target API
- [ ] Custom install prompt

### Phase 3 (Optional)
- [ ] Native wrapper (Capacitor)
- [ ] Publish to app stores
- [ ] Deep linking
- [ ] Biometric auth

## 💡 Tips

- Test on real devices early
- Keep bundle size minimal
- Monitor Lighthouse score
- Test on slow 3G network
- Handle offline gracefully
- Use touch-friendly UI

## 📈 Metrics

Track these after deployment:
- Install rate
- Lighthouse score
- Load time
- Offline usage
- User engagement
- Error rate

## 🤝 Contributing

1. Fork the repo
2. Create feature branch
3. Test on mobile
4. Submit PR

## 📄 License

MIT License - see [LICENSE.md](LICENSE.md)

## 🙏 Credits

- Built with Go, Gin, PostgreSQL
- UI with Tailwind CSS
- PWA with Service Workers
- Deployed on Fly.io/Render/Railway

---

## 🚀 Get Started Now!

```bash
# Clone repo
git clone https://github.com/your-username/evolipia-radar.git
cd evolipia-radar

# Deploy to mobile
deploy-mobile.bat    # Windows
./deploy-mobile.sh   # Linux/Mac

# Follow the prompts
# Done in 5 minutes! 🎉
```

---

**Made with ❤️ for the AI/ML community**

📱 Install now and stay updated with the latest AI/ML trends!
