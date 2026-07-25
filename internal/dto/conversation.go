package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Conversation DTOs ---

type CreateDirectConversationRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

type ConversationParticipant struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	AvatarURL   string    `json:"avatar_url"`
}

type LastMessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	SenderID  uuid.UUID `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationResponse struct {
	ID          uuid.UUID                `json:"id"`
	Participant ConversationParticipant  `json:"participant"`
	LastMessage *LastMessageResponse     `json:"last_message"`
	UnreadCount int64                    `json:"unread_count"`
}

// --- Message DTOs ---

type SendMessageRequest struct {
	Body string `json:"body" binding:"required,min=1,max=5000"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	SenderID  uuid.UUID `json:"sender_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CursorPaginationResponse struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
