package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CollaborationRequest represents an application to collaborate on a discussion topic or project.
type CollaborationRequest struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID            *uuid.UUID     `gorm:"type:uuid;index"`
	DiscussionID         *uuid.UUID     `gorm:"type:uuid;index"`
	ApplicantID          uuid.UUID      `gorm:"type:uuid;index;not null"`
	ProposedContribution string         `gorm:"type:text;not null"`
	Status               string         `gorm:"type:varchar(20);default:'PENDING';not null"` // PENDING, ACCEPTED, REJECTED
	RejectionReason      *string        `gorm:"type:text"`
	CreatedAt            time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`

	// Relations
	Project    *Project    `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Discussion *Discussion `gorm:"foreignKey:DiscussionID;constraint:OnDelete:SET NULL"`
	Applicant  User        `gorm:"foreignKey:ApplicantID;constraint:OnDelete:CASCADE"`
}
