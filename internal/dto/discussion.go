package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDiscussionRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description"`
	Excerpt     string `json:"excerpt"`
	Category    string `json:"category" binding:"required"`
	Tags        string `json:"tags"`
}

func (r *CreateDiscussionRequest) GetDescription() string {
	if r.Description != "" {
		return r.Description
	}
	return r.Excerpt
}

type UpdateDiscussionRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description"`
	Excerpt     string `json:"excerpt"`
	Category    string `json:"category" binding:"required"`
	Tags        string `json:"tags"`
	Status      string `json:"status" binding:"omitempty,oneof=open resolved closed"`
}

func (r *UpdateDiscussionRequest) GetDescription() string {
	if r.Description != "" {
		return r.Description
	}
	return r.Excerpt
}

type DiscussionCreator struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"fullName"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url"`
}

type DiscussionReplyResponse struct {
	ID        uuid.UUID         `json:"id"`
	Content   string            `json:"content"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	Creator   DiscussionCreator `json:"creator"`
}

type DiscussionRepliesData struct {
	DiscussionID uuid.UUID                 `json:"discussion_id"`
	TopicTitle   string                    `json:"topic_title"`
	TotalReplies int64                     `json:"total_replies"`
	Replies      []DiscussionReplyResponse `json:"replies"`
}

type DiscussionResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Excerpt            string     `json:"excerpt,omitempty"`
	Category           string     `json:"category"`
	Tags               string     `json:"tags"`
	Status             string     `json:"status"`
	IsInCollaboration  bool       `json:"isInCollaboration"`
	Replies            int64      `json:"replies"`
	UpvoteCount        int64      `json:"upvote_count"`
	HasUpvoted         bool       `json:"has_upvoted"`
	Time               time.Time  `json:"time"`
	CreatedAt          time.Time  `json:"createdAt"`
	Author             string     `json:"author"`
	Role               string     `json:"role"`
	SourceDiscussionID *uuid.UUID `json:"sourceDiscussionId"`
}

type DiscussionDetailResponse struct {
	ID                 uuid.UUID                 `json:"id"`
	Title              string                    `json:"title"`
	Description        string                    `json:"description"`
	Excerpt            string                    `json:"excerpt,omitempty"`
	Category           string                    `json:"category"`
	Tags               string                    `json:"tags"`
	Status             string                    `json:"status"`
	IsInCollaboration  bool                      `json:"isInCollaboration"`
	Replies            int64                     `json:"replies"`
	UpvoteCount        int64                     `json:"upvote_count"`
	HasUpvoted         bool                      `json:"has_upvoted"`
	Time               time.Time                 `json:"time"`
	CreatedAt          time.Time                 `json:"createdAt"`
	Author             string                    `json:"author"`
	Role               string                    `json:"role"`
	SourceDiscussionID *uuid.UUID                `json:"sourceDiscussionId"`
	Creator            DiscussionCreator         `json:"creator"`
	RepliesList        []DiscussionReplyResponse `json:"repliesList"`
}

type CreateReplyRequest struct {
	Content  string     `json:"content" binding:"required,min=3"`
	ParentID *uuid.UUID `json:"parent_id"`
}
