package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Connection DTOs ---

type SendConnectionRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

type ConnectionUserResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"fullName"`
	Role      string    `json:"role"`
	Affiliation string  `json:"affiliation"`
	AvatarURL string    `json:"avatar_url"`
}

type ConnectionResponse struct {
	ID        uuid.UUID              `json:"id"`
	User      ConnectionUserResponse `json:"user"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
}

type ConnectionRequestResponse struct {
	ID        uuid.UUID              `json:"id"`
	User      ConnectionUserResponse `json:"user"` // requester for incoming, addressee for outgoing
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
}

type UserSuggestionResponse struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"fullName"`
	Role           string    `json:"role"`
	Affiliation    string    `json:"affiliation"`
	Expertise      string    `json:"expertise"`
	AvatarURL      string    `json:"avatar_url"`
	MutualCount    int64     `json:"mutual_count"`
}
