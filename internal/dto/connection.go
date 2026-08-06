package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Connection DTOs ---

type SendConnectionRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" form:"target_user_id" binding:"required"`
}

type ConnectionUserResponse struct {
	ID          uuid.UUID `json:"id" form:"id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	Afiliasi    string    `json:"afiliasi" form:"afiliasi"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

type ConnectionResponse struct {
	ID        uuid.UUID              `json:"id" form:"id"`
	User      ConnectionUserResponse `json:"user" form:"user"`
	Status    string                 `json:"status" form:"status"`
	CreatedAt time.Time              `json:"created_at" form:"created_at"`
}

type ConnectionRequestResponse struct {
	ID        uuid.UUID              `json:"id" form:"id"`
	User      ConnectionUserResponse `json:"user" form:"user"` // requester for incoming, addressee for outgoing
	Status    string                 `json:"status" form:"status"`
	CreatedAt time.Time              `json:"created_at" form:"created_at"`
}

type UserSuggestionResponse struct {
	ID             uuid.UUID `json:"id" form:"id"`
	NamaLengkap    string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role           string    `json:"role" form:"role"`
	Afiliasi       string    `json:"afiliasi" form:"afiliasi"`
	BidangKeahlian string    `json:"bidang_keahlian" form:"bidang_keahlian"`
	AvatarURL      string    `json:"avatar_url" form:"avatar_url"`
	MutualCount    int64     `json:"mutual_count" form:"mutual_count"`
}
