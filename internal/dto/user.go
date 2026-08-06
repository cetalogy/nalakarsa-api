package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	NamaLengkap    string `json:"nama_lengkap" form:"nama_lengkap"`
	GelarDepan     string `json:"gelar_depan" form:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang" form:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" form:"afiliasi"`
	Lokasi         string `json:"lokasi" form:"lokasi"`
	BidangKeahlian string `json:"bidang_keahlian" form:"bidang_keahlian"`
	Industry       string `json:"industry" form:"industry"`
	Bio            string `json:"bio" form:"bio"`
	Mission        string `json:"mission" form:"mission"`
	AvatarURL      string `json:"avatar_url" form:"avatar_url"`
	ViewCount      int    `json:"view_count" form:"view_count"`
}

type UserProfileResponse struct {
	ID        uuid.UUID       `json:"id" form:"id"`
	Email     string          `json:"email" form:"email"`
	Role      string          `json:"role" form:"role"`
	CreatedAt time.Time       `json:"created_at" form:"created_at"`
	Profile   ProfileResponse `json:"profile" form:"profile"`
}

// UserProfileStatsResponse adds aggregate stats for profile view
type UserProfileStatsResponse struct {
	UserProfileResponse
	Stats ProfileStats `json:"stats" form:"stats"`
}

type ProfileStats struct {
	ConnectionCount int64 `json:"connection_count" form:"connection_count"`
	ProjectCount    int64 `json:"project_count" form:"project_count"`
	DiscussionCount int64 `json:"discussion_count" form:"discussion_count"`
}

type UpdateProfileRequest struct {
	NamaLengkap    string `json:"nama_lengkap" form:"nama_lengkap" binding:"required"`
	GelarDepan     string `json:"gelar_depan" form:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang" form:"gelar_belakang"`
	Afiliasi       string `json:"afiliasi" form:"afiliasi" binding:"required"`
	Lokasi         string `json:"lokasi" form:"lokasi" binding:"required"`
	BidangKeahlian string `json:"bidang_keahlian" form:"bidang_keahlian" binding:"required"`
	Industry       string `json:"industry" form:"industry"`
	Bio            string `json:"bio" form:"bio"`
	Mission        string `json:"mission" form:"mission"`
	AvatarURL      string `json:"avatar_url" form:"avatar_url"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page" form:"current_page"`
	TotalPages  int   `json:"total_pages" form:"total_pages"`
	TotalItems  int64 `json:"total_items" form:"total_items"`
	Limit       int   `json:"limit" form:"limit"`
}

// APIResponse is the standard response envelope
type APIResponse struct {
	Data    interface{}         `json:"data,omitempty" form:"data"`
	Meta    *PaginationResponse `json:"meta,omitempty" form:"meta"`
	Message string              `json:"message" form:"message"`
}

// APIErrorResponse is the standard error envelope
type APIErrorResponse struct {
	Error APIErrorDetail `json:"error" form:"error"`
}

type APIErrorDetail struct {
	Code    string            `json:"code" form:"code"`
	Message string            `json:"message" form:"message"`
	Fields  map[string]string `json:"fields,omitempty" form:"fields"`
}
