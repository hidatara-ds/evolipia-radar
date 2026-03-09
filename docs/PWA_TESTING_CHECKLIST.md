# 📱 PWA Testing Checklist

Gunakan checklist ini untuk memastikan PWA berjalan sempurna di HP.

## ✅ Pre-Deployment Checks

### 1. Manifest.json
- [ ] File `web/manifest.json` ada dan valid
- [ ] `name` dan `short_name` sudah diisi
- [ ] `start_url` benar (biasanya `/`)
- [ ] `display: "standalone"` untuk fullscreen
- [ ] Icons 192x192 dan 512x512 ada di `/assets/`
- [ ] `theme_color` dan `background_color` sesuai brand

### 2. Service Worker
- [ ] File `web/sw.js` ada
- [ ] Service worker registered di `web/index.html`
- [ ] Cache strategy sudah benar (cache-first untuk assets, network-first untuk API)

### 3. HTTPS
- [ ] Aplikasi serve via HTTPS (wajib untuk PWA)
- [ ] SSL certificate valid
- [ ] No mixed content warnings

### 4. Icons & Assets
- [ ] Icon 192x192: `/assets/icon.png`
- [ ] Icon 512x512: `/assets/icon1.png`
- [ ] Maskot: `/assets/maskot1.png`
- [ ] Favicon di `<head>`
- [ ] Apple touch icon untuk iOS

### 5. Meta Tags
- [ ] `<meta name="viewport">` ada
- [ ] `<meta name="theme-color">` sesuai manifest
- [ ] `<meta name="description">` untuk SEO
- [ ] `<link rel="manifest">` pointing ke manifest.json

---

## 🧪 Testing di Desktop (Chrome DevTools)

### Lighthouse Audit
1. Buka Chrome DevTools (F12)
2. Tab "Lighthouse"
3. Select "Progressive Web App"
4. Click "Analyze page load"

**Target Score: 100/100**

Checklist Lighthouse:
- [ ] ✅ Installable
- [ ] ✅ PWA optimized
- [ ] ✅ Fast and reliable
- [ ] ✅ Works offline
- [ ] ✅ Configured for custom splash screen

### Application Tab
1. DevTools > Application

**Manifest:**
- [ ] Manifest loads without errors
- [ ] Icons preview correctly
- [ ] All fields populated

**Service Workers:**
- [ ] Service worker registered
- [ ] Status: "activated and running"
- [ ] Update on reload works
- [ ] Offline mode works (toggle offline in DevTools)

**Storage:**
- [ ] Cache Storage shows cached files
- [ ] LocalStorage works (settings, chat history)
- [ ] No quota errors

---

## 📱 Testing di HP - Android

