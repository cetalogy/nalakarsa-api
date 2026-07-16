# API Reference - Nalakarsa

Dokumen ini menjelaskan spesifikasi API RESTful untuk **Nalakarsa** platform backend.

---

## Format Respon Standar (Standard Response Format)

Semua respon API mengembalikan format JSON yang konsisten.

### Respon Sukses (Success Response)
```json
{
  "success": true,
  "message": "Pesan sukses detail",
  "data": {} // Bisa berupa object atau array
}
```

### Respon Error (Error Response)
```json
{
  "success": false,
  "message": "Deskripsi error utama",
  "errors": [
    "Field email wajib diisi",
    "Field password minimal berisi 6 karakter"
  ] // Detail error validasi (jika ada)
}
```

---

## 1. Authentication Endpoints

### 1.1 Register
Mendaftarkan akun baru sekaligus membuat data profil.
* **URL**: `/api/v1/auth/register`
* **Method**: `POST`
* **Authentication**: None (Public)
* **Request Body (JSON)**:
  * `email` (string, required, email format)
  * `password` (string, required, min 6 characters)
  * `role` (string, required, value: `akademisi`, `praktisi`, `profesional`)
  * `nama_lengkap` (string, required)
  * `gelar_depan` (string, optional)
  * `gelar_belakang` (string, optional)
  * `afiliasi` (string, required)
  * `lokasi` (string, required)
  * `bidang_keahlian` (string, required)
  * `bio_misi` (string, optional)
* **Response (201 Created)**:
  ```json
  {
    "success": true,
    "message": "Registrasi berhasil",
    "data": {
      "id": "e30b6510-1ab2-4b2b-bbd8-4903c73dc2d5",
      "email": "dosen@nalakarsa.id",
      "role": "akademisi"
    }
  }
  ```

### 1.2 Login
Melakukan login untuk mendapatkan Access Token dan Refresh Token.
* **URL**: `/api/v1/auth/login`
* **Method**: `POST`
* **Authentication**: None (Public)
* **Request Body (JSON)**:
  * `email` (string, required)
  * `password` (string, required)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Login berhasil",
    "data": {
      "access_token": "eyJhbGciOi...",
      "refresh_token": "eyJhbGciOi...",
      "expires_in": 3600,
      "user": {
        "id": "e30b6510-1ab2-4b2b-bbd8-4903c73dc2d5",
        "email": "dosen@nalakarsa.id",
        "role": "akademisi"
      }
    }
  }
  ```

### 1.3 Refresh Token
Memperpanjang masa aktif sesi menggunakan Refresh Token.
* **URL**: `/api/v1/auth/refresh`
* **Method**: `POST`
* **Authentication**: None (Public - mengirimkan refresh token di body)
* **Request Body (JSON)**:
  * `refresh_token` (string, required)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Token berhasil diperbarui",
    "data": {
      "access_token": "eyJhbGciOi...",
      "refresh_token": "eyJhbGciOi...",
      "expires_in": 3600
    }
  }
  ```

### 1.4 Logout
Mencabut Refresh Token dari server.
* **URL**: `/api/v1/auth/logout`
* **Method**: `POST`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `refresh_token` (string, required)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Logout berhasil"
  }
  ```

---

## 2. User & Profile Endpoints

### 2.1 Get Current User Profile
Mendapatkan data lengkap profil sendiri setelah terautentikasi.
* **URL**: `/api/v1/users/me`
* **Method**: `GET`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Profil berhasil dimuat",
    "data": {
      "id": "e30b6510-1ab2-4b2b-bbd8-4903c73dc2d5",
      "email": "dosen@nalakarsa.id",
      "role": "akademisi",
      "profile": {
        "nama_lengkap": "Dr. Ir. Budi Santoso",
        "gelar_depan": "Dr.",
        "gelar_belakang": "M.T.",
        "afiliasi": "Universitas Indonesia",
        "lokasi": "Depok, Indonesia",
        "bidang_keahlian": "Kecerdasan Buatan & IoT",
        "bio_misi": "Menghubungkan riset akademik ke industri nyata.",
        "avatar_url": ""
      }
    }
  }
  ```

