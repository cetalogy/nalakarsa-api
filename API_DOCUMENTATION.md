# Nalakarsa API Documentation untuk Frontend

Dokumen ini berisi panduan lengkap integrasi API Nalakarsa untuk kebutuhan Frontend. Semua endpoint mengembalikan format standar (kecuali jika terjadi error fatal di luar aplikasi).

## Base URL
Semua endpoint berawalan dengan: `/api/v1`

## Standard Response Format
Setiap response (berhasil atau gagal yang ditangani sistem) akan dibungkus dalam format berikut:

```json
{
  "data": { ... },       // Payload response (bisa berupa object atau array)
  "meta": {              // Meta informasi (opsional, biasanya untuk pagination)
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "limit": 10
  },
  "message": "Sukses"    // Pesan keterangan
}
```

## Error Response Format
Jika terjadi error (misalnya 400 Bad Request atau 404 Not Found), strukturnya akan seperti ini:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Data tidak valid",
    "fields": {
      "email": "Email harus berformat valid",
      "password": "Password minimal 8 karakter"
    }
  }
}
```

---

## 1. Authentication (Public)
Endpoint: `/auth`

### 1.1 Register
- **Endpoint**: `POST /auth/register`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "role": "Akademisi", // enum: "Akademisi", "Praktisi", "Profesional"
    "nama_lengkap": "John Doe",
    "gelar_depan": "Dr.",
    "gelar_belakang": "S.Kom",
    "afiliasi": "Universitas Nalakarsa",
    "lokasi": "Jakarta",
    "bidang_keahlian": "Software Engineering",
    "industry": "Technology",
    "bio": "Dosen dan Peneliti",
    "mission": "Mengembangkan ekosistem digital",
    "avatar_url": "https://..."
  }
  ```

### 1.2 Login
- **Endpoint**: `POST /auth/login`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response Data**:
  ```json
  {
    "access_token": "jwt-token-here",
    "refresh_token": "refresh-token-here",
    "expires_in": 3600,
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "role": "Akademisi"
    }
  }
  ```

### 1.3 Refresh Token
- **Endpoint**: `POST /auth/refresh`
- **Body**:
  ```json
  {
    "refresh_token": "your-refresh-token"
  }
  ```

### 1.4 Lupa Password
- **Endpoint**: `POST /auth/forgot-password`
- **Body**:
  ```json
  {
    "email": "user@example.com"
  }
  ```

### 1.5 Reset Password
- **Endpoint**: `POST /auth/reset-password`
- **Body**:
  ```json
  {
    "token": "reset-token-dari-email",
    "new_password": "newpassword123"
  }
  ```

### 1.6 Logout (Protected)
- **Endpoint**: `POST /auth/logout`
- **Headers**: `Authorization: Bearer <token>`

---

## 2. Users & Profile

### 2.1 List Users (Public)
- **Endpoint**: `GET /users`

### 2.2 Get Public Profile (Public)
- **Endpoint**: `GET /users/:id`

### 2.3 Get My Profile (Protected)
- **Endpoint**: `GET /users/me`
- **Headers**: `Authorization: Bearer <token>`

