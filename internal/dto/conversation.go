package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type SendMessageRequest struct {
	Text               string `json:"text" binding:"max=5000"`
	Body               string `json:"body"`
	AttachmentPath     string `json:"attachment_path"`
	AttachmentName     string `json:"attachment_name"`
	AttachmentMimeType string `json:"attachment_mime_type"`
	AttachmentSize     int64  `json:"attachment_size"`
	AttachmentType     string `json:"attachment_type"`
}

type AttachmentUploadResponse struct {
	Path     string `json:"path"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

type MessageResponse struct {
	ID         uuid.UUID           `json:"id"`
	Sender     string              `json:"sender"`
	Text       string              `json:"text"`
	Time       time.Time           `json:"time"`
	Attachment *AttachmentResponse `json:"attachment,omitempty"`
}

type AttachmentResponse struct {
	Path     string `json:"path"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

type CursorPaginationResponse struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

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
	Text               string `json:"text"`
	Content            string `json:"content"`
	AttachmentPath     string `json:"attachment_path"`
	AttachmentName     string `json:"attachment_name"`
	AttachmentMimeType string `json:"attachment_mime_type"`
	AttachmentSize     int64  `json:"attachment_size"`
	AttachmentType     string `json:"attachment_type"`
}

type GroupMessageResponse struct {
	ID              uuid.UUID           `json:"id"`
	GroupChatID     uuid.UUID           `json:"groupChatId"`
	SenderID        *uuid.UUID          `json:"senderId,omitempty"`
	SenderName      string              `json:"senderName,omitempty"`
	SenderRole      string              `json:"senderRole,omitempty"`
	SenderAvatar    string              `json:"senderAvatar,omitempty"`
	IsSystemMessage bool                `json:"isSystemMessage"`
	Content         string              `json:"content"`
	Attachment      *AttachmentResponse `json:"attachment,omitempty"`
	CreatedAt       time.Time           `json:"createdAt"`
}
