package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDiscussionRequest struct {
	Title    string `json:"title" form:"title" binding:"required,min=5"`
	Content  string `json:"content" form:"content" binding:"required,min=10"`
	Category string `json:"category" form:"category" binding:"required"`
	Tags     string `json:"tags" form:"tags"`
}

type UpdateDiscussionRequest struct {
	Title    string `json:"title" form:"title" binding:"required,min=5"`
	Content  string `json:"content" form:"content" binding:"required,min=10"`
	Category string `json:"category" form:"category" binding:"required"`
	Tags     string `json:"tags" form:"tags"`
	Status   string `json:"status" form:"status" binding:"omitempty,oneof=open resolved closed"`
}

type DiscussionCreator struct {
	ID          uuid.UUID `json:"id" form:"id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

type DiscussionReplyResponse struct {
	ID        uuid.UUID         `json:"id" form:"id"`
	Content   string            `json:"content" form:"content"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty" form:"parent_id"`
	CreatedAt time.Time         `json:"created_at" form:"created_at"`
	Creator   DiscussionCreator `json:"creator" form:"creator"`
}

type DiscussionResponse struct {
	ID          uuid.UUID         `json:"id" form:"id"`
	Title       string            `json:"title" form:"title"`
	Content     string            `json:"content" form:"content"`
	Category    string            `json:"category" form:"category"`
	Tags        string            `json:"tags" form:"tags"`
	Status      string            `json:"status" form:"status"`
	ReplyCount  int64             `json:"reply_count" form:"reply_count"`
	UpvoteCount int64             `json:"upvote_count" form:"upvote_count"`
	HasUpvoted  bool              `json:"has_upvoted" form:"has_upvoted"`
	CreatedAt   time.Time         `json:"created_at" form:"created_at"`
	Creator     DiscussionCreator `json:"creator" form:"creator"`
}

type DiscussionDetailResponse struct {
	ID          uuid.UUID                 `json:"id" form:"id"`
	Title       string                    `json:"title" form:"title"`
	Content     string                    `json:"content" form:"content"`
	Category    string                    `json:"category" form:"category"`
	Tags        string                    `json:"tags" form:"tags"`
	Status      string                    `json:"status" form:"status"`
	ReplyCount  int64                     `json:"reply_count" form:"reply_count"`
	UpvoteCount int64                     `json:"upvote_count" form:"upvote_count"`
	HasUpvoted  bool                      `json:"has_upvoted" form:"has_upvoted"`
	CreatedAt   time.Time                 `json:"created_at" form:"created_at"`
	Creator     DiscussionCreator         `json:"creator" form:"creator"`
	Replies     []DiscussionReplyResponse `json:"replies" form:"replies"`
}

type CreateReplyRequest struct {
	Content  string     `json:"content" form:"content" binding:"required,min=3"`
	ParentID *uuid.UUID `json:"parent_id" form:"parent_id"`
}
