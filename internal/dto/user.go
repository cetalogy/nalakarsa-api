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
	Industry       string `json:"industry"`
	Bio            string `json:"bio"`
	Mission        string `json:"mission"`
	AvatarURL      string `json:"avatar_url"`
	ViewCount      int    `json:"view_count"`
}

type UserProfileResponse struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	Role      string          `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	Profile   ProfileResponse `json:"profile"`
}

// UserProfileStatsResponse adds aggregate stats for profile view
type UserProfileStatsResponse struct {
	UserProfileResponse
	Stats ProfileStats `json:"stats"`
}

type ProfileStats struct {
	ConnectionCount int64 `json:"connection_count"`
	ProjectCount    int64 `json:"project_count"`
	DiscussionCount int64 `json:"discussion_count"`
}

type UpdateProfileRequest struct {
	NamaLengkap    string `json:"nama_lengkap" binding:"required"`
	GelarDepan     string `json:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" binding:"required"`
	Lokasi         string `json:"lokasi" binding:"required"`
	BidangKeahlian string `json:"bidang_keahlian" binding:"required"`
	Industry       string `json:"industry"`
	Bio            string `json:"bio"`
	Mission        string `json:"mission"`
	AvatarURL      string `json:"avatar_url"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}

// APIResponse is the standard response envelope
type APIResponse struct {
	Data       interface{}         `json:"data,omitempty"`
	Meta       *PaginationResponse `json:"meta,omitempty"`
	Message    string              `json:"message"`
}

// APIErrorResponse is the standard error envelope
type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
