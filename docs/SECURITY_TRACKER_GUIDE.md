# 🛡️ Repository Security, Anti-Copy & Telemetry Guide

Dokumen ini menjelaskan sistem pelacakan (tracker), deteksi copas/cloning, dan notifikasi otomatis email/webhook yang dipasang pada repository **Evolipia Radar**.

---

## 📌 Fitur Utama Tracker & Anti-Theft System

1. **Visitor & Hit Count Badges (README.md)**
   - Menghitung secara otomatis berapa kali halaman `README.md` dan repository ini dilihat/diakses oleh publik.
   - Menggunakan SVG Counter Badge yang terupdate secara real-time.

2. **Canarytoken Email Alert (Web Bug Beacon)**
   - Piksel/gambar pelacak tersembunyi (*web bug*) yang tertanam di `README.md` & `docs/README.md`.
   - **Cara Kerja**: Setiap kali seseorang membuka, membaca, atau meng-copy repository (baik di GitHub, VS Code, maupun Git GUI lokal), browser/markdown viewer akan mengunduh gambar piksel pelacak.
   - **Hasil**: Server Canarytoken langsung **mengirimkan EMAIL OTOMATIS** ke pemilik repository berisi:
     - 🌐 Alamat IP Pengakses / Copier
     - 📍 Lokasi Geografis & ISP
     - 💻 User-Agent / Client (GitHub / VS Code / Browser)
     - 🕒 Timestamp Waktu Pengaksesan

3. **GitHub Action Automation (`repo-tracker.yml`)**
   - Berjalan secara otomatis setiap ada event `fork`, `watch/star`, atau jadwal berkala (`schedule` 6 jam).
   - Mengambil data statistik dari GitHub Traffic API (`clones` & `views`).
   - Mengirim notifikasi otomatis via Webhook (Discord / Telegram / Email Canary).

---

## 🚀 Cara Mengaktifkan Email Notifikasi Otomatis (Thinkst Canarytoken)

Untuk menerima **email otomatis** setiap kali ada yang membaca/meng-copy repo Anda:

### Langkah 1: Buat Canarytoken Gratis
1. Buka [https://canarytokens.org](https://canarytokens.org).
2. Pada menu **"Select your token"**, pilih **"Web bug / URL"**.
3. Masukkan **Alamat Email Anda** (tempat notifikasi akan dikirimkan).
4. Masukkan catatan pengenal pada *Comment*, misalnya: `Evolipia Radar Repo Copy Alert`.
5. Klik **"Create my Canarytoken"**.

### Langkah 2: Pasang Token di `README.md`
1. Salin **Web bug URL** atau **Image URL** yang dihasilkan Canarytokens (misal: `https://canarytokens.com/static/tags/terms/xxxxxx/index.html` atau URL `http://canarytokens.com/xxxxxx/submit.png`).
2. Masukkan URL tersebut pada `README.md` pada bagian badge pelacak:
   ```markdown
   <img src="https://canarytokens.com/xxxxxx/submit.png" width="1" height="1" alt="telemetry-beacon" />
   ```
3. *(Opsional)* Simpan URL Canarytoken ke **GitHub Repository Secrets**:
   - Buka Repository Settings -> **Secrets and variables** -> **Actions** -> **New repository secret**.
   - Name: `CANARY_TOKEN_URL`
   - Value: `https://canarytokens.com/xxxxxx/submit.png`

---

## 📊 Notifikasi Webhook Tambahan (Discord / Telegram / Slack)

Jika Anda ingin notifikasi langsung ke Discord / Telegram saat ada orang yang menge-fork atau clone repo:
1. Buat Webhook di Discord Server (Channel Settings -> Integrations -> Webhooks).
2. Tambahkan Secret di GitHub Repository:
   - Name: `DISCORD_WEBHOOK_URL`
   - Value: `https://discord.com/api/webhooks/YOUR_WEBHOOK_URL`
3. Workflow `.github/workflows/repo-tracker.yml` akan secara otomatis mengirim pesan ke Discord setiap ada aktivitas fork/clone!

---

## 🔒 Peringatan Keamanan & Hak Cipta

Semua file markdown dan kode di repository ini dilindungi oleh watermark jejak digital dan telemetri otomatis. Segera laporkan jika ditemukan penyalahgunaan tanpa izin.
