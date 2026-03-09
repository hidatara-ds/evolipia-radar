# Fix Row Level Security (RLS) - Supabase

Security Advisor mendeteksi beberapa tables belum enable RLS. Mari kita fix!

## ⚠️ Tables yang Perlu RLS

Dari screenshot Security Advisor:
- `public.sources`
- `public.items` ✅ (sudah enable)
- `public.signals`
- `public.scores`
- `public.summaries`
- `public.fetch_runs`
- `public.scrape_logs`

## 🔒 Enable RLS untuk Semua Tables

Copy-paste SQL berikut di **Supabase SQL Editor**, lalu klik **Run**:

```sql
-- Enable RLS untuk semua tables
ALTER TABLE sources ENABLE ROW LEVEL SECURITY;
ALTER TABLE items ENABLE ROW LEVEL SECURITY;
ALTER TABLE signals ENABLE ROW LEVEL SECURITY;
ALTER TABLE scores ENABLE ROW LEVEL SECURITY;
ALTER TABLE summaries ENABLE ROW LEVEL SECURITY;
ALTER TABLE fetch_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE scrape_logs ENABLE ROW LEVEL SECURITY;
```

## 📋 Create RLS Policies

### 1. Public Read Access (untuk Flutter App)

Tables yang perlu public read access:

```sql
-- Items: Allow anonymous read
CREATE POLICY "Allow anonymous read items" ON items
    FOR SELECT USING (true);

-- Scores: Allow anonymous read
CREATE POLICY "Allow anonymous read scores" ON scores
    FOR SELECT USING (true);

-- Summaries: Allow anonymous read
CREATE POLICY "Allow anonymous read summaries" ON summaries
    FOR SELECT USING (true);

-- Signals: Allow anonymous read
CREATE POLICY "Allow anonymous read signals" ON signals
    FOR SELECT USING (true);

-- Sources: Allow anonymous read (only enabled sources)
CREATE POLICY "Allow anonymous read enabled sources" ON sources
    FOR SELECT USING (enabled = true);
```

### 2. Service Role Access (untuk Worker)

Worker menggunakan `service_role` key yang **bypass RLS**, jadi tidak perlu policy khusus.

Tapi untuk monitoring, kita bisa allow read untuk scrape_logs dan fetch_runs:

```sql
-- Scrape logs: Allow anonymous read (untuk monitoring)
CREATE POLICY "Allow anonymous read scrape_logs" ON scrape_logs
    FOR SELECT USING (true);

-- Fetch runs: Allow anonymous read (untuk monitoring)
CREATE POLICY "Allow anonymous read fetch_runs" ON fetch_runs
    FOR SELECT USING (true);
```

### 3. Write Access (hanya untuk authenticated/service_role)

Secara default, jika tidak ada policy untuk INSERT/UPDATE/DELETE, maka hanya `service_role` yang bisa write. Ini sudah aman!

Tapi jika mau explicit:

```sql
-- Deny all writes from anonymous users (optional, sudah default behavior)
-- Hanya service_role (worker) yang bisa write

-- Items: Deny anonymous insert/update/delete
CREATE POLICY "Deny anonymous write items" ON items
    FOR ALL USING (auth.role() = 'authenticated');

-- Scores: Deny anonymous write
CREATE POLICY "Deny anonymous write scores" ON scores
    FOR ALL USING (auth.role() = 'authenticated');

-- Summaries: Deny anonymous write
CREATE POLICY "Deny anonymous write summaries" ON summaries
    FOR ALL USING (auth.role() = 'authenticated');

-- Signals: Deny anonymous write
CREATE POLICY "Deny anonymous write signals" ON signals
    FOR ALL USING (auth.role() = 'authenticated');

-- Sources: Deny anonymous write
CREATE POLICY "Deny anonymous write sources" ON sources
    FOR ALL USING (auth.role() = 'authenticated');

-- Scrape logs: Deny anonymous write
CREATE POLICY "Deny anonymous write scrape_logs" ON scrape_logs
    FOR ALL USING (auth.role() = 'authenticated');

-- Fetch runs: Deny anonymous write
CREATE POLICY "Deny anonymous write fetch_runs" ON fetch_runs
    FOR ALL USING (auth.role() = 'authenticated');
```

