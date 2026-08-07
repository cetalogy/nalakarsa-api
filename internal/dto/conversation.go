package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Conversation DTOs ---

type CreateDirectConversationRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

type StartChatRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type ConversationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Avatar      string    `json:"avatar"`
	LastMessage string    `json:"lastMessage"`
}

// --- Message DTOs ---

type SendMessageRequest struct {
	Text string `json:"text" binding:"required,min=1,max=5000"`
	Body string `json:"body"` // Fallback for older clients if needed
}

type MessageResponse struct {
	ID     uuid.UUID `json:"id"`
	Sender string    `json:"sender"` // "me" or "them"
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}

type CursorPaginationResponse struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
