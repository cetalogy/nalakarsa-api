package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Role         string         `gorm:"type:varchar(50);not null"` // akademisi, praktisi, profesional
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations
	Profile Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
