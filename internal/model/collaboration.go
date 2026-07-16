package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Collaboration struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null"`
	Title        string         `gorm:"type:varchar(255);not null"`
	Description  string         `gorm:"type:text;not null"`
	RoleRequired string         `gorm:"type:varchar(50);index;not null"` // akademisi, praktisi, profesional
	Status       string         `gorm:"type:varchar(50);index;default:'open';not null"` // open, in_progress, closed
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations
	User         User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Applications []Application `gorm:"foreignKey:CollaborationID;constraint:OnDelete:CASCADE"`
}
