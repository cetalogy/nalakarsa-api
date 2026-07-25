# Kebutuhan Backend Nalakarsa

Dokumen ini merangkum kebutuhan backend berdasarkan progres frontend Nalakarsa saat ini. Tujuannya adalah menjadi checklist implementasi dan acuan kontrak antara tim Frontend (FE) dan Backend (BE).

## 1. Kondisi Saat Ini

Frontend sudah memiliki halaman dan alur berikut:

- Registrasi, login, logout, dan edit profil.
- Daftar dan pembuatan topik diskusi.
- Daftar peluang/proyek kolaborasi dan aksi menghubungi inisiator.
- Dasbor ringkasan proyek, pesan, jaringan, dan notifikasi.
- Jaringan pengguna: koneksi, permintaan masuk, rekomendasi, terima, abaikan, dan kirim permintaan.
- Chat satu lawan satu dan memulai chat dari halaman kolaborasi.
- Halaman profil dengan statistik koneksi, proyek, dan diskusi.

Saat ini seluruh data masih berasal dari `localStorage`, konstanta dummy, atau state Pinia. Belum tersedia API, autentikasi server, database, upload file, maupun komunikasi realtime.

## 2. Keputusan Teknis yang Harus Disiapkan

- [ ] Tentukan stack backend. Rekomendasi awal: Node.js + TypeScript dengan NestJS/Express/Fastify.
- [ ] Tentukan database relasional. Rekomendasi: PostgreSQL.
- [ ] Tentukan ORM dan sistem migrasi. Contoh: Prisma atau Drizzle.
- [ ] Tentukan metode autentikasi: access token pendek + refresh token aman, atau session berbasis cookie.
- [ ] Tentukan penyimpanan avatar dan lampiran: object storage kompatibel S3.
- [ ] Tentukan realtime chat/notifikasi: WebSocket atau Socket.IO.
- [ ] Tentukan layanan email untuk verifikasi akun dan reset password.
- [ ] Siapkan environment `development`, `staging`, dan `production`.
- [ ] Sepakati format respons, error, pagination, filter, dan penamaan field dengan FE.
- [ ] Buat dokumentasi OpenAPI/Swagger sebagai sumber kontrak API.

## 3. Prioritas Implementasi

### P0 - Fondasi Backend

- [ ] Buat repository atau folder backend sesuai keputusan arsitektur proyek.
- [ ] Siapkan konfigurasi environment dan validasi variabel environment saat aplikasi dimulai.
- [ ] Siapkan koneksi database, migration, dan seed data untuk akun/demo UI.
- [ ] Buat struktur modul: `auth`, `users`, `profiles`, `discussions`, `projects`, `connections`, `conversations`, `notifications`.
- [ ] Terapkan validasi request pada seluruh endpoint.
- [ ] Terapkan middleware autentikasi dan otorisasi berbasis pemilik data/anggota proyek.
- [ ] Terapkan format respons API yang konsisten.
- [ ] Terapkan global error handler tanpa membocorkan stack trace atau data sensitif.
- [ ] Tambahkan logging terstruktur dengan request/correlation ID.
- [ ] Konfigurasikan CORS hanya untuk origin FE yang diizinkan.
- [ ] Tambahkan endpoint `GET /health` untuk health check aplikasi dan database.
- [ ] Buat spesifikasi OpenAPI dan contoh request/response.

Format respons yang disarankan:

```json
{
  "data": {},
  "meta": {},
  "message": "Operasi berhasil"
}
```

