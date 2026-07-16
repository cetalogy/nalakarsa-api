package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discussion struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null"`
	Title     string         `gorm:"type:varchar(255);not null"`
	Content   string         `gorm:"type:text;not null"`
	Category  string         `gorm:"type:varchar(100);index;not null"`
	Tags      string         `gorm:"type:varchar(255);default:'';not null"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:DiscussionID;constraint:OnDelete:CASCADE"`
}
