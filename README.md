# Nalakarsa Backend - Collaboration Ecosystem

Nalakarsa Backend adalah prototipe backend ekosistem kolaborasi yang production-ready. Proyek ini menghubungkan **Akademisi**, **Praktisi**, dan **Profesional** dalam satu platform interaktif yang dilengkapi dengan forum diskusi, data direktori pengguna, dan sistem pengajuan kolaborasi proyek.

Sistem ini diimplementasikan menggunakan bahasa **Golang** dengan **Clean Architecture**, **Repository Pattern**, **PostgreSQL** (menggunakan **GORM**), **JWT Authentication** (dengan Refresh Token Rotation), input validation, structured logging, unit testing, dan kontainerisasi **Docker**.

---

## Fitur Utama

1. **Role-Based Authentication**: Pendaftaran dan masuk akun dengan memilih peran: Akademisi, Praktisi, atau Profesional. Kata sandi dienkripsi menggunakan Bcrypt.
2. **Access & Refresh Token (Double Token System)**: Sesi aman menggunakan access token berdurasi pendek dan refresh token rotasi tersimpan di DB.
3. **Direktori Profil Pengguna**: Menampilkan info kualifikasi profesional, afiliasi institusi, dan bidang keahlian yang dapat dicari/disaring.
4. **Ruang Diskusi (Discussion Forum)**: Kategori forum terstruktur dengan pencarian, filter kategori/peran, serta komentar.
5. **Matchmaking Proyek Kolaborasi**: Pengajuan proyek oleh inisiator dan pengiriman lamaran (apply) dari pengguna peran lain dengan kualifikasi yang cocok.

---

## Persyaratan (Requirements)

Sebelum memulai, pastikan Anda telah menginstal:
* **Golang** (versi 1.21 ke atas)
* **PostgreSQL** (versi 14 ke atas)
* **Docker & Docker Compose** (jika ingin menjalankan dengan kontainer)
* **Git**

---

## Folder Structure

Aplikasi ini menggunakan standar **Clean Architecture** Go:
```
cmd/api/               # Entrypoint aplikasi (main.go)
internal/
  ├── config/          # Parser konfigurasi .env
  ├── database/        # Koneksi database & auto-migration GORM
  ├── dto/             # Data Transfer Objects & validasi binding
  ├── handler/         # Controller HTTP (Gin)
  ├── middleware/      # Middleware JWT Auth, CORS, dan structured logger
  ├── model/           # Entitas tabel GORM database
  ├── repository/      # Kueri database (Repository Pattern)
  ├── routes/          # Mapping endpoint API
  ├── service/         # Logika Bisnis (Usecase layer)
  └── utils/           # Helper enkripsi, JWT, dan pagination
docs/                  # Dokumentasi API (openapi.yaml)
seed/                  # Database seeder dummy data
```

---

## Cara Install & Konfigurasi

### 1. Clone Project
```bash
git clone <repository-url>
cd nalakarsa-backend
```

### 2. Setup Environment
Salin template konfigurasi `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi kredensial PostgreSQL dan JWT Key Anda di dalam file `.env`.

### 3. Install Dependencies
```bash
go mod download
```

---

## Menjalankan Aplikasi

Ada dua cara untuk menjalankan aplikasi ini:

### Metode A: Menjalankan Secara Lokal

1. **Jalankan Database Seeder** (Opsional - untuk mengisi data demonstrasi ke PostgreSQL lokal):
   *Pastikan database PostgreSQL Anda sudah menyala dan nama DB sesuai konfigurasi `.env`.*
   ```bash
   go run cmd/api/main.go --seed
   ```
2. **Jalankan Server API**:
   ```bash
   go run cmd/api/main.go
   ```
   Aplikasi akan berjalan di http://localhost:8080.

### Metode B: Menjalankan Menggunakan Docker Compose (Sangat Direkomendasikan)

Metode ini otomatis mengunduh PostgreSQL, menginisialisasi database, dan menjalankan backend di dalam container.
```bash
docker-compose up -build -d
```
Untuk menghentikan container:
```bash
docker-compose down
```

---

## Menjalankan Unit Test

Kami menyediakan pengujian unit untuk layer bisnis (Service) menggunakan mock repository:
```bash
go test ./internal/service/... -v
```

---

## Dokumentasi API & Swagger

Dokumentasi API lengkap didefinisikan menggunakan standar OpenAPI 3.0. Anda dapat membaca referensi lengkap endpoint di:
* [API_REFERENCE.md](file:///d:/NALAKARSA/nalakarsa-backend/API_REFERENCE.md)
* Spesifikasi OpenAPI YAML di [docs/openapi.yaml](file:///d:/NALAKARSA/nalakarsa-backend/docs/openapi.yaml) (Bisa di-import ke Postman atau Swagger UI).
