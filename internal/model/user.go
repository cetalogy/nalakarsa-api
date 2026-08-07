package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash    string         `gorm:"type:varchar(255);not null"`
	Role            string         `gorm:"type:varchar(50);not null"`                          // akademisi, praktisi, profesional
	SystemRole      string         `gorm:"type:varchar(20);not null;default:'user'"`            // user, moderator, admin
	Status          string         `gorm:"type:varchar(20);not null;default:'active'"`          // active, suspended
	EmailVerifiedAt *time.Time     `gorm:"type:timestamptz"`
	CreatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Profile fields
	FirstName   string  `gorm:"type:varchar(100);not null"`
	MiddleName  *string `gorm:"type:varchar(100)"`
	LastName    string  `gorm:"type:varchar(100);not null"`
	FullName    string  `gorm:"type:varchar(255);not null"`
	PrefixTitle string  `gorm:"type:varchar(100);default:'';not null"`
	SuffixTitle string  `gorm:"type:varchar(100);default:'';not null"`
	Affiliation string  `gorm:"type:varchar(255);not null"`
	Location    string  `gorm:"type:varchar(255);not null"`
	Expertise   string  `gorm:"type:varchar(255);not null"`
	Industry    string  `gorm:"type:varchar(255);default:'';not null"`
	Bio         string  `gorm:"type:text;default:'';not null"`
	Mission     string  `gorm:"type:text;default:'';not null"`
	AvatarURL   string  `gorm:"type:varchar(255);default:'';not null"`
	ViewCount   int     `gorm:"type:int;default:0;not null"`
}
