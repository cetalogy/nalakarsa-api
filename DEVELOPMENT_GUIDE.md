# Development Guide - Nalakarsa Backend

Panduan ini bertujuan untuk membantu pengembang memahami arsitektur proyek Nalakarsa Backend, standar penulisan kode, serta alur penambahan fitur baru (model, repository, service, handler, dan endpoint).

---

## 1. Arsitektur Aplikasi

Proyek ini menerapkan **Clean Architecture** (atau sering dikenal sebagai Onion/Hexagonal Architecture) dikombinasikan dengan **Repository Pattern**. Arsitektur ini bertujuan memisahkan kode berdasarkan tanggung jawabnya (Separation of Concerns) sehingga mempermudah proses pengujian, pemeliharaan, dan penggantian pustaka eksternal (seperti database driver atau HTTP router).

### Flow Request & Dependency Rule
Ketergantungan kode hanya berjalan satu arah ke dalam (dari luar ke dalam). Layer terdalam tidak tahu apa-apa tentang teknologi di layer terluar.

```
+-------------------------------------------------------------+
|                     1. USER INTERACTION                     |
|            (Postman / Swagger / Frontend Client)            |
+------------------------------+------------------------------+
                               | (HTTP Request)
                               v
+-------------------------------------------------------------+
|                2. HTTP HANDLER (Gin Controller)             |
|   - Tempat parsing JSON/Query Parameter                     |
|   - Melakukan validasi awal via Struct Binding DTO          |
+------------------------------+------------------------------+
                               | (DTO Struct)
                               v
+-------------------------------------------------------------+
|                  3. SERVICE LAYER (Business Logic)          |
|   - Berisi aturan bisnis (business logic)                   |
|   - Verifikasi hak akses pengguna secara kualitatif         |
|   - Tidak peduli HTTP protocol atau framework               |
+------------------------------+------------------------------+
                               | (Model Struct / Params)
                               v
+-------------------------------------------------------------+
|            4. REPOSITORY LAYER (Data Access - GORM)         |
|   - Bertanggung jawab melakukan kueri ke database           |
|   - Terdiri atas Interface dan Implementasi konkret         |
+------------------------------+------------------------------+
                               | (SQL Query)
                               v
+-------------------------------------------------------------+
|                         5. DATABASE                         |
|                         (PostgreSQL)                        |
+-------------------------------------------------------------+
```

---

## 2. Alur Menambah Fitur / Endpoint Baru

Untuk menambahkan fitur baru (misalnya fitur **Notifikasi**), ikuti langkah-langkah terstruktur berikut:

### Langkah A: Membuat Model Database Baru
1. Buat file model di dalam `internal/model/notifikasi.go`.
2. Definisikan struktur tabel menggunakan tag GORM:
   ```go
   package model
   
   import (
       "time"
       "github.com/google/uuid"
   )
   
   type Notifikasi struct {
       ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
       UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
       Pesan     string    `gorm:"type:varchar(255);not null"`
       IsRead    bool      `gorm:"default:false;not null"`
       CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
   }
   ```
