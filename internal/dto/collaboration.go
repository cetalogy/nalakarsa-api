package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCollaborationRequest struct {
	Title        string `json:"title" binding:"required,min=5"`
	Description  string `json:"description" binding:"required,min=15"`
	RoleRequired string `json:"role_required" binding:"required,oneof=akademisi praktisi profesional"`
}

type UpdateCollaborationRequest struct {
	Title        string `json:"title" binding:"required,min=5"`
	Description  string `json:"description" binding:"required,min=15"`
	RoleRequired string `json:"role_required" binding:"required,oneof=akademisi praktisi profesional"`
	Status       string `json:"status" binding:"required,oneof=open in_progress closed"`
}

type CollaborationOwner struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	Afiliasi    string    `json:"afiliasi"`
	AvatarURL   string    `json:"avatar_url"`
}

type CollaborationResponse struct {
	ID           uuid.UUID          `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	RoleRequired string             `json:"role_required"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	Owner        CollaborationOwner `json:"owner"`
}

type ApplyCollaborationRequest struct {
	Message string `json:"message" binding:"required,min=10"`
}

type ApplicantResponse struct {
	ID          uuid.UUID `json:"id"`
	NamaLengkap string    `json:"nama_lengkap"`
	Role        string    `json:"role"`
	Afiliasi    string    `json:"afiliasi"`
	Lokasi      string    `json:"lokasi"`
	AvatarURL   string    `json:"avatar_url"`
}

type CollaborationApplicationResponse struct {
	ID        uuid.UUID         `json:"id"`
	Message   string            `json:"message"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Applicant ApplicantResponse `json:"applicant"`
}

type UpdateApplicationRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected"`
}
