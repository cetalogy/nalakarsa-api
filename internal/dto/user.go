package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	FullName    string    `json:"fullName"`
	PrefixTitle string    `json:"prefixTitle"`
	SuffixTitle string    `json:"suffixTitle"`
	FirstName   string    `json:"firstName"`
	MiddleName  *string   `json:"middleName"`
	LastName    string    `json:"lastName"`
	Affiliation string    `json:"affiliation"`
	Location    string    `json:"location"`
	Expertise   string    `json:"expertise"`
	Mission     string    `json:"mission"`
	Industry    string    `json:"industry"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserProfileStatsResponse adds aggregate stats for profile view
type UserProfileStatsResponse struct {
	UserResponse
	Stats ProfileStats `json:"stats"`
}

type ProfileStats struct {
	ConnectionCount int64 `json:"connection_count"`
	ProjectCount    int64 `json:"project_count"`
	DiscussionCount int64 `json:"discussion_count"`
}

type MyProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int64            `json:"total"`
}

type UserStatsResponse struct {
	ConnectionCount int64 `json:"connection_count"`
	ProjectCount    int64 `json:"project_count"`
	DiscussionCount int64 `json:"discussion_count"`
	ViewCount       int64 `json:"view_count"`
}

type UpdateProfileRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	MiddleName     string `json:"middleName"`
	LastName       string `json:"lastName" binding:"required"`
	FullName       string `json:"fullName"` // usually ignored in request, derived from parts
	PrefixTitle    string `json:"prefixTitle"`
	SuffixTitle    string `json:"suffixTitle"`
	Affiliation    string `json:"affiliation" binding:"required"`
	Location       string `json:"location" binding:"required"`
	Expertise      string `json:"expertise" binding:"required"`
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
	Data    interface{}         `json:"data,omitempty"`
	Meta    *PaginationResponse `json:"meta,omitempty"`
	Message string              `json:"message"`
}

// APIErrorResponse is the standard error envelope
type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	// Use `details` as the canonical error details field.
	Details map[string]string `json:"details,omitempty"`
}
