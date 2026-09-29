package email

import "embed"

// Templates berisi seluruh template HTML email yang digunakan oleh modul aplikasi.
// Modul pemilik fitur bertanggung jawab merender dan mengirim emailnya.
//
//go:embed templates/*.html
var Templates embed.FS
