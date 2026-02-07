# Install Go di Windows

Panduan install Go untuk Windows (Git Bash / MINGW64).

## 🎯 Cara Install Go di Windows

### Opsi 1: Download dari Official Website (Recommended)

1. **Download Go:**
   - Buka: https://go.dev/dl/
   - Download file `.msi` untuk Windows (contoh: `go1.21.x.windows-amd64.msi`)

2. **Install:**
   - Double-click file `.msi` yang sudah didownload
   - Ikuti installer (Next, Next, Install)
   - Default install location: `C:\Program Files\Go`

3. **Verify Install:**
   - Buka **PowerShell baru** atau **CMD baru** (bukan Git Bash)
   ```powershell
   go version
   ```
   - Harusnya muncul: `go version go1.21.x windows/amd64`

### Opsi 2: Install via Chocolatey (jika sudah punya)

```powershell
# Run PowerShell as Administrator
choco install golang
```

### Opsi 3: Install via Scoop (jika sudah punya)

```powershell
scoop install go
```

## ⚙️ Setup PATH (Jika Belum Otomatis)

Jika `go version` masih error setelah install:

1. **Cek Go sudah terinstall:**
   ```powershell
   # Di PowerShell
   Test-Path "C:\Program Files\Go\bin\go.exe"
   # Harusnya return: True
   ```

2. **Add ke PATH (jika belum):**
   - Buka **System Properties** → **Environment Variables**
   - Di **System Variables**, cari `Path`
   - Klik **Edit** → **New**
   - Tambahkan: `C:\Program Files\Go\bin`
   - Klik **OK** di semua window

3. **Restart Terminal:**
   - Tutup semua terminal/PowerShell/Git Bash
   - Buka terminal baru
   - Test lagi: `go version`

## 🔧 Setup untuk Git Bash

Setelah Go terinstall, setup untuk Git Bash:

1. **Cek Go di PATH:**
   ```bash
   # Di Git Bash
   echo $PATH
   # Harusnya ada: /c/Program Files/Go/bin
   ```

2. **Jika tidak ada, tambahkan ke `.bashrc` atau `.bash_profile`:**
   ```bash
   # Edit file
   nano ~/.bashrc
   # atau
   notepad ~/.bashrc
   
   # Tambahkan baris ini:
   export PATH="/c/Program Files/Go/bin:$PATH"
   
   # Save dan reload
   source ~/.bashrc
   ```

3. **Verify di Git Bash:**
   ```bash
   go version
   ```

## ✅ Verify Setup Lengkap

Setelah install, test semua:

```bash
# 1. Go version
go version
# Expected: go version go1.21.x windows/amd64

# 2. Go environment
go env GOPATH
go env GOROOT

# 3. Test compile
go run --help
```

## 🚀 Setelah Go Terinstall

Kembali ke project dan lanjutkan:

```bash
# 1. Navigate ke project
cd /d/papaengineer/evolipia-radar

# 2. Download dependencies
go mod download

# 3. Verify
go mod verify
```

## 🐛 Troubleshooting

### Problem: "go: command not found" di Git Bash

**Solution 1: Gunakan PowerShell/CMD**
- Go biasanya sudah di PATH untuk PowerShell/CMD
- Gunakan PowerShell untuk Go commands

**Solution 2: Fix PATH di Git Bash**
```bash
# Tambahkan ke ~/.bashrc
export PATH="/c/Program Files/Go/bin:$PATH"
export PATH="/c/Program Files/Go/bin:$PATH"

# Reload
source ~/.bashrc
```

### Problem: Go version berbeda di Git Bash vs PowerShell

**Solution:**
- Gunakan PowerShell untuk Go commands (lebih reliable di Windows)
- Atau fix PATH di Git Bash seperti di atas

### Problem: Permission denied saat install

**Solution:**
- Run installer sebagai Administrator
- Right-click → Run as Administrator

## 📝 Quick Reference

**Install Go:**
1. Download dari https://go.dev/dl/
2. Install `.msi` file
3. Restart terminal
4. Test: `go version`

**Setup Project:**
```bash
cd /d/papaengineer/evolipia-radar
go mod download
```

---

**Setelah Go terinstall, lanjutkan ke [LOCAL_SETUP.md](LOCAL_SETUP.md)**
