# Roadmap & Future Tasks - Nalakarsa Backend

Daftar berikut adalah rekomendasi roadmap pengembangan selanjutnya untuk meningkatkan sistem Nalakarsa Backend dari prototipe menjadi produk berskala besar:

---

## Short-Term Tasks (Pengembangan Dekat)

- [ ] **Rate Limiting Middleware**: Tambahkan middleware pembatas request (misal: tollbooth / go-redis-rate) untuk melindungi API dari serangan brute force pada endpoint login/register.
- [ ] **Penyimpanan Media (Profile Avatar)**: Hubungkan API untuk pengunggahan foto profil langsung ke cloud object storage (seperti AWS S3, Cloudinary, atau MinIO) daripada sekadar menyimpan URL mentah.
- [ ] **Email Verification**: Implementasikan modul verifikasi email pendaftaran dan fitur reset password via token OTP/link email (menggunakan SMTP / Sendgrid).

---

## Medium-Term Tasks (Skala Menengah)

- [ ] **Redis Caching**: Integrasikan Redis Cache pada list direktori pengguna dan daftar diskusi terpopuler untuk mempercepat latency baca serta memangkas beban query PostgreSQL.
- [ ] **OAuth2 Social Login**: Tambahkan opsi pendaftaran dan login cepat satu-klik menggunakan akun pihak ketiga seperti Google dan GitHub.
- [ ] **Algoritma Matchmaking Pintar**: Kembangkan algoritma pencocokan cerdas yang otomatis mengirimkan rekomendasi tawaran kolaborasi berdasarkan kecocokan tag keahlian dan afiliasi pengguna.

---

## Long-Term Tasks (Skala Besar)

- [ ] **Real-time Chat & WebSockets**: Tambahkan fitur chat antar-pengguna langsung di dalam platform untuk mempermudah inisiasi kolaborasi setelah lamaran diterima.
- [ ] **Notifikasi Real-time**: Kirim push notification instan kepada pembuat proyek jika ada pelamar baru, dan kirim notifikasi jika topik diskusi dikomentari.
- [ ] **Pemisahan Unit & Integration Test Terisolasi**: Konfigurasikan CI/CD pipeline yang menjalankan unit test otomatis dan integration test menggunakan database test ephemeral (`testcontainers-go`).
- [ ] **Centralized Logging & Monitoring**: Hubungkan log `slog` JSON ke ElasticSearch/Kibana atau Prometheus/Grafana untuk mempermudah monitoring kesehatan server di production.
