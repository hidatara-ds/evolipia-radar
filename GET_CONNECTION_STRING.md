# Cara Mendapatkan Connection String Supabase yang Benar

## ❌ Error: "no such host"

Error ini terjadi karena hostname Supabase tidak valid. Kemungkinan:
1. Project Supabase belum dibuat
2. Connection string salah/typo
3. Project Supabase sudah dihapus

## ✅ Cara Mendapatkan Connection String yang Benar

### Step 1: Buka Supabase Dashboard

1. Login ke [supabase.com](https://supabase.com)
2. Pilih project Anda (atau buat baru jika belum ada)

### Step 2: Get Connection String

1. Di sidebar kiri, klik **Settings** (icon gear ⚙️)
2. Klik **Database**
3. Scroll ke bagian **Connection string**
4. Pilih tab **URI** (bukan Nodejs atau lainnya)
5. Copy connection string yang muncul

**Format yang benar:**
```
postgresql://postgres.[PROJECT-REF]:[YOUR-PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres
```

Atau format lama:
```
postgresql://postgres:[YOUR-PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
```

### Step 3: Verify Project Reference

**Project Reference** adalah string unik untuk project Anda. Contoh:
- `axxbcfnrlfnjyfanaogb` ← Ini yang Anda pakai (mungkin tidak valid)
- `xyzabcdefghijklmnop` ← Contoh valid

**Cara cek Project Reference:**
1. Di Supabase Dashboard → **Settings** → **General**
2. Lihat **Reference ID**
3. Pastikan sama dengan yang ada di connection string

### Step 4: Update .env File

Edit file `.env` dengan connection string yang benar:

```bash
# File: .env
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
MAX_FETCH_BYTES=2000000
FETCH_TIMEOUT_SECONDS=8
```

**Contoh dengan data real:**
```bash
DATABASE_URL=postgresql://postgres:eRmnQG8QCxiblkWe@db.xyzabcdefghijklmnop.supabase.co:5432/postgres
```

### Step 5: Test Connection

```bash
# Test dengan worker
go run ./cmd/worker
```

Expected output:
```
========================================
Worker started at: 2026-03-09T23:12:01+07:00
========================================
Starting ingestion...
```

## 🔍 Troubleshooting

### Error: "no such host"

**Penyebab:**
- Hostname tidak valid
- Project Supabase tidak ada
- Typo di connection string

**Solusi:**
1. Verify project exists di Supabase Dashboard
2. Copy connection string lagi dari Settings → Database
3. Pastikan tidak ada typo
4. Pastikan project tidak paused (free tier auto-pause setelah 1 minggu inactive)

### Error: "password authentication failed"

**Penyebab:**
- Password salah
- Password berubah

**Solusi:**
1. Reset password di Supabase Dashboard:
   - Settings → Database → Database password → Reset password
2. Copy password baru
3. Update `.env` file

### Error: "connection refused"

**Penyebab:**
- Port salah
- Firewall blocking

**Solusi:**
1. Pastikan port `5432` (bukan `6543` atau lainnya)
2. Check firewall/antivirus
3. Try dari network lain

### Project Paused (Free Tier)

Supabase free tier auto-pause setelah 1 minggu tidak aktif.

**Cara unpause:**
1. Buka Supabase Dashboard
2. Pilih project
3. Klik **Restore** atau **Resume**
4. Tunggu ~2 menit
5. Try connection lagi

## 📝 Checklist

Sebelum run worker, pastikan:

- [ ] Project Supabase exists dan active (not paused)
- [ ] Connection string copied dari Dashboard (Settings → Database → URI)
- [ ] Password correct (no typo)
- [ ] Project Reference ID correct
- [ ] File `.env` updated dengan connection string yang benar
- [ ] Migrations sudah di-run (via SQL Editor)
- [ ] Default sources sudah di-seed

## 🎯 Quick Test Connection

Buat file `test-connection.go` untuk test connection:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5"
)

func main() {
    // Get from .env or hardcode untuk testing
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        fmt.Println("DATABASE_URL not set")
        os.Exit(1)
    }
    
    fmt.Printf("Testing connection to: %s\n", maskPassword(connStr))
    
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ Connection failed: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close(context.Background())
    
    var result int
    err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ Query failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("✅ Connection successful!")
    
    // Test tables exist
    var count int
    err = conn.QueryRow(context.Background(), 
        "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ Table check failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("✅ Found %d tables in database\n", count)
}

func maskPassword(connStr string) string {
    // Simple masking for display
    return "postgresql://postgres:***@db.xxx.supabase.co:5432/postgres"
}
```

Run:
```bash
# Load .env
export $(cat .env | xargs)

# Test
go run test-connection.go
```

## 📞 Need Help?

Jika masih error:
1. Screenshot error message
2. Screenshot Supabase Dashboard → Settings → Database (hide password!)
3. Verify project status (active/paused)
4. Check project region (Singapore, US, EU, etc)

## 🔗 Useful Links

- [Supabase Dashboard](https://supabase.com/dashboard)
- [Supabase Database Settings](https://supabase.com/dashboard/project/_/settings/database)
- [Supabase Connection Docs](https://supabase.com/docs/guides/database/connecting-to-postgres)
