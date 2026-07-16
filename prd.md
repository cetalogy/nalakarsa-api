# Product Requirements Document (PRD)

## 1. Ringkasan Produk
Sistem ini adalah halaman web prototipe ekosistem kolaborasi bernama **Nalakarsa**. Tujuannya adalah menghubungkan **Akademisi**, **Praktisi**, dan **Profesional** dalam satu platform untuk:
- Kolaborasi riset dan praktik
- Pendaftaran dan autentikasi pengguna
- Penyajian profil pengguna dan ruang diskusi
- Menampilkan fitur unggulan, testimonial, dan visualisasi pengalaman pengguna

## 2. Sasaran Produk
- Memudahkan pengguna memilih peran dan memahami fungsi platform
- Menyediakan proses pendaftaran dan login sederhana
- Menyajikan halaman profil dan status sesi pengguna aktif
- Menyediakan akses ke halaman diskusi dan fitur kolaborasi utama
- Menciptakan pengalaman visual yang menarik, responsif, dan intuitif

## 3. Pemangku Kepentingan
- Pengguna akhir: akademisi, praktisi, profesional, mahasiswa, peneliti
- Tim produk/desain: pemilik brand Nalakarsa
- Tim pengembang: front-end HTML/CSS/JS
- Sponsor / organisasi: pihak yang ingin memperkuat ekosistem kolaborasi

## 4. Persona Pengguna
1. **Dosen / Peneliti**
   - Ingin mempublikasikan ide riset dan mencari mitra praktik
2. **Pengusaha / Praktisi**
   - Ingin menemukan proyek riset yang dapat dikomersialkan
3. **Spesialis / Profesional**
   - Ingin berkontribusi dalam produk digital atau kolaborasi teknis

## 5. Fitur Utama
### 5.1 Landing Page / Home
- Intro splash screen dengan animasi teks dan efek visual
- Navigasi utama: Beranda, Diskusi, Kolaborasi
- Panel pilihan peran: Akademisi, Praktisi, Profesional
- Form login dan daftar dengan status peran terpilih
- Carousel profil anggota dan testimonial pengguna
- Tombol CTA untuk mendorong pendaftaran

### 5.2 Autentikasi & Akun
- Pendaftaran akun dengan data:
  - nama lengkap
  - gelar depan/belakang
  - afiliasi dan lokasi
  - bidang keahlian
  - bio/misi
  - email dan kata sandi
- Login dengan email dan password
- Penyimpanan akun dan sesi di `localStorage`
- Tampilan mode pengguna terautentikasi dalam halaman utama

### 5.3 Halaman Profil Pengguna
- Menampilkan data profil pengguna
- Foto profil, nama, jabatan, afiliasi, statistik ringkas
- Tombol logout
- Akses kembali ke profil setelah login

### 5.4 Halaman Diskusi
- Ruang diskusi terpisah di `diskusi.html`
- Menyediakan sistem filter dan pencarian sederhana
- Menampilkan daftar topik diskusi dan konten ringkas
- Visual ambient untuk pengalaman premium

### 5.5 Navigasi & Interaksi
- Navbar responsif dengan toggle hamburger di mobile
- Scroll smooth saat navigasi antar section
- Dropdown profil di navbar saat pengguna login
- Interaksi tombol CTA dan carousel

## 6. Alur Pengguna Utama
1. Pengguna membuka `homepage.html`
2. Melihat intro splash dan masuk ke halaman utama
3. Memilih peran terbaik
4. Melakukan pendaftaran atau login
5. Setelah login diarahkan ke `profil.html`
6. Mengakses fitur, testimonial, dan halaman diskusi
7. Bisa logout dan kembali ke status guest

## 7. Kebutuhan Fungsional
- FR1: Menampilkan splash intro saat pertama kali membuka halaman
- FR2: Memungkinkan pemilihan peran dengan visual terpilih
- FR3: Menyimpan data pendaftaran pengguna di browser
- FR4: Menyediakan login dengan validasi email/password
- FR5: Menampilkan halaman profil pengguna setelah login
- FR6: Menyediakan halaman diskusi dengan topik dan filter
- FR7: Navbar harus responsif di desktop dan mobile

## 8. Kebutuhan Non-Fungsional
- NFR1: Halaman tampil responsif pada semua ukuran layar
- NFR2: Desain konsisten dengan brand warna emas-hijau
- NFR3: Animasi dan transisi berjalan halus
- NFR4: Data sesi pengguna dikelola secara lokal di browser
- NFR5: Performa halaman tetap ringan dan tidak berat

## 9. Batasan Saat Ini
- Autentikasi hanya berbasis browser lokal (`localStorage`), belum backend
- Data diskusi dan profil sebagian besar statis
- Fitur kolaborasi belum tersambung ke sistem real-time
- Tidak ada manajemen hak akses server-side

## 10. Kriteria Sukses
- Pengguna dapat mendaftar dan login tanpa error
- Pengguna melihat profil dan akses diskusi berfungsi
- Role selection dan CTA berjalan sesuai ekspektasi
- Tampilan tetap rapi di layar kecil dan besar
- Sistem memberikan kesan ekosistem kolaborasi profesional

## 11. Rekomendasi Pengembangan Selanjutnya
- Integrasi backend untuk autentikasi dan database pengguna
- Diskusi dinamis dengan posting dan komentar
- Dashboard proyek kolaborasi terintegrasi
- Notifikasi real-time dan pesan antar pengguna
- Manajemen peran dan hak akses pengguna