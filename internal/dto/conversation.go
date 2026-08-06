package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Conversation DTOs ---

type CreateDirectConversationRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" form:"target_user_id" binding:"required"`
}

type ConversationParticipant struct {
	ID          uuid.UUID `json:"id" form:"id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

type LastMessageResponse struct {
	ID        uuid.UUID `json:"id" form:"id"`
	Body      string    `json:"body" form:"body"`
	SenderID  uuid.UUID `json:"sender_id" form:"sender_id"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
}

type ConversationResponse struct {
	ID          uuid.UUID               `json:"id" form:"id"`
	Participant ConversationParticipant `json:"participant" form:"participant"`
	LastMessage *LastMessageResponse    `json:"last_message" form:"last_message"`
	UnreadCount int64                   `json:"unread_count" form:"unread_count"`
}

// --- Message DTOs ---

type SendMessageRequest struct {
	Body string `json:"body" form:"body" binding:"required,min=1,max=5000"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id" form:"id"`
	Body      string    `json:"body" form:"body"`
	SenderID  uuid.UUID `json:"sender_id" form:"sender_id"`
	Status    string    `json:"status" form:"status"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
}

type CursorPaginationResponse struct {
	NextCursor string `json:"next_cursor,omitempty" form:"next_cursor"`
	HasMore    bool   `json:"has_more" form:"has_more"`
}
