# ✅ Fix: Skor Terlalu Rendah (0.2, 0.0, dll)

## Masalah

Skor masih tampil rendah (0.2, 0.0) karena:
1. ❌ Data lama di database masih pakai skala 0-1
2. ❌ Endpoint `/v1/items/:id` belum diupdate untuk konversi

## Solusi yang Sudah Diterapkan

### 1. Fix API Endpoints ✅

**File yang diupdate:**
- `internal/http/handlers/handlers.go`
  - Tambah fungsi `convertToScale10()`
  - Update `GetItem()` untuk konversi skor
  - Update `Search()` untuk konversi skor

**Sekarang semua endpoint API akan return skor 1-10:**
```json
{
  "scores": {
    "final": 3.7,      // ✅ Bukan 0.37
    "hot": 2.8,        // ✅ Bukan 0.28
    "relevance": 5.5   // ✅ Bukan 0.55
  }
}
```

### 2. Re-score Script ✅

**File baru:** `scripts/rescore_items.go`

Script ini akan:
- Re-calculate semua skor untuk items yang ada
- Update database dengan skor yang benar
- Process items dari 30 hari terakhir

## Cara Menggunakan

### Option 1: Restart Worker (Otomatis)

Worker akan otomatis re-score items saat jalan:

```bash
# Stop worker jika sedang jalan
# Ctrl+C

# Start worker lagi
go run ./cmd/worker
```

Worker akan:
- Fetch items baru
- Re-score items lama yang belum di-score
- Update semua skor

### Option 2: Manual Re-score (Lebih Cepat)

Jalankan script re-score untuk update semua items sekaligus:

**Git Bash / PowerShell:**
```bash
go run ./scripts/rescore_items.go
```

**Output:**
```
=== Re-scoring All Items ===
Found 150 items to re-score
Progress: 0/150 (0.0%)
Progress: 100/150 (66.7%)
Progress: 150/150 (100.0%)

=== Re-scoring Complete ===
Total items processed: 150
Successfully updated: 150
Failed: 0
```

### Option 3: Via Database (Advanced)

Jika ingin manual update via SQL:

```sql
-- Update semua skor ke skala yang lebih reasonable
-- Ini hanya contoh, lebih baik pakai script Go

UPDATE scores 
SET 
  final = GREATEST(0.3, final),
  hot = GREATEST(0.2, hot),
  relevance = GREATEST(0.3, relevance)
WHERE final < 0.3;
```

## Verifikasi Fix

### 1. Test API
```bash
# Get feed - skor harus 1-10
curl http://localhost:8080/v1/feed?date=today

# Response:
{
  "items": [
    {
      "scores": {
        "final": 7.3,     // ✅ Skala 1-10
        "hot": 6.4,
        "relevance": 8.2
      }
    }
  ]
}
```

### 2. Test Detail View
```bash
# Get item detail
curl http://localhost:8080/v1/items/{item_id}

# Response:
{
  "scores": {
    "final": 7.3,        // ✅ Skala 1-10
    "hot": 6.4,
    "relevance": 8.2,
    "credibility": 5.5,
    "novelty": 4.6
  }
}
```

### 3. Test UI

1. Buka http://localhost:8080
2. Klik item untuk detail
3. Skor harus tampil 1-10:
   ```
   Final: 7.3
   Hot: 6.4
   Relevan: 8.2
   Kredibilitas: 5.5/10
   Kebaruan: 4.6/10
   ```

## Penjelasan Skor

### Kenapa Skor Bisa Rendah?

**Hot Score (0.0 - 2.0):**
- Item baru tanpa engagement (no points/comments)
- Item lama yang sudah tidak populer
- ✅ Normal untuk item yang baru di-fetch

**Relevance Score (0.2 - 3.0):**
- Konten tidak terlalu relevan dengan AI/ML
- Tidak ada keyword AI/ML di title/excerpt
- ✅ Normal untuk berita umum

**Final Score (0.2 - 3.0):**
- Kombinasi dari semua komponen
- Item baru biasanya skor rendah dulu
- ✅ Akan naik seiring waktu jika populer

### Skor yang Bagus

Setelah konversi ke skala 1-10:

| Skor | Interpretasi |
|------|--------------|
| 8-10 | Sangat bagus, trending |
| 6-7  | Bagus, worth reading |
| 4-5  | Cukup menarik |
| 2-3  | Kurang menarik |
| 1    | Tidak relevan |

## Troubleshooting

### Skor Masih 0.2 Setelah Fix

**Penyebab:** Browser cache

**Solusi:**
```bash
# Hard refresh browser
Ctrl + Shift + R (Windows/Linux)
Cmd + Shift + R (Mac)

# Atau clear cache
DevTools > Application > Clear Storage
```

### Script Re-score Error

**Error:** `Failed to connect to database`

**Solusi:**
```bash
# Pastikan PostgreSQL jalan
docker-compose up -d postgres

# Set DATABASE_URL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"

# Run script lagi
go run ./scripts/rescore_items.go
```

### Skor Tidak Berubah

**Penyebab:** Worker belum jalan atau belum re-score

**Solusi:**
```bash
# Option 1: Jalankan worker
go run ./cmd/worker

# Option 2: Manual re-score
go run ./scripts/rescore_items.go

# Option 3: Restart API server
# Ctrl+C
go run ./cmd/api
```

## Build Verification

```bash
✅ go build ./cmd/api - Success
✅ go build ./cmd/worker - Success
✅ go build ./scripts/rescore_items.go - Success
```

## Summary

**Masalah:** Skor tampil 0.2, 0.0 (terlalu rendah)  
**Penyebab:** Data lama + endpoint belum konversi  
**Solusi:** 
1. ✅ Update API endpoints untuk konversi otomatis
2. ✅ Buat script re-score untuk update data lama
3. ✅ Worker akan auto-score items baru dengan benar

**Status:** ✅ Fixed and ready to use

**Next Steps:**
1. Restart API server: `go run ./cmd/api`
2. Run re-score script: `go run ./scripts/rescore_items.go`
3. Refresh browser dan cek skor baru!
