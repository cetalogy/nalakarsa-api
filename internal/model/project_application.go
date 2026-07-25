package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectApplication struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_proj_applicant;index;not null"`
	ApplicantID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_proj_applicant;index;not null"`
	Message     string         `gorm:"type:text;not null"`
	Status      string         `gorm:"type:varchar(20);default:'pending';not null"` // pending, accepted, rejected
	CreatedAt   time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relations
	Applicant User `gorm:"foreignKey:ApplicantID;constraint:OnDelete:CASCADE"`
}
