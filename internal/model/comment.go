package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DiscussionID uuid.UUID      `gorm:"type:uuid;index;not null"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null"`
	Content      string         `gorm:"type:text;not null"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
