# Project Structure - Nalakarsa Backend

Dokumen ini menjelaskan struktur folder dan pembagian modul dalam pengembangan backend **Nalakarsa** menggunakan **Clean Architecture** di Go.

---

## Folder Layout

Struktur direktori utama proyek dirancang sebagai berikut:

```
├── cmd/
│   └── api/
│       └── main.go          # Entry point utama aplikasi
├── internal/
│   ├── config/              # Manajemen variabel lingkungan (.env)
│   ├── database/            # Inisialisasi DB, migrasi, dan seeders GORM
│   ├── dto/                 # Data Transfer Objects untuk request/response & validasi
│   ├── handler/             # Handler HTTP per modul (Gin Controllers)
│   ├── middleware/          # JWT Auth, CORS, Logger, Recovery
│   ├── model/               # Entitas model database GORM
│   ├── repository/          # Repository per modul (Repository Pattern)
│   ├── routes/              # Register route per modul
│   ├── service/              # Service per modul (Usecase Layer)
│   └── utils/               # Utilitas (hash bcrypt, JWT generation, dll)
├── pkg/                     # Kode utilitas publik (opsional jika digunakan oleh project luar)
├── docs/                    # Dokumentasi OpenAPI/Swagger (dihasilkan otomatis oleh Swag)
├── migration/               # File SQL schema (jika ada raw migration)
├── seed/                    # Data seeder dummy untuk testing local
├── scripts/                 # Shell script / Makefile helper untuk otomatisasi
├── tests/                   # Integrasi test suite
├── .env.example             # Template file environment
├── docker-compose.yml       # Konfigurasi orkestrasi container backend & db
├── Dockerfile               # Konfigurasi build container aplikasi Go (Multi-stage)
├── go.mod                   # File dependency Go modules
└── go.sum                   # Checksum dependency Go modules
```

---

## Penjelasan Detail Folder & File

### 1. `cmd/api/main.go`
Merupakan pintu masuk (entrypoint) program. Tugas utamanya adalah:
* Membaca konfigurasi menggunakan `internal/config`.
* Membuka koneksi database melalui `internal/database`.
* Menginisialisasi semua layer dependensi (Repository -> Service -> Handler) menggunakan Manual Dependency Injection.
* Mendaftarkan rute API.
* Menjalankan HTTP server (Gin) pada port yang ditentukan.

### 2. `internal/config/`
Mengatur konfigurasi aplikasi menggunakan library (seperti `spf13/viper` atau standard library dengan `.env` parser). Menyimpan struct konfigurasi seperti database URL, JWT secret key, JWT expiration, port server, dan environment state (development/production).

### 3. `internal/database/`
* Mengelola pembuatan koneksi ke PostgreSQL menggunakan GORM.
* Mengatur proses migrasi schema database secara otomatis (*GORM AutoMigrate*) atau manual.
* Menyediakan fungsi pemicu data seeder awal.

### 4. `internal/model/`
Mendefinisikan entitas database yang dipetakan langsung ke tabel PostgreSQL menggunakan tag GORM (contoh: `gorm:"primaryKey"`, `gorm:"uniqueIndex"`).
File utama meliputi:
* `user.go`: Entitas User dengan validasi peran.
* `profile.go`: Entitas Profile detail 1-to-1 dengan user.
* `refresh_token.go`: Menyimpan data token refresh aktif.
* `discussion.go` & `comment.go`: Model untuk forum diskusi.
* `collaboration.go` & `application.go`: Model untuk penawaran dan lamaran kolaborasi.

### 5. `internal/dto/`
Data Transfer Object (DTO) digunakan untuk memisahkan data internal database dengan data yang dikirimkan/diterima oleh client.
* Memiliki tag `binding` (seperti `binding:"required,email"`) untuk validasi input request otomatis menggunakan Validator.
* Menyusun format response terstandarisasi untuk menyaring data sensitif (seperti membuang `password_hash` dari output).

### 6. `internal/repository/`
Menerapkan **Repository Pattern**. Layer ini adalah satu-satunya bagian kode yang berinteraksi langsung dengan database melalui GORM.
* Setiap repositori didefinisikan sebagai *interface* (misal: `UserRepository` interface) dan memiliki implementasi struct konkret (misal: `pgUserRepository`).
* Hal ini mempermudah proses mocking pada pengujian unit.
* Implementasi dipisah ke subfolder modul, misalnya `repository/user/`,
  `repository/project/`, dan `repository/discussion/`.

### 7. `internal/service/`
Layer logika bisnis (sering disebut *Usecase*).
* Layer ini memproses aturan bisnis sistem (contoh: mengecek keunikan email sebelum registrasi, membandingkan password dengan bcrypt, memvalidasi apakah pelamar kolaborasi memiliki peran yang sesuai).
* Tidak mengetahui detail database atau HTTP (tidak bergantung pada Gin atau SQL langsung).
* Setiap modul memiliki subfolder sendiri, misalnya `service/auth/` dan
  `service/project/`. Dependency repository di-import langsung dari modul terkait.

### 8. `internal/handler/`
Controller HTTP yang bertindak sebagai jembatan antara protokol HTTP dan service.
* Mengambil input request dari JSON/query param.
* Memvalidasi request menggunakan DTO.
* Memanggil fungsi logika bisnis di layer Service.
* Mengembalikan HTTP Status Code dan JSON response standar.
* Setiap modul memiliki subfolder handler sendiri dan hanya menerima service dari
  modul terkait melalui constructor.

### 9. `internal/middleware/`
Menyimpan fungsi interseptor request:
* `AuthMiddleware`: Memvalidasi header Authorization `Bearer <token>` menggunakan JWT.
* `CORSMiddleware`: Mengizinkan akses antar domain.
* `LoggerMiddleware`: Mencatat setiap request masuk ke dalam logger terstruktur (`slog`).

### 10. `internal/routes/`
`routes.NewRouter` hanya menyiapkan engine, middleware global, dan health check.
Pemetaan endpoint dipisah per modul melalui fungsi `RegisterRoutes`, misalnya
`routes/auth`, `routes/project`, dan `routes/discussion`. Semua fungsi register
dipanggil eksplisit dari `main.go` setelah dependency repository, service, dan
handler selesai dibuat. Dengan begitu route module tidak lagi terkumpul dalam
satu file besar.

### 11. `internal/utils/`
Berisi helper kecil yang dapat digunakan secara horizontal, seperti pengolah JWT (Generate, Validate), password hasher (Bcrypt), format standard pagination, dan penyeragaman format respon JSON.