### 2.2 Update Current User Profile
Memperbarui informasi profil pribadi.
* **URL**: `/api/v1/users/me`
* **Method**: `PUT`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `nama_lengkap` (string, required)
  * `gelar_depan` (string, optional)
  * `gelar_belakang` (string, optional)
  * `afiliasi` (string, required)
  * `lokasi` (string, required)
  * `bidang_keahlian` (string, required)
  * `bio_misi` (string, optional)
  * `avatar_url` (string, optional)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Profil berhasil diperbarui",
    "data": {
      "nama_lengkap": "Dr. Ir. Budi Santoso, M.T.",
      "afiliasi": "Universitas Indonesia",
      "lokasi": "Jakarta, Indonesia"
    }
  }
  ```

### 2.3 List Users (Directory)
Mendapatkan daftar profil seluruh pengguna dengan pencarian dan filter peran.
* **URL**: `/api/v1/users`
* **Method**: `GET`
* **Authentication**: None (Public)
* **Query Parameters**:
  * `search` (string, optional - cari nama/afiliasi/keahlian)
  * `role` (string, optional - filter `akademisi`, `praktisi`, `profesional`)
  * `page` (int, default `1`)
  * `limit` (int, default `10`)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Daftar pengguna berhasil dimuat",
    "data": [
      {
        "id": "e30b6510-1ab2-4b2b-bbd8-4903c73dc2d5",
        "email": "dosen@nalakarsa.id",
        "role": "akademisi",
        "profile": {
          "nama_lengkap": "Dr. Ir. Budi Santoso",
          "afiliasi": "Universitas Indonesia",
          "bidang_keahlian": "IoT"
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 48,
      "limit": 10
    }
  }
  ```

---

## 3. Discussion Endpoints

### 3.1 List Discussions
* **URL**: `/api/v1/discussions`
* **Method**: `GET`
* **Authentication**: None (Public)
* **Query Parameters**:
  * `search` (string, optional - cari judul/isi)
  * `category` (string, optional - filter kategori)
  * `role` (string, optional - filter berdasarkan peran pembuat)
  * `sort` (string, default `newest`, value: `newest`, `oldest`)
  * `page` (int, default `1`)
  * `limit` (int, default `10`)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Daftar diskusi berhasil dimuat",
    "data": [
      {
        "id": "99e1966a-11bc-4e5c-9c98-468df9f201aa",
        "title": "Arsitektur Microservices di Go",
        "content": "Menurut rekan-rekan bagaimanakah performa...",
        "category": "Tech",
        "tags": "golang,microservices",
        "created_at": "2026-07-13T20:30:00Z",
        "creator": {
          "id": "a1b2c3d4...",
          "nama_lengkap": "Setyo Nugroho",
          "role": "profesional"
        }
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 2,
      "total_items": 15
    }
  }
  ```

### 3.2 Get Discussion Detail with Comments
* **URL**: `/api/v1/discussions/:id`
* **Method**: `GET`
* **Authentication**: None (Public)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Detail diskusi berhasil dimuat",
    "data": {
      "id": "99e1966a-11bc-4e5c-9c98-468df9f201aa",
      "title": "Arsitektur Microservices di Go",
      "content": "Menurut rekan-rekan bagaimanakah performa...",
      "category": "Tech",
      "tags": "golang,microservices",
      "created_at": "2026-07-13T20:30:00Z",
      "creator": {
        "id": "a1b2c3d4...",
        "nama_lengkap": "Setyo Nugroho",
        "role": "profesional"
      },
      "comments": [
        {
          "id": "f5f5f5f5...",
          "content": "Sangat bagus, pastikan menggunakan gRPC untuk komunikasi inter-service.",
          "created_at": "2026-07-13T20:35:00Z",
          "creator": {
            "id": "e30b6510...",
            "nama_lengkap": "Dr. Ir. Budi Santoso",
            "role": "akademisi"
          }
        }
      ]
    }
  }
  ```

### 3.3 Create Discussion Topic
* **URL**: `/api/v1/discussions`
* **Method**: `POST`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `title` (string, required, min 5 chars)
  * `content` (string, required, min 10 chars)
  * `category` (string, required)
  * `tags` (string, optional - comma separated)
* **Response (201 Created)**:
  ```json
  {
    "success": true,
    "message": "Diskusi berhasil dibuat",
    "data": {
      "id": "99e1966a-11bc-4e5c-9c98-468df9f201aa"
    }
  }
  ```

### 3.4 Update Discussion Topic
* **URL**: `/api/v1/discussions/:id`
* **Method**: `PUT`
* **Authentication**: JWT Bearer Token (Authenticated - Owner Only)
* **Request Body (JSON)**:
  * `title` (string, required)
  * `content` (string, required)
  * `category` (string, required)
  * `tags` (string, optional)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Diskusi berhasil diperbarui"
  }
  ```

### 3.5 Delete Discussion Topic
* **URL**: `/api/v1/discussions/:id`
* **Method**: `DELETE`
* **Authentication**: JWT Bearer Token (Authenticated - Owner Only)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Diskusi berhasil dihapus"
  }
  ```