## 🎯 Recommended: Simple & Secure

Untuk use case kita (Flutter app read-only, worker write via service_role), cukup:

```sql
-- Enable RLS untuk semua tables
ALTER TABLE sources ENABLE ROW LEVEL SECURITY;
ALTER TABLE items ENABLE ROW LEVEL SECURITY;
ALTER TABLE signals ENABLE ROW LEVEL SECURITY;
ALTER TABLE scores ENABLE ROW LEVEL SECURITY;
ALTER TABLE summaries ENABLE ROW LEVEL SECURITY;
ALTER TABLE fetch_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE scrape_logs ENABLE ROW LEVEL SECURITY;

-- Allow anonymous READ untuk tables yang dibutuhkan Flutter app
CREATE POLICY "Allow anonymous read" ON items FOR SELECT USING (true);
CREATE POLICY "Allow anonymous read" ON scores FOR SELECT USING (true);
CREATE POLICY "Allow anonymous read" ON summaries FOR SELECT USING (true);
CREATE POLICY "Allow anonymous read" ON signals FOR SELECT USING (true);
CREATE POLICY "Allow anonymous read" ON sources FOR SELECT USING (enabled = true);
CREATE POLICY "Allow anonymous read" ON scrape_logs FOR SELECT USING (true);
CREATE POLICY "Allow anonymous read" ON fetch_runs FOR SELECT USING (true);

-- Write access: Secara default hanya service_role yang bisa write (sudah aman)
-- Tidak perlu policy tambahan untuk write
```

## ✅ Verify RLS Enabled

```sql
-- Check RLS status
SELECT 
    schemaname,
    tablename,
    rowsecurity as rls_enabled
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY tablename;

-- Expected: All tables should have rls_enabled = true

-- Check policies
SELECT 
    schemaname,
    tablename,
    policyname,
    cmd as command,
    qual as using_expression
FROM pg_policies 
WHERE schemaname = 'public'
ORDER BY tablename, policyname;

-- Expected: Should see policies for each table
```

## 🧪 Test dari Flutter App

Setelah enable RLS, test query dari Flutter:

```dart
// Should work (anonymous read allowed)
final items = await Supabase.instance.client
  .from('items')
  .select('*, scores(*), summaries(*)')
  .limit(10);

print('Items: ${items.length}'); // Should return data

// Should fail (anonymous write not allowed)
try {
  await Supabase.instance.client
    .from('items')
    .insert({'title': 'Test', 'url': 'https://test.com'});
} catch (e) {
  print('Write blocked: $e'); // Expected: Permission denied
}
```

## 🔐 Security Best Practices

### 1. Service Role Key (Worker)

Worker menggunakan `service_role` key yang **bypass RLS**:

```bash
# GitHub Secret: SUPABASE_DB_URL
# Format: postgresql://postgres:[PASSWORD]@db.xxx.supabase.co:5432/postgres
# Ini pakai service_role credentials, bypass RLS
```

### 2. Anon Key (Flutter App)

Flutter app menggunakan `anon` key yang **respect RLS**:

```dart
// Flutter app pakai anon key
await Supabase.initialize(
  url: 'https://xxx.supabase.co',
  anonKey: 'eyJhbGc...', // Anon key (public, safe to expose)
);

// Queries akan respect RLS policies
```

### 3. Never Expose Service Role Key

❌ **JANGAN** pakai `service_role` key di Flutter app!
✅ **HANYA** pakai `anon` key di Flutter app
✅ `service_role` key **HANYA** untuk backend (GitHub Actions worker)

## 📝 Summary

**Enable RLS:**
```sql
ALTER TABLE [table_name] ENABLE ROW LEVEL SECURITY;
```

**Allow Read:**
```sql
CREATE POLICY "Allow anonymous read" ON [table_name]
    FOR SELECT USING (true);
```

**Deny Write:**
- Secara default sudah deny (tidak perlu policy tambahan)
- Hanya `service_role` yang bisa write

**Result:**
- ✅ Flutter app bisa read data (via anon key)
- ✅ Flutter app tidak bisa write data (protected)
- ✅ Worker bisa read & write (via service_role)
- ✅ Security Advisor warnings hilang