### 2.4 Update My Profile (Protected)
- **Endpoint**: `PATCH /users/me`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "nama_lengkap": "John Doe",
    "gelar_depan": "",
    "gelar_belakang": "",
    "afiliasi": "Univ A",
    "lokasi": "Bandung",
    "bidang_keahlian": "AI",
    "industry": "Edu",
    "bio": "...",
    "mission": "...",
    "avatar_url": ""
  }
  ```

### 2.5 Upload Avatar (Protected)
- **Endpoint**: `POST /users/me/avatar`
- **Headers**: `Authorization: Bearer <token>`, `Content-Type: multipart/form-data`
- **Body**: File upload via form-data dengan key `avatar` (atau disesuaikan dengan handler upload).

### 2.6 User Suggestions (Protected)
- **Endpoint**: `GET /users/suggestions`
- **Headers**: `Authorization: Bearer <token>`
- **Deskripsi**: Mendapatkan daftar pengguna yang direkomendasikan untuk koneksi.

---

## 3. Discussions / Forum

### 3.1 List Discussions (Public)
- **Endpoint**: `GET /discussions`

### 3.2 Detail Discussion (Public)
- **Endpoint**: `GET /discussions/:id` (Termasuk list replies)

### 3.3 Create Discussion (Protected)
- **Endpoint**: `POST /discussions`
- **Body**:
  ```json
  {
    "title": "Judul Diskusi",
    "content": "Isi detail diskusi...",
    "category": "Technology",
    "tags": "golang, react, api"
  }
  ```

### 3.4 Update Discussion (Protected)
- **Endpoint**: `PATCH /discussions/:id`
- **Body**: Sama seperti create, ditambah opsi `"status": "open|resolved|closed"`.

### 3.5 Delete Discussion (Protected)
- **Endpoint**: `DELETE /discussions/:id`

### 3.6 Add Reply (Protected)
- **Endpoint**: `POST /discussions/:id/replies`
- **Body**:
  ```json
  {
    "content": "Ini balasan saya",
    "parent_id": "uuid-optional-jika-balas-komentar-lain"
  }
  ```

### 3.7 Delete Reply (Protected)
- **Endpoint**: `DELETE /discussions/:id/replies/:reply_id`

### 3.8 Upvote Discussion (Protected)
- **Endpoint**: `POST /discussions/:id/votes`

### 3.9 Unvote Discussion (Protected)
- **Endpoint**: `DELETE /discussions/:id/votes`

---

## 4. Connections (Protected)

### 4.1 List My Connections
- **Endpoint**: `GET /connections`

### 4.2 List Connection Requests (Pending)
- **Endpoint**: `GET /connections/requests`

### 4.3 Send Connection Request
- **Endpoint**: `POST /connections/requests`
- **Body**:
  ```json
  {
    "target_user_id": "uuid-user-tujuan"
  }
  ```

### 4.4 Accept Request
- **Endpoint**: `PATCH /connections/requests/:id/accept`

### 4.5 Reject Request
- **Endpoint**: `PATCH /connections/requests/:id/reject`

### 4.6 Cancel Request (Batalkan request yg sudah dikirim)
- **Endpoint**: `DELETE /connections/requests/:id`

### 4.7 Remove Connection (Hapus dari daftar teman)
- **Endpoint**: `DELETE /connections/:userId`

---

## 5. Projects (Kolaborasi)

### 5.1 List Projects (Public)
- **Endpoint**: `GET /projects`

### 5.2 Project Detail (Public)
- **Endpoint**: `GET /projects/:id`

### 5.3 List Project Members (Public)
- **Endpoint**: `GET /projects/:id/members`

### 5.4 Create Project (Protected)
- **Endpoint**: `POST /projects`
- **Body**:
  ```json
  {
    "title": "Nama Project",
    "description": "Deskripsi project minimal 15 karakter",
    "category": "Penelitian",
    "role_required": "Praktisi", // Opsional
    "funding_status": "Didanai",
    "location": "Jakarta",
    "deadline": "2026-12-31T00:00:00Z" // Opsional
  }
  ```

### 5.5 Update Project (Protected)
- **Endpoint**: `PATCH /projects/:id`
- **Body**: Sama seperti create, ditambah `"status": "draft|open|in_review|active|completed|archived"` dan `"progress": 0-100`.

### 5.6 Delete Project (Protected)
- **Endpoint**: `DELETE /projects/:id`

### 5.7 Apply to Project (Protected)
- **Endpoint**: `POST /projects/:id/applications`
- **Body**:
  ```json
  {
    "message": "Saya tertarik bergabung karena..."
  }
  ```

### 5.8 List Applications in Project (Protected - For Owner)
- **Endpoint**: `GET /projects/:id/applications`

### 5.9 Accept/Reject Application (Protected - For Owner)
- **Endpoint**: `PATCH /projects/:id/applications/:applicationId`
- **Body**:
  ```json
  {
    "status": "accepted" // atau "rejected"
  }
  ```

### 5.10 Create Milestone (Protected)
- **Endpoint**: `POST /projects/:id/milestones`
- **Body**:
  ```json
  {
    "title": "Fase 1 Selesai",
    "due_at": "2026-10-10T00:00:00Z", // Opsional
    "assignee_id": "uuid-member" // Opsional
  }
  ```

### 5.11 Update Milestone (Protected)
- **Endpoint**: `PATCH /projects/:id/milestones/:milestoneId`
- **Body**:
  ```json
  {
    "title": "Fase 1 Selesai",
    "due_at": "2026-10-10T00:00:00Z",
    "status": "pending|in_progress|completed",
    "assignee_id": "uuid-member"
  }
  ```

---

## 6. Conversations (Chat)

### 6.1 List Conversations
- **Endpoint**: `GET /conversations`

### 6.2 Get or Create Direct Chat
- **Endpoint**: `POST /conversations/direct`
- **Body**:
  ```json
  {
    "target_user_id": "uuid-lawan-bicara"
  }
  ```

### 6.3 List Messages in Conversation
- **Endpoint**: `GET /conversations/:id/messages`
- **Note**: Support cursor pagination. Response Meta akan mengandung `NextCursor` dan `HasMore`.

### 6.4 Send Message
- **Endpoint**: `POST /conversations/:id/messages`
- **Body**:
  ```json
  {
    "body": "Halo, salam kenal!"
  }
  ```

### 6.5 Mark Conversation as Read
- **Endpoint**: `PATCH /conversations/:id/read`

---

## 7. Notifications (Protected)

### 7.1 List Notifications
- **Endpoint**: `GET /notifications`

### 7.2 Mark Notification as Read
- **Endpoint**: `PATCH /notifications/:id/read`

### 7.3 Mark All as Read
- **Endpoint**: `PATCH /notifications/read-all`

---

## 8. Dashboard (Protected)

### 8.1 Get Dashboard Data
- **Endpoint**: `GET /dashboard`
- **Response Data**:
  ```json
  {
    "active_projects_count": 2,
    "unread_messages_count": 5,
    "active_connections": 10,
    "my_projects": [
      {
        "id": "uuid",
        "title": "Project A",
        "role": "Owner",
        "progress": 50,
        "status": "active",
        "next_milestone": "Fase 2",
        "deadline": "2026-12-31T00:00:00Z"
      }
    ],
    "recent_activity": [
      // Format sama seperti List Notifications
    ]
  }
  ```

---

## Panduan Umum Headers
Untuk setiap endpoint yang bertanda **(Protected)**, Frontend diwajibkan mengirimkan token otentikasi di HTTP Headers:
```http
Authorization: Bearer <access_token>
```
Jika `access_token` expired (biasanya melempar status `401 Unauthorized`), FE bisa melakukan silent refresh ke endpoint `POST /auth/refresh` dengan mengirimkan `refresh_token`, lalu ulangi request aslinya dengan `access_token` yang baru.
