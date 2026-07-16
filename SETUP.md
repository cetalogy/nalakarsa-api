# Setup Guide - Nalakarsa Backend

Panduan lengkap instalasi dari nol hingga aplikasi **Nalakarsa Backend** berjalan sukses di lokal Anda.

---

## 1. Instalasi Prasyarat (Prerequisites)

### A. Git (Version Control)
* **Windows**: Unduh installer di [git-scm.com](https://git-scm.com/) lalu ikuti petunjuk wizard.
* **Mac (Homebrew)**: `brew install git`
* **Linux (Ubuntu/Debian)**: `sudo apt update && sudo apt install git -y`

### B. Golang (Runtime Go)
* Unduh paket installer Go stabil versi terbaru di [go.dev/dl](https://go.dev/dl/).
* Setelah instalasi, buka terminal/command prompt lalu verifikasi instalasi:
  ```bash
  go version
  ```

### C. PostgreSQL (Database)
* Unduh di [postgresql.org/download](https://www.postgresql.org/download/) sesuai sistem operasi Anda.
* Selama proses instalasi, tetapkan password default untuk user `postgres` (misalnya: `postgres`).
* Buka PGAdmin atau psql client, lalu buat database kosong bernama `nalakarsa`:
  ```sql
  CREATE DATABASE nalakarsa;
  ```

### D. Docker & Docker Compose (Pilihan Kontainerisasi)
* Unduh dan instal **Docker Desktop** dari [docker.com](https://www.docker.com/).
* Pastikan daemon Docker menyala di komputer Anda.

### E. Air (Hot Reload Development - Opsional)
Air digunakan untuk melakukan build otomatis setiap kali Anda melakukan perubahan kode tanpa perlu restart server manual.
* Instal via Go CLI:
  ```bash
  go install github.com/air-verse/air@latest
  ```

---

## 2. Pemasangan Proyek (Project Setup)

### A. Clone Project
Clone repositori ke dalam folder komputer lokal Anda:
```bash
git clone <repository-url> nalakarsa-backend
cd nalakarsa-backend
```

### B. Konfigurasi File Lingkungan (.env)
1. Salin template `.env.example` menjadi `.env`:
   ```bash
   cp .env.example .env
   ```
2. Buka file `.env` menggunakan teks editor Anda, sesuaikan variabel database Anda:
   ```env
   PORT=8080
   ENV=development
   
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=YOUR_POSTGRES_PASSWORD
   DB_NAME=nalakarsa
   ```

### C. Unduh Dependencies
```bash
go mod tidy
```

---

## 3. Menjalankan Database Seeder & Migrasi

GORM akan secara otomatis membuat seluruh skema tabel pada saat server pertama kali dijalankan (*GORM AutoMigrate*).

Namun, jika Anda ingin mengisi database dengan data dummy awal (mock user, diskusi, & proyek kolaborasi) secara instan:
```bash
go run cmd/api/main.go --seed
```
*Catatan: Parameter `--seed` akan membersihkan database terlebih dahulu dan mengisinya dengan 3 akun testing dengan password default `password123`:*
* *Akademisi*: `dosen@nalakarsa.id`
* *Praktisi*: `pengusaha@nalakarsa.id`
* *Profesional*: `engineer@nalakarsa.id`

---

## 4. Menjalankan Aplikasi

### Metode A: Menjalankan Lokal (Biasa)
```bash
go run cmd/api/main.go
```
API server akan mendengarkan request di http://localhost:8080.

### Metode B: Menjalankan Lokal (Menggunakan Air - Hot Reload)
Cukup jalankan command berikut di direktori utama:
```bash
air
```

### Metode C: Menjalankan dengan Docker Compose
Jika menggunakan docker compose, Anda tidak perlu menginstal PostgreSQL lokal atau menjalankan migrasi manual.
```bash
docker-compose up --build -d
```
Docker compose akan otomatis:
1. Menyiapkan database PostgreSQL di port `5432`.
2. Melakukan pengecekan kesehatan (*health check*) DB.
3. Mem-build image Go backend dan menghubungkannya ke DB.
4. Menjalankan migrasi database otomatis.

---

## 5. Build Produksi (Production Build)

Untuk menghasilkan binary mandiri yang siap disebarkan ke server cloud:
```bash
go build -o nalakarsa_app ./cmd/api
```
Untuk menjalankan binary hasil build:
```bash
./nalakarsa_app
```
*(Pastikan file `.env` berada di folder yang sama saat Anda menjalankan binary ini di server produksi).*
