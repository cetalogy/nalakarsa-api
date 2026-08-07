package dto

import (
	"time"

	"github.com/google/uuid"
)

// --- Project DTOs (replaces Collaboration) ---

type CreateProjectRequest struct {
	Title         string     `json:"title" binding:"required,min=5"`
	Description   string     `json:"description" binding:"required,min=15"`
	Category      string     `json:"category" binding:"required"`
	Needs         string     `json:"needs" binding:"omitempty,oneof=Akademisi Praktisi Profesional"`
	FundingStatus string     `json:"funding_status"`
	Location      string     `json:"location"`
	Deadline      *time.Time `json:"deadline"`
}

type CreateCollaborationRequest struct {
	DiscussionID uuid.UUID `json:"discussionId" binding:"required"`
}

type UpdateProjectRequest struct {
	Title         string     `json:"title" binding:"required,min=5"`
	Description   string     `json:"description" binding:"required,min=15"`
	Category      string     `json:"category" binding:"required"`
	Status        string     `json:"status" binding:"required,oneof=draft open in_review active completed archived"`
	Needs         string     `json:"needs" binding:"omitempty,oneof=Akademisi Praktisi Profesional"`
	FundingStatus string     `json:"funding_status"`
	Location      string     `json:"location"`
	Deadline      *time.Time `json:"deadline"`
	Progress      *int       `json:"progress" binding:"omitempty,min=0,max=100"`
}

type ProjectResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	Category           string     `json:"category"`
	Status             string     `json:"status"`
	Needs              string     `json:"needs"`
	FundingStatus      string     `json:"funding_status"`
	Location           string     `json:"location"`
	Deadline           *time.Time `json:"deadline"`
	Progress           int        `json:"progress"`
	CreatedAt          time.Time  `json:"created_at"`
	Initiator          string     `json:"initiator"`
	SourceDiscussionID *uuid.UUID `json:"sourceDiscussionId"`
}

type ProjectDetailResponse struct {
	ProjectResponse
	Members    []ProjectMemberResponse    `json:"members"`
	Milestones []ProjectMilestoneResponse `json:"milestones"`
}

// --- Project Application ---

type ApplyProjectRequest struct {
	Message string `json:"message" binding:"required,min=10"`
}

type ApplicantResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"fullName"`
	Role     string    `json:"role"`
	Afiliasi string    `json:"afiliasi"`
	Lokasi   string    `json:"lokasi"`
	AvatarURL string    `json:"avatar_url"`
}

type ProjectApplicationResponse struct {
	ID        uuid.UUID         `json:"id"`
	Message   string            `json:"message"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Applicant ApplicantResponse `json:"applicant"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected"`
}

// --- Project Member ---

type ProjectMemberResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FullName  string    `json:"fullName"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	JoinedAt  time.Time `json:"joined_at"`
	AvatarURL string    `json:"avatar_url"`
}

// --- Project Milestone ---

type CreateMilestoneRequest struct {
	Title      string     `json:"title" binding:"required"`
	DueAt      *time.Time `json:"due_at"`
	AssigneeID *uuid.UUID `json:"assignee_id"`
}

type UpdateMilestoneRequest struct {
	Title      string     `json:"title" binding:"required"`
	DueAt      *time.Time `json:"due_at"`
	Status     string     `json:"status" binding:"required,oneof=pending in_progress completed"`
	AssigneeID *uuid.UUID `json:"assignee_id"`
}

type ProjectMilestoneResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	DueAt       *time.Time `json:"due_at"`
	Status      string     `json:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