3. Tambahkan struct model tersebut ke daftar migrasi otomatis di [internal/database/db.go](file:///d:/NALAKARSA/nalakarsa-backend/internal/database/db.go):
   ```go
   err = db.AutoMigrate(
       &model.User{},
       // ...,
       &model.Notifikasi{}, // Daftarkan di sini
   )
   ```

### Langkah B: Membuat DTO (Data Transfer Object)
1. Buat file DTO baru atau tambahkan di file relevan `internal/dto/notifikasi.go` untuk request body dan response API:
   ```go
   package dto
   
   type CreateNotifikasiRequest struct {
       Pesan string `json:"pesan" binding:"required"`
   }
   ```

### Langkah C: Membuat Repository Layer
1. Buat interface dan implementasi database di `internal/repository/notifikasi_repository.go`:
   ```go
   package repository
   
   import (
       "nalakarsa/internal/model"
       "gorm.io/gorm"
   )
   
   type NotifikasiRepository interface {
       Create(notif *model.Notifikasi) error
   }
   
   type pgNotifikasiRepository struct {
       db *gorm.DB
   }
   
   func NewNotifikasiRepository(db *gorm.DB) NotifikasiRepository {
       return &pgNotifikasiRepository{db: db}
   }
   
   func (r *pgNotifikasiRepository) Create(notif *model.Notifikasi) error {
       return r.db.Create(notif).Error
   }
   ```

### Langkah D: Membuat Service Layer (Business Logic)
1. Buat interface dan implementasi di `internal/service/notifikasi_service.go`:
   ```go
   package service
   
   import (
       "nalakarsa/internal/dto"
       "nalakarsa/internal/model"
       "nalakarsa/internal/repository"
       "github.com/google/uuid"
   )
   
   type NotifikasiService interface {
       Kirim(userID uuid.UUID, req dto.CreateNotifikasiRequest) error
   }
   
   type notifikasiService struct {
       notifRepo repository.NotifikasiRepository
   }
   
   func NewNotifikasiService(repo repository.NotifikasiRepository) NotifikasiService {
       return &notifikasiService{notifRepo: repo}
   }
   
   func (s *notifikasiService) Kirim(userID uuid.UUID, req dto.CreateNotifikasiRequest) error {
       notif := &model.Notifikasi{
           UserID: userID,
           Pesan:  req.Pesan,
       }
       return s.notifRepo.Create(notif)
   }
   ```

### Langkah E: Membuat Handler Layer (HTTP Controller)
1. Buat di `internal/handler/notifikasi_handler.go`:
   ```go
   package handler
   
   import (
       "net/http"
       "nalakarsa/internal/dto"
       "nalakarsa/internal/service"
       "nalakarsa/internal/utils"
       "github.com/gin-gonic/gin"
       "github.com/google/uuid"
   )
   
   type NotifikasiHandler struct {
       service service.NotifikasiService
   }
   
   func NewNotifikasiHandler(s service.NotifikasiService) *NotifikasiHandler {
       return &NotifikasiHandler{service: s}
   }
   
   func (h *NotifikasiHandler) Kirim(c *gin.Context) {
       userID := c.MustGet("user_id").(uuid.UUID)
       var req dto.CreateNotifikasiRequest
       if err := c.ShouldBindJSON(&req); err != nil {
           utils.ErrorJSONResponse(c, http.StatusBadRequest, "Invalid input", []string{err.Error()})
           return
       }
       if err := h.service.Kirim(userID, req); err != nil {
           utils.ErrorJSONResponseWithMessage(c, http.StatusInternalServerError, err.Error())
           return
       }
       utils.JSONResponse(c, http.StatusCreated, "Notifikasi terkirim", nil, nil)
   }
   ```

### Langkah F: Hubungkan Rute & Daftarkan Dependensi
1. Daftarkan dan hubungkan di entrypoint [cmd/api/main.go](file:///d:/NALAKARSA/nalakarsa-backend/cmd/api/main.go):
   ```go
   // Repositories
   notifRepo := repository.NewNotifikasiRepository(db)
   // Services
   notifService := service.NewNotifikasiService(notifRepo)
   // Handlers
   notifHandler := handler.NewNotifikasiHandler(notifService)
   ```
2. Daftarkan path endpoint rute di [internal/routes/routes.go](file:///d:/NALAKARSA/nalakarsa-backend/internal/routes/routes.go):
   ```go
   protected.POST("/notifications", notifHandler.Kirim)
   ```

---

## 3. Best Practice & Coding Standard

* **Tangani Semua Error**: Jangan pernah mengabaikan return value error. Selalu log error menggunakan structured logging (`slog`) atau kembalikan ke client dengan response JSON standar.
* **Interface Segregation**: Selalu definisikan interface untuk layer Repository dan Service untuk mendukung decoupling dan mock testing.
* **Gunakan Validation Binding**: Selalu tambahkan tag `binding:"required"` (atau validasi custom) pada model input request DTO untuk menangkap error input sedini mungkin sebelum memproses database.
* **Gunakan Transaction untuk Multi-write**: Jika Anda menulis ke lebih dari satu tabel dalam satu alur bisnis (contoh: pembuatan `User` dan `Profile` saat registrasi), bungkus kueri tersebut dalam transaksi `db.Transaction(func(tx *gorm.DB) error { ... })` agar integritas data tetap terjaga.
