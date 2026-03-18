# Vercel Setup Guide - Database Integration

## Problem
API `/api/news` gagal load karena tidak bisa akses file `news.json` di Vercel serverless environment.

## Solution
Mengubah API untuk membaca langsung dari Neon.tech PostgreSQL database.

## Setup Steps

### 1. Add Environment Variable di Vercel

1. Buka [Vercel Dashboard](https://vercel.com/dashboard)
2. Pilih project **evolipia-radar**
3. Klik **Settings** (tab di atas)
4. Klik **Environment Variables** (menu di kiri)
5. Klik **Add New**
6. Isi form:
   - **Name**: `DATABASE_URL`
   - **Value**: `postgresql://neondb_owner:npg_ntTN8wojqf3R@ep-falling-grass-a1dfoa60-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require`
   - **Environment**: Centang semua (Production, Preview, Development)
7. Klik **Save**

### 2. Redeploy

Setelah menambahkan environment variable:

**Option A - Automatic:**
- Vercel akan otomatis redeploy dalam beberapa menit

**Option B - Manual:**
1. Klik tab **Deployments**
2. Klik titik tiga (...) di deployment terakhir
3. Klik **Redeploy**
4. Pilih **Use existing Build Cache** (lebih cepat)
5. Klik **Redeploy**

### 3. Verify Deployment

Setelah deployment selesai:

1. Buka deployment logs di Vercel
2. Cari log dari `/api/news` endpoint
3. Seharusnya muncul log seperti:
   ```
   ✅ Returning X news items
   ```

4. Test API di browser:
   ```
   https://your-domain.vercel.app/api/news
   ```

## Architecture Changes

### Before:
```
GitHub Actions → Scrape → Save to data/news.json → Commit to Git
Vercel API → Read news.json file → Return JSON
```

### After:
```
GitHub Actions → Scrape → Save to Neon.tech Database
Vercel API → Query Neon.tech Database → Return JSON
```

## Benefits

✅ No file path issues in Vercel serverless environment
✅ Real-time data from database
✅ Better performance with database indexing
✅ Scalable solution
✅ Can filter by topic directly in database query

## Troubleshooting

### Error: "Database configuration missing"
- DATABASE_URL belum ditambahkan di Vercel environment variables
- Atau environment variable belum ter-apply (perlu redeploy)

### Error: "Failed to connect to database"
- Check DATABASE_URL format benar
- Check Neon.tech database masih aktif
- Check network/firewall settings

### Error: "Failed to load news"
- Check database ada data (run GitHub Actions scraper)
- Check database schema sudah benar
- Check logs di Vercel untuk detail error

## Next Steps

1. ✅ Code sudah di-push ke GitHub (commit: bc7f7e5)
2. ⏳ **ANDA HARUS**: Tambahkan DATABASE_URL di Vercel
3. ⏳ **ANDA HARUS**: Redeploy di Vercel
4. ✅ Test API endpoint

## Contact

Jika masih ada masalah setelah setup, cek:
- Vercel deployment logs
- Browser console untuk error messages
- Network tab untuk API response