### Chrome (Android)
1. **Buka URL di Chrome**
   - [ ] Page loads correctly
   - [ ] No console errors (inspect via chrome://inspect)

2. **Install Prompt**
   - [ ] "Add to Home Screen" banner muncul otomatis ATAU
   - [ ] Menu (⋮) > "Install app" / "Add to Home Screen" available

3. **Install & Launch**
   - [ ] Klik install
   - [ ] Icon muncul di home screen dengan nama & icon yang benar
   - [ ] Splash screen muncul saat launch (background_color + icon)
   - [ ] App buka dalam standalone mode (no browser bar)

4. **Functionality**
   - [ ] Navigation works (bottom nav)
   - [ ] Feed loads
   - [ ] Search works
   - [ ] Item detail opens
   - [ ] Chat AI works (dengan API key)
   - [ ] Settings saves (dark mode, API key)

5. **Offline Mode**
   - [ ] Matikan internet
   - [ ] Buka app dari home screen
   - [ ] Cached pages still load
   - [ ] Offline indicator muncul (optional)

6. **Performance**
   - [ ] Smooth scrolling
   - [ ] Fast page transitions
   - [ ] No lag on interactions
   - [ ] Pull-to-refresh works

### Samsung Internet (Android)
- [ ] Repeat steps 1-6 di Samsung Internet
- [ ] Install prompt works
- [ ] Standalone mode works

---

## 📱 Testing di HP - iOS (Safari)

### Safari (iOS)
1. **Buka URL di Safari**
   - [ ] Page loads correctly
   - [ ] No errors in console (Settings > Safari > Advanced > Web Inspector)

2. **Install Process**
   - [ ] Tap Share button (⬆️)
   - [ ] "Add to Home Screen" option available
   - [ ] Tap "Add to Home Screen"
   - [ ] Edit name if needed
   - [ ] Tap "Add"

3. **Launch**
   - [ ] Icon muncul di home screen
   - [ ] Tap icon
   - [ ] Splash screen muncul (optional di iOS)
   - [ ] App buka fullscreen (no Safari UI)

4. **Functionality**
   - [ ] All features work (same as Android checklist)
   - [ ] Touch gestures work
   - [ ] Keyboard doesn't cover input fields

5. **iOS-Specific**
   - [ ] Status bar color matches theme
   - [ ] Safe area insets respected (notch/home indicator)
   - [ ] No bounce scroll issues

---

## 🔧 Common Issues & Fixes

### "Add to Home Screen" tidak muncul

**Android:**
- Pastikan HTTPS aktif
- Cek manifest.json valid (DevTools > Application > Manifest)
- Service worker harus registered & active
- Buka di Chrome (bukan in-app browser)

**iOS:**
- Harus buka di Safari (bukan Chrome iOS atau in-app browser)
- HTTPS wajib
- Manifest.json harus valid

### Service Worker tidak register

```javascript
// Cek di console
navigator.serviceWorker.getRegistrations().then(regs => console.log(regs));

// Force re-register
navigator.serviceWorker.register('/web/sw.js', { scope: '/' });
```

### Icons tidak muncul

- Cek path di manifest.json: `/assets/icon.png` (absolute path)
- Pastikan file exists dan accessible
- Clear cache & hard reload (Ctrl+Shift+R)

### Offline mode tidak jalan

- Cek service worker cache strategy
- Pastikan files di-cache di install event
- Test di DevTools > Application > Service Workers > Offline

### App tidak fullscreen (browser bar masih ada)

- Cek `display: "standalone"` di manifest.json
- Uninstall & reinstall app
- Pastikan buka dari home screen icon (bukan bookmark)

---

## 📊 Performance Benchmarks

Target metrics:

| Metric | Target | Tool |
|--------|--------|------|
| First Contentful Paint | < 1.5s | Lighthouse |
| Time to Interactive | < 3.5s | Lighthouse |
| Speed Index | < 3.0s | Lighthouse |
| Lighthouse PWA Score | 100/100 | Lighthouse |
| Bundle Size | < 500KB | Network tab |

---

## 🚀 Pre-Launch Checklist

Sebelum launch ke production:

- [ ] All tests passed (Android + iOS)
- [ ] Lighthouse PWA score 100/100
- [ ] Tested on multiple devices (min 3 devices)
- [ ] Tested on slow 3G network
- [ ] Offline mode works
- [ ] Analytics integrated (optional)
- [ ] Error tracking setup (Sentry, optional)
- [ ] API endpoints production-ready
- [ ] Database migrations run
- [ ] Environment variables set
- [ ] Health check endpoint works
- [ ] CORS configured (if needed)

---

## 📝 Testing Log Template

```
Date: ___________
Tester: ___________

Device: ___________
OS: ___________
Browser: ___________

✅ Install prompt works
✅ Icon & splash screen correct
✅ Standalone mode works
✅ All features functional
✅ Offline mode works
✅ Performance acceptable

Issues found:
- 
- 

Notes:
- 
```

---

## 🎯 Next Steps After Testing

1. **Fix any issues** found during testing
2. **Re-test** on affected devices
3. **Monitor** analytics & error logs post-launch
4. **Iterate** based on user feedback
5. **Consider** native app wrapper (Capacitor) if needed

---

## 📚 Resources

- [PWA Checklist](https://web.dev/pwa-checklist/)
- [Lighthouse CI](https://github.com/GoogleChrome/lighthouse-ci)
- [PWA Builder](https://www.pwabuilder.com/)
- [Can I Use - PWA](https://caniuse.com/?search=pwa)
- [iOS PWA Limitations](https://firt.dev/ios-14.5/)