Format error yang disarankan:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Data tidak valid",
    "fields": {
      "email": "Format email tidak valid"
    }
  }
}
```

### P1 - Autentikasi dan Profil

- [ ] Registrasi dengan `fullName`, `email`, `password`, `affiliation`, dan `role`.
- [ ] Batasi nilai `role` menjadi `Akademisi`, `Praktisi`, atau `Profesional`.
- [ ] Pastikan email unik dan dinormalisasi menjadi huruf kecil.
- [ ] Hash password menggunakan Argon2id atau bcrypt dengan konfigurasi yang aman.
- [ ] Login menggunakan email dan password.
- [ ] Buat endpoint refresh session/token.
- [ ] Logout dan invalidasi refresh token/session.
- [ ] Endpoint untuk mengambil profil pengguna yang sedang login.
- [ ] Endpoint untuk memperbarui nama, afiliasi, peran, bio, lokasi, bidang industri, dan misi pengembangan.
- [ ] Upload, validasi tipe/ukuran, dan penggantian avatar.
- [ ] Profil publik berdasarkan ID/slug pengguna.
- [ ] Statistik profil: jumlah koneksi, proyek, diskusi, pengikut, mengikuti, dan jumlah dilihat.
- [ ] Tambahkan verifikasi email.
- [ ] Tambahkan alur lupa dan reset password.
- [ ] Tambahkan rate limit pada login, register, reset password, dan verifikasi email.

Endpoint minimal:

| Method | Endpoint | Fungsi | Akses |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | Membuat akun | Publik |
| `POST` | `/api/v1/auth/login` | Login | Publik |
| `POST` | `/api/v1/auth/refresh` | Memperbarui sesi/token | Publik dengan refresh token |
| `POST` | `/api/v1/auth/logout` | Mengakhiri sesi | Login |
| `POST` | `/api/v1/auth/forgot-password` | Meminta reset password | Publik |
| `POST` | `/api/v1/auth/reset-password` | Mengganti password | Publik dengan token |
| `GET` | `/api/v1/users/me` | Profil pengguna aktif | Login |
| `PATCH` | `/api/v1/users/me` | Edit profil | Login |
| `POST` | `/api/v1/users/me/avatar` | Upload avatar | Login |
| `GET` | `/api/v1/users/:id` | Profil publik | Publik/Login sesuai kebijakan |

### P1 - Diskusi

- [ ] Mengambil daftar diskusi dengan pagination.
- [ ] Mendukung pencarian berdasarkan judul dan nama penulis.
- [ ] Mendukung filter kategori dan peran penulis.
- [ ] Membuat topik diskusi baru.
- [ ] Mengambil detail diskusi.
- [ ] Memperbarui dan menghapus diskusi oleh pemilik atau moderator.
- [ ] Menambahkan, mengubah, dan menghapus balasan.
- [ ] Upvote/unvote diskusi dan balasan dengan aturan satu vote per pengguna.
- [ ] Mengembalikan jumlah balasan dan upvote secara efisien.
- [ ] Menyediakan status topik seperti `open`, `resolved`, dan `closed` bila dibutuhkan.
- [ ] Membuat notifikasi saat diskusi dibalas atau menerima upvote.

Endpoint minimal:

| Method | Endpoint | Fungsi | Akses |
|---|---|---|---|
| `GET` | `/api/v1/discussions` | Daftar, cari, filter, pagination | Publik |
| `POST` | `/api/v1/discussions` | Membuat topik | Login |
| `GET` | `/api/v1/discussions/:id` | Detail topik dan ringkasan balasan | Publik |
| `PATCH` | `/api/v1/discussions/:id` | Edit topik | Pemilik/Moderator |
| `DELETE` | `/api/v1/discussions/:id` | Hapus topik | Pemilik/Moderator |
| `GET` | `/api/v1/discussions/:id/replies` | Daftar balasan | Publik |
| `POST` | `/api/v1/discussions/:id/replies` | Membalas topik | Login |
| `POST` | `/api/v1/discussions/:id/votes` | Upvote | Login |
| `DELETE` | `/api/v1/discussions/:id/votes` | Membatalkan upvote | Login |

Query daftar yang disarankan:

```text
GET /api/v1/discussions?q=iot&category=Riset&role=Akademisi&page=1&limit=20&sort=latest
```

### P1 - Jaringan Pengguna

- [ ] Daftar koneksi pengguna.
- [ ] Daftar permintaan koneksi masuk dan terkirim.
- [ ] Kirim permintaan koneksi.
- [ ] Terima atau abaikan/tolak permintaan koneksi.
- [ ] Batalkan permintaan yang telah dikirim.
- [ ] Hapus koneksi.
- [ ] Cegah koneksi ke akun sendiri dan permintaan duplikat.
- [ ] Hitung koneksi mutual.
- [ ] Buat rekomendasi pengguna berdasarkan peran, afiliasi, bidang, dan koneksi mutual.
- [ ] Sediakan pencarian pengguna untuk menemukan calon mitra.
- [ ] Buat notifikasi untuk permintaan baru dan permintaan yang diterima.

Endpoint minimal:

| Method | Endpoint | Fungsi |
|---|---|---|
| `GET` | `/api/v1/connections` | Daftar koneksi saya |
| `GET` | `/api/v1/connections/requests?type=incoming` | Permintaan masuk |
| `GET` | `/api/v1/connections/requests?type=outgoing` | Permintaan terkirim |
| `POST` | `/api/v1/connections/requests` | Kirim permintaan dengan `targetUserId` |
| `PATCH` | `/api/v1/connections/requests/:id/accept` | Terima permintaan |
| `PATCH` | `/api/v1/connections/requests/:id/reject` | Tolak/abaikan permintaan |
| `DELETE` | `/api/v1/connections/requests/:id` | Batalkan permintaan |
| `DELETE` | `/api/v1/connections/:userId` | Hapus koneksi |
| `GET` | `/api/v1/users/suggestions` | Rekomendasi koneksi |
| `GET` | `/api/v1/users?q=nama&role=Akademisi` | Cari pengguna |

### P2 - Kolaborasi dan Manajemen Proyek

- [ ] Daftar proyek/peluang kolaborasi dengan pagination, pencarian, filter status, dan kategori.
- [ ] Detail proyek yang dapat dibagikan melalui URL.
- [ ] Membuat, mengubah, mempublikasikan, dan mengarsipkan proyek.
- [ ] Tentukan lifecycle status, misalnya `draft`, `open`, `in_review`, `active`, `completed`, `archived`.
- [ ] Simpan inisiator, organisasi, deskripsi, kebutuhan peran, kategori, lokasi, tenggat, dan status pendanaan.
- [ ] Pengguna dapat mengajukan minat/partisipasi pada proyek.
- [ ] Inisiator dapat menerima atau menolak pengajuan anggota.
- [ ] Kelola anggota proyek dan perannya.
- [ ] Simpan progres proyek dalam rentang 0-100.
- [ ] Kelola milestone/target, tenggat, status, dan penanggung jawab.
- [ ] Upload dokumen proposal/lampiran dengan kontrol akses.
- [ ] Tombol "Hubungi Inisiator" membuat atau membuka conversation dengan inisiator, bukan mencari berdasarkan nama.
- [ ] Buat notifikasi untuk pengajuan, keputusan pengajuan, perubahan status, dan tenggat.

Endpoint minimal:

| Method | Endpoint | Fungsi |
|---|---|---|
| `GET` | `/api/v1/projects` | Daftar proyek/peluang |
| `POST` | `/api/v1/projects` | Membuat proyek |
| `GET` | `/api/v1/projects/:id` | Detail proyek |
| `PATCH` | `/api/v1/projects/:id` | Memperbarui proyek |
| `DELETE` | `/api/v1/projects/:id` | Arsip/hapus proyek sesuai kebijakan |
| `POST` | `/api/v1/projects/:id/applications` | Mengajukan partisipasi |
| `GET` | `/api/v1/projects/:id/applications` | Daftar pengajuan untuk inisiator |
| `PATCH` | `/api/v1/projects/:id/applications/:applicationId` | Terima/tolak pengajuan |
| `GET` | `/api/v1/projects/:id/members` | Daftar anggota |
| `POST` | `/api/v1/projects/:id/milestones` | Membuat milestone |
| `PATCH` | `/api/v1/projects/:id/milestones/:milestoneId` | Memperbarui milestone |

### P2 - Chat Realtime

- [ ] Membuat atau mengambil percakapan langsung berdasarkan pasangan `userId`, bukan nama pengguna.
- [ ] Daftar percakapan diurutkan berdasarkan pesan terakhir.
- [ ] Mengambil pesan dengan cursor pagination.
- [ ] Mengirim pesan dan menyimpannya secara persisten.
- [ ] Status pesan: `sent`, `delivered`, dan `read`.
- [ ] Jumlah pesan belum dibaca per percakapan dan total.
- [ ] Event realtime untuk pesan baru, status dibaca, dan pembaruan daftar percakapan.
- [ ] Validasi bahwa pengirim adalah anggota percakapan.
- [ ] Batasi panjang pesan dan lakukan sanitasi konten.
- [ ] Siapkan dukungan lampiran sebagai fase lanjutan.
- [ ] Tentukan kebijakan edit/hapus pesan dan retensi data.

Endpoint minimal:

| Method | Endpoint | Fungsi |
|---|---|---|
| `GET` | `/api/v1/conversations` | Daftar percakapan |
| `POST` | `/api/v1/conversations/direct` | Membuat/mengambil chat dengan `targetUserId` |
| `GET` | `/api/v1/conversations/:id/messages` | Riwayat pesan dengan cursor |
| `POST` | `/api/v1/conversations/:id/messages` | Mengirim pesan |
| `PATCH` | `/api/v1/conversations/:id/read` | Tandai pesan telah dibaca |

Event realtime yang disarankan:

```text
message.created
message.delivered
message.read
conversation.updated
notification.created
```

### P2 - Dasbor dan Notifikasi

- [ ] Endpoint agregasi dasbor agar FE tidak perlu memanggil banyak endpoint hanya untuk kartu ringkasan.
- [ ] Ringkasan jumlah proyek aktif, pesan belum dibaca, dan koneksi aktif.
- [ ] Daftar proyek berjalan milik pengguna beserta peran, progres, milestone berikutnya, dan tenggat.
- [ ] Daftar aktivitas/notifikasi terbaru.
- [ ] Tandai satu notifikasi sebagai dibaca.
- [ ] Tandai seluruh notifikasi sebagai dibaca.
- [ ] Link/metadata notifikasi mengarah ke resource terkait.
- [ ] Push notifikasi baru melalui koneksi realtime.

Endpoint minimal:

| Method | Endpoint | Fungsi |
|---|---|---|
| `GET` | `/api/v1/dashboard` | Statistik dan ringkasan pengguna |
| `GET` | `/api/v1/notifications` | Daftar notifikasi |
| `PATCH` | `/api/v1/notifications/:id/read` | Tandai satu sebagai dibaca |
| `PATCH` | `/api/v1/notifications/read-all` | Tandai semua sebagai dibaca |

### P3 - Moderasi dan Administrasi

- [ ] Role sistem terpisah dari role profesi: `user`, `moderator`, dan `admin`.
- [ ] Pelaporan pengguna, diskusi, balasan, proyek, dan pesan bila masuk ruang lingkup moderasi.
- [ ] Antrean laporan untuk moderator.
- [ ] Suspend/aktifkan akun dengan audit trail.
- [ ] Moderasi atau arsip konten tanpa hard delete langsung.
- [ ] Audit log untuk tindakan admin dan perubahan data sensitif.
- [ ] Dashboard metrik dasar: pengguna aktif, diskusi, proyek, koneksi, dan laporan.

## 4. Rancangan Entitas Database Minimum

| Entitas | Field penting |
|---|---|
| `users` | `id`, `email`, `password_hash`, `system_role`, `email_verified_at`, `status`, timestamps |
| `profiles` | `user_id`, `full_name`, `professional_role`, `affiliation`, `bio`, `location`, `industry`, `mission`, `avatar_url`, `view_count` |
| `refresh_tokens` atau `sessions` | `id`, `user_id`, `token_hash`, `expires_at`, `revoked_at`, metadata perangkat |
| `discussions` | `id`, `author_id`, `title`, `body`, `category`, `status`, timestamps, `deleted_at` |
| `discussion_replies` | `id`, `discussion_id`, `author_id`, `parent_id` opsional, `body`, timestamps, `deleted_at` |
| `discussion_votes` | `user_id`, `discussion_id`/`reply_id`, `created_at`, unique constraint |
| `projects` | `id`, `owner_id`, `title`, `description`, `category`, `status`, `funding_status`, `location`, `deadline`, `progress`, timestamps |
| `project_members` | `project_id`, `user_id`, `role`, `status`, `joined_at` |
| `project_applications` | `id`, `project_id`, `applicant_id`, `message`, `status`, timestamps |
| `project_milestones` | `id`, `project_id`, `title`, `due_at`, `status`, `assignee_id`, `completed_at` |
| `connections` | `id`, `requester_id`, `addressee_id`, `status`, timestamps |
| `conversations` | `id`, `type`, `last_message_at`, timestamps |
| `conversation_members` | `conversation_id`, `user_id`, `last_read_message_id`, `joined_at` |
| `messages` | `id`, `conversation_id`, `sender_id`, `body`, `status`, timestamps, `deleted_at` |
| `notifications` | `id`, `user_id`, `type`, `actor_id`, `resource_type`, `resource_id`, `payload`, `read_at`, `created_at` |
| `attachments` | `id`, `owner_id`, `resource_type`, `resource_id`, `storage_key`, `mime_type`, `size`, timestamps |

Catatan constraint penting:

- [ ] Gunakan UUID/ULID atau strategi ID yang konsisten.
- [ ] Unique index pada email yang sudah dinormalisasi.
- [ ] Unique constraint untuk satu koneksi per pasangan pengguna.
- [ ] Unique constraint untuk direct conversation per pasangan pengguna.
- [ ] Index pada foreign key, `created_at`, status, kategori, dan kolom pencarian/filter utama.
- [ ] Gunakan soft delete untuk konten yang membutuhkan audit/moderasi.
- [ ] Semua timestamp disimpan dalam UTC dan dikonversi di FE.

## 5. Kontrak Data yang Dibutuhkan Frontend

Jangan mengirim nilai waktu siap tampil seperti `"2 jam yang lalu"` dari database. Kirim timestamp ISO agar FE dapat melakukan lokalisasi.

Contoh ringkas user:

```json
{
  "id": "usr_123",
  "fullName": "Dr. Budi Santoso",
  "email": "budi@example.com",
  "role": "Akademisi",
  "affiliation": "Universitas Nalakarsa",
  "avatarUrl": null
}
```

Contoh item diskusi:

```json
{
  "id": "dsc_123",
  "title": "Kolaborasi Riset AI untuk Pertanian",
  "body": "Isi diskusi",
  "category": "Teknologi",
  "author": {
    "id": "usr_123",
    "fullName": "Dr. Budi Santoso",
    "role": "Akademisi",
    "avatarUrl": null
  },
  "replyCount": 24,
  "upvoteCount": 10,
  "hasUpvoted": false,
  "createdAt": "2026-07-19T08:00:00.000Z"
}
```

Contoh item percakapan:

```json
{
  "id": "cnv_123",
  "participant": {
    "id": "usr_456",
    "fullName": "Siti Rahma, M.Sc",
    "role": "Profesional",
    "avatarUrl": null
  },
  "lastMessage": {
    "id": "msg_123",
    "body": "Proposal sudah saya periksa.",
    "senderId": "usr_456",
    "createdAt": "2026-07-19T08:00:00.000Z"
  },
  "unreadCount": 1
}
```

## 6. Keamanan dan Kualitas

- [ ] Password tidak pernah dikembalikan oleh API atau dicatat di log.
- [ ] Refresh token disimpan sebagai hash; cookie memakai `HttpOnly`, `Secure`, dan `SameSite` yang sesuai jika memilih cookie.
- [ ] Terapkan rate limiting dan proteksi brute force.
- [ ] Validasi authorization pada level resource, bukan hanya memastikan pengguna sudah login.
- [ ] Sanitasi teks dan cegah stored XSS pada diskusi, profil, proyek, dan chat.
- [ ] Validasi upload berdasarkan MIME sebenarnya, ukuran, extension, dan scanning bila diperlukan.
- [ ] Gunakan query terparameterisasi melalui ORM/query builder.
- [ ] Simpan secret hanya di environment/secret manager.
- [ ] Tambahkan audit log untuk login sensitif dan tindakan admin.
- [ ] Siapkan backup database dan prosedur restore.
- [ ] Tambahkan unit test untuk business rule dan integration test untuk API utama.
- [ ] Tambahkan test authorization untuk mencegah akses data milik pengguna lain.
- [ ] Tambahkan lint, type-check, migration check, dan test ke CI.

## 7. Integrasi Frontend yang Perlu Dikoordinasikan

- [ ] FE membuat `src/services/api.js` atau modul service per domain.
- [ ] Ganti pemeriksaan auth dari `localStorage.isLoggedIn` dengan validasi session/token dari backend.
- [ ] Ganti store `auth.js`, `discussion.js`, `network.js`, dan `chat.js` dari data lokal ke pemanggilan API.
- [ ] Tambahkan state `loading`, `error`, empty state, dan pagination di setiap store.
- [ ] Tambahkan interceptor untuk refresh token dan penanganan `401`.
- [ ] Gunakan ID pengguna untuk membuka chat; jangan mencocokkan kontak berdasarkan nama.
- [ ] Sinkronkan enum role, kategori diskusi, status proyek, dan tipe notifikasi.
- [ ] Sepakati base URL melalui `VITE_API_BASE_URL`.
- [ ] Sepakati strategi WebSocket auth, reconnect, dan deduplikasi pesan.
- [ ] Hapus data dummy secara bertahap setelah setiap modul API stabil.

Catatan progres FE yang perlu ditangani saat integrasi:

- Halaman dasbor sudah mengharapkan `activeProjects`, tetapi sumber data proyek aktif belum terhubung.
- Statistik profil dan dasbor masih berupa angka statis.
- Kolaborasi saat ini langsung membuka chat berdasarkan nama inisiator; backend harus menggunakan `initiatorId`.
- Diskusi baru saat ini hanya memiliki judul dan kategori; BE sebaiknya meminta `body`/deskripsi agar detail diskusi bermakna.
- Chat belum persisten, belum memiliki unread count, dan belum mendukung realtime.

## 8. Urutan Sprint yang Disarankan

### Sprint BE 1 - Fondasi dan Auth

- [ ] Setup aplikasi, database, migration, seed, logging, dan Swagger.
- [ ] Register, login, refresh, logout, dan profil saya.
- [ ] Integrasi auth pertama dengan FE.

### Sprint BE 2 - Diskusi dan Jaringan

- [ ] CRUD diskusi, balasan, vote, pencarian, dan pagination.
- [ ] Koneksi, request masuk/keluar, dan rekomendasi pengguna.
- [ ] Notifikasi dasar untuk diskusi dan koneksi.

### Sprint BE 3 - Kolaborasi dan Dasbor

- [ ] CRUD proyek, pengajuan, anggota, milestone, dan progres.
- [ ] Endpoint agregasi dasbor dan statistik profil.
- [ ] Notifikasi proyek.

### Sprint BE 4 - Chat Realtime dan Hardening

- [ ] Conversation, message persistence, unread count, dan read receipt.
- [ ] WebSocket untuk chat dan notifikasi.
- [ ] Upload avatar/lampiran, security hardening, observability, dan load test dasar.

## 9. Definition of Done Backend

Sebuah fitur BE dianggap selesai jika:

- [ ] Endpoint dan business rule sudah diimplementasikan.
- [ ] Request dan response sesuai OpenAPI yang disepakati.
- [ ] Validasi input, autentikasi, dan authorization sudah diterapkan.
- [ ] Migration dan seed terkait tersedia.
- [ ] Unit/integration test utama lulus.
- [ ] Error case sudah terdokumentasi dan dapat ditangani FE.
- [ ] Tidak ada data sensitif di response atau log.
- [ ] Endpoint sudah diuji pada environment staging.
- [ ] FE berhasil mengintegrasikan happy path dan error state.
- [ ] Dokumentasi deployment/environment diperbarui.

## 10. Deliverable Awal untuk Tim BE

- [ ] Diagram ERD final.
- [ ] OpenAPI v1 untuk Auth, User/Profile, Discussion, Connection, Project, Chat, Dashboard, dan Notification.
- [ ] Daftar enum dan business rule yang disetujui bersama.
- [ ] Migration awal dan seed akun/demo.
- [ ] Postman/Bruno collection atau akses Swagger UI.
- [ ] File contoh environment tanpa secret.
- [ ] URL API development/staging.
- [ ] Panduan menjalankan backend secara lokal.
- [ ] Panduan migration, seed, test, dan deployment.
