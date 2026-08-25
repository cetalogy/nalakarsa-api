package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Conversation DTOs ---

type CreateDirectConversationRequest struct {
	TargetUserID string `json:"target_user_id" binding:"required"`
}

type StartChatRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type ConversationResponse struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Role               string    `json:"role"`
	Avatar             string    `json:"avatar"`
	LastMessage        string    `json:"lastMessage"`
	LastMessageText    string    `json:"last_message,omitempty"`
	LastMessagePreview string    `json:"last_message_preview,omitempty"`
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

// --- Group Chat DTOs (FE Contract Specification) ---

type GroupChatMemberResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

type GroupChatResponse struct {
	ID              uuid.UUID                 `json:"id"`
	TopicID         *uuid.UUID                `json:"topicId,omitempty"`
	ProjectID       *uuid.UUID                `json:"projectId,omitempty"`
	Title           string                    `json:"title"`
	Badge           string                    `json:"badge"`
	LastMessage     *string                   `json:"lastMessage,omitempty"`
	LastMessageTime *time.Time                `json:"lastMessageTime,omitempty"`
	CreatedAt       time.Time                 `json:"createdAt"`
	Members         []GroupChatMemberResponse `json:"members"`
}

type SendGroupMessageRequest struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}

type GroupMessageResponse struct {
	ID              uuid.UUID  `json:"id"`
	GroupChatID     uuid.UUID  `json:"groupChatId"`
	SenderID        *uuid.UUID `json:"senderId,omitempty"`
	SenderName      string     `json:"senderName,omitempty"`
	SenderRole      string     `json:"senderRole,omitempty"`
	SenderAvatar    string     `json:"senderAvatar,omitempty"`
	IsSystemMessage bool       `json:"isSystemMessage"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"createdAt"`
}
