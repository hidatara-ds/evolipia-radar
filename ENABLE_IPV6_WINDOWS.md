# Enable IPv6 di Windows untuk Supabase

## Masalah

Supabase database `db.axxbcfnrlfnjyfanaogb.supabase.co` hanya resolve ke IPv6 address, tapi Windows Anda belum enable IPv6.

## ✅ Solusi: Enable IPv6 (5 Menit)

### Cara 1: Via Network Adapter Settings (Recommended)

1. **Buka Control Panel**
   - Tekan `Win + R`
   - Ketik: `ncpa.cpl`
   - Enter

2. **Pilih Network Adapter**
   - Right-click pada adapter yang aktif (WiFi atau Ethernet)
   - Pilih **Properties**

3. **Enable IPv6**
   - Cari **"Internet Protocol Version 6 (TCP/IPv6)"**
   - **Centang** checkbox-nya (jika belum)
   - Klik **OK**

4. **Restart Network**
   - Disable adapter (right-click → Disable)
   - Tunggu 5 detik
   - Enable lagi (right-click → Enable)

5. **Test**
   ```bash
   # Test IPv6 connectivity
   ping -6 google.com
   
   # Test Supabase hostname
   ping -6 db.axxbcfnrlfnjyfanaogb.supabase.co
   ```

### Cara 2: Via Command Prompt (Alternatif)

```cmd
# Buka CMD as Administrator
# Tekan Win + X → Command Prompt (Admin)

# Check IPv6 status
netsh interface ipv6 show interface

# Enable IPv6 (jika disabled)
netsh interface ipv6 install

# Reset IPv6 stack
netsh int ipv6 reset

# Restart network adapter
ipconfig /release
ipconfig /renew
```

### Cara 3: Via PowerShell (Alternatif)

```powershell
# Buka PowerShell as Administrator

# Enable IPv6
Enable-NetAdapterBinding -Name "Wi-Fi" -ComponentID ms_tcpip6

# Atau untuk Ethernet
Enable-NetAdapterBinding -Name "Ethernet" -ComponentID ms_tcpip6

# Restart adapter
Restart-NetAdapter -Name "Wi-Fi"
```

## 🧪 Test Connection Setelah Enable IPv6

```bash
# Test DNS resolve
nslookup db.axxbcfnrlfnjyfanaogb.supabase.co

# Test connection
export DATABASE_URL="postgresql://postgres:eRmnQG8QCxiblkWe@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres"
go run test-connection.go

# Run worker
go run ./cmd/worker
```

## 🔄 Alternatif: Deploy ke GitHub Actions (Recommended!)

Jika IPv6 tetap tidak bisa di-enable (misalnya ISP tidak support), **deploy ke GitHub Actions saja**!

GitHub Actions runner **support IPv6**, jadi worker akan jalan normal di cloud.

### Setup GitHub Actions:

1. **Add Secret**
   - GitHub repo → Settings → Secrets → Actions
   - Name: `SUPABASE_DB_URL`
   - Value: `postgresql://postgres:eRmnQG8QCxiblkWe@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres`

2. **Push Branch**
   ```bash
   git push origin deploy-supabase
   ```

3. **Manual Trigger**
   - GitHub → Actions tab
   - Pilih "Scheduled News Scraper"
   - Klik "Run workflow"

4. **Check Logs**
   - Worker akan jalan di GitHub Actions (support IPv6)
   - Check logs untuk verify success

## 📱 Flutter App

Untuk Flutter app, **tidak ada masalah**! Flutter app akan jalan di:
- **Android/iOS**: Support IPv6 by default
- **Web**: Browser support IPv6
- **Desktop**: Bisa pakai Supabase SDK yang handle connection

Flutter pakai **Supabase SDK** yang otomatis handle connection (bukan direct PostgreSQL), jadi tidak ada masalah IPv6.

## 🎯 Rekomendasi

**Untuk Development Lokal:**
- Enable IPv6 di Windows (Cara 1 di atas)

**Untuk Production:**
- Worker jalan di GitHub Actions (sudah support IPv6)
- Flutter app pakai Supabase SDK (tidak perlu direct PostgreSQL)

## ❓ Troubleshooting

### IPv6 Masih Tidak Bisa

Jika setelah enable IPv6 masih tidak bisa:

1. **Check ISP Support**
   ```bash
   # Test IPv6 connectivity
   ping -6 ipv6.google.com
   ```
   
   Jika gagal, ISP Anda mungkin tidak support IPv6.

2. **Gunakan VPN**
   - Beberapa VPN support IPv6
   - Atau pakai Cloudflare WARP (gratis, support IPv6)

3. **Deploy ke Cloud**
   - GitHub Actions (recommended)
   - Atau deploy worker ke Heroku/Railway/Fly.io

### Check IPv6 Status

```bash
# Windows
ipconfig /all | findstr "IPv6"

# Bash
ip -6 addr show
```

## 💡 Kesimpulan

**Pilihan Terbaik:**
1. ✅ Enable IPv6 di Windows (untuk local testing)
2. ✅ Deploy worker ke GitHub Actions (untuk production)
3. ✅ Flutter app pakai Supabase SDK (tidak perlu worry tentang IPv6)

Dengan setup ini:
- Worker jalan otomatis 3x/hari di GitHub Actions
- Flutter app query langsung ke Supabase via SDK
- Tidak perlu maintain server sendiri
