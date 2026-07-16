package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Application struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CollaborationID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_collab_user;not null"`
	UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_collab_user;index;not null"`
	Message         string         `gorm:"type:text;not null"`
	Status          string         `gorm:"type:varchar(50);default:'pending';not null"` // pending, accepted, rejected
	CreatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