---

## 4. Comment Endpoints

### 4.1 Add Comment to Discussion
* **URL**: `/api/v1/discussions/:id/comments`
* **Method**: `POST`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `content` (string, required, min 3 chars)
* **Response (201 Created)**:
  ```json
  {
    "success": true,
    "message": "Komentar berhasil ditambahkan",
    "data": {
      "id": "f5f5f5f5-8cb6-4d56-b09e-7164a6d91244",
      "content": "Sangat bagus, pastikan menggunakan gRPC...",
      "created_at": "2026-07-13T20:35:00Z"
    }
  }
  ```

### 4.2 Delete Comment
* **URL**: `/api/v1/discussions/:id/comments/:comment_id`
* **Method**: `DELETE`
* **Authentication**: JWT Bearer Token (Authenticated - Owner Only)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Komentar berhasil dihapus"
  }
  ```

---

## 5. Collaboration Endpoints

### 5.1 List Collaborations
* **URL**: `/api/v1/collaborations`
* **Method**: `GET`
* **Authentication**: None (Public)
* **Query Parameters**:
  * `search` (string, optional - cari judul/isi)
  * `role_required` (string, optional - filter `akademisi`, `praktisi`, `profesional`)
  * `status` (string, optional - filter `open`, `in_progress`, `closed`)
  * `page` (int, default `1`)
  * `limit` (int, default `10`)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Daftar kolaborasi berhasil dimuat",
    "data": [
      {
        "id": "cb104278-df35-430b-8d07-2a54b38d3810",
        "title": "Uji Coba Lapangan Sensor IoT Pertanian",
        "description": "Dibutuhkan praktisi pertanian untuk...",
        "role_required": "praktisi",
        "status": "open",
        "created_at": "2026-07-13T20:40:00Z",
        "owner": {
          "id": "e30b6510...",
          "nama_lengkap": "Dr. Ir. Budi Santoso",
          "role": "akademisi"
        }
      }
    ]
  }
  ```

### 5.2 Create Collaboration Proposal
* **URL**: `/api/v1/collaborations`
* **Method**: `POST`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `title` (string, required, min 5 chars)
  * `description` (string, required, min 15 chars)
  * `role_required` (string, required, value: `akademisi`, `praktisi`, `profesional`)
* **Response (201 Created)**:
  ```json
  {
    "success": true,
    "message": "Proposal kolaborasi berhasil dibuat",
    "data": {
      "id": "cb104278-df35-430b-8d07-2a54b38d3810"
    }
  }
  ```

### 5.3 Apply to Collaboration (Submit Interest)
* **URL**: `/api/v1/collaborations/:id/apply`
* **Method**: `POST`
* **Authentication**: JWT Bearer Token (Authenticated)
* **Request Body (JSON)**:
  * `message` (string, required, min 10 chars)
* **Response (201 Created)**:
  ```json
  {
    "success": true,
    "message": "Lamaran kolaborasi berhasil dikirim",
    "data": {
      "id": "a9a9a9a9-468c-4a30-9b0d-b108cc681a2e"
    }
  }
  ```

### 5.4 Get Collaboration Applicants (Owner Only)
Melihat daftar pelamar untuk proyek kolaborasi yang dibuat oleh pengguna aktif.
* **URL**: `/api/v1/collaborations/:id/applications`
* **Method**: `GET`
* **Authentication**: JWT Bearer Token (Authenticated - Owner of project only)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Daftar pelamar berhasil dimuat",
    "data": [
      {
        "id": "a9a9a9a9-468c-4a30-9b0d-b108cc681a2e",
        "message": "Saya memiliki lahan uji coba 2 hektar di Bogor dan tertarik...",
        "status": "pending",
        "created_at": "2026-07-13T20:50:00Z",
        "applicant": {
          "id": "c7c7c7c7...",
          "nama_lengkap": "Hendra Wijaya",
          "role": "praktisi",
          "afiliasi": "PT Tani Maju"
        }
      }
    ]
  }
  ```

### 5.5 Update Application Status (Owner Only)
Menerima atau menolak lamaran kolaborasi.
* **URL**: `/api/v1/collaborations/:id/applications/:app_id`
* **Method**: `PUT`
* **Authentication**: JWT Bearer Token (Authenticated - Owner of project only)
* **Request Body (JSON)**:
  * `status` (string, required, value: `accepted`, `rejected`)
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Status lamaran kolaborasi berhasil diperbarui"
  }
  ```
