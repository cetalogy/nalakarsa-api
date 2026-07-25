package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDiscussionRequest struct {
	Title    string `json:"title" binding:"required,min=5"`
	Content  string `json:"content" binding:"required,min=10"`
	Category string `json:"category" binding:"required"`
	Tags     string `json:"tags"`
}

type UpdateDiscussionRequest struct {
	Title    string `json:"title" binding:"required,min=5"`
	Content  string `json:"content" binding:"required,min=10"`
	Category string `json:"category" binding:"required"`
	Tags     string `json:"tags"`
	Status   string `json:"status" binding:"omitempty,oneof=open resolved closed"`
}

type DiscussionCreator struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	AvatarURL   string    `json:"avatar_url"`
}

type DiscussionReplyResponse struct {
	ID        uuid.UUID         `json:"id"`
	Content   string            `json:"content"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	Creator   DiscussionCreator `json:"creator"`
}

type DiscussionResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Category    string            `json:"category"`
	Tags        string            `json:"tags"`
	Status      string            `json:"status"`
	ReplyCount  int64             `json:"reply_count"`
	UpvoteCount int64             `json:"upvote_count"`
	HasUpvoted  bool              `json:"has_upvoted"`
	CreatedAt   time.Time         `json:"created_at"`
	Creator     DiscussionCreator `json:"creator"`
}

type DiscussionDetailResponse struct {
	ID          uuid.UUID                 `json:"id"`
	Title       string                    `json:"title"`
	Content     string                    `json:"content"`
	Category    string                    `json:"category"`
	Tags        string                    `json:"tags"`
	Status      string                    `json:"status"`
	ReplyCount  int64                     `json:"reply_count"`
	UpvoteCount int64                     `json:"upvote_count"`
	HasUpvoted  bool                      `json:"has_upvoted"`
	CreatedAt   time.Time                 `json:"created_at"`
	Creator     DiscussionCreator         `json:"creator"`
	Replies     []DiscussionReplyResponse `json:"replies"`
}

type CreateReplyRequest struct {
	Content  string     `json:"content" binding:"required,min=3"`
	ParentID *uuid.UUID `json:"parent_id"`
}
