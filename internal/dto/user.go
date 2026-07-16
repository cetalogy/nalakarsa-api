package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	NamaLengkap    string `json:"nama_lengkap"`
	GelarDepan     string `json:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi"`
	Lokasi         string `json:"lokasi"`
	BidangKeahlian string `json:"bidang_keahlian"`
	BioMisi        string `json:"bio_misi"`
	AvatarURL      string `json:"avatar_url"`
}

type UserProfileResponse struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	Role      string          `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	Profile   ProfileResponse `json:"profile"`
}

type UpdateProfileRequest struct {
	NamaLengkap    string `json:"nama_lengkap" binding:"required"`
	GelarDepan     string `json:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" binding:"required"`
	Lokasi         string `json:"lokasi" binding:"required"`
	BidangKeahlian string `json:"bidang_keahlian" binding:"required"`
	BioMisi        string `json:"bio_misi"`
	AvatarURL      string `json:"avatar_url"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}

type APIResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data,omitempty"`
	Errors     []string            `json:"errors,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}
