package discussion

import (
	"time"

	"nalakarsa/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscussionReply struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DiscussionID uuid.UUID      `gorm:"type:uuid;index;not null"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null"`
	ParentID     *uuid.UUID     `gorm:"type:uuid;index"` // optional, for nested replies
	Content      string         `gorm:"type:text;not null"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations
	User user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
