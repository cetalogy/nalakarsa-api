package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMilestone struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID  `gorm:"type:uuid;index;not null"`
	Title       string     `gorm:"type:varchar(255);not null"`
	DueAt       *time.Time `gorm:"type:timestamptz"`
	Status      string     `gorm:"type:varchar(20);default:'pending';not null"` // pending, in_progress, completed
	AssigneeID  *uuid.UUID `gorm:"type:uuid;index"`
	CompletedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Assignee *User `gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL"`
}
