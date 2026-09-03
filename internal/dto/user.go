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
	Bio         string    `json:"bio"`
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
	Total    int64             `json:"total"`
}

type UserStatsResponse struct {
	ConnectionCount int64 `json:"connection_count"`
	FollowersCount  int64 `json:"followers_count"`
	FollowingCount  int64 `json:"following_count"`
	ProjectCount    int64 `json:"project_count"`
	DiscussionCount int64 `json:"discussion_count"`
	ViewCount       int64 `json:"view_count"`
}

type ToggleFollowResponse struct {
	Message      string    `json:"message"`
	IsFollowing  bool      `json:"isFollowing"`
	TargetUserID uuid.UUID `json:"targetUserId"`
}

type FollowUserItemResponse struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"fullName"`
	Role        string    `json:"role"`
	Affiliation string    `json:"affiliation"`
	AvatarURL   string    `json:"avatar_url"`
	IsFollowing bool      `json:"isFollowing,omitempty"`
}

type UpdateProfileRequest struct {
	FirstName   *string `json:"firstName"`
	MiddleName  *string `json:"middleName"`
	LastName    *string `json:"lastName"`
	FullName    *string `json:"fullName"`
	PrefixTitle *string `json:"prefixTitle"`
	SuffixTitle *string `json:"suffixTitle"`
	Affiliation *string `json:"affiliation"`
	Location    *string `json:"location"`
	Expertise   *string `json:"expertise"`
	Industry    *string `json:"industry"`
	Bio         *string `json:"bio" binding:"max=500"`
	AvatarURL   *string `json:"avatar_url"`
}

type PaginationResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}

// APIResponse is the standard response envelope
type APIResponse struct {
	Data      interface{}         `json:"data,omitempty"`
	Meta      *PaginationResponse `json:"meta,omitempty"`
	Message   string              `json:"message"`
	Timestamp string              `json:"timestamp,omitempty"`
}

// APIErrorResponse is the standard error envelope
type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp string            `json:"timestamp,omitempty"`
}
