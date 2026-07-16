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
}

type DiscussionCreator struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	AvatarURL   string    `json:"avatar_url"`
}

type DiscussionCommentResponse struct {
	ID        uuid.UUID         `json:"id"`
	Content   string            `json:"content"`
	CreatedAt time.Time         `json:"created_at"`
	Creator   DiscussionCreator `json:"creator"`
}

type DiscussionResponse struct {
	ID        uuid.UUID         `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Category  string            `json:"category"`
	Tags      string            `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
	Creator   DiscussionCreator `json:"creator"`
}

type DiscussionDetailResponse struct {
	ID        uuid.UUID                   `json:"id"`
	Title     string                      `json:"title"`
	Content   string                      `json:"content"`
	Category  string                      `json:"category"`
	Tags      string                      `json:"tags"`
	CreatedAt time.Time                   `json:"created_at"`
	Creator   DiscussionCreator           `json:"creator"`
	Comments  []DiscussionCommentResponse `json:"comments"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=3"`
}
