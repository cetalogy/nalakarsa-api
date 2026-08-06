package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Project DTOs (replaces Collaboration) ---

type CreateProjectRequest struct {
	Title         string     `json:"title" form:"title" binding:"required,min=5"`
	Description   string     `json:"description" form:"description" binding:"required,min=15"`
	Category      string     `json:"category" form:"category" binding:"required"`
	RoleRequired  string     `json:"role_required" form:"role_required" binding:"omitempty,oneof=Akademisi Praktisi Profesional"`
	FundingStatus string     `json:"funding_status" form:"funding_status"`
	Location      string     `json:"location" form:"location"`
	Deadline      *time.Time `json:"deadline" form:"deadline"`
}

type UpdateProjectRequest struct {
	Title         string     `json:"title" form:"title" binding:"required,min=5"`
	Description   string     `json:"description" form:"description" binding:"required,min=15"`
	Category      string     `json:"category" form:"category" binding:"required"`
	Status        string     `json:"status" form:"status" binding:"required,oneof=draft open in_review active completed archived"`
	RoleRequired  string     `json:"role_required" form:"role_required" binding:"omitempty,oneof=Akademisi Praktisi Profesional"`
	FundingStatus string     `json:"funding_status" form:"funding_status"`
	Location      string     `json:"location" form:"location"`
	Deadline      *time.Time `json:"deadline" form:"deadline"`
	Progress      *int       `json:"progress" form:"progress" binding:"omitempty,min=0,max=100"`
}

type ProjectOwner struct {
	ID          uuid.UUID `json:"id" form:"id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	Afiliasi    string    `json:"afiliasi" form:"afiliasi"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

type ProjectResponse struct {
	ID            uuid.UUID    `json:"id" form:"id"`
	Title         string       `json:"title" form:"title"`
	Description   string       `json:"description" form:"description"`
	Category      string       `json:"category" form:"category"`
	Status        string       `json:"status" form:"status"`
	RoleRequired  string       `json:"role_required" form:"role_required"`
	FundingStatus string       `json:"funding_status" form:"funding_status"`
	Location      string       `json:"location" form:"location"`
	Deadline      *time.Time   `json:"deadline" form:"deadline"`
	Progress      int          `json:"progress" form:"progress"`
	CreatedAt     time.Time    `json:"created_at" form:"created_at"`
	Owner         ProjectOwner `json:"owner" form:"owner"`
}

type ProjectDetailResponse struct {
	ProjectResponse
	Members    []ProjectMemberResponse    `json:"members" form:"members"`
	Milestones []ProjectMilestoneResponse `json:"milestones" form:"milestones"`
}

// --- Project Application ---

type ApplyProjectRequest struct {
	Message string `json:"message" form:"message" binding:"required,min=10"`
}

type ApplicantResponse struct {
	ID          uuid.UUID `json:"id" form:"id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	Afiliasi    string    `json:"afiliasi" form:"afiliasi"`
	Lokasi      string    `json:"lokasi" form:"lokasi"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

type ProjectApplicationResponse struct {
	ID        uuid.UUID         `json:"id" form:"id"`
	Message   string            `json:"message" form:"message"`
	Status    string            `json:"status" form:"status"`
	CreatedAt time.Time         `json:"created_at" form:"created_at"`
	Applicant ApplicantResponse `json:"applicant" form:"applicant"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" form:"status" binding:"required,oneof=accepted rejected"`
}

// --- Project Member ---

type ProjectMemberResponse struct {
	ID          uuid.UUID `json:"id" form:"id"`
	UserID      uuid.UUID `json:"user_id" form:"user_id"`
	NamaLengkap string    `json:"nama_lengkap" form:"nama_lengkap"`
	Role        string    `json:"role" form:"role"`
	Status      string    `json:"status" form:"status"`
	JoinedAt    time.Time `json:"joined_at" form:"joined_at"`
	AvatarURL   string    `json:"avatar_url" form:"avatar_url"`
}

// --- Project Milestone ---

type CreateMilestoneRequest struct {
	Title      string     `json:"title" form:"title" binding:"required"`
	DueAt      *time.Time `json:"due_at" form:"due_at"`
	AssigneeID *uuid.UUID `json:"assignee_id" form:"assignee_id"`
}

type UpdateMilestoneRequest struct {
	Title      string     `json:"title" form:"title" binding:"required"`
	DueAt      *time.Time `json:"due_at" form:"due_at"`
	Status     string     `json:"status" form:"status" binding:"required,oneof=pending in_progress completed"`
	AssigneeID *uuid.UUID `json:"assignee_id" form:"assignee_id"`
}

type ProjectMilestoneResponse struct {
	ID          uuid.UUID  `json:"id" form:"id"`
	Title       string     `json:"title" form:"title"`
	DueAt       *time.Time `json:"due_at" form:"due_at"`
	Status      string     `json:"status" form:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id" form:"assignee_id"`
	CompletedAt *time.Time `json:"completed_at" form:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" form:"created_at"`
}
