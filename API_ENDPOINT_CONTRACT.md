# Kontrak Endpoint API (Ditetapkan)

Mulai sekarang, kontrak backend yang harus dipakai FE adalah sebagai berikut:

## 1) Base URL

- Wajib pakai prefix: `/api/v1`.

## 2) Format error standar (tunggal)

Semua error response harus mengikuti envelope:

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "Validation failed",
    "details": {
      "email": "Email harus valid"
    }
  }
}
```

- `code`: kode error internal API.
- `message`: pesan ringkas.
- `details`: map informasi detail.

Jika sebelumnya ada `fields`, sekarang memakai `details`.

## 3) Endpoint homepage landing

- `GET /api/v1/homepage/hero`
- `GET /api/v1/homepage/sections`
- `GET /api/v1/homepage/testimonials` (opsional)

## 4) Catatan

- Implementasi helper error:
  - `internal/utils/response.go`
- DTO error:
  - `internal/dto/user.go`
- Dokumen detail referensi juga ada di:
  - `docs/md/FE/API_ENDPOINT_CONTRACT.md`
